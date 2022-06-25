package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	_, err := run()

	if err != nil {
		t.Error("Failed run()")
	}

	// switch v := (reflect.TypeOf(db)){
	// case driver.DB:
	// 	// do nothing
	// default:
	// 	t.Errorf("Type is not DB driver, its of type %s", v)
	// }
}
