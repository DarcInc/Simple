package reflex

import (
	"fmt"
	"reflect"
)

// DependencyManager manages the data used to store the injectable dependencies.
type DependencyManager struct {
	// Todo: Cache constructed dependencies
	// Todo: Some way of getting a singleton or goroutine reference to a dependency manager
	// For example, we start a goroutine that encapsulates the dependency management and
	// the reference to the Dependency Manager is a channel.
	guts map[string]interface{}
}

// RegisterFactory registers factory functions that take multiple arguments, and the dependency manager
// and return a constructed dependency.
func (dm *DependencyManager) RegisterFactory(name string, factory interface{}) {
	dm.guts[name] = factory
}

// Register registers either a concrete value (e.g. an int), a type that serves as a template for
// creating a dependency, or a function that takes a dependency manager and returns a constructed
// dependency.  The dependency can be looked up by its name.
func (dm *DependencyManager) Register(name string, item interface{}) {
	dm.guts[name] = item
}

// MakeInstance allows passing multiple parameters to a constructor that builds an object.
func (dm DependencyManager) MakeInstance(name string, params ...interface{}) interface{} {
	someFactory := dm.guts[name]

	t, v := reflect.TypeOf(someFactory), reflect.ValueOf(someFactory)
	if t.Kind() != reflect.Func {
		panic("you failed to pass a function!")
	}

	callArray := make([]reflect.Value, len(params))
	for i, v := range params {
		callArray[i] = reflect.ValueOf(v)
	}

	result := v.Call(callArray)[0]
	fmt.Printf("Type of result: %v\n", result.Type())

	return result.Interface()
}

func (dm DependencyManager) setByName(item reflect.Value, fieldName string, val interface{}) {
	field := item.Elem().FieldByName(fieldName)

	switch {
	case reflect.TypeOf(val).Kind() == field.Kind():
		field.Set(reflect.ValueOf(val))
	case reflect.TypeOf(val).Kind() == reflect.Ptr:
		field.Set(reflect.ValueOf(val).Elem())
	case field.Kind() == reflect.Ptr:
		field.Set(reflect.ValueOf(&val))
	}
}

func (dm DependencyManager) setField(item reflect.Value, field reflect.StructField) {
	injectFrom, ok := field.Tag.Lookup("inject")
	var v interface{}
	if ok {
		v = dm.GetInstance(injectFrom)
	} else {
		v = dm.GetInstance(field.Name)
	}
	dm.setByName(item, field.Name, v)
}

// GetInstance returns the constructed object that is our dependency.
func (dm DependencyManager) GetInstance(name string) interface{} {
	anInstance, ok := dm.guts[name]
	if !ok {
		return nil
	}

	t, v := reflect.TypeOf(anInstance), reflect.ValueOf(anInstance)

	if t.Kind() == reflect.Func {
		results := v.Call([]reflect.Value{reflect.ValueOf(dm)})
		return results[0].Interface()
	} else if t.Kind() == reflect.Struct {
		result := reflect.New(t)
		for i := 0; i < t.NumField(); i++ {
			dm.setField(result, t.Field(i))
		}
		return result.Interface()
	} else {
		return anInstance
	}
}
