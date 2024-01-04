import pytest
import git
from jinja2 import Template
from tag import Tag
from typing import Union
import yaml
import jsonschema
import json
import datetime
import os

class ResultAccumulator:
    def __init__(self, 
                 tags_fp: str, 
                 tags_schema_fp: str, 
                 template_fp: str,
                 historic_tests_fp: str, 
                 reports_dir: str,
                 pages_repo_dir: str,
                 pages_branch: str,
                 git_username: str,
                 git_email: str) -> None:
        # Generate tag database from yaml file
        self.__parse_args(tags_fp, tags_schema_fp)
        self.tag_submissions: dict[str, any] = {}  # tag_id -> value cache
        self.error_submissions: list[str] = []  # list of cached errors
        self.historic_tests_fp = historic_tests_fp
        self.template_fp = template_fp
        self.reports_dir = reports_dir
        self.pages_repo_dir = pages_repo_dir
        self.pages_branch = pages_branch
        self.git_username = git_username
        self.git_email = git_email
        self.all_tags_passing = True

    def submit_tag(self, tag_id: str, value: any) -> Union[bool, KeyError]:
        """On tag submissions"""
        try:
            is_passing: bool = self.tag_db[tag_id].is_passing(value)
            if not is_passing:
                self.all_tags_passing = False
            self.tag_submissions[tag_id] = value
            return is_passing, None
        except KeyError as e:
            self.all_tags_passing = False
            return False, e
        
    def get_all_tags(self):
        return self.tag_db
    
    def get_all_errors(self):
        return self.error_submissions

    def submit_error(self, error: str) -> int:
        """On error submissions"""
        self.error_submissions.append(error)

        return len(self.error_submissions)

    def __parse_args(self, tag_file_path: str, schema_file_path: str) -> None:
        self.tag_db: dict[str, Tag] = {}  # tag_id -> Tag
        with open(tag_file_path) as f:
            tags_yaml = yaml.load(f, Loader=yaml.FullLoader)

            self.__validate_tags(tag_file_path, schema_file_path)

            for tag_id, tag_info in tags_yaml.items():
                if tag_id in self.tag_db:
                    # TODO: throw exception here. Ideally this check is done
                    # inside of the JSON schema. could do this by making the
                    # yaml a list instead since JSON schema can check for
                    # unique entries only for arrays, not dicts (idk why)
                    return

                newTag = Tag(
                    description=tag_info.get("description", ""),
                    comp_op=tag_info.get("compareOp", ""),
                    type=tag_info.get("type", ""),
                    lower_limit=tag_info.get("lowerLimit", None),
                    upper_limit=tag_info.get("upperLimit", None),
                    expected_val=tag_info.get("expectedVal", None),
                    unit=tag_info.get("unit", ""),
                )

                self.tag_db[tag_id] = newTag

    def __validate_tags(self, tags_file_path, schema_file_path) -> bool:
        # Load YAML schema
        with open(schema_file_path, "r") as schema_file:
            schema = yaml.safe_load(schema_file)

        # Load YAML data
        with open(tags_file_path, "r") as tags_file:
            tags_data = yaml.safe_load(tags_file)

        # Validate against the schema
        jsonschema.validate(tags_data, schema)

    @staticmethod
    def load_template_from_file(file_path) -> Template:
        """Loads a jinja template from a file"""
        with open(file_path, "r") as file:
            template_content = file.read()
        return Template(template_content)

    def __generate_test_file(self, tag_ids: list) -> [str, str]:
        """Creates file from tag submissions, runs tests on it, then deletes it"""
        template = self.load_template_from_file(self.template_fp)
        test_file_content = template.render(
            tag_db=self.tag_db,
            tag_submissions=self.tag_submissions,
            tag_ids=tag_ids,
            error_submissions=self.error_submissions,
        )

        test_file_path = "results/server/temp_test_file.py"
        with open(test_file_path, "w") as file:
            file.write(test_file_content)

        dt = datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
        html_file_name = f"hil_report_{dt}.html"
        html_cli_arg = f"--html={self.reports_dir}/{html_file_name}"
        pytest.main(
            ["-v", "--showlocals", "--durations=10", html_cli_arg, test_file_path],
            plugins=[self],
        )

        os.remove(test_file_path)

        return dt

    def generate_and_run_tests(self, test_id: str) -> bool:
        """On CompleteTest submissions
        - Wrapper for generate_test_file and pytest.main.
        the full RA tests w/ reports are run with this one

        returns overall pass/fail"""
        tag_ids = list(self.tag_submissions.keys())
        date_time = self.__generate_test_file(tag_ids)

        has_errors = len(self.error_submissions) > 0
        overall_pass_fail = self.all_tags_passing and (not has_errors)

        self.__update_historic_tests(test_id, date_time, overall_pass_fail)
        self.__push_to_github_pages(test_id)

        # reset cached submissions
        self.tag_submissions = {}
        self.error_submissions = []
        self.all_tags_passing = True

        return overall_pass_fail

    def __update_historic_tests(self, test_id: str, dt: str, test_passed: bool) -> None:
        # Replace the following line with your logic to determine test pass/fail
        new_test = {
            "testId": test_id,
            "testPassed": test_passed,
            "date": dt
        }

        # Load existing tests from the JSON file
        with open(self.historic_tests_fp, 'r') as file:
            existing_tests = json.load(file)

        # Prepend the new test to the existing tests
        existing_tests.insert(0, new_test)

        # Save the updated list back to the JSON file
        with open(self.historic_tests_fp, 'w') as file:
            json.dump(existing_tests, file, indent=2)

    def __push_to_github_pages(self, test_id) -> None:
        repo = git.Repo(self.pages_repo_dir)

        # Set Git user name and email for the repository
        repo.config_writer().set_value("user", "name", self.git_username).release()
        repo.config_writer().set_value("user", "email", self.git_email).release()

        # Add all changes to the index
        repo.git.add('*')

        # Commit changes
        repo.index.commit(test_id)

        # Fetch updates from the remote repository
        repo.remotes.origin.fetch()

        # Ensure we are on the branch that needs to be updated
        repo.heads[self.pages_branch].checkout()

        # Rebase changes
        repo.git.rebase('origin/' + self.pages_branch)

        # Push the changes, forcing the update to the remote branch
        repo.git.push()
        return
