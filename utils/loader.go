package utils

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadFromYaml(file string) ([]string, error) {
	s := make([]string, 0, 20)

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &s)
	return s, err
}
