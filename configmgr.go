package supercmd

import "github.com/pkg/errors"

// ConfigMgr provides wraps Config to provide more useful access methods.
type ConfigMgr struct {
	Config
}

// GetString gets an optional string from the config.
func (c *ConfigMgr) GetString(key string) (string, bool, error) {
	v, ok := c.GetValue(key)
	if !ok {
		return "", false, nil
	}
	str, ok := v.(string)
	if !ok {
		return "", false, errors.Errorf("expected a string for key '%s', got '%s' (%T)")
	}
	return str, true, nil
}

// RequireString gets an required string from the config.
func (c *ConfigMgr) RequireString(key string) (string, error) {
	str, ok, err := c.GetString(key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.Errorf("key '%s' is required")
	}
	return str, nil
}
