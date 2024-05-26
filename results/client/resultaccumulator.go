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
	ra.tagDb, err = loadTestsFromYAML(tagsFilepath)

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
		// Nil Check Before Type Assertion: Check if tagInfo is not nil before trying to access it
		if tagInfo == nil {
			return nil, fmt.Errorf("nil tag info for tag %s", tagID)
		}
		fmt.Println("tagInfo ", tagInfo, "\n")
		// infoMap, ok := tagInfo.(map[interface{}]interface{})
		// if !ok {
		// 	fmt.Println("ok ", ok, "\n")
		// 	return nil, fmt.Errorf("invalid tag info format for tag %s: %T", tagID, tagInfo)
		// }

		// fmt.Println("infoMap ", infoMap, "\n")

		if m, ok := data.(map[string]interface{}); ok {

			compareOp, _ := m["compareOp"].(string)     // Extract string
			description, _ := m["description"].(string) // Extract string
			expectedVal, _ := m["expectedVal"].(bool)   // Extract bool
			typeStr, _ := m["type"].(string)            // Extract string
			unit, _ := m["unit"].(string)               // Extract string

			fmt.Println("Compare Op:", compareOp)
			fmt.Println("Description:", description)
			fmt.Println("Expected Val:", expectedVal)
			fmt.Println("Type:", typeStr)
			fmt.Println("Unit:", unit)
		} else {
			fmt.Println("Error: Data is not in the expected map format")
		}

		// test := Test{
		// 	ID: uuid.New(),
		// }

		// if description, ok := infoMap["description"].(string); ok {
		// 	test.Description = description
		// }

		// if compOp, ok := infoMap["compareOp"].(string); ok {
		// 	test.CompOp = compOp
		// }

		// if unit, ok := infoMap["unit"].(string); ok {
		// 	test.Unit = unit
		// }

		// if expectedVal, ok := infoMap["expectedVal"]; ok {
		// 	switch v := expectedVal.(type) {
		// 	case bool:
		// 		test.Value = v
		// 		test.Type = "bool"
		// 	case string:
		// 		test.ExpectedValue = v
		// 	case int:
		// 		test.ExpectedValue = strconv.Itoa(v)
		// 	case float64:
		// 		test.ExpectedValue = strconv.FormatFloat(v, 'f', -1, 64) // Convert float64 to string with full precision
		// 	default:
		// 		return nil, fmt.Errorf("invalid type for expectedVal in tag %s", tagID)
		// 	}
		// }

		// if upperLimit, ok := infoMap["upperLimit"]; ok {
		// 	test.UpperLimit = getStringFromInterface(upperLimit)
		// }

		// if lowerLimit, ok := infoMap["lowerLimit"]; ok {
		// 	test.LowerLimit = getStringFromInterface(lowerLimit)
		// }

		// testMap[tagID] = test
	}
	fmt.Println("testMap", testMap)
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
