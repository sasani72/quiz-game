package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Load(path string) *Config {
	var k = koanf.New(".")

	k.Load(confmap.Provider(defaultConfig, "."), nil)

	k.Load(file.Provider(path), yaml.Parser())

	k.Load(env.Provider("GAMEAPP_", ".", func(s string) string {
		str := strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "GAMEAPP_")), "_", ".", -1)
		return strings.Replace(str, "..", "_", -1)
	}), nil)

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
