package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"regexp"
	"strings"
)

type Command struct {
	Name      string
	SkipLines int
	Keys      []string
	Regex     string
}

type Commands struct {
	List map[string]Command `toml:"command"`
}

func ParseToJson(command string, output string) string {
	result, err := Parse(command, output)

	if err != nil {
		return fmt.Sprintf(`{"status": "error", "message": "%s"}`, err)
	}

	resultJson, _ := json.Marshal(result)
	return string(resultJson)
}

func Parse(command string, output string) ([]map[string]string, error) {
	commands, err := loadCommandsConfig()

	if err != nil {
		return nil, errors.New("Error loading command config file")
	}

	fmt.Printf("Commands: %q \n", commands.List)
	fmt.Printf("Command: '%s' \n", command)
	fmt.Printf("Output: %s \n", output)
	fmt.Printf("CommandConfig: %q \n",  commands.List[command])

	outputLines := strings.Split(output, "\n")

	if commandConfig, ok := commands.List[command]; ok {
		resultCount := len(outputLines) - commandConfig.SkipLines - 1 // -1 is for last line always empty
		resultMap := make([]map[string]string, resultCount)

		re := regexp.MustCompile(commandConfig.Regex)
		for i, line := range outputLines[commandConfig.SkipLines : len(outputLines)-1] {
			results := re.FindAllStringSubmatch(line, -1)

			if results != nil {
				tmpMap := make(map[string]string)
				for j, key := range commandConfig.Keys {

					if r := results[j][0]; r != "" && key != "" {
						tmpMap[key] = r
					}
				}
				resultMap[i] = tmpMap
			}
		}

		return resultMap, nil
	}

	return nil, errors.New("Command Not Found")
}

func loadCommandsConfig() (Commands, error) {
	var commands Commands
	if _, err := toml.DecodeFile("parser/commands.toml", &commands); err != nil {
		return Commands{}, err
	}

	return commands, nil
}
