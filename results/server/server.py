import grpc
from concurrent import futures
from proto import results_pb2_grpc
from proto import results_pb2
from result_accumulator import ResultAccumulator


class TagTunnel(results_pb2_grpc.TagTunnel):
    def __init__(self):
        self.ra = ResultAccumulator(
            "./results/server/good_tags.yaml",
            "./results/server/schema/tags_schema.json")

    def SubmitTag(self, request, context):
        # get the value of the oneof entry that is actually set
        value = getattr(request, request.WhichOneof("data"))
        success, _ = self.ra.submit_tag(request.tag, value)
        return results_pb2.SubmitTagReply(success=success)

    def CompleteTest(self, request, context):
        self.ra.generate_and_run_tests()
        return results_pb2.CompleteTestReply(success=True)


def serve():
    port = "8080"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    results_pb2_grpc.add_TagTunnelServicer_to_server(TagTunnel(), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
