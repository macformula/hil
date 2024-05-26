package results

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

const (
	tagsFilepath       = "../../results/server/tags.yaml"
	tagsSchemaFilepath = "../../results/server/schema/tags_schema.json"
)

type Test struct {
	ID            uuid.UUID `json:"id"`
	Description   string    `json:"description"`
	Value         bool      `json:"final_value"`
	CompOp        string    `json:"comp_op"`
	UpperLimit    string    `json:"upper_limit"`
	LowerLimit    string    `json:"lower_limit"`
	ExpectedValue string    `json:"expected_val"`
	Type          bool      `json:"type"`
	Unit          string    `json:"unit"`
}

type ResultAccumulator struct {
	tagDb            map[string]Test        // tag_id -> Tag
	tagSubmissions   map[string]interface{} // tag_id -> value cache
	errorSubmissions []string               // list of cached errors
	allTagsPassing   bool
}

func (ra *ResultAccumulator) NewResultAccumulator() error {
	// 1. Validate tags.yaml against tags_scheme.json
	err := validateTags(tagsFilepath, tagsSchemaFilepath)
	if err != nil {
		return err
	}

	return nil
}

func validateTags(tagsFilepath, schemaFilepath string) error {
	tagsData, err := loadYAML(tagsFilepath)
	if err != nil {
		return err
	}

	// Convert tagsData from YAML to JSON compatible map[string]interface{}
	jsonData, err := convertYAMLtoJSON(tagsData)
	if err != nil {
		return err
	}

	absSchemaPath, _ := filepath.Abs(schemaFilepath)
	schemaURI := "file://" + absSchemaPath
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)

	documentLoader := gojsonschema.NewGoLoader(jsonData)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("error during validation: %v", err)
	}
	fmt.Println("result ", result)

	if !result.Valid() {
		var errorMessages []string
		for _, desc := range result.Errors() {
			errorMessages = append(errorMessages, desc.String())
		}
		return fmt.Errorf("tags.yaml does not conform to the schema:\n%s", errorMessages)
	}

	return nil
}

// Helper function to convert YAML data to JSON compatible format
func convertYAMLtoJSON(yamlData interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(yamlData)
	if err != nil {
		return nil, fmt.Errorf("error converting YAML to JSON: %v", err)
	}

	var jsonResult interface{}
	err = json.Unmarshal(jsonData, &jsonResult)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return jsonResult, nil
}
func loadYAML(relativeFilepath string) (interface{}, error) {
	// Get the absolute path of the YAML file
	absFilepath, err := filepath.Abs(relativeFilepath)
	if err != nil {
		return nil, fmt.Errorf("error resolving absolute path for %s: %v", relativeFilepath, err)
	}

	// Read the YAML file contents using the absolute path
	data, err := os.ReadFile(absFilepath)
	if err != nil {
		return nil, err
	}
	var out interface{}
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// func getOrDefault(data map[interface{}]interface{}, key string, defaultValue interface{}) interface{} {
// 	if val, ok := data[key]; ok {
// 		return val
// 	}
// 	return defaultValue
// }
