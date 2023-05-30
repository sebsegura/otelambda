package env

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
)

type Vars struct {
	Env     string `env:"ENV" default:"dev" json:"env"`
	Debug   bool   `env:"DEBUG" default:"false" json:"debug"`
	BaseURL string `env:"BASE_URL" default:"" json:"base_url"`
}

var (
	_v    *Vars
	_once sync.Once
)

func GetVars() *Vars {
	_once.Do(func() {
		setup()
	})
	return _v
}

func setup() {
	_v = &Vars{}

	ref := reflect.ValueOf(_v).Elem()
	if err := loadVars(ref, ref.Type()); err != nil {
		log.Println(err)
	}
}

func loadVar(field reflect.Value, fieldDef reflect.StructField) error {
	var (
		err          error
		envField     = fieldDef.Tag.Get("env")
		defaultValue = fieldDef.Tag.Get("default")
		value        = os.Getenv(envField)
	)

	if len(value) == 0 {
		value = defaultValue
	} else {
		switch field.Type().Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				log.Println("invalid config")
			} else {
				field.SetBool(boolValue)
			}
		}
	}

	return err
}

func loadVars(value reflect.Value, valueType reflect.Type) error {
	var err error
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Type().Kind() == reflect.Struct {
			err = loadVars(field, field.Type())
		} else {
			err = loadVar(field, valueType.Field(i))
		}
	}
	return err
}
