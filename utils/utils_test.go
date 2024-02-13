package utils

import (
	"reflect"
	"testing"
)

func TestBase64(t *testing.T) {
	test1 := Base64()
	test2 := Base64()

	if test1 == test2 {
		t.Error("not unique")
	}
}

func TestMergeMaps(t *testing.T) {
	map1 := map[string]string{"TestAttr": "TestData"}
	map2 := map[string]string{"Test2": "TestData2"}
	map3 := map[string]string{"Test3": "TestData3"}

	mapArr := []map[string]string{map1, map2, map3}
	control := map[string]string{"TestAttr": "TestData", "Test2": "TestData2", "Test3": "TestData3"}

	test := MergeMaps(mapArr)
	eq := reflect.DeepEqual(test, control)
	if !eq {
		t.Error(eq)
	}
}
