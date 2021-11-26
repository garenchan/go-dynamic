package dynamic

import (
	"errors"
	"reflect"
)

// ErrNoSuchMethod no such method error.
var ErrNoSuchMethod = errors.New("No such method")

// Call dynamic call method of object by name.
func Call(obj interface{}, methodName string, args ...interface{}) (result []interface{}, err error) {
	// prevent panic crash
	defer func() {
		reason := recover()

		// avoid overwriting real error
		if reason != nil && err == nil {
			err = interfaceToError(reason)
		}
	}()

	// get method of object by name
	method := methodByName(obj, methodName)

	// no such method
	if method == nil {
		err = ErrNoSuchMethod
		return
	}

	// call method
	in := argsToValues(args)
	out := method.Call(in)

	// convert return values to []interface{}
	result = valuesToResult(out)
	return
}

// interfaceToError convert interface to error.
func interfaceToError(i interface{}) (err error) {
	if i == nil {
		return
	}

	switch x := i.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = errors.New("Unknown error")
	}

	return
}

// methodByName return method of object by name.
func methodByName(obj interface{}, methodName string) *reflect.Value {
	value := reflect.ValueOf(obj)
	method := value.MethodByName(methodName)

	// It means that the method does not exist in the object
	if !method.IsValid() {
		return nil
	}

	return &method
}

// argsToValues converts function arguments.
func argsToValues(args []interface{}) []reflect.Value {
	values := make([]reflect.Value, len(args))

	for index, arg := range args {
		values[index] = reflect.ValueOf(arg)
	}

	return values
}

// valuesToResult converts the return values of a method.
func valuesToResult(values []reflect.Value) []interface{} {
	result := make([]interface{}, len(values))

	for index, value := range values {
		result[index] = value.Interface()
	}

	return result
}
