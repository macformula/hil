import grpc
from concurrent import futures
from repo_handler import RepoHandler
from result_accumulator import ResultAccumulator
from generated import results_pb2_grpc
from tag_tunnel import TagTunnel

tags_fp="./results/server/tags.yaml"
tags_schema_fp="./results/server/schema/tags_schema.json"
template_fp="./results/server/pytest_template.py.j2"
historic_tests_fp="../macfe-hil.github.io/historic_tests.json"
reports_dir="../macfe-hil.github.io/reports"
pages_repo_dir="../macfe-hil.github.io"
pages_branch="main"
git_username="macformularacing"
git_email="macformulaelectric@gmail.com"
rp_server_address = "localhost:31763"


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    # Insta
    repo_handler = RepoHandler(pages_repo_dir=pages_repo_dir,
                               pages_branch=pages_branch,
                               git_username=git_username,
                               git_email=git_email,
                               )
    
    ra = ResultAccumulator(tags_fp=tags_fp,
                           tags_schema_fp=tags_schema_fp,
                           template_fp=template_fp,
                           historic_tests_fp=historic_tests_fp,
                           reports_dir=reports_dir,
                           repo_handler=repo_handler,
                           )

    results_pb2_grpc.add_TagTunnelServicer_to_server(TagTunnel(result_accumulator=ra), server)
    server.add_insecure_port(rp_server_address)
    server.start()
    print("Listening on " + rp_server_address)
    server.wait_for_termination()

if __name__ == "__main__":
    serve()