"""
Validates tags from yaml using pytest
"""
import pytest
import yaml
import traceback

with open("server/tags.yml") as f:
    tags_yaml = yaml.load(f, Loader=yaml.FullLoader)


class TestTags:
    def __init__(self, tags):
        self.tags_yaml = tags_yaml
        try:
            self.validate_tags(self.tags_yaml)
        except Exception:
            print("Error in tags.yml")
            print(traceback.format_exc())

        self.run_tests()

    # @pytest.mark.parametrize("tag", tags.keys())
    # def test_tag(self, tag):
    #     self.__validate_tags(tags[tag])

    @staticmethod
    def validate_tags(tags: dict):
        for key, value in tags.items():
            if value["greaterThan"] > value["lessThan"]:
                raise ValueError("greaterThan value must be less than lessThan value")
            if value["type"] not in ["float", "int", "bool", "string"]:
                raise ValueError("type must be float, int, bool or string")

    @staticmethod
    def validate_tags2(tags):
        for key, value in tags.items():
            if value["greaterThan"] > value["lessThan"]:
                raise ValueError("greaterThan value must be less than lessThan value")
            if value["greaterThan"] == value["lessThan"]:
                raise ValueError(
                    "greaterThan value must be different than lessThan value"
                )
            if value["type"] not in ["float", "int", "bool", "string"]:
                raise ValueError("type must be float, int, bool or string")
            if value["type"] == "bool":
                if value["greaterThan"] not in [0, 1]:
                    raise ValueError("greaterThan value must be 0 or 1")
                if value["lessThan"] not in [0, 1]:
                    raise ValueError("lessThan value must be 0 or 1")
            if value["type"] == "string":
                if not isinstance(value["greaterThan"], str):
                    raise ValueError("greaterThan value must be a string")
                if not isinstance(value["lessThan"], str):
                    raise ValueError("lessThan value must be a string")
            if value["type"] == "int":
                if not isinstance(value["greaterThan"], int):
                    raise ValueError("greaterThan value must be an integer")
                if not isinstance(value["lessThan"], int):
                    raise ValueError("lessThan value must be an integer")
            if value["type"] == "float":
                if not isinstance(value["greaterThan"], float):
                    raise ValueError("greaterThan value must be a float")
                if not isinstance(value["lessThan"], float):
                    raise ValueError("lessThan value must be a float")

    def run_tests(self):
        pytest.main()


if __name__ == "__main__":
    TestTags(tags_yaml)
