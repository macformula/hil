package results

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

const (
	tagsFilepath       = "../server/tags.yaml"
	tagsSchemaFilepath = "../server/schema/tags_schema.json"
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
		fmt.Println("err  ", err)
		return err
	}
	fmt.Println("yaml ", tagsData)

	// schemaLoader := gojsonschema.NewReferenceLoader(schemaFilepath)

	// // Load the tags data into a JSON schema loader (since the library expects JSON)
	// documentLoader := gojsonschema.NewGoLoader(tagsData)

	// // Perform the validation
	// result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	// if err != nil {
	// 	return fmt.Errorf("error during validation: %v", err)
	// }

	// // Check the validation result
	// if !result.Valid() {
	// 	var errorMessages []string
	// 	for _, desc := range result.Errors() {
	// 		errorMessages = append(errorMessages, desc.String())
	// 	}
	// 	return fmt.Errorf("tags.yaml does not conform to the schema:\n%s", errorMessages)
	// }
	// fmt.Println("result yaml validation", result)
	return nil
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
