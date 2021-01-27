package configs

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func New(fileName string) Config {
	if fileName == "" {
		panic("Not found: fileName is empty")
	}

	workDir, err := os.Getwd()
	workDir = filepath.ToSlash(workDir) // for window os
	if err != nil {
		panic(err)
	}

	configPath := ""
	moduleDir := []string{"ratelimiter"}
	for _, dir := range moduleDir {
		if strings.Contains(workDir, dir) {
			path := strings.Split(workDir, dir)
			configPath = path[0] + dir
			break
		}
	}

	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName(fileName)
	vp.AddConfigPath(configPath)
	if err := vp.ReadInConfig(); err != nil {
		panic(err)
	}

	cfg := new(Config)
	option := func(c *mapstructure.DecoderConfig) { c.TagName = "configs" }
	if err := vp.Unmarshal(cfg, option); err != nil {
		panic(err)
	}
	return *cfg
}

type Config struct {
	Port    string  `configs:"port"`
	Limiter Limiter `configs:"rate_limiter"`
	Redis   Redis   `configs:"redis"`
}

type Limiter struct {
	MaxLimitCount             int32 `configs:"max_limit_count"`
	ResetCountIntervalSeconds int64 `configs:"reset_count_interval"`
}

func (l *Limiter) ResetCountInterval() time.Duration {
	return time.Duration(l.ResetCountIntervalSeconds) * time.Second
}

type Redis struct {
	Host     string `configs:"host"`
	Port     string `configs:"port"`
	Password string `configs:"password"`
}

func (r *Redis) Address() string {
	return r.Host + ":" + r.Port
}
