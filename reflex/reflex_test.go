package reflex

import "testing"

func TestDependencyManager_GetInstanceSimpleTypes(t *testing.T) {
	dm := DependencyManager{guts: make(map[string]interface{})}

	var shortInt int16 = 20
	var regularInt int = 21
	var bigInt int64 = 22
	var someString string = "Hello World"
	var shortFloat float32 = 2.3
	var regularFloat float64 = 3.4
	var someRune rune = 'A'
	var someBool bool = true

	dm.Register("shortInt", shortInt)
	dm.Register("regularInt", regularInt)
	dm.Register("bigInt", bigInt)
	dm.Register("someString", someString)
	dm.Register("shortFloat", shortFloat)
	dm.Register("regularFloat", regularFloat)
	dm.Register("someRune", someRune)
	dm.Register("someBool", someBool)

	if val, ok := dm.GetInstance("shortInt").(int16); !ok || val != shortInt {
		t.Errorf("Expected short integer but got %v %d\n", ok, val)
	}

	if val, ok := dm.GetInstance("regularInt").(int); !ok || val != regularInt {
		t.Errorf("Expected integer but got %v %d\n", ok, val)
	}

	if val, ok := dm.GetInstance("bigInt").(int64); !ok || val != bigInt {
		t.Errorf("Expected big int but got %v %d\n", ok, val)
	}

	if val, ok := dm.GetInstance("someString").(string); !ok || val != someString {
		t.Errorf("Expected some string but got %v %s\n", ok, val)
	}

	if val, ok := dm.GetInstance("shortFloat").(float32); !ok || val != shortFloat {
		t.Errorf("Expected short float but got %v %f\n", ok, val)
	}

	if val, ok := dm.GetInstance("regularFloat").(float64); !ok || val != regularFloat {
		t.Errorf("Expected regular float but got %v %f\n", ok, val)
	}

	if val, ok := dm.GetInstance("someRune").(rune); !ok || val != someRune {
		t.Errorf("Expected rune but got %v %d\n", ok, val)
	}

	if val, ok := dm.GetInstance("someBool").(bool); !ok || val != someBool {
		t.Errorf("Expected bool but got %v %v\n", ok, val)
	}
}
