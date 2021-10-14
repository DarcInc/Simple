package reflex

import (
	"fmt"
	"reflect"
)

type Reflex struct {
	guts  map[string]interface{}
	types map[string]reflect.Type
}

func NewReflex() Reflex {
	return Reflex{
		guts:  make(map[string]interface{}),
		types: make(map[string]reflect.Type),
	}
}

func (dm *Reflex) Register(name string, item interface{}) {
	if _, ok := item.(reflect.Type); ok {
		dm.types[name] = item.(reflect.Type)
	} else {
		dm.guts[name] = item
	}
}

func (dm Reflex) setByName(item reflect.Value, fieldName string, val interface{}) {
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

func (dm Reflex) setField(item reflect.Value, field reflect.StructField) {
	injectFrom, ok := field.Tag.Lookup("inject")
	var v interface{}
	var found bool
	if ok {
		v, found = dm.Get(injectFrom)
	} else {
		v, found = dm.Get(field.Name)
	}
	if !found {
		dm.setByName(item, field.Name, reflect.Zero(field.Type))
	} else {
		dm.setByName(item, field.Name, v)
	}
}

func (dm Reflex) constructFromType(aType reflect.Type) (interface{}, bool) {
	result := reflect.New(aType)
	for i := 0; i < aType.NumField(); i++ {
		dm.setField(result, aType.Field(i))
	}
	return result.Elem().Interface(), true
}

func (dm Reflex) returnValue(anInstance interface{}) (interface{}, bool) {
	t, v := reflect.TypeOf(anInstance), reflect.ValueOf(anInstance)

	if t.Kind() == reflect.Func {
		results := v.Call([]reflect.Value{reflect.ValueOf(dm)})
		return results[0].Interface(), results[1].Bool()
	} else if t.Kind() == reflect.Struct {
		return anInstance, true
	} else {
		return anInstance, true
	}
}

func (dm Reflex) Get(name string) (interface{}, bool) {
	anInstance, ok := dm.guts[name]
	if aType, hasType := dm.types[name]; !ok && hasType {
		return dm.constructFromType(aType)
	}

	return dm.returnValue(anInstance)
}

func (dm Reflex) MustGet(name string) interface{} {
	someAsset, ok := dm.Get(name)
	if !ok {
		panic(fmt.Sprintf("failed to find a registered value for %s", name))
	}

	return someAsset
}

func (dm Reflex) Inject(someType reflect.Type) interface{} {
	result, _ := dm.constructFromType(someType)
	return result
}
