from result_accumulator import ResultAccumulator


class RATest:
    """Singleton type test class for ResultAccumulator
    This takes place of the main process (submit tags and run tests)
    Note this same class is used within the jinja template to refer to the
    same submissions here
    TODO: submit tags and test them through gRPC"""

    _ra = None

    @classmethod
    def get_ra(cls):
        if cls._ra is None:
            cls.initialize_ra()
        return cls._ra

    @classmethod
    def initialize_ra(cls):
        TAG_FILE_PATH = "./results/server/good_tags.yaml"
        SCHEMA_FILE_PATH = "./results/server/schema/tags_schema.json"

        cls._ra = ResultAccumulator(TAG_FILE_PATH, SCHEMA_FILE_PATH)

        b, err = cls._ra.submit_tag("PV003", "Hello")

        # Should fail
        b, err = cls._ra.submit_tag("PV001", 1)

        # Should pass
        # b, err = cls._ra.submit_tag("PV001", 96)


if __name__ == "__main__":
    RATest.initialize_ra()
    ra = RATest.get_ra()
    ra.generate_and_run_tests()
