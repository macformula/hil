from tag import Tag
from typing import Union
import yaml
import jsonschema


class ResultAccumulator:
    def __init__(self, tag_file_path: str, schema_file_path: str) -> None:
        self.__parse_args(tag_file_path, schema_file_path)
        self.tag_submissions = {}

    def submit_tag(self, tag_id: str, value: any) ->Union[bool, KeyError]:
        try:
            is_passing = self.tag_db[tag_id].is_passing(value)
            self.tag_submissions[tag_id] = value
            return is_passing, None
        except KeyError as e:
            return False, e
    
    def trigger_pytest(self):
        for tag_id, tag in self.tag_db:
            assert tag.is_passing(self.tag_submissions[tag_id])

    def __parse_args(self, tag_file_path: str, schema_file_path: str) -> None:
            self.tag_db = {}
            with open(tag_file_path) as f:
                tags_yaml = yaml.load(f, Loader=yaml.FullLoader)
                
                self.__validate_tags(tag_file_path, schema_file_path)
                
                for tag_id, tag_info in tags_yaml.items():
                    if tag_id in self.tag_db:
                        #TODO: throw exception here. Ideally this check is done
                        # inside of the JSON schema. could do this by making the 
                        # yaml a list instead since JSON schema can check for
                        # unique entries only for arrays, not dicts (idk why)
                        return

                    newTag = Tag(
                        description=tag_info.get("description", ""),
                        comp_op=tag_info.get("compareOp", ""),
                        type=tag_info.get("type", ""),
                        lower_limit=tag_info.get("lowerLimit", None),
                        upper_limit=tag_info.get("upperLimit", None),
                        expected_val=tag_info.get("expectedVal", None),
                        unit=tag_info.get("unit", ""),
                    )

                    self.tag_db[tag_id] = newTag

    
    def __validate_tags(self, tags_file_path, schema_file_path) -> bool:
        # Load YAML schema
        with open(schema_file_path, 'r') as schema_file:
            schema = yaml.safe_load(schema_file)

        # Load YAML data
        with open(tags_file_path, 'r') as tags_file:
            tags_data = yaml.safe_load(tags_file)

        # Validate against the schema
        jsonschema.validate(tags_data, schema)
