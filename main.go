package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func writeStringToFile(filePath, content string) error {
	if filePath == "" {
		return errors.New("No path provided!")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

func transformIfSpecialEnv(envKeyValuePair EnvKeyValuePair) (EnvKeyValuePair, error) {
	if envKeyValuePair.Key == "__INPUT_FILE__" {
		log.Println(" (i) Special key: __INPUT_FILE__")
		usr, err := user.Current()
		if err != nil {
			return EnvKeyValuePair{}, err
		}
		tmpFolderPath := filepath.Join(usr.HomeDir, "bitrise/tmp")
		if err := os.MkdirAll(tmpFolderPath, 0777); err != nil {
			return EnvKeyValuePair{}, err
		}
		stepInputStoreFilePath := filepath.Join(tmpFolderPath, "step_input_store")
		if err := writeStringToFile(stepInputStoreFilePath, envKeyValuePair.Value); err != nil {
			return EnvKeyValuePair{}, err
		}
		envKeyValuePair.Value = stepInputStoreFilePath
	}
	return envKeyValuePair, nil
}

func filterEnvironmentKeyValuePairs(envKeyValuePair []EnvKeyValuePair) []EnvKeyValuePair {
	filteredPairs := []EnvKeyValuePair{}

	for _, aPair := range envKeyValuePair {
		if aPair.Key == "" {
			log.Println("[i] Key is missing - won't add it to the environment. Value: ", aPair.Value)
			continue
		}
		if aPair.Value == "" {
			log.Printf("[i] Value is missing - won't add it to the environment (default value will be used by the Step) (Key: %s)\n", aPair.Key)
			continue
		}

		aPair, err := transformIfSpecialEnv(aPair)
		if err != nil {
			log.Printf("[i] Failed to convert special Env - ignored (Key: %s | Value: %s)\n", aPair.Key, aPair.Value)
			continue
		}
		filteredPairs = append(filteredPairs, aPair)
	}

	return filteredPairs
}

func runStepWithAdditionalEnvironment(commandPath string, envsToAdd []EnvKeyValuePair) error {
	commandDir := filepath.Dir(commandPath)
	commandName := filepath.Base(commandPath)
	c := exec.Command("bash", commandName)

	envLength := len(envsToAdd)
	if envLength > 0 {
		envStringPairs := make([]string, len(envsToAdd), len(envsToAdd))
		for idx, aEnvPair := range envsToAdd {
			envStringPairs[idx] = aEnvPair.ToStringWithExpand()
		}
		c.Env = append(os.Environ(), envStringPairs...)
	}
	c.Dir = commandDir
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
	decodedStepCommand = ExpandPath(decodedStepCommand)
	decodedStepEnvPairs, err := decodeCombinedEnvs(encodedCombinedStepEnvs)
	if err != nil {
		return err
	}

	filteredEnvPairs := filterEnvironmentKeyValuePairs(decodedStepEnvPairs)

	fmt.Println("Perform: ", decodedStepCommand, filteredEnvPairs)
	return runStepWithAdditionalEnvironment(decodedStepCommand, filteredEnvPairs)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var (
		flagEncodedStepPath         = flag.String("steppath", "", "[REQUIRED] step's path (base64 encoded)")
		flagEncodedCombinedStepEnvs = flag.String("stepenvs", "", "[REQUIRED] step's encoded-combined environment key-value pairs")
	)

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
