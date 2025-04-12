package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/fatih/structs"
	"github.com/jeremywohl/flatten"
	"github.com/spf13/viper"

	"github.com/ccbrown/cloud-snitch/backend/api"
	"github.com/ccbrown/cloud-snitch/backend/app"
)

type Config struct {
	App app.Config
	API api.Config
}

// Takes a pointer to a struct and ensures that all fields which are pointers to structs are
// initialized to pointers to their zero value.
func allocStructPointers(obj any) {
	v := reflect.ValueOf(obj)
	for i := 0; i < v.Elem().NumField(); i++ {
		f := v.Elem().Field(i)
		if f.Type().Kind() == reflect.Pointer && f.Type().Elem().Kind() == reflect.Struct {
			newObj := reflect.New(f.Type().Elem())
			allocStructPointers(newObj.Interface())
			f.Set(newObj)
		} else if f.Type().Kind() == reflect.Struct {
			allocStructPointers(f.Addr().Interface())
		}
	}
}

// This loads config values from environment variables, e.g. API_PROXYCOUNT. And any environment
// variables with values of the form "secret:NAME_OR_ARN" are fetched from the AWS secrets manager.
func LoadConfigEnvVariables() error {
	return loadConfigEnvVariablesImpl(os.Getenv, viper.Set)
}

func coerceStringToVar(v string, dest any) (any, error) {
	switch dest.(type) {
	case []byte:
		return base64.StdEncoding.DecodeString(v)
	case []string:
		return strings.Split(v, ","), nil
	default:
		return v, nil
	}
}

func loadConfigEnvVariablesImpl(getenv func(string) string, setvar func(string, any)) error {
	var cfg Config
	allocStructPointers(&cfg)
	confMap := structs.Map(cfg)

	flat, err := flatten.Flatten(confMap, "", flatten.DotStyle)
	if err != nil {
		return err
	}

	var sm *secretsmanager.Client

	for key, dest := range flat {
		envVar := strings.ReplaceAll(strings.ToUpper(key), ".", "_")

		// Load the value from the secrets manager if needed.
		if v := getenv(envVar); strings.HasPrefix(v, "secret:") {
			secretId := v[7:]

			if sm == nil {
				if awsConfig, err := config.LoadDefaultConfig(context.Background()); err != nil {
					return fmt.Errorf("error loading aws config: %w", err)
				} else {
					sm = secretsmanager.NewFromConfig(awsConfig)
				}
			}

			if r, err := sm.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
				SecretId: &secretId,
			}); err != nil {
				return fmt.Errorf("error getting secret for %v: %w", envVar, err)
			} else if r.SecretString == nil {
				return fmt.Errorf("secret for %v is not a string", envVar)
			} else if v, err := coerceStringToVar(*r.SecretString, dest); err != nil {
				return fmt.Errorf("error coercing secret for %v: %w", envVar, err)
			} else {
				setvar(key, v)
			}
		} else if v != "" {
			if v, err := coerceStringToVar(v, dest); err != nil {
				return fmt.Errorf("error coercing %v: %w", envVar, err)
			} else {
				setvar(key, v)
			}
		}
	}

	return nil
}
