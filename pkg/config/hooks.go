package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

// TimeLocationDecodeHook is a function to transform time.Location values.
func TimeLocationDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		var timeLocation *time.Location
		if t != reflect.TypeOf(timeLocation) {
			return data, nil
		}

		str, ok := data.(string)
		if !ok {
			return data, nil
		}

		return time.LoadLocation(str)
	}
}

// CMQTypeDecodeHook is a function to transform CMQType values.
func CMQTypeDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		var ct CMQType
		if t != reflect.TypeOf(ct) {
			return data, nil
		}

		str, ok := data.(string)
		if !ok {
			return data, nil
		}

		if c, ok := map[string]CMQType{
			CMQNatsStreaming.String(): CMQNatsStreaming,
			CMQJetStream.String():     CMQJetStream,
			CMQNats.String():          CMQNats,
		}[str]; ok {
			return c, nil
		}

		return nil, fmt.Errorf("invalid cmq type")
	}
}
