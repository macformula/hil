import pytest
from jinja2 import Template
from tag import Tag
from typing import Union
import yaml
import jsonschema
import json
import datetime
import os

# File paths
TEMPLATE_FILE_PATH = "./results/server/tm.py.jinja"


class ResultAccumulator:
    def __init__(self, tag_fp: str, schema_fp: str, historic_tests_fp: str) -> None:
        # Generate tag database from yaml file
        self.__parse_args(tag_fp, schema_fp)
        self.tag_submissions: dict[str, any] = {}  # tag_id -> value cache
        self.error_submissions: list[str] = []  # list of cached errors
        self.historic_tests_fp = historic_tests_fp
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
        template = self.load_template_from_file(TEMPLATE_FILE_PATH)
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
        html_file_name = f"pytest_report_{dt}.html"
        html_cli_arg = f"--html=logs/{html_file_name}"
        pytest.main(
            ["-v", "--showlocals", "--durations=10", test_file_path, html_cli_arg],
            plugins=[self],
        )

        os.remove(test_file_path)

        return html_file_name, dt

    def generate_and_run_tests(self, test_id: str) -> bool:
        """On CompleteTest submissions
        - Wrapper for generate_test_file and pytest.main.
        the full RA tests w/ reports are run with this one

        returns overall pass/fail"""
        tag_ids = list(self.tag_submissions.keys())
        html_file_name, dt = self.__generate_test_file(tag_ids)
        
        has_errors = len(self.error_submissions) > 0
        overall_pass_fail = self.all_tags_passing and (not has_errors)

        self.__update_historic_tests(test_id, dt)
        self.__generate_index_html(html_file_name)
        self.__commit_to_github_pages()

        # reset cached submissions
        self.tag_submissions = {}
        self.error_submissions = []
        self.all_tags_passing = True

        return overall_pass_fail

    def __update_historic_tests(self, test_id: str, dt: str) -> None:
        # Replace the following line with your logic to determine test pass/fail
        new_test = {
            "testId": test_id,
            "testPassed": self.overall_pass_fail(),
            "date": dt
        }

        # Load existing tests from the JSON file
        with open(self.historic_test_fp, 'r') as file:
            existing_tests = json.load(file)

        # Append the new test to the existing tests
        existing_tests.append(new_test)

        # Save the updated list back to the JSON file
        with open(self.historic_test_fp, 'w') as file:
            json.dump(existing_tests, file, indent=2)

    def __generate_index_html(self, html_file_name):
        #TODO: Implement me @itshady
        return

    def __commit_to_github_pages(self):
        #TODO: Implement me
        return

    def overall_pass_fail(self) -> bool:
        # Implement me
        return False
