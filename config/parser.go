package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Struct that has the yml config file
// MUST be updated every time you add a new entry in the yml file
type Conf struct {
	Storage struct {
		Hostname string `yaml:"hostname"`
		Port     int32  `yaml:"port"`
	}
}

// Parse the config file using the filename parameter and return the configuration as a Conf object
func ParseConfig(filename string) (*Conf, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error while trying to read %s file\n", filename)
		return nil, err
	}

	var conf Conf
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Error while parsing the %s yaml file (Unmarshal)\n", filename)
		return nil, err
	}

	return &conf, nil
}

func Parse(filename string, conf *interface{}) error {
	ymlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error while trying to read %s file\n", filename)
		return err
	}

	err = yaml.Unmarshal(ymlFile, conf)
	return err
}
