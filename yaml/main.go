package main

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"strconv"
	"strings"
)

//go:embed ingress.yaml
var ingressyaml []byte

func main() {
	fmt.Println(string(ingressyaml))
	fmt.Printf("%T\n", ingressyaml)
	fmt.Println(reflect.TypeOf(ingressyaml))
	//	var yamlInterface map[string](interface{})
	var yamlInterface interface{}
	err := yaml.Unmarshal(ingressyaml, &yamlInterface)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", yamlInterface)
	fmt.Println(reflect.TypeOf(yamlInterface))
	metadata := yamlInterface.(map[string]interface{})["metadata"]
	fmt.Printf("%T\n", metadata)
	fmt.Println(reflect.TypeOf(metadata))
	annotations := metadata.(map[string]interface{})
	fmt.Printf("%T\n", annotations)
	fmt.Println(reflect.TypeOf(annotations))

	rules := yamlInterface.(map[string]interface{})["spec"].(map[string]interface{})["rules"].([]interface{})
	fmt.Printf("%T\n", rules)
	fmt.Println(reflect.TypeOf(rules))
	fmt.Println(rules)
	fmt.Println("PROCESSING")
	// The case where it's only a list?
	// ProcessNode(rules, []string{"[0]"})
	keys := splitKeysCorrectly("spec.rules[0]")
	returnedfield, err := ProcessNode(yamlInterface, keys)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("returned field: %v\n", returnedfield)
}

func splitKeysCorrectly(keys string) []string {
	keylist := strings.FieldsFunc(keys, func(r rune) bool {
		return r == '.' || r == '['
	})
	return keylist
}

func ProcessNode(yamlInput interface{}, keys []string) (interface{}, error) {
	currentKey := keys[0]
	fmt.Printf("current key: %v\n", currentKey)

	if len(keys) == 0 {
		fmt.Println("no more keys")
		return nil, nil
	}
	switch yamlNodeType := yamlInput.(type) {
	case map[string]interface{}:
		fmt.Printf("%T\n", yamlInput)
		fmt.Println(reflect.TypeOf(yamlInput))
		fmt.Println(yamlNodeType)
		fmt.Println(yamlInput.(map[string]interface{})[currentKey])
		nextYamlInput := yamlInput.(map[string]interface{})[currentKey]
		if len(keys) == 1 {
			fmt.Println("no more keys")
			fmt.Println(nextYamlInput)
			return nextYamlInput, nil
		}
		return ProcessNode(nextYamlInput, keys[1:])
	case []interface{}:
		fmt.Printf("%T\n", yamlInput)
		fmt.Println(reflect.TypeOf(yamlInput))
		fmt.Println(yamlNodeType)
		currentIndex, err := strconv.Atoi(currentKey)
		if err != nil {
			fmt.Printf("cannot convert %v to integer: %v", currentKey, err)
		}
		nextYamlInput := yamlInput.([]interface{})[currentIndex]
		if len(keys) == 1 {
			fmt.Println("no more keys")
			fmt.Println(nextYamlInput)
			return nextYamlInput, nil
		}
		return ProcessNode(nextYamlInput, keys[1:])
	default:
		fmt.Println("error now hahaha")
		fmt.Printf("%T\n", yamlInput)
		fmt.Println(reflect.TypeOf(yamlInput))
		fmt.Println(yamlNodeType)
	}
	return nil, nil

}
