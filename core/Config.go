package core

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/apex/log"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/go-yaml/yaml"
)

type templateData struct {
	Env map[string]string
}

type inValid struct {
	Name string
	Type string
}

//服务端配置数据结构
type AppConfig struct {
	RunMode  bool
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	//RenamedC int   `yaml:"c"`
	//D        []int `yaml:",flow"`
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Type     string
	User     string
	Password string
	Host     string
	DbName   string
	Charset  string
}

type RedisConfig struct {
	Host     string
	Password string
	Db       int
}

var (
	//https://studygolang.com/articles/4490
	debugMode  = flag.Bool("d", false, "debug mode")
	configPath = flag.String("c", "app.yaml", "config path")
	//
	appConfig *AppConfig
)

//加载服务端配置
func init() {
	err := LoadConfig(*configPath, &appConfig, false)
	if nil != err {
		log.Fatal(err.Error())
	}
	//
	appConfig.RunMode = *debugMode
}

func GetAppConfig() *AppConfig {
	return appConfig
}

func GetServerConfig() ServerConfig {
	return appConfig.Server
}

func GetDatabaseConfig() DatabaseConfig {
	return appConfig.Database
}

func GetRedisConfig() RedisConfig {
	return appConfig.Redis
}

func LoadConfig(configPath string, configStruct interface{}, valid bool) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	file, err = substitute(file)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(file, configStruct); err != nil {
		return err
	}

	if valid {
		if err = validate(configStruct); err != nil {
			return err
		}
	}

	return nil
}

func substitute(in []byte) ([]byte, error) {
	t, err := template.New("config").Parse(string(in))
	if err != nil {
		return nil, err
	}

	data := &templateData{
		Env: make(map[string]string),
	}

	values := os.Environ()
	for _, val := range values {
		keyval := strings.SplitN(val, "=", 2)
		if len(keyval) != 2 {
			continue
		}
		data.Env[keyval[0]] = keyval[1]
	}

	buffer := &bytes.Buffer{}
	if err = t.Execute(buffer, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func validate(object interface{}) error {
	if valid := validateValue(object); valid != nil {
		return errors.New(fmt.Sprintf("Missing required config field: %v of type %s", valid.Name, valid.Type))
	}
	return nil
}

func validateValue(object interface{}) *inValid {
	objType := reflect.TypeOf(object)
	objValue := reflect.ValueOf(object)
	// If object is a nil interface value, TypeOf returns nil.
	if objType == nil {
		// Don't validate nil interfaces
		return nil
	}

	switch objType.Kind() {
	case reflect.Ptr:
		// If the ptr is nil
		if objValue.IsNil() {
			return &inValid{Type: objType.String()}
		}
		// De-reference the ptr and pass the object to validate
		return validateValue(objValue.Elem().Interface())
	case reflect.Struct:
		for idx := 0; idx < objValue.NumField(); idx++ {
			if valid := validateValue(objValue.Field(idx).Interface()); valid != nil {
				field := objType.Field(idx)
				// Capture sub struct names
				if valid.Name != "" {
					field.Name = field.Name + "." + valid.Name
				}

				// If our field is a pointer and it's pointing to an object
				if field.Type.Kind() == reflect.Ptr && !objValue.Field(idx).IsNil() {
					// The optional doesn't apply because our field does exist
					// instead the de-referenced object failed validation
					if field.Tag.Get("config") == "optional" {
						return &inValid{Name: field.Name, Type: valid.Type}
					}
				}
				// If the field is optional, don't invalidate
				if field.Tag.Get("config") != "optional" {
					return &inValid{Name: field.Name, Type: valid.Type}
				}
			}
		}
	// no way to tell if boolean or integer fields are provided or not
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool, reflect.Interface,
		reflect.Func:
		return nil
	default:
		if objValue.Len() == 0 {
			return &inValid{Type: objType.Name()}
		}
	}
	return nil
}
