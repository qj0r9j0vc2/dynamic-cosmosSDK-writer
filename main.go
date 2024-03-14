package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

type CustomConfig struct {
	// ConfigToml specifies parameter(s) in config.toml
	ConfigToml Entry `yaml:"configToml"`
	// ConfigToml specifies parameter(s) in config.toml
	AppToml Entry `yaml:"appToml"`
}

type Entry interface{}

func MergeWithEntryList(entryList []Entry, parameterMap map[string]interface{}) map[string]interface{} {
	mapInterface := make(map[string]interface{})
	for _, entry := range entryList {
		mapInterface = mergeMaps(mapInterface, mergeInterfaceConverter(entry, parameterMap))
	}
	return mapInterface
}

func mergeInterfaceConverter(entry interface{}, parameterMap map[string]interface{}) map[string]interface{} {
	var (
		convertValue interface{}
		stop         = false
	)

	convert, ok := entry.(map[string]interface{})
	if ok {
		for key, value := range convert {
			for paramKey, paramValue := range parameterMap {
				if paramKey == key {
					if _, hasChild := value.(map[string]interface{}); !hasChild {
						convert[key] = paramValue
					} else {
						pv, pvOk := paramValue.(map[string]interface{})
						if pvOk {
							t := mergeInterfaceConverter(value, pv)
							convert[key] = t
						}
					}
					stop = true
					delete(parameterMap, paramKey)
				}
			}

			if stop {
				stop = false
			} else if convertValue, ok = value.(string); ok {
				convert[key] = convertValue
			} else if convertValue, ok = value.(int); ok {
				convert[key] = convertValue
			} else if convertValue, ok = value.(bool); ok {
				convert[key] = convertValue
			} else {
				convert[key] = interfaceConverter(value)
			}
		}

		return convert
	} else if reflect.TypeOf(entry).Kind() == reflect.Slice {
		s := reflect.ValueOf(entry)
		for i := 0; i < s.Len(); i++ {
			entry = interfaceConverter(s.Index(i).Interface())
		}
		return entry.(map[string]interface{})
	} else {
		return map[string]interface{}{}
	}
}

// mergeMaps merge maps as args.
// If maps' key duplicated, it recovered by latest key - value.
func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for key, value := range m {
			result[key] = value
		}
	}
	return result
}

func interfaceConverter(entry interface{}) map[string]interface{} {
	var convertValue interface{}

	convert, ok := entry.(map[string]interface{})
	if ok {
		for key, value := range convert {
			if convertValue, ok = value.(string); ok {
				convert[key] = convertValue
			} else if convertValue, ok = value.(int); ok {
				convert[key] = convertValue
			} else if convertValue, ok = value.(bool); ok {
				convert[key] = convertValue
			} else {
				convert[key] = interfaceConverter(value)
			}
		}
		return convert
	} else if reflect.TypeOf(entry).Kind() == reflect.Slice {
		s := reflect.ValueOf(entry)
		for i := 0; i < s.Len(); i++ {
			entry = interfaceConverter(s.Index(i).Interface())
		}
		return entry.(map[string]interface{})
	} else {
		return map[string]interface{}{}
	}
}

func parse() {
	f, _ := os.ReadFile("out/config.toml")
	var (
		targetConfig interface{}
		readConfig   CustomConfig
		err          error
	)
	_ = toml.Unmarshal(f, &targetConfig)

	e := Entry(targetConfig)

	readConfigBytes, _ := os.ReadFile("config.yaml")

	err = yaml.Unmarshal(readConfigBytes, &readConfig)
	if err != nil {
		panic(err)
	}

	readConfigTomlConfig, ok := readConfig.ConfigToml.(map[string]interface{})
	if !ok {
		fmt.Printf("Nothing\n")
		return
	}
	formatedConfig := MergeWithEntryList([]Entry{e}, readConfigTomlConfig)

	out, _ := toml.Marshal(formatedConfig)
	testFile, _ := os.Create("test-config.toml")
	testFile.Write(out)

}

func main() {
	parse()
}
