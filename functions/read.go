package functions

import (
	"log"
	// "path"
	"gopkg.in/yaml.v3"
)


type Yaml struct {
	Name string `yaml:"name"`
	Github string `yaml:"github"`
	Location string `yaml:"location"`
	Install string `yaml:"install"`
	Build string `yaml:"build"`
	Run string `yaml:"run"`
	Port string `yaml:"port"`
}

func ReadYaml(file []byte) (Yaml,error)  {

	var fileContent Yaml
	if err := yaml.Unmarshal(file, &fileContent); err != nil {
		log.Fatalf("Failed to read app.yaml. Error %v\n",err)
		return Yaml{}, err
	}

	return fileContent, nil
}