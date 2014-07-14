package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	flagEncodedStepPath         = flag.String("steppath", "", "[REQUIRED] step's path (base64 encoded)")
	flagEncodedCombinedStepEnvs = flag.String("stepenvs", "", "[REQUIRED] step's encoded-combined environment key-value pairs")
)

func transformIfSpecialEnvPair(envKeyValuePair EnvKeyValuePair) EnvKeyValuePair {
	// if envKeyValuePair.Key == "__INPUT_FILE__" {

	// }
	return envKeyValuePair
}

func filterEnvironmentKeyValuePairs(envKeyValuePair []EnvKeyValuePair) []EnvKeyValuePair {
	filteredPairs := []EnvKeyValuePair{}

	for _, aPair := range envKeyValuePair {
		if aPair.Key == "" {
			log.Println("[i] Key is missing - won't add it to the environment. Value: ", aPair.Value)
			continue
		}
		if aPair.Value == "" {
			log.Printf("[i] Value is missing - won't add it to the environment (Key: %s)\n", aPair.Key)
			continue
		}

		aPair = transformIfSpecialEnvPair(aPair)
		filteredPairs = append(filteredPairs, aPair)
	}

	return filteredPairs
}

func runCommandWithAdditionalEnvironment(commandPath string, envsToAdd []EnvKeyValuePair) error {
	c := exec.Command(commandPath)

	envLength := len(envsToAdd)
	if envLength > 0 {
		envStringPairs := make([]string, len(envsToAdd), len(envsToAdd))
		for idx, aEnvPair := range envsToAdd {
			envStringPairs[idx] = aEnvPair.String()
		}
		c.Env = append(os.Environ(), envStringPairs...)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func runCommandWithArgs(command string, cmdArgs ...string) error {
	c := exec.Command(command, cmdArgs...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func perform(encodedStepPath, encodedCombinedStepEnvs string) error {
	if encodedStepPath == "" {
		return errors.New("No Step Path provided")
	}

	decodedStepCommand, err := decodeSingleValue(encodedStepPath)
	if err != nil {
		return err
	}
	decodedStepEnvPairs, err := decodeCombinedEnvs(encodedCombinedStepEnvs)
	if err != nil {
		return err
	}

	filteredEnvPairs := filterEnvironmentKeyValuePairs(decodedStepEnvPairs)

	fmt.Println("Perform: ", decodedStepCommand, filteredEnvPairs)
	if err := runCommandWithArgs("chmod", "+x", decodedStepCommand); err != nil {
		return err
	}
	return runCommandWithAdditionalEnvironment(decodedStepCommand, filteredEnvPairs)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *flagEncodedStepPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := perform(*flagEncodedStepPath, *flagEncodedCombinedStepEnvs); err != nil {
		log.Fatal(err)
	}
}
