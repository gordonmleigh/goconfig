package supercmd

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config provides access to configuration settings.
type Config interface {
	GetValue(key string) (interface{}, bool)
}

// FallbackConfig is an implementation of the Config interface which checks
// the first config for a value first, then falls back to the subsequent configs
// until a value is found.
type FallbackConfig struct {
	Configs []Config
}

// FallbackConfig implements the Config interface.
var _ Config = &FallbackConfig{}

// MakeFallbackConfig makes a fallback config.
func MakeFallbackConfig(cfg ...Config) Config {
	return &FallbackConfig{
		Configs: cfg,
	}
}

// GetValue gets the value associated with the given key.
func (c *FallbackConfig) GetValue(key string) (interface{}, bool) {
	for _, cfg := range c.Configs {
		v, ok := cfg.GetValue(key)
		if ok {
			return v, true
		}
	}
	return nil, false
}

// AddFirst creates a new config with the given config looked up first.
func (c *FallbackConfig) AddFirst(cfg Config) *FallbackConfig {
	return &FallbackConfig{Configs: append([]Config{cfg}, c.Configs...)}
}

// AddLast creates a new config with the given config looked up last.
func (c *FallbackConfig) AddLast(cfg Config) Config {
	return &FallbackConfig{
		Configs: append(append([]Config{}, c.Configs...), cfg),
	}
}

// tieredMap is a map of maps containing config.
type tieredMap map[string]interface{}

// tieredMap implements Config interface.
var _ Config = new(tieredMap)

// FromTieredMap reads config settings from a map of values, where each value
// might be a map that has child values.
func FromTieredMap(m map[string]interface{}) Config {
	return tieredMap(m)
}

// GetValue gets a value from the config.
func (m tieredMap) GetValue(key string) (interface{}, bool) {
	var v interface{} = map[string]interface{}(m)

	for _, k := range strings.Split(key, ".") {
		if k == "" {
			return nil, false
		}

		haveValue := false

		if level, ok := v.(map[string]interface{}); ok {
			v, haveValue = level[k]
		} else if level, ok := v.(map[interface{}]interface{}); ok {
			// special case for YAML deserialising
			v, haveValue = level[k]
		} else if r := reflect.ValueOf(v); r.Kind() == reflect.Array || r.Kind() == reflect.Slice {
			i, err := strconv.Atoi(k)
			if err != nil || i < 0 || i >= r.Len() {
				return nil, false
			}
			v, haveValue = r.Index(i).Interface(), true
		}

		if !haveValue {
			return nil, false
		}
	}

	return v, true
}

// flatMap is a flat map of values by dot-seperated key.
type flatMap map[string]interface{}

// flatMap implements config.
var _ Config = new(flatMap)

// GetValue gets a value from the config.
func (m flatMap) GetValue(key string) (interface{}, bool) {
	v, ok := (map[string]interface{}(m))[key]
	return v, ok
}

// FromJSON reads config settings from JSON.
func FromJSON(r io.Reader) (Config, error) {
	d := json.NewDecoder(r)
	m := make(map[string]interface{})
	err := d.Decode(&m)
	if err != nil {
		return nil, err
	}
	return FromTieredMap(m), nil
}

// FromYAML reads config settings from YAML.
func FromYAML(r io.Reader) (Config, error) {
	d := yaml.NewDecoder(r)
	m := make(map[string]interface{})
	err := d.Decode(&m)
	if err != nil {
		return nil, err
	}
	return FromTieredMap(m), nil
}

// FromArgs reads config settings from commandline arguments.
func FromArgs(args []string) (Config, error) {
	m := make(map[string]interface{})

	for _, f := range args {
		if f == "--" {
			continue
		}
		assign := strings.SplitN(strings.TrimLeft(f, "-"), "=", 2)
		if len(assign) != 2 {
			return nil, errors.Errorf("expected assignment, got: %s", f)
		}
		m[assign[0]] = assign[1]
	}

	return flatMap(m), nil
}

// envConfig reads config from environment variables.
type envConfig struct {
	prefix  string
	lookup  func(string) (string, bool)
	list    func() []string
	convert func(string) string
}

// FromEnv reads config from environment variables.
func FromEnv(prefix string) (Config, error) {
	return &envConfig{
		prefix:  prefix,
		lookup:  os.LookupEnv,
		convert: DotToSnake,
	}, nil
}

// Get gets a value from the config.
func (e *envConfig) GetValue(key string) (interface{}, bool) {
	return e.lookup(e.convert(key))
}

// DotToSnake converts dot-seperated identifiers to underscore-seperated
// identifiers.
func DotToSnake(str string) string {
	// make mutable copy of str
	chars := []rune(str)
	for i, c := range chars {
		if c == '.' {
			chars[i] = '_'
		}
	}
	return string(chars)
}
