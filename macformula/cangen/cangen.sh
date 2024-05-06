#!/bin/bash

# Function to display help menu
display_help() {
    echo "Usage: $0 <directory>"
    echo "   <directory>: Directory path containing DBC files."
    echo ""
    echo "Example:"
    echo "   $0 /path/to/dbc/files"
}

# Check if an argument is provided
if [ -z "$1" ]; then
  echo "Error: Please provide a directory path containing DBC files."
  display_help
  exit 1
fi

# Check if the provided directory exists
if [ ! -d "$1" ]; then
  echo "Error: Directory '$1' not found."
  display_help
  exit 1
fi

# Store the provided directory path
DBC_DIR="$1"

# Array to store temporary folders
temp_folders=()

DBC_GEN_DIR="./temp"

# Loop through all files ending with .dbc in the directory
for dbc_file in "${DBC_DIR}"/*.dbc; do
  # Extract filename without extension
  filename=$(basename "${dbc_file}" .dbc)

  # Create a folder named "<filename>can"
  temp_dbc_folder="${DBC_GEN_DIR}/${filename}can"
  mkdir -p "$temp_dbc_folder"

  # Copy the dbc file to the created folder
  cp "${dbc_file}" "${temp_dbc_folder}/"

  # Add folder to the array
  temp_folders+=("${temp_dbc_folder}")
done

# Call the go generate command
go run -v -mod=mod go.einride.tech/can/cmd/cantool generate ${DBC_GEN_DIR} .

# Delete the temporary folders
for folder in "${temp_folders[@]}"; do
  rm -rf "${DBC_GEN_DIR}"
done

echo "DBC files processed, folders generated, and deleted after processing."
