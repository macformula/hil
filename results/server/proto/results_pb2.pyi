from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SubmitTagRequest(_message.Message):
    __slots__ = ["tag", "value_str", "value_int", "value_float", "value_bool"]
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
    def __init__(
        self,
        tag: _Optional[str] = ...,
        value_str: _Optional[str] = ...,
        value_int: _Optional[int] = ...,
        value_float: _Optional[float] = ...,
        value_bool: bool = ...,
    ) -> None: ...

class SubmitTagReply(_message.Message):
    __slots__ = ["success"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    def __init__(self, success: bool = ...) -> None: ...

class CompleteTestReply(_message.Message):
    __slots__ = ["success"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    def __init__(self, success: bool = ...) -> None: ...
