package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App    app
	Svc    Service
	Http   HTTPServerConfig
	S3     S3
	Crypto Crypto
}

var Cfg *Config

func Load() {
	viper.GetViper().AutomaticEnv()
	// App Defaults
	viper.SetDefault("ENV", "dev")
	viper.SetDefault("LOG_LEVEL", "debug")
	// S3 Provider Default
	viper.SetDefault("S3_PROVIDER", "S3")
	// ARGON 2 Default Config
	// Config for Generating ID:
	// parallelism  memory  rounds  time
	// 12           32      32      0.210900
	viper.SetDefault("ARGON2_KEY_MEMORY", 64)
	viper.SetDefault("ARGON2_KEY_ROUNDS", 12)
	viper.SetDefault("ARGON2_KEY_PARALLELISM", 16)
	// Config for Generating Key:
	// parallelism  memory  rounds  time
	// 16           64      12      0.249461
	viper.SetDefault("ARGON2_ID_MEMORY", 32)
	viper.SetDefault("ARGON2_ID_ROUNDS", 32)
	viper.SetDefault("ARGON2_ID_PARALLELISM", 12)
	// Web View Defaults
	viper.SetDefault("HTTP_LISTEN_PORT", 8080)
	viper.SetDefault("HTTP_LISTEN_HOST", "127.0.0.1")
	viper.SetDefault("HTTP_BASE_URL", "http://localhost:8080")
	viper.SetDefault("HTTP_RETURN_FORMAT", "raw")
	viper.SetDefault("HTTP_API_ENDPOINT", "/snippets")
	viper.SetDefault("HTTP_CORS_LIST", "http://localhost:*")
	cfg := &Config{
		App: app{
			env: viper.GetString("ENV"),
			ll:  viper.GetString("LOG_LEVEL"),
		},
		Svc: Service{
			Overrides: make(map[string]string),
		},
		S3: S3{
			Bucket: viper.GetString("S3_BUCKET"),
			Provider: viper.GetString("S3_PROVIDER"),
		},
		Crypto: Crypto{
			Salt: []byte(viper.GetString("SALT")),
			ARGON2Key: ARGON2Config{
				Parallelism: uint8(viper.GetUint("ARGON2_KEY_PARALLELISM")),
				Memory:      viper.GetUint32("ARGON2_KEY_MEMORY") * 1024,
				Rounds:      viper.GetUint32("ARGON2_KEY_ROUNDS"),
			},
			ARGON2ID: ARGON2Config{
				Parallelism: uint8(viper.GetUint("ARGON2_ID_PARALLELISM")),
				Memory:      viper.GetUint32("ARGON2_ID_MEMORY") * 1024,
				Rounds:      viper.GetUint32("ARGON2_ID_ROUNDS"),
			},
		},
		Http: HTTPServerConfig{
			host:         viper.GetString("HTTP_LISTEN_HOST"),
			port:         viper.GetInt("HTTP_LISTEN_PORT"),
			CORS:         viper.GetString("HTTP_CORS_LIST"),
			baseURL:      viper.GetString("HTTP_BASE_URL"),
			Enpoint:      viper.GetString("HTTP_API_ENDPOINT"),
			returnFormat: viper.GetString("HTTP_RETURN_FORMAT"),
		},
	}

	overrides := viper.GetString("OVERRIDES")
	for _, override := range strings.Split(overrides, ",") {
		entry := strings.Split(override, ":")
		if len(entry) == 2 {
			cfg.Svc.Overrides[entry[0]] = entry[1]
		}
	}

	if len(cfg.Crypto.Salt) == 0 {
		panic("config missing: salt not provided")
	}

	Cfg = cfg
}
