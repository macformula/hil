import grpc
from concurrent import futures
from generated import results_pb2_grpc
from tag_tunnel import TagTunnel

tags_fp="./results/server/tags.yaml"
tags_schema_fp="./results/server/schema/tags_schema.json"
historic_tests_fp="./results/server/historic_tests.json"



def serve():
    port = "8080"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    results_pb2_grpc.add_TagTunnelServicer_to_server(TagTunnel(tags_fp=tags_fp, 
                                                               tags_schema_fp=tags_schema_fp, 
                                                               historic_tests_fp=historic_tests_fp), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    serve()