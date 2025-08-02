package config

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// LoadConfig Load config from file
func LoadConfig[T any](
	envPrefix, configFilename string, builtinConfig []byte, c T,
	decodeHookFuncs ...mapstructure.DecodeHookFunc,
) (T, error) {
	cPtr := &c
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if err := v.ReadConfig(bytes.NewReader(builtinConfig)); err != nil {
		return *cPtr, fmt.Errorf("failed on config initialization: %v", err)
	}

	if configFilename != "" {
		v.SetConfigFile(configFilename)
		if err := v.MergeInConfig(); err != nil {
			return *cPtr, fmt.Errorf("failed on config `%s` merging: %w", configFilename, err)
		}
		log.Printf("config file [%s] opened and merged successfully\n", configFilename)
	}

	hooks := []mapstructure.DecodeHookFunc{
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		TimeLocationDecodeHook(),
		CMQTypeDecodeHook(),
	}
	hooks = append(hooks, decodeHookFuncs...)
	err := v.Unmarshal(cPtr, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(hooks...)))
	if err != nil {
		return *cPtr, fmt.Errorf("failed on config `%s` unmarshal: %w", configFilename, err)
	}

	return *cPtr, nil
}
