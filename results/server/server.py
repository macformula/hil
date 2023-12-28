import grpc
from concurrent import futures
from proto import results_pb2_grpc


class TagTunnel(results_pb2_grpc.TagTunnel):
    def SubmitTag(self, request, context):
        pass

    def CompleteTest(self, request, context):
        pass


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
