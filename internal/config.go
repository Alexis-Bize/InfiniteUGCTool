package internal

import (
	"bufio"
	"embed"
	"strings"
)

type Config struct {
	Name 	string
	Version	string
	GitHub	string
}

var config Config

func LoadConfig(f embed.FS) error {
	file, err := f.Open("config.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, "\"")
			switch key {
			case "APP_NAME":
				config.Name = value
			case "APP_VERSION":
				config.Version = value
			case "APP_GITHUB":
				config.GitHub = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return config
}
