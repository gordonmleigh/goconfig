package supercmd

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTieredMap(t *testing.T) {
	example := func(cfg Config, key, desc string, value interface{}) {
		Convey(fmt.Sprintf("When GetValue is called for a %s value", desc), func() {
			v, ok := cfg.GetValue(key)
			Convey("Then the value is returned", func() {
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, value)
			})
		})
	}
	Convey("Given a map of values", t, func() {
		cfg := FromTieredMap(map[string]interface{}{
			"string":  "hello world",
			"integer": 42,
			"section": map[string]interface{}{
				"foo": "bar",
			},
			"array": []map[string]interface{}{
				map[string]interface{}{
					"key":   1,
					"value": "one",
				},
				map[string]interface{}{
					"key":   2,
					"value": "two",
				},
			},
		})

		example(cfg, "string", "string", "hello world")
		example(cfg, "integer", "integer", 42)
		example(cfg, "section.foo", "nested", "bar")
		example(cfg, "array.0.key", "array nested", 1)
	})
}

func TestJSON(t *testing.T) {
	example := func(cfg Config, key, desc string, value interface{}) {
		Convey(fmt.Sprintf("When GetValue is called for a %s value", desc), func() {
			v, ok := cfg.GetValue(key)
			Convey("Then the value is returned", func() {
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, value)
			})
		})
	}
	Convey("Given a JSON config", t, func() {
		jsonCfg := `{
			"string": "hello world",
			"integer": 42,
			"section": {
				"foo": "bar"
			},
			"array": [
				{"key": 1, "value": "one"},
				{"key": 2, "value": "two"}
			]
		}`
		cfg, err := FromJSON(strings.NewReader(jsonCfg))
		So(err, ShouldBeNil)

		example(cfg, "string", "string", "hello world")
		example(cfg, "integer", "integer", 42)
		example(cfg, "section.foo", "nested", "bar")
		example(cfg, "array.0.key", "array nested", 1)
	})
}

func TestYAML(t *testing.T) {
	example := func(cfg Config, key, desc string, value interface{}) {
		Convey(fmt.Sprintf("When GetValue is called for a %s value", desc), func() {
			v, ok := cfg.GetValue(key)
			Convey("Then the value is returned", func() {
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, value)
			})
		})
	}
	Convey("Given a YAML config", t, func() {
		jsonCfg := `string: "hello world"
integer: 42
section:
  foo: bar
array:
- key: 1
  value: one 
- key: 2
  value: two`
		cfg, err := FromYAML(strings.NewReader(jsonCfg))
		So(err, ShouldBeNil)

		example(cfg, "string", "string", "hello world")
		example(cfg, "integer", "integer", 42)
		example(cfg, "section.foo", "nested", "bar")
		example(cfg, "array.0.key", "array nested", 1)
	})
}
