#!/bin/bash

set -e  # Exit on errors

# Initialize verbose flag
VERBOSE=false

# Get current date
CURRENT_DATE=$(date +"%l:%M%PM - %b %-d, %Y")

# Get Git information
GIT_COMMIT=$(git rev-parse --short HEAD)
IS_DIRTY=$(git status --porcelain | wc -l)

if [[ $IS_DIRTY -ne 0 ]]; then
   DIRTY_VS_CLEAN="dirty"
else
   DIRTY_VS_CLEAN="clean"
fi

# Build flags for ldflags
LDFLAGS="-X 'main.Date=$CURRENT_DATE' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.DirtyVsClean=${DIRTY_VS_CLEAN}'"

# Parse command-line arguments
while getopts ":e:v" opt; do
  case ${opt} in
    e ) ENV_FILE="$OPTARG"
        ;;
    v ) VERBOSE=true
        ;;
    \? ) echo "Usage: $0 -e <env_file> [-v]"
         exit 1
        ;;
  esac
done
shift $((OPTIND -1))

if [ -z "$ENV_FILE" ]; then
  echo "ERROR: Environment file not provided."
  exit 1
fi

# Check if environment file exists
if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: Environment file '$ENV_FILE' not found."
  exit 1
fi

source "$ENV_FILE"  # Load environment variables from file

# Set default values for GOOS and GOARCH if not provided in the environment file
if [ -n "$GOOS" ]; then
  GOOS_ARG="GOOS=$GOOS"
fi

if [ -n "$GOARCH" ]; then
  GOARCH_ARG="GOARCH=$GOARCH"
fi

if [ -n "$GOARM" ]; then
  GOARM_ARG="GOARM=$GOARM"
fi

# Build the Go program with the specified GOOS and GOARCH, if provided
if [ -n "$GOOS_ARG" ] || [ -n "$GOARCH_ARG" ] || [ -n "$GOARM_ARG" ]; then
  if [ "$VERBOSE" = true ]; then
    echo "Building Go program with specified GOOS and GOARCH..."
  fi
  BUILD_CMD="$GOOS_ARG $GOARCH_ARG $GOARM_ARG go build -ldflags=\"$LDFLAGS\" -o hilapp main.go"
else
  if [ "$VERBOSE" = true ]; then
    echo "Building Go program..."
  fi
  BUILD_CMD="go build -ldflags="$LDFLAGS" -o hilapp main.go"
fi

if [ "$VERBOSE" = true ]; then
  echo "Executing build command: $BUILD_CMD"
fi

eval "$BUILD_CMD"

echo ""
echo "-------------------------------"
echo ""

# Check if REMOTE_TARGET is set in the environment file
if [ -n "$REMOTE_TARGET" ]; then
  if [ -n "$REMOTE_USER" ] && [ -n "$REMOTE_HOST" ] && [ -n "$BIN_PATH_REMOTE" ]; then
    if [ "$VERBOSE" = true ]; then
      echo "Syncing binary to remote target..."
    fi
    # Sync binary to remote target
    RSYNC_CMD="rsync -avz -e ssh hilapp '${REMOTE_USER}@${REMOTE_HOST}:${BIN_PATH_REMOTE}'"
    if [ "$VERBOSE" = true ]; then
      echo "Executing rsync command: $RSYNC_CMD"
    fi
    eval "$RSYNC_CMD"
    echo "Binary synced to remote target."
  else
    echo "ERROR: Remote user, host, or path not provided in environment file."
  fi
else
  echo "INFO: No remote target specified in environment file. Keeping binary locally."
fi

echo ""
echo "-------------------------------"
echo ""

if [ -n "$REMOTE_TARGET" ]; then
  # Check if RESULTS_SERVER_PATH_LOCAL and RESULTS_PI_PATH are set in the environment file
  if [ -n "$RESULTS_SERVER_PATH_LOCAL" ] && [ -n "$RESULTS_SERVER_PATH_REMOTE" ]; then
    if [ "$VERBOSE" = true ]; then
      echo "Syncing results server folder..."
    fi
    RSYNC_RESULTS_CMD="rsync -avz -e ssh $RESULTS_SERVER_PATH_LOCAL ${REMOTE_USER}@${REMOTE_HOST}:${RESULTS_SERVER_PATH_REMOTE}"
    if [ "$VERBOSE" = true ]; then
      echo "Executing rsync command: $RSYNC_RESULTS_CMD"
    fi
    eval "$RSYNC_RESULTS_CMD"
    echo "Results server folder synced."
  else
    echo "INFO: Results server local folder and/or remote path not provided in environment file. Skipping sync."
  fi
fi

echo ""
echo "-------------------------------"
echo ""

if [ -n "$REMOTE_TARGET" ]; then
  # Check if CONFIG_YAML_PATH_LOCAL and CONFIG_YAML_PATH_REMOTE are set in the environment file
  if [ -n "$CONFIG_YAML_PATH_LOCAL" ] && [ -n "$CONFIG_YAML_PATH_REMOTE" ]; then
    if [ "$VERBOSE" = true ]; then
      echo "Syncing config.yaml file..."
    fi
    RSYNC_CONFIG_CMD="rsync -e ssh $CONFIG_YAML_PATH_LOCAL ${REMOTE_USER}@${REMOTE_HOST}:${CONFIG_YAML_PATH_REMOTE}"
    if [ "$VERBOSE" = true ]; then
      echo "Executing rsync command: $RSYNC_CONFIG_CMD"
    fi
    eval "$RSYNC_CONFIG_CMD"
    echo "config.yaml file synced."
  else
    echo "INFO: config.yaml local file and/or remote path not provided in environment file. Skipping sync."
  fi
fi

echo ""
echo "-------------------------------"
echo ""

if [ -n "$REMOTE_TARGET" ]; then
  # Check if TAGS_YAML_PATH_LOCAL and TAGS_YAML_PATH_REMOTE are set in the environment file
  if [ -n "$TAGS_YAML_PATH_LOCAL" ] && [ -n "$TAGS_YAML_PATH_REMOTE" ]; then
    if [ "$VERBOSE" = true ]; then
      echo "Syncing tags.yaml file..."
    fi
    RSYNC_TAGS_CMD="rsync -e ssh $TAGS_YAML_PATH_LOCAL ${REMOTE_USER}@${REMOTE_HOST}:${TAGS_YAML_PATH_REMOTE}"
    if [ "$VERBOSE" = true ]; then
      echo "Executing rsync command: $RSYNC_TAGS_CMD"
    fi
    eval "$RSYNC_TAGS_CMD"
    echo "tags.yaml file synced."
  else
    echo "INFO: tags.yaml local file and/or remote path not provided in environment file. Skipping sync."
  fi
fi

echo ""
echo "-------------------------------"
echo ""

echo "Successfully built and synced all targets."
echo ""


