package results

import (
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
		fmt.Println("err 2 ", err)
		return err
	}
	absSchemaPath, _ := filepath.Abs(schemaFilepath)
	schemaURI := "file://" + absSchemaPath
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)

	documentLoader := gojsonschema.NewGoLoader(tagsData)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println("err ", err)
		return nil
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

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
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

	// Unmarshal YAML data into a generic interface{}
	var out interface{}
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	// Convert the potentially nested map[interface{}]interface{} to map[string]interface{}
	out = convert(out)

	// Return the unmarshalled data (not the JSON string)
	return out, nil
}

// func getOrDefault(data map[interface{}]interface{}, key string, defaultValue interface{}) interface{} {
// 	if val, ok := data[key]; ok {
// 		return val
// 	}
// 	return defaultValue
// }
