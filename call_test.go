package dynamic

import (
	"encoding/json"
	"errors"
	"testing"
)

// rpcServer rpc server
type rpcServer struct {
	eps []endpoint
}

// endpoint rpc endpoint
type endpoint interface {
}

// request rpc request
type request struct {
	MethodName string
	Args       []interface{}
}

// newRPCServer create a rpc server
func newRPCServer(endpoints ...endpoint) *rpcServer {
	rs := &rpcServer{}
	rs.eps = endpoints

	return rs
}

// newRequestString use for create a request string
func newRequestString(methodName string, args ...interface{}) (string, error) {
	req := request{MethodName: methodName, Args: args}
	bytes, err := json.Marshal(req)

	return string(bytes), err
}

// Handle handle a request
func (rs *rpcServer) Handle(reqString string) (string, error) {
	var req request
	if err := json.Unmarshal([]byte(reqString), &req); err != nil {
		return "", err
	}

	for _, ep := range rs.eps {
		result, err := Call(ep, req.MethodName, req.Args...)

		if err == nil {
			if bytes, err := json.Marshal(result); err != nil {
				return "", err
			} else {
				return string(bytes), nil
			}
		} else if err != ErrNoSuchMethod {
			return "", err
		}
	}

	return "", ErrNoSuchMethod
}

// myEndpoint my endpoint for test
type myEndpoint struct {
}

// Add float64 addition
func (ep *myEndpoint) Add(a, b float64) float64 {
	return a + b
}

// Sub float64 subtraction
func (ep *myEndpoint) Sub(a, b float64) float64 {
	return a - b
}

func TestCall(t *testing.T) {
	ep := &myEndpoint{}
	server := newRPCServer(ep)

	var tests = []struct {
		methodName   string
		args         []interface{}
		expectResult string
		expectErr    error
	}{
		{"Add", []interface{}{1, 2}, "[3]", nil},
		{"Sub", []interface{}{2, 1}, "[1]", nil},
		{"Add", []interface{}{1, 2, 3}, "", errors.New("reflect: Call with too many input arguments")},
		{"Add", []interface{}{1}, "", errors.New("reflect: Call with too few input arguments")},
		{"NotImplemented", []interface{}{}, "", ErrNoSuchMethod},
	}

	for _, test := range tests {
		reqString, _ := newRequestString(test.methodName, test.args...)
		actualResult, actualErr := server.Handle(reqString)

		if actualErr != test.expectErr && (
			actualErr != nil && test.expectErr != nil && actualErr.Error() != test.expectErr.Error()) {
			t.Errorf("Expect error %v, but got %v", test.expectErr, actualErr)
		}

		if actualErr == nil && actualResult != test.expectResult {
			t.Errorf("Expect result %v, but got %v", test.expectResult, actualResult)
		}
	}
}
