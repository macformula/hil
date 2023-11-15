package dispatcher

type FakeTestData struct {
	Tags         map[string]string
	ExpectedData map[string]interface{}
	ActualData   map[string]interface{}
}

func NewFakeTestData() *FakeTestData {
	return &FakeTestData{
		Tags: map[string]string{
			"tag1": "description for tag 1",
			"tag2": "description for tag 2",
		},
		ExpectedData: map[string]interface{}{
			"tag1": "expected value for tag 1",
			"tag2": 42,
		},
		ActualData: map[string]interface{}{
			"tag1": "actual value for tag 1",
			"tag2": 42,
		},
	}
}
