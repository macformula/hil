import sys
sys.path.append("./results/server")
from result_accumulator import ResultAccumulator

TAG_FILE_PATH = "./results/server/good_tags.yaml"
SCHEMA_FILE_PATH = "./results/server/schema/tags_schema.json"

class RATest:
    """Test class to run ResultAccumulator tests without server.
    Change jinja path to demo: TEMPLATE_FILE_PATH = "./results/server/demo/demo.py.jinja"""

    _ra = None

    @classmethod
    def get_ra(cls):
        if cls._ra is None:
            cls.initialize_ra()
        return cls._ra

    @classmethod
    def initialize_ra(cls):

        cls._ra = ResultAccumulator(TAG_FILE_PATH, SCHEMA_FILE_PATH)

        b, err = cls._ra.submit_tag("PV003", "Hello")

        # Should fail
        b, err = cls._ra.submit_tag("PV001", 1)
        # cls._ra.submit_error("Error1")
        # Should pass
        # b, err = cls._ra.submit_tag("PV001", 96)


if __name__ == "__main__":
    RATest.initialize_ra()
    ra = RATest.get_ra()
    ra.generate_and_run_tests()
