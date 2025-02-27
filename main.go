package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	cmd := exec.Command("cmd", "/c", cfg.TaniumInstaller)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return fmt.Errorf("failed to execute replace.exe: %s, output: %s, error: %v", cmd.String(), output, err)
	}
	return nil
}

func removeAllFilesInExeDir() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	files, err := os.ReadDir(exeDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(exeDir, file.Name())
		err2300 := os.Remove(filePath)
		if err2300 != nil {
			fmt.Println("Error occurred removing file " + file.Name())
		}
	}

	return nil
}
func updateScheduleTemplate(templatePath, outputPath string) error {
	// Read the template file
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}

	// Get the current time and calculate the replacement times
	startTime := time.Now().Add(2 * time.Minute).Format(time.RFC3339)
	endTime := time.Now().AddDate(1, 0, 0).Format(time.RFC3339)

	// Get the path to the executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Replace the placeholders in the template
	updatedData := strings.ReplaceAll(string(templateData), "##REPLACE START##", startTime)
	updatedData = strings.ReplaceAll(updatedData, "##REPLACE END##", endTime)
	updatedData = strings.ReplaceAll(updatedData, "##REPLACE FILE PATH##", exePath)

	// Write the updated data to the output file
	err = os.WriteFile(outputPath, []byte(updatedData), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}

func createScheduledTaskFromXML(xmlFilePath string) error {
	cmd := exec.Command("schtasks", "/create", "/tn", "YourTaskName", "/xml", xmlFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create scheduled task: %v, output: %s", err, output)
	}
	fmt.Println("Scheduled task created successfully")
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
		return
	} else {
		fmt.Println("Registry entries added successfully")
	}

	err232 := updateScheduleTemplate("template.xml", "schedule.xml")
	if err232 != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Schedule updated successfully")
	}
	err300 := createScheduledTaskFromXML("schedule.xml")
	if err300 != nil {
		fmt.Printf("Error: %v\n", err)
	}

	removeAllFilesInExeDir()

	var input string
	fmt.Print("Enter your input: ")
	fmt.Scanln(&input)
	fmt.Printf("You entered: %s\n", input)
}
