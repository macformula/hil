import git


class RepoHandler:

    def __init__(
        self, pages_repo_dir: str, pages_branch: str, git_username: str, git_email: str
    ) -> None:
        self.pages_repo_dir = pages_repo_dir
        self.pages_branch = pages_branch
        self.git_username = git_username
        self.git_email = git_email

    def push_to_github_pages(self, test_id) -> None:
        print("START OF PUSH")
        repo = git.Repo(self.pages_repo_dir)

        # Set Git user name and email for the repository
        repo.config_writer().set_value("user", "name", self.git_username).release()
        repo.config_writer().set_value("user", "email", self.git_email).release()

        # Add all changes to the index
        repo.git.add("*")

        # Commit changes
        repo.index.commit(test_id)

        # Fetch updates from the remote repository
        repo.remotes.origin.fetch()

        # Ensure we are on the branch that needs to be updated
        repo.heads[self.pages_branch].checkout()

        # Rebase changes
        repo.git.rebase("origin/" + self.pages_branch)

        # Push the changes, forcing the update to the remote branch
        repo.git.push()
        print("END OF PUSH")

        return
