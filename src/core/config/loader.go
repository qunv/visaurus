package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/qunv/visaurus/core/utils"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"reflect"
	"strings"
)

type Loader interface {
	Bind(properties ...Properties) error
}

type ViperLoader struct {
	viper          *viper.Viper
	option         Option
	groupedConfig  map[string]interface{}
	decodeHookFunc mapstructure.DecodeHookFunc
	validate       *validator.Validate
}

func NewLoader(option Option, properties []Properties) (Loader, error) {
	setDefaultOption(&option)
	reader, err := NewDefaultProfileReader(option.ConfigPaths, option.ConfigFormat, option.KeyDelimiter)
	if err != nil {
		return nil, err
	}
	vi, err := loadViper(reader, option, properties)
	if err != nil {
		return nil, err
	}
	return &ViperLoader{
		viper:         vi,
		option:        option,
		groupedConfig: groupPropertiesConfig(vi, properties, option),
		decodeHookFunc: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			MapStructurePlaceholderValueHook(),
		),
		validate: validator.New(),
	}, nil
}

func (l *ViperLoader) Bind(propertiesList ...Properties) error {
	for _, props := range propertiesList {
		propsName := reflect.TypeOf(props).String()
		// Run pre-binding life cycle
		if propsPreBind, ok := props.(PropertiesPreBinding); ok {
			if err := propsPreBind.PreBinding(); err != nil {
				return err
			}
		}

		if err := l.decodeWithDefaults(props); err != nil {
			return errors.WithMessage(err,
				fmt.Sprintf("[Visaurus-error] Error when decode config key [%s] to [%s]", props.Prefix(), propsName))
		}

		if err := l.validateProps(props); err != nil {
			return errors.WithMessage(err,
				fmt.Sprintf("[Visaurus-error] Error when validate properties [%s]", propsName))
		}

		// Run post-binding life cycle
		if propsPostBind, ok := props.(PropertiesPostBinding); ok {
			if err := propsPostBind.PostBinding(); err != nil {
				return err
			}
		}
		l.option.DebugFunc("[Visaurus-debug] Properties [%s] was loaded with prefix [%s]", propsName, props.Prefix())
	}
	return nil
}

func (l *ViperLoader) validateProps(props Properties) error {
	return l.validate.Struct(props)
}

func (l ViperLoader) decodeWithDefaults(props Properties) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       l.decodeHookFunc,
		WeaklyTypedInput: true,
		Result:           props,
	})
	if err != nil {
		return err
	}
	loadedCfMap, ok := l.groupedConfig[normalizeKey(props.Prefix())].(map[string]interface{})
	if !ok {
		return errors.New("loaded config inside prefix is not a map")
	}
	if err := decoder.Decode(loadedCfMap); err != nil {
		return errors.WithMessage(err, "cannot decode props")
	}

	// Set default value if its missing
	if err := defaults.Set(props); err != nil {
		return errors.WithMessage(err, "cannot set default")
	}

	propsMap := make(map[string]interface{})
	if err := mapstructure.Decode(props, &propsMap); err != nil {
		return errors.WithMessage(err, "cannot decode props to map")
	}

	// Because mapstructure cannot decode struct inside slice, so convert to yaml and unmarshal again
	propsMapBytes, err := yaml.Marshal(propsMap)
	if err != nil {
		return errors.WithMessage(err, "cannot encode propsMap to yaml")
	}
	var defaultProps map[string]interface{}
	if err := yaml.Unmarshal(propsMapBytes, &defaultProps); err != nil {
		return errors.WithMessage(err, "cannot unmarshal propsMapBytes")
	}

	newPropsMap := utils.MergeCaseInsensitiveMaps(loadedCfMap, defaultProps)
	if err := decoder.Decode(newPropsMap); err != nil {
		return errors.New("cannot decode props again")
	}
	return nil
}

func loadViper(reader ProfileReader, option Option, propertiesList []Properties) (*viper.Viper, error) {
	option.DebugFunc("[Visaurus-debug] Loading active profiles [%s] in paths [%s] with format [%s]",
		strings.Join(option.ActiveProfiles, ", "), strings.Join(option.ConfigPaths, ", "), option.ConfigFormat)

	vi := viper.NewWithOptions(viper.KeyDelimiter(option.KeyDelimiter))
	vi.SetEnvKeyReplacer(strings.NewReplacer(option.KeyDelimiter, "_"))
	vi.AutomaticEnv()

	if err := discoverActiveProfiles(vi, reader, option); err != nil {
		return nil, err
	}

	if err := discoverEnvKeys(vi, option, propertiesList); err != nil {
		return nil, err
	}
	return vi, nil
}

// discoverEnvKeys Discover env keys for multiple properties at once
func discoverEnvKeys(vi *viper.Viper, option Option, propertiesList []Properties) error {
	for _, props := range propertiesList {
		propsName := reflect.TypeOf(props).String()
		propsMap := make(map[string]interface{})
		if err := mapstructure.Decode(props, &propsMap); err != nil {
			return errors.WithMessage(err, "cannot decode props to map")
		}

		defaultMap := convertSliceToNestedMap(strings.Split(normalizeKey(props.Prefix()), option.KeyDelimiter), propsMap, nil)
		for key, env := range buildEnvKeys(defaultMap, option.KeyDelimiter, "_", "", "") {
			if err := vi.BindEnv(key, env); err != nil {
				return fmt.Errorf("[Visaurus-error] Error when build env keys properties [%s]: %v", propsName, err)
			}
		}
		option.DebugFunc("[Visaurus-debug] Default value was discovered for properties [%s]", propsName)
	}
	return nil
}

// discoverActiveProfiles Discover values for multiple active profiles at once
func discoverActiveProfiles(vi *viper.Viper, reader ProfileReader, option Option) error {
	debugPaths := strings.Join(option.ConfigPaths, ", ")
	for _, activeProfile := range option.ActiveProfiles {
		cfMap, err := reader.Read(activeProfile)
		if err != nil {
			return err
		}
		if err := vi.MergeConfigMap(cfMap); err != nil {
			return fmt.Errorf("[Visaurus-error] Error when read active profile [%s] in paths [%s]: %v",
				activeProfile, debugPaths, err)
		}
		option.DebugFunc("[Visaurus-debug] Active profile [%s] was loaded", activeProfile)
	}
	return nil
}

func groupPropertiesConfig(vi *viper.Viper, propertiesList []Properties, option Option) map[string]interface{} {
	allSettings := vi.AllSettings()
	group := make(map[string]interface{})
	for _, props := range propertiesList {
		normalizedPrefix := normalizeKey(props.Prefix())
		m := utils.DeepSearchInMap(allSettings, normalizedPrefix, option.KeyDelimiter)
		correctedVal, exists := correctSliceValues(vi, normalizedPrefix, option.KeyDelimiter, m)
		if exists {
			group[normalizedPrefix] = correctedVal
		} else {
			group[normalizedPrefix] = m
		}
	}
	return group
}

func correctSliceValues(vi *viper.Viper, prefix string, delim string, val interface{}) (interface{}, bool) {
	if slice, ok := val.([]interface{}); ok {
		for k, v := range slice {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%d", prefix, delim, k), delim, v)
			if exists {
				slice[k] = correctedVal
			}
		}
	} else if m, ok := val.(map[interface{}]interface{}); ok {
		for k, v := range m {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%s", prefix, delim, k), delim, v)
			if exists {
				m[k] = correctedVal
			}
		}
	} else if m, ok := val.(map[string]interface{}); ok {
		for k, v := range m {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%s", prefix, delim, k), delim, v)
			if exists {
				m[k] = correctedVal
			}
		}
	} else {
		if correctedVal := vi.Get(prefix); correctedVal != nil {
			return correctedVal, true
		}
	}
	return nil, false
}

func convertSliceToNestedMap(paths []string, endVal interface{}, inMap map[interface{}]interface{}) map[interface{}]interface{} {
	if inMap == nil {
		inMap = map[interface{}]interface{}{}
	}
	if len(paths) == 0 {
		return inMap
	}
	if len(paths) == 1 {
		inMap[paths[0]] = endVal
		return inMap
	}
	inMap[paths[0]] = convertSliceToNestedMap(paths[1:], endVal, map[interface{}]interface{}{})
	return inMap
}

func buildEnvKeys(data map[interface{}]interface{}, keyDelim string, envDelim, baseKey string, baseEnv string) map[string]string {
	keyEnvMap := make(map[string]string)
	if data == nil {
		return keyEnvMap
	}
	for key, val := range data {
		keyStr, ok := key.(string)
		if !ok {
			continue
		}
		kPath := keyStr
		ePath := keyStr
		if baseKey != "" {
			kPath = baseKey + keyDelim + keyStr
			ePath = baseEnv + envDelim + keyStr
		}
		ePath = strings.ToUpper(ePath)
		switch valT := val.(type) {
		case map[string]interface{}:
			if len(valT) == 0 {
				keyEnvMap[kPath] = ePath
				break
			}
			mapI := make(map[interface{}]interface{})
			for k1, v1 := range valT {
				mapI[k1] = v1
			}
			for sk, sv := range buildEnvKeys(mapI, keyDelim, envDelim, kPath, ePath) {
				keyEnvMap[sk] = sv
			}
			break
		case map[interface{}]interface{}:
			if len(valT) == 0 {
				keyEnvMap[kPath] = ePath
				break
			}
			for sk, sv := range buildEnvKeys(valT, keyDelim, envDelim, kPath, ePath) {
				keyEnvMap[sk] = sv
			}
			break
		default:
			keyEnvMap[kPath] = ePath
			break
		}
	}
	return keyEnvMap
}

func normalizeKey(key string) string {
	return strings.ToLower(key)
}
