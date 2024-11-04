#!/bin/bash

# Author: Teghveer Singh Ateliey
# Date: 2024-10-30

# Check args
if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <commit-hash> <repo-path> <project> <platform>"
    exit 2
fi

# Assign inputs to variables
COMMIT_HASH=$1
RACECAR_PATH=$2
PROJECT=$3
PLATFORM=$4

# Checks if directory exists
FW_PATH=$RACECAR_PATH/firmware
if [ ! -d "$FW_PATH" ]; then
    echo "Error: Firmware directory not found within repository path '$RACECAR_PATH'."
    exit 2
fi

# Navigate to the racecar/firmware directory path
cd "$FW_PATH"

# Fetch latest changes and check out the specified commit
git fetch --all
git checkout "$COMMIT_HASH"

# Build command
if make PROJECT="$PROJECT" PLATFORM="$PLATFORM" clean build; then
    echo "Build completed successfully!!"
else
    echo "Build failed :("
fi

# Return to the main branch
git checkout main