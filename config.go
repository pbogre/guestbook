package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

type Config struct {
    Title               string  `yaml:"title"`
    GlobalRateLimit     int     `yaml:"globalRateLimit"`
    GlobalBurstLimit    int     `yaml:"globalBurstLimit"`
    EntriesPerPage      int     `yaml:"entriesPerPage"`
}

var GuestbookConfig Config

func loadConfig(path string) {
    data, err := os.ReadFile(path)
    if err != nil {
        log.Fatal(err)
        return
    }

    if err := yaml.Unmarshal(data, &GuestbookConfig); err != nil {
        log.Fatal(err)
    }

    if err := validateConfig(GuestbookConfig); err != nil {
        log.Fatal(err)
    }
}

// to validate, make sure all entries have a value
func validateConfig(config interface{}) error {
    v := reflect.ValueOf(config)
    t := reflect.TypeOf(config)

    if v.Kind() == reflect.Pointer {
        v = v.Elem()
        t = t.Elem()
    }

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        if field.IsZero() {
            return fmt.Errorf("Missing or zero value for field '%s'", fieldType.Name)
        }
    }

    return nil
}
