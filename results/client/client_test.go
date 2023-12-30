package client

import (
	"testing"
)

type testCase[T Data] struct {
	tag   string
	value T
	want  bool
}

func TestCompleteTest(t *testing.T) {
	// Cannot reduce this to an array because of the generic typing (without adding a lot of extra code).
	test1 := testCase[float32]{
		tag:   "PV001",
		value: 94,
		want:  false,
	}
	test2 := testCase[string]{
		tag:   "PV003",
		value: "Hello",
		want:  true,
	}

	want := true
	testSubmitTag(test1, t)
	testSubmitTag(test2, t)
	success, err := CompleteTest()
	if err != nil {
		t.Errorf("CompleteTest: got error %s, expected %t", err.Error(), want)
	} else if success != want {
		t.Errorf("CompleteTest: got %t, expected %t", success, want)
	}
}

func testSubmitTag[T Data](test testCase[T], t *testing.T) {
	success, err := SubmitTag(test.tag, test.value)
	if err != nil {
		t.Errorf("{\"%s\":%v}: got error %s, expected %t", test.tag, test.value, err.Error(), test.want)
	} else if success != test.want {
		t.Errorf("{\"%s\":%v}: got %t, expected %t", test.tag, test.value, success, test.want)
	}
}
