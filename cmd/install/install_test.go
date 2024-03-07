package install

import (
	"testing"
)

func TestGenerateVariables(t *testing.T) {
	test := map[string]string{"base64": "func(base64)", "nofunction": "6446"}
	varMap := generateVariables(test)

	val, ok := varMap["base64"]

	if !ok {
		t.Error("Base64 key does not exist")
	}

	if val == "func(base64)" {
		t.Error("base64 function not called")
	}

	val, ok = varMap["nofunction"]

	if !ok {
		t.Error("nofunction key does not exist")
	}

	if val != test["nofunction"] {
		t.Error("Value not copied correctly", val)
	}
}
