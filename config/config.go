package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AuthContexts struct {
	AuthContext map[string]string `yaml:"auth-contexts"`
}

func GetDoTokens() map[string]string {
	yamlFile, err := os.ReadFile("/Users/adamzidiker/Library/Application Support/doctl/config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	authContexts := AuthContexts{}

	if err := yaml.Unmarshal(yamlFile, &authContexts); err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	return authContexts.AuthContext
}

func GetDoTokensName() []string {
	do_tokens := GetDoTokens()
	names := make([]string, 0, len(do_tokens))
	for k := range do_tokens {
		names = append(names, k)
	}
	return names
}

func GetDoToken(context string) string {
	return GetDoTokens()[context]
}
