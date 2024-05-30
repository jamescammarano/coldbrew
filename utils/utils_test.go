package utils

import (
	"reflect"
	"testing"
)

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

func TestRandomString(t *testing.T) {
	test := RandomString(8)
	println(test)
}
