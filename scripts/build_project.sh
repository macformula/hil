# Author: Teghveer Singh Ateliey
# Date: 2024-10-30

#!/bin/bash

# Check args
if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <commit-hash> <repo-path> <project> <platform>"
    exit 1
fi

# Assign inputs to variables
COMMIT_HASH=$1
REPO_PATH=$2
PROJECT=$3
PLATFORM=$4

# Checks if directory exists
if [ ! -d "$REPO_PATH" ]; then
    echo "Error: Repository path '$REPO_PATH' does not exist."
    exit
fi

# Navigate to the racecar/firmware directory path
cd "$REPO_PATH" || exit

# Activate the Python virtual environment
source .env/Scripts/activate

# Fetch latest changes and check out the specified commit
git fetch --all
git checkout "$COMMIT_HASH"

# Build command
if make PROJECT="$PROJECT" PLATFORM="$PLATFORM" build; then
    echo "Build completed successfully!!"
else
    echo "Build failed :("
    exit
fi

# Return to the main branch
git checkout main