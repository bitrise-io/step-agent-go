package main

import (
	"encoding/base64"
	"fmt"
)

func decode(encodedContent string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(encodedContent)
	if err != nil {
		fmt.Println("Failed for input: ", encodedContent)
		return "", err
	}

	return string(bytes), nil
}

func decodeArgs(encodedArgs []string) ([]string, error) {
	theDecodedArgs := make([]string, len(encodedArgs), len(encodedArgs))
	for idx, anEncodedArg := range encodedArgs {
		fmt.Println("anEncodedArg: ", anEncodedArg)
		aDecodedArg, err := decode(anEncodedArg)
		fmt.Println("aDecodedArg: ", aDecodedArg)
		if err != nil {
			return []string{}, err
		}
		theDecodedArgs[idx] = aDecodedArg
	}
	return theDecodedArgs, nil
}

func perform() {

}

func main() {
	fmt.Println("main")
}
