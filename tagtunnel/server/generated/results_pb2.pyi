from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SubmitTagRequest(_message.Message):
    __slots__ = ("tag", "value_str", "value_int", "value_float", "value_bool")
    TAG_FIELD_NUMBER: _ClassVar[int]
    VALUE_STR_FIELD_NUMBER: _ClassVar[int]
    VALUE_INT_FIELD_NUMBER: _ClassVar[int]
    VALUE_FLOAT_FIELD_NUMBER: _ClassVar[int]
    VALUE_BOOL_FIELD_NUMBER: _ClassVar[int]
    tag: str
    value_str: str
    value_int: int
    value_float: float
    value_bool: bool
    def __init__(self, tag: _Optional[str] = ..., value_str: _Optional[str] = ..., value_int: _Optional[int] = ..., value_float: _Optional[float] = ..., value_bool: bool = ...) -> None: ...

class SubmitTagResponse(_message.Message):
    __slots__ = ("success", "error", "is_passing")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    IS_PASSING_FIELD_NUMBER: _ClassVar[int]
    success: bool
    error: str
    is_passing: bool
    def __init__(self, success: bool = ..., error: _Optional[str] = ..., is_passing: bool = ...) -> None: ...

class CompleteTestRequest(_message.Message):
    __slots__ = ("test_id", "sequence_name", "push_report_to_github")
    TEST_ID_FIELD_NUMBER: _ClassVar[int]
    SEQUENCE_NAME_FIELD_NUMBER: _ClassVar[int]
    PUSH_REPORT_TO_GITHUB_FIELD_NUMBER: _ClassVar[int]
    test_id: str
    sequence_name: str
    push_report_to_github: bool
    def __init__(self, test_id: _Optional[str] = ..., sequence_name: _Optional[str] = ..., push_report_to_github: bool = ...) -> None: ...

class CompleteTestResponse(_message.Message):
    __slots__ = ("test_passed",)
    TEST_PASSED_FIELD_NUMBER: _ClassVar[int]
    test_passed: bool
    def __init__(self, test_passed: bool = ...) -> None: ...

class SubmitErrorRequest(_message.Message):
    __slots__ = ("error",)
    ERROR_FIELD_NUMBER: _ClassVar[int]
    error: str
    def __init__(self, error: _Optional[str] = ...) -> None: ...

class SubmitErrorResponse(_message.Message):
    __slots__ = ("error_count",)
    ERROR_COUNT_FIELD_NUMBER: _ClassVar[int]
    error_count: int
    def __init__(self, error_count: _Optional[int] = ...) -> None: ...

class EnumerateErrorsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class EnumerateErrorsResponse(_message.Message):
    __slots__ = ("errors",)
    ERRORS_FIELD_NUMBER: _ClassVar[int]
    errors: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, errors: _Optional[_Iterable[str]] = ...) -> None: ...

class EnumerateTagsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class EnumerateTagsResponse(_message.Message):
    __slots__ = ("tags",)
    TAGS_FIELD_NUMBER: _ClassVar[int]
    tags: _containers.RepeatedCompositeFieldContainer[Tag]
    def __init__(self, tags: _Optional[_Iterable[_Union[Tag, _Mapping]]] = ...) -> None: ...

class Tag(_message.Message):
    __slots__ = ("tag_id", "description", "comp_operator", "upper_limit", "lower_limit", "expected_val_str", "expected_val_int", "expected_val_float", "expected_val_bool")
    TAG_ID_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    COMP_OPERATOR_FIELD_NUMBER: _ClassVar[int]
    UPPER_LIMIT_FIELD_NUMBER: _ClassVar[int]
    LOWER_LIMIT_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_VAL_STR_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_VAL_INT_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_VAL_FLOAT_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_VAL_BOOL_FIELD_NUMBER: _ClassVar[int]
    tag_id: str
    description: str
    comp_operator: str
    upper_limit: float
    lower_limit: float
    expected_val_str: str
    expected_val_int: int
    expected_val_float: float
    expected_val_bool: bool
    def __init__(self, tag_id: _Optional[str] = ..., description: _Optional[str] = ..., comp_operator: _Optional[str] = ..., upper_limit: _Optional[float] = ..., lower_limit: _Optional[float] = ..., expected_val_str: _Optional[str] = ..., expected_val_int: _Optional[int] = ..., expected_val_float: _Optional[float] = ..., expected_val_bool: bool = ...) -> None: ...
