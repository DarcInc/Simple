//Package reflex
// This package provides dependency assistance services stuff.
//
// Copyright 2021 Paul C. Hoehne
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
// OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
// THE POSSIBILITY OF SUCH DAMAGE.
package reflex

import (
	"fmt"
	"reflect"
)

type Reflex struct {
	guts  map[string]interface{}
	types map[string]reflect.Type
}

// NewReflex constructs a new reflex.  It is up to the user to decide how to
// manage references to a global object or to other classes.
func NewReflex() Reflex {
	return Reflex{
		guts:  make(map[string]interface{}),
		types: make(map[string]reflect.Type),
	}
}

// Register a new type or value with the reflex.
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

// Get returns a value and a bool for a previously registered value, type, or function.
// If there is no matching type, value, or function, then Get returns `nil` and `false`.
func (dm Reflex) Get(name string) (interface{}, bool) {
	anInstance, ok := dm.guts[name]
	if aType, hasType := dm.types[name]; !ok && hasType {
		return dm.constructFromType(aType)
	}

	return dm.returnValue(anInstance)
}

// MustGet returns a value for a given name.  If no such value has been registered, it
// panics.  This is useful for key values that are required for the system to run.
// Otherwise, uset `Get` and test the values as appropriate for a safer experience.
func (dm Reflex) MustGet(name string) interface{} {
	someAsset, ok := dm.Get(name)
	if !ok {
		panic(fmt.Sprintf("failed to find a registered value for %s", name))
	}

	return someAsset
}

// Inject takes a type and constructs an object based on those injected values.
func (dm Reflex) Inject(someType reflect.Type) interface{} {
	result, _ := dm.constructFromType(someType)
	return result
}
