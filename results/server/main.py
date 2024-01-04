import grpc
from concurrent import futures
from result_accumulator import ResultAccumulator
from generated import results_pb2_grpc
from tag_tunnel import TagTunnel

tags_fp="./results/server/tags.yaml"
tags_schema_fp="./results/server/schema/tags_schema.json"
template_fp="./results/server/tm.py.j2"
historic_tests_fp="../macfe-hil.github.io/historic_tests.json"
reports_dir="../macfe-hil.github.io/reports"
pages_repo_dir="../macfe-hil.github.io"
pages_branch="main"
git_username="macformularacing"
git_email="macformulaelectric@gmail.com"

def serve():
    port = "8080"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    ra = ResultAccumulator(tags_fp=tags_fp,
                           tags_schema_fp=tags_schema_fp,
                           template_fp=template_fp,
                           historic_tests_fp=historic_tests_fp,
                           reports_dir=reports_dir,
                           pages_repo_dir=pages_repo_dir,
                           pages_branch=pages_branch,
                           git_username=git_username,
                           git_email=git_email)

    results_pb2_grpc.add_TagTunnelServicer_to_server(TagTunnel(result_accumulator=ra), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    serve()