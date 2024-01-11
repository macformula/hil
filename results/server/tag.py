"""
MACFE HIL Test Tag
"""

""" 
TODO: Check the type inside of the __test_<comp_op> functions
      and call a different function __test_<comp_op>_array for 
      array types.
"""


class Tag:
    def __init__(
        self,
        description: str,
        type: str,
        comp_op: str,
        lower_limit: any = 0,
        upper_limit: any = 0,
        expected_val: any = 0,
        unit: str = "Unitless",
    ) -> None:
        self.description = description
        self.comp_op = comp_op
        self.upper_limit = upper_limit
        self.lower_limit = lower_limit
        self.expected_val = expected_val
        self.type = type
        self.unit = unit

    def is_passing(self, value: any) -> bool:
        "Determines whether the value passes the given test tag"
        if self.comp_op == "GELE":
            return self.__test_gele(value)
        elif self.comp_op == "GTLT":
            return self.__test_gtlt(value)
        elif self.comp_op == "EQ":
            return self.__test_equal(value)
        elif self.comp_op == "GT":
            return self.__test_gt(value)
        elif self.comp_op == "GE":
            return self.__test_ge(value)
        elif self.comp_op == "LT":
            return self.__test_lt(value)
        elif self.comp_op == "LE":
            return self.__test_le(value)
        elif self.comp_op == "LOG":
            return True
        else:
            return False

    def __test_gele(self, value: any) -> bool:
        """Greater than or equal to and less than or equal to"""
        return (value >= self.lower_limit) and (value <= self.upper_limit)

    def __test_gtlt(self, value: any):
        """Greater than and less than"""
        return (value > self.lower_limit) and (value < self.upper_limit)

    def __test_equal(self, value: any):
        """Check if the value is equal to the expected value"""
        return value == self.expected_val

    def __test_gt(self, value: any):
        """Check if the value is greater than the lower limit"""
        return value > self.lower_limit

    def __test_ge(self, value: any):
        """Check if the value is greater than or equal to the lower limit"""
        return value >= self.lower_limit

    def __test_lt(self, value: any):
        """Check if the value is less than the upper limit"""
        return value < self.upper_limit

    def __test_le(self, value: any):
        """Check if the value is less than or equal to the upper limit"""
        return value <= self.upper_limit

    def __str__(self) -> str:
        return f"""
            Description: {self.description}
            Comparison Operator: {self.comp_op}
            Upper Limit: {self.upper_limit}
            Lower Limit: {self.lower_limit}
            Expected Value: {self.expected_val}
            Type: {self.type}
            Unit: {self.unit}
        """
