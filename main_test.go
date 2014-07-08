package main

import (
	"encoding/base64"
	"testing"
)

func Test_decode(t *testing.T) {
	t.Log("should decode")

	testContString := "any + old & data"
	if decodedString, err := decode(base64.StdEncoding.EncodeToString([]byte(testContString))); err != nil || decodedString != testContString {
		t.Error("decode result doesn't match the expected.\n Expected: %s\n Got:%s\n Error: %v", testContString, decodedString, err)
	}

	testInvalidContString := "this should not be a valid input if not encoded"
	if _, err := decode(testInvalidContString); err == nil {
		t.Error("expected to return an error for an invalid input")
	}
}

func Test_decodeArgs(t *testing.T) {
	t.Log("should decode an array of args")

	testContentStrings := []string{"cont 1", "cont 2"}
	if _, err := decodeArgs(testContentStrings); err == nil {
		t.Error("expected to return an error for an un-encoded input")
	}

	testEncodedStrings := []string{
		base64.StdEncoding.EncodeToString([]byte(testContentStrings[0])),
		base64.StdEncoding.EncodeToString([]byte(testContentStrings[1])),
	}
	decodedArgs, err := decodeArgs(testEncodedStrings)
	if err != nil {
		t.Error("returned error:", err)
	}
	if decodedArgs[0] != "cont 1" || decodedArgs[1] != "cont 2" {
		t.Error("result doesn't match the expected.\n Expected: %v\n Got:%v\n Error: %v", testContentStrings, decodedArgs, err)
	}
}

func Test_perform(t *testing.T) {
	t.Log("should perform")
}
