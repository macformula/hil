from generated import results_pb2_grpc
from generated import results_pb2
from result_accumulator import ResultAccumulator

class TagTunnel(results_pb2_grpc.TagTunnel):
    def __init__(self, result_accumulator: ResultAccumulator):
        self.ra = result_accumulator
    
    def CompleteTest(self, request, context):
        test_passed = self.ra.generate_and_run_tests(request.test_id)

        return results_pb2.CompleteTestResponse(test_passed=test_passed)

    def EnumerateErrors(self, request, context):
        errors = self.ra.get_all_errors()

        return results_pb2.EnumerateErrorsResponse(errors=errors)

    def EnumerateTags(self, request, context):
        tags = self.ra.get_all_tags()

        # Create a list to store Tag messages
        proto_tags = []

        for tag_id, tag in tags.items():
            # Convert each Tag instance to a Tag message
            proto_tag = results_pb2.Tag(
                tag_id=tag_id,
                description=tag.description,
                comp_operator=tag.comp_op,
                upper_limit=tag.upper_limit,
                lower_limit=tag.lower_limit,
            )

            # Determine the expected value based on the type attribute
            if tag.comp_op == "EQ":
                if tag.type == "string":
                    proto_tag.expected_val_str = str(tag.expected_val)
                elif "int" in tag.type:
                    proto_tag.expected_val_int = int(tag.expected_val)
                elif "float" in tag.type:
                    proto_tag.expected_val_float = float(tag.expected_val)
                elif tag.type == "bool":
                    proto_tag.expected_val_bool = bool(tag.expected_val)
                else:
                    proto_tag.expected_val_str = "ERROR: Unexpected Type"
            else:
               proto_tag.expected_val_str = "" 

            # Add the Tag message to the list
            proto_tags.append(proto_tag)

        return results_pb2.EnumerateTagsResponse(tags=proto_tags)

    def SubmitError(self, request, context):
        error_count = self.ra.submit_error(request.error)
        
        return results_pb2.SubmitErrorResponse(error_count=error_count)

    def SubmitTag(self, request, context):
        # get the value of the oneof entry that is actually set
        value = getattr(request, request.WhichOneof("data"))
        is_passing, e = self.ra.submit_tag(request.tag, value)
        if e != None:
            return results_pb2.SubmitTagResponse(success=False, error=str(e), is_passing=is_passing)
        else:
            return results_pb2.SubmitTagResponse(success=True, error="", is_passing=is_passing)
    

    