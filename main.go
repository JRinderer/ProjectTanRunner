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
	TaniumInstaller  string      `json:"TaniumInstallerName"`
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
	fmt.Println("Starting")
	for _, tag := range cfg.CustomTags {
		//cmdStr := fmt.Sprintf(`reg add "HKLM\SOFTWARE\WOW6432Node\Tanium\Tanium Client\Sensor Data\Tags" /v %s /t REG_SZ /d "%s"`, tag.Name, tag.Description)
		//cmdStr := fmt.Sprintf(`reg add / "HKLM\SOFTWARE\WOW6432Node\Python\PyLauncher" /v %s /t REG_SZ /d "%s"`, tag.Name, tag.Description)
		fmt.Println("Adding custom tag %s\n", tag.Name)
		cmd := exec.Command("reg", "add", `HKLM\SOFTWARE\WOW6432Node\Tanium\Tanium Client\Sensor Data\Tags`, "/v", tag.Name, "/t", "REG_SZ", "/d", tag.Description)

		//cmd := exec.Command("cmd", "/c", cmdStr)
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))
		if err != nil {
			return fmt.Errorf("failed to execute command: %s, output: %s, error: %v", "registry", output, err)
		}
	}

	return nil
}

func runeExe(cfg Config) error {
	cmd := exec.Command("cmd", "/c", cfg.LocationInitFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute replace.exe: %s, output: %s, error: %v", cmd.String(), output, err)
	}
	return nil
}

func main() {
	config, err := readConfig("conf.json")

	if err != nil {
		fmt.Printf("failed to read config: %v", err)
		return
	}

	err23 := runeExe(*config)

	if err23 != nil {
		fmt.Printf("failed to launch Tanium installer " + config.TaniumInstaller)
		return
	}

	err2 := addRegistryEntries(*config)
	if err2 != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Registry entries added successfully")
	}
}
