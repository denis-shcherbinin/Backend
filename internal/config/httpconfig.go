package config

import (
	"github.com/spf13/viper"
	"time"
)

const (
	httpPort               = "http.port"
	httpReadTimeout        = "http.read_timout"
	httpWriteTimeout       = "http.write_timeout"
	httpMaxHeaderMegabytes = "http.max_header_bytes"
)

type HTTPConfig struct {
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

func httpConfigInit() *HTTPConfig {
	return &HTTPConfig{
		Port:               viper.GetString(httpPort),
		ReadTimeout:        viper.GetDuration(httpReadTimeout),
		WriteTimeout:       viper.GetDuration(httpWriteTimeout),
		MaxHeaderMegabytes: viper.GetInt(httpMaxHeaderMegabytes),
	}
}
