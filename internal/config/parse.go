package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

type Config struct {
	Service struct {
		RefreshInterval int      `yaml:"refresh_interval"`
		Chains          []string `yaml:"supported_chains"`
	}
	Db struct {
		Name string `yaml:"name"`
	}
	Telegram struct {
		Token string `yaml:"token"`
	}
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Println("Current working directory:", dir)
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("no .env file found")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file, err = replaceEnvVars(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ParseCLI() (string, error) {
	var path string

	flag.StringVar(&path, "config", "./config.yaml", "path to config file")
	flag.Parse()

	s, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("cannot get stat for %s", path)
	}
	if s.IsDir() {
		return "", fmt.Errorf("'%s' is a directory, not a normal file", path)
	}

	return path, nil
}

func replaceEnvVars(input []byte) ([]byte, error) {
	envVarRegexp := regexp.MustCompile(`\$\{(\w+)\}`)
	return envVarRegexp.ReplaceAllFunc(input, func(match []byte) []byte {
		key := string(match[2 : len(match)-1])
		return []byte(os.Getenv(key))
	}), nil
}
