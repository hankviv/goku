package gson

import (
	"fmt"
	"reflect"
)

func populateStructReflect(in interface{}) error {
	val := reflect.ValueOf(in)
	if val.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("you must pass in a pointer")
	}
	elmv := val.Elem()
	if elmv.Type().Kind() != reflect.Struct {
		return fmt.Errorf("you must pass in a pointer to a struct")
	}

	fval := elmv.FieldByName("B")
	fval.SetInt(42)

	return nil
}
