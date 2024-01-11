# Project Name

## Important Notice Regarding Generated gRPC Python Code

### Issue Description

If you encounter an import error related to `results_pb2_grpc.py` in the generated gRPC Python code (specifically on line 5), you may need to replace that line with the following:

```python
from . import results_pb2 as results__pb2
```

This is a workaround fix for now. For more details please refer to [Protobuf Issue #1491](https://github.com/protocolbuffers/protobuf/issues/1491).

