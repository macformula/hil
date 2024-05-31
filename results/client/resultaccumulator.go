package results

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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
	Type          string    `json:"type"`
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
	value, err := loadTestsFromYAML(tagsFilepath)
	fmt.Println("value: ", value)
	if err != nil {
		fmt.Println("err load yaml ", err)
		return err
	}
	return nil
}

func loadTestsFromYAML(filepath string) (map[string]Test, error) {
	tagData, err := loadYAML(filepath)
	if err != nil {
		return nil, fmt.Errorf("invalid tags data format in %s", filepath)
	}

	tags, ok := tagData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid tags data format in %s", filepath)
	}

	testMap := make(map[string]Test) 
	for tagID, tagInfo := range tags {
		if tagInfo == nil {
			return nil, fmt.Errorf("nil tag info for tag %s", tagID)
		}

		test := Test{}
		if m, ok := tagInfo.(map[string]interface{}); ok {
			test.CompOp = m["compareOp"].(string)        
			test.Description = m["description"].(string) 
			test.Type = m["type"].(string)               
			test.Unit = m["unit"].(string)               
			test.UpperLimit, _ = m["upperLimit"].(string)
			test.LowerLimit, _ = m["lowerLimit"].(string)

			// Handle expectedVal carefully
			expectedVal, exists := m["expectedVal"]
			if exists {
				// Only assign if the key exists in the map
				test.ExpectedValue, _ = expectedVal.(string)
			}

			// fmt.Println("Compare Op:", test.CompOp)
			// fmt.Println("Description:", test.Description)
			// fmt.Println("Expected Val:", test.ExpectedValue) // Might be empty if not in YAML
			// fmt.Println("Type:", test.Type)
			// fmt.Println("Unit:", test.Unit)
		} else {
			fmt.Println("Error: Data is not in the expected map format")
		}
		testMap[tagID] = test
	}
	fmt.Println("end of load")
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

func getStringFromInterface(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64) // Convert float64 to string with full precision
	case bool:
		return strconv.FormatBool(v) // Convert boolean to string ("true" or "false")
	case nil:
		return "" // Handle nil values, returning an empty string
	default:
		return fmt.Sprintf("%v", v) // Fallback: convert to string representation for other types
	}
}
