import yaml
from jsonschema import validate

def validate_tags(tags_file_path, schema_file_path) -> bool:
    # Load YAML schema
    with open(schema_file_path, 'r') as schema_file:
        schema = yaml.safe_load(schema_file)

    # Load YAML data
    with open(tags_file_path, 'r') as tags_file:
        tags_data = yaml.safe_load(tags_file)

    # Validate the data against the schema
    try:
        validate(tags_data, schema)
        return True
    except Exception as e:
        print(f"Validation error: {e}")
        return False

# Paths to your YAML files
tags_schema_yaml_path = "tags_schema.json"
tags_yaml_path_good = "good_tags.yaml"
tags_yaml_path_bad = "bad_tags.yaml"

# Validate the YAML data against the schema
if validate_tags(tags_yaml_path_good, tags_schema_yaml_path):
    print("Passed test for good yaml")
else:
    print("Failed test for good yaml")
if not validate_tags(tags_yaml_path_bad, tags_schema_yaml_path):
    print("Passed test for bad yaml")
else:
    print("Failed test for bad yaml")


