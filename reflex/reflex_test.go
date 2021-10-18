package reflex

import "testing"

func TestDependencyManager_GetInstanceSimpleTypes(t *testing.T) {
	dm := Reflex{guts: make(map[string]interface{})}

	var shortInt int16 = 20
	var regularInt = 21
	var bigInt int64 = 22
	var someString = "Hello World"
	var shortFloat float32 = 2.3
	var regularFloat = 3.4
	var someRune = 'A'

	dm.Register("shortInt", shortInt)
	dm.Register("regularInt", regularInt)
	dm.Register("bigInt", bigInt)
	dm.Register("someString", someString)
	dm.Register("shortFloat", shortFloat)
	dm.Register("regularFloat", regularFloat)
	dm.Register("someRune", someRune)
	dm.Register("someBool", true)

	if val, ok := dm.MustGet("shortInt").(int16); !ok || val != shortInt {
		t.Errorf("Expected short integer but got %v %d\n", ok, val)
	}

	if val, ok := dm.MustGet("regularInt").(int); !ok || val != regularInt {
		t.Errorf("Expected integer but got %v %d\n", ok, val)
	}

	if val, ok := dm.MustGet("bigInt").(int64); !ok || val != bigInt {
		t.Errorf("Expected big int but got %v %d\n", ok, val)
	}

	if val, ok := dm.MustGet("someString").(string); !ok || val != someString {
		t.Errorf("Expected some string but got %v %s\n", ok, val)
	}

	if val, ok := dm.MustGet("shortFloat").(float32); !ok || val != shortFloat {
		t.Errorf("Expected short float but got %v %f\n", ok, val)
	}

	if val, ok := dm.MustGet("regularFloat").(float64); !ok || val != regularFloat {
		t.Errorf("Expected regular float but got %v %f\n", ok, val)
	}

	if val, ok := dm.MustGet("someRune").(rune); !ok || val != someRune {
		t.Errorf("Expected rune but got %v %d\n", ok, val)
	}

	if val, ok := dm.MustGet("someBool").(bool); !ok || !val {
		t.Errorf("Expected bool but got %v %v\n", ok, val)
	}
}

func TestDependencyManager_GetInstanceSlices(t *testing.T) {
	dm := Reflex{guts: make(map[string]interface{})}

	stringSlice := []string{"foo", "bar", "baz"}

	dm.Register("stringSlice", stringSlice)

	if val, ok := dm.MustGet("stringSlice").([]string); !ok {
		t.Error("Failed to get a string slice")
	} else {
		for i := range stringSlice {
			if val[i] != stringSlice[i] {
				t.Errorf("Expected %s to equal %s", val[i], stringSlice[i])
			}
		}
	}
}

func TestDependencyManager_GetInstanceBasicStruct(t *testing.T) {
	type BasicStruct struct {
		Foo int
		Bar string
		Baz float64
	}

	dm := Reflex{guts: make(map[string]interface{})}
	dm.Register("Basic", BasicStruct{Foo: 100, Bar: "some bar", Baz: 3.4})

	if val, ok := dm.MustGet("Basic").(BasicStruct); !ok {
		t.Error("Expected to get back a type of basic struct")
	} else {
		if val.Foo != 100 || val.Bar != "some bar" || val.Baz != 3.4 {
			t.Errorf("Expected 100, 'some bar', 3.4 but got %d '%s' %f", val.Foo, val.Bar, val.Baz)
		}
	}
}
