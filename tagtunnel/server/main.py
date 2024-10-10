import grpc
import argparse
import os
import sys
import yaml
from concurrent import futures
from repo_handler import RepoHandler
from result_accumulator import ResultAccumulator
from generated import results_pb2_grpc
from tag_tunnel import TagTunnel

RESULT_PROCESSOR_KEY = "resultProcessor"
SERVER_ADDRESS_KEY = "serverAddress"
TAGS_FP_KEY = "tagsFilePath"
HISTORIC_TESTS_FP_KEY = "historicTestsFilePath"
REPORTS_DIR_KEY = "reportsDir"
GITHUB_PAGES_REPO_KEY = "githubPagesRepoDir"
GITHUB_PAGES_BRANCH_KEY = "githubPagesBranch"
GITHUB_USERNAME_KEY = "githubPagesUsername"
GITHUB_EMAIL_KEY = "githubPagesEmail"

TAGS_SCHEMA_FILE_PATH = "./schema/tags_schema.json"
TEMPLATE_FILE_PATH = "./pytest_template.py.j2"
TEMP_TEST_FILE_PATH = "./temp_test_file.py" #Note: this fle gets generated and deleted on the fly.

def read_config(config_fp: str):
    try:
        with open(config_fp, "r") as file:
            config = yaml.safe_load(file)
        return config
    except Exception as e:
        print(f"Error reading configuration file: {e}")
        sys.exit(1)


def serve(rpConfig: dict):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    # Insta
    repo_handler = RepoHandler(
        pages_repo_dir=rpConfig[GITHUB_PAGES_REPO_KEY],
        pages_branch=rpConfig[GITHUB_PAGES_BRANCH_KEY],
        git_username=rpConfig[GITHUB_USERNAME_KEY],
        git_email=rpConfig[GITHUB_EMAIL_KEY],
    )

    ra = ResultAccumulator(
        tags_fp=rpConfig[TAGS_FP_KEY],
        tags_schema_fp=TAGS_SCHEMA_FILE_PATH,
        template_fp=TEMPLATE_FILE_PATH,
        historic_tests_fp=rpConfig[HISTORIC_TESTS_FP_KEY],
        temp_test_fp=TEMP_TEST_FILE_PATH,
        reports_dir=rpConfig[REPORTS_DIR_KEY],
        repo_handler=repo_handler,
    )

    results_pb2_grpc.add_TagTunnelServicer_to_server(
        TagTunnel(result_accumulator=ra), server
    )
    server.add_insecure_port(rpConfig[SERVER_ADDRESS_KEY])
    server.start()
    print("Listening on " + rpConfig[SERVER_ADDRESS_KEY])
    server.wait_for_termination()


def main():
    parser = argparse.ArgumentParser(description="Starts the result processor server.")

    parser.add_argument("--config", type=str, help="Path to the configuration file")

    args = parser.parse_args()

    # Get the directory of the current script
    script_dir = os.path.dirname(os.path.realpath(__file__))

    # Change the current working directory to the script's directory
    os.chdir(script_dir)

    config_fp = args.config
    if not config_fp:
        print(
            'Error: Configuration file path is required. Use --config="/path/to/config".'
        )
        sys.exit(1)

    config = read_config(config_fp)

    if RESULT_PROCESSOR_KEY not in config:
        print(f'Error: Configuration file must contain key "{RESULT_PROCESSOR_KEY}".')
        sys.exit(1)

    serve(rpConfig=config[RESULT_PROCESSOR_KEY])


if __name__ == "__main__":
    main()
