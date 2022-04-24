package grpc

import "testing"

func createClient() {

}

func Test_GrpcGet(t *testing.T) {
	Assert(t, "Hello", 1, "Bonjour")
}

func Assert(t *testing.T, expected any, actual any, message string) {
	if expected != actual {
		t.Errorf(`%s: Expected "%v" but got "%v"`, message, expected, actual)
	}
}
