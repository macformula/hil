import pytest
from jinja2 import Template
from tag import Tag
from typing import Union
import yaml
import jsonschema
import datetime
import os

# File paths
TEMPLATE_FILE_PATH = "./results/server/tm.py.jinja"


class ResultAccumulator:
    def __init__(self, tag_file_path: str, schema_file_path: str) -> None:
        # Generate tag database from yaml file
        self.__parse_args(tag_file_path, schema_file_path)
        self.tag_submissions: dict[str, any] = {} # tag_id -> value cache
        self.error_submissions: list[str] = [] # list of cached errors

    def submit_tag(self, tag_id: str, value: any) -> Union[bool, KeyError]:
        """On tag submissions"""
        try:
            is_passing: bool = self.tag_db[tag_id].is_passing(value)
            self.tag_submissions[tag_id] = value
            return is_passing, None
        except KeyError as e:
            return False, e
        
    def submit_error(self, error: str) -> None:
        """On error submissions"""
        self.error_submissions.append(error)

    def __parse_args(self, tag_file_path: str, schema_file_path: str) -> None:
        self.tag_db: dict[str, Tag] = {} # tag_id -> Tag
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

    def __generate_test_file(self, tag_ids: list) -> None:
        """Creates file from tag submissions, runs tests on it, then deletes it"""
        template = self.load_template_from_file(TEMPLATE_FILE_PATH)
        test_file_content = template.render(
            tag_db=self.tag_db, tag_submissions=self.tag_submissions, tag_ids=tag_ids
        )

        test_file_path = "results/server/temp_test_file.py"
        with open(test_file_path, "w") as file:
            file.write(test_file_content)

        dt = datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
        html = f"--html=logs/pytest_report_{dt}.html"
        pytest.main(
            ["-v", "--showlocals", "--durations=10", test_file_path, html],
            plugins=[self],
        )

        os.remove(test_file_path)

    def generate_and_run_tests(self) -> bool:
        """On CompleteTest submissions
        - Wrapper for generate_test_file and pytest.main.
        the full RA tests w/ reports are run with this one

        returns overall pass/fail"""
        tag_ids = list(self.tag_submissions.keys())
        self.__generate_test_file(tag_ids)
        has_errors = len(self.error_submissions) > 0
        # reset cached submissions
        self.tag_submissions = {}
        self.error_submissions = []

        return has_errors
    
