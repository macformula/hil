import pytest
from jinja2 import Template
from tag import Tag
from repo_handler import RepoHandler
from typing import Union
import yaml
import jsonschema
import json
import datetime
import os


class ResultAccumulator:
    def __init__(
        self,
        tags_fp: str,
        tags_schema_fp: str,
        template_fp: str,
        historic_tests_fp: str,
        temp_test_fp: str,
        reports_dir: str,
        repo_handler: RepoHandler,
    ) -> None:
        # Generate tag database from yaml file
        self.__parse_args(tags_fp, tags_schema_fp)
        self.tag_submissions: dict[str, any] = {}  # tag_id -> value cache
        self.error_submissions: list[str] = []  # list of cached errors
        self.historic_tests_fp = historic_tests_fp
        self.template_fp = template_fp
        self.temp_test_fp = temp_test_fp
        self.reports_dir = reports_dir
        self.repo_handler = repo_handler
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

        with open(self.temp_test_fp, "w") as file:
            file.write(test_file_content)

        dt = datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
        html_file_name = f"hil_report_{dt}.html"
        html_cli_arg = f"--html={self.reports_dir}/{html_file_name}"
        pytest.main(
            ["-v", "--showlocals", "--durations=10", html_cli_arg, self.temp_test_fp],
            plugins=[self],
        )

        os.remove(self.temp_test_fp)

        return dt

    def generate_and_run_tests(self, test_id: str, sequence_name: str, push_to_github: bool) -> bool:
        """On CompleteTest submissions
        - Wrapper for generate_test_file and pytest.main.
        the full RA tests w/ reports are run with this one

        returns overall pass/fail"""
        tag_ids = list(self.tag_submissions.keys())
        date_time = self.__generate_test_file(tag_ids)

        has_errors = len(self.error_submissions) > 0
        overall_pass_fail = self.all_tags_passing and (not has_errors)

        self.__update_historic_tests(test_id, sequence_name, date_time, overall_pass_fail)
        if push_to_github:
            self.repo_handler.push_to_github_pages(test_id)

        # reset cached submissions
        self.tag_submissions = {}
        self.error_submissions = []
        self.all_tags_passing = True

        return overall_pass_fail

    def __update_historic_tests(self, test_id: str, sequence_name: str, dt: str, test_passed: bool) -> None:
        ymd, hms = dt.split("_")

        # Replace the following line with your logic to determine test pass/fail
        new_test = {
            "testId": test_id,
            "sequenceName": sequence_name,
            "testPassed": test_passed,
            "date": ymd,
            "time": hms,
        }

        # Load existing tests from the JSON file
        with open(self.historic_tests_fp, "r") as file:
            existing_tests = json.load(file)

        # Prepend the new test to the existing tests
        existing_tests.insert(0, new_test)

        # Save the updated list back to the JSON file
        with open(self.historic_tests_fp, "w") as file:
            json.dump(existing_tests, file, indent=2)
