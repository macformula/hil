package flow

// Tag is a way to identify expected bounds for a given result.
type Tag struct {
	// ID is the tag identifier.
	ID string `yaml:"tagId"`
	// Description describes the result that will be posted to this tag.
	Description string `yaml:"description"`
}
