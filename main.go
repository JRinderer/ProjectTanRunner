package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type CustomTag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Config struct {
	CustomTags       []CustomTag `json:"custom_tags"`
	LocationInitFile string      `json:"location_init_file"`
}

func readConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &config, nil
}

func addRegistryEntries(cfg Config) error {
	for _, tag := range cfg.CustomTags {
		cmdStr := fmt.Sprintf(`reg add "HKLM\SOFTWARE\WOW6432Node\Tanium\Tanium Client\Sensor Data\Tags" /v %s /t REG_SZ /d "%s"`, tag.Name, tag.Description)
		cmd := exec.Command("cmd", "/c", cmdStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to execute command: %s, output: %s, error: %v", cmdStr, output, err)
		}
	}

	return nil
}

func main() {

	config, err := readConfig("conf.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Location Init File: %s\n", config.LocationInitFile)
	for _, tag := range config.CustomTags {
		fmt.Printf("Tag Name: %s, Description: %s\n", tag.Name, tag.Description)
	}

	err2 := addRegistryEntries(*config)
	if err2 != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Registry entries added successfully")
	}
}
