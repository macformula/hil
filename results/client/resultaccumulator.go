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

type Test struct { // These are Tags, not calling this tag right now because it might cause confusion
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
	// 2. Convert tag.yaml into Test structures
	//ra.tagDb, err = loadTestsFromYAML(tagsFilepath)
	err = loadTestsFromYAML(tagsFilepath)
	if err != nil {
		fmt.Println("err load yaml ", err)
		return err
	}
	return nil
}

func loadTestsFromYAML(filepath string) (map[string]Test, error) {
	tagData, err := loadYAML(filepath)
	if err != nil {
		fmt.Printf("err load yaml in", err)
		return nil, fmt.Errorf("invalid tags data format in %s", filepath)
	}

	// // Type assertion to ensure tagData is a map[string]interface{}
	tags, ok := tagData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid tags data format in %s", filepath)
	}

	testMap := make(map[string]Test)
	for tagID, tagInfo := range tags {
		infoMap, ok := tagInfo.(map[interface{}]interface{})
		if !ok {
			fmt.Printf("ok in \n", ok)
			return nil, fmt.Errorf("invalid tag info format for tag %s", tagID)
		}
		fmt.Println("||||", infoMap, "||||", testMap, "||||", tagID, "\n\n")
		// // Create a new Test struct with defaults and override with values from tagInfo
		// test := Test{
		// 	ID:            uuid.Nil, // Generate a unique ID
		// 	Description:   getOrDefault(infoMap, "description", "").(string),
		// 	CompOp:        getOrDefault(infoMap, "compareOp", "").(string),
		// 	Type:          false, // Assuming "type" is a bool, set default to false
		// 	UpperLimit:    getOrDefault(infoMap, "upperLimit", "0").(string),
		// 	LowerLimit:    getOrDefault(infoMap, "lowerLimit", "0").(string),
		// 	ExpectedValue: getOrDefault(infoMap, "expectedVal", "0").(string),
		// 	Unit:          getOrDefault(infoMap, "unit", "Unitless").(string),
		// }

		// // If value type is boolean then we want the final_value as a boolean
		// if test.Type == "bool" {
		// 	test.Value = getOrDefault(infoMap, "expectedVal", "false").(bool)
		// }

		// testMap[tagID] = test
	}

	return testMap, nil
}

func validateTags(tagsFilepath, schemaFilepath string) error {
	tagsData, err := loadYAML(tagsFilepath)
	if err != nil {
		fmt.Println("err 2", err)
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

	if !result.Valid() {
		var errorMessages []string
		for _, desc := range result.Errors() {
			errorMessages = append(errorMessages, desc.String())
		}
		fmt.Println("errorMessages ", errorMessages)
		return nil
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

func getOrDefault(data map[interface{}]interface{}, key string, defaultValue interface{}) interface{} {
	if val, ok := data[key]; ok {
		return val
	}
	return defaultValue
}
