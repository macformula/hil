"""
Tests individual tag

Flags:
    - <tag-name> <value>
"""
import pytest
import sys
import yaml


class Tag:
    def __init__(self) -> None:
        self.__parse_args()
        self.run_tests()

    def test_less_than(self):
        assert self.value < self.less_than

    def test_greater_than(self):
        assert self.value > self.greater_than

    def __parse_args(self):
        if len(sys.argv) < 1:
            raise Exception("Tag and value argument not specified")
        if len(sys.argv) < 2:
            raise Exception("Value argument not specified")

        tag_name = sys.argv[1]

        with open("server/tags.yml") as f:
            tags_yaml = yaml.load(f, Loader=yaml.FullLoader)
            print(tags_yaml)
            # print(tag_name)
        if tag_name not in tags_yaml:
            raise Exception("Invalid Tag name not found in yaml")
        else:
            print("Parsing Tag Parameters..")

            self.type = tags_yaml[tag_name]["type"]
            self.value = sys.argv[2]
            try:
                self.value = eval(f"{self.type}({self.value})")
            except Exception as e:
                print(f"Error converting value to {self.type}: {e}")

            self.less_than = tags_yaml[tag_name]["lessThan"]
            self.greater_than = tags_yaml[tag_name]["greaterThan"]
            self.func_name = tags_yaml[tag_name]["function_name"]

    def run_tests(self):
        pytest.main()
        # pass
        self.test_greater_than()
        self.test_less_than()
        # pytest.main()

        # test_class = type('TestDynamic', (TestTag,), {})
        # setattr(test_class, 'test_dynamic', lambda self: None)
        # pytest.main(['-k', 'test_dynamic'])

    def __str__(self) -> str:
        return f"""less than: {self.less_than}\nGreater than:{self.greater_than}"""


if __name__ == "__main__":
    t = Tag()
    print(t)
