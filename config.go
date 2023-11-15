package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func loadEnvStr(key string, res *string) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	*res = s
}

func loadEnvUint(key string, res *uint) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	si, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	*res = uint(si)
}

type pgConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     uint   `yaml:"port" json:"port"`
	DBName   string `yaml:"db_name" json:"db_name"`
	SslMode  string `yaml:"ssl_mode" json:"ssl_mode"`
	Password string `yaml:"password" json:"password"`
	Username string `yaml:"user" json:"user"`
}

func (p *pgConfig) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d database=%s sslmode=%s user=%s password=%s", p.Host, p.Port, p.DBName, p.SslMode, p.Username, p.Password)
}

func defaultPgConfig() pgConfig {
	return pgConfig{
		Host:     "localhost",
		Port:     5433,
		DBName:   "postgres",
		SslMode:  "disable",
		Password: "mysecret",
		Username: "postgres",
	}
}

func (p *pgConfig) loadFromEnv() {
	loadEnvStr("DB_HOST", &p.Host)
	loadEnvUint("DB_PORT", &p.Port)
	loadEnvStr("DB_NAME", &p.DBName)
	loadEnvStr("DB_SSL", &p.SslMode)
	loadEnvStr("DB_PASSWORD", &p.Password)
	loadEnvStr("DB_USER", &p.Username)
}

type listenConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         uint   `yaml:"port" json:"port"`
	ReadTimeout  uint   `yaml:"read_to" json:"read_to"`
	WriteTimeout uint   `yaml:"write_to" json:"write_to"`
	IdleTimeout  uint   `yaml:"idle_to" json:"idle_to"`
}

func (l listenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

func defaultListenConfig() listenConfig {
	return listenConfig{
		Host:         "127.0.0.1",
		Port:         8080,
		ReadTimeout:  25,
		WriteTimeout: 25,
		IdleTimeout:  300,
	}
}

func (l *listenConfig) loadFromEnv() {
	loadEnvStr("LISTEN_HOST", &l.Host)
	loadEnvUint("LISTEN_PORT", &l.Port)
	loadEnvUint("LISTEN_READ_TIMEOUT", &l.ReadTimeout)
	loadEnvUint("LISTEN_WRITE_TIMEOUT", &l.WriteTimeout)
	loadEnvUint("LISTEN_IDLE_TIMEOUT", &l.IdleTimeout)
}

type jwtConfig struct {
	Secret         string `yaml:"secret" json:"secret"`
	RefreshExpTime uint   `yaml:"refresh_exp" json:"refresh_exp"`
	AccessExpTime  uint   `yaml:"access_exp" json:"access_exp"`
}

func defaultJwtConfig() jwtConfig {
	return jwtConfig{
		Secret:         "mysecret",
		RefreshExpTime: 7,
		AccessExpTime:  15,
	}
}

func (p *jwtConfig) loadFromEnv() {
	loadEnvStr("JWT_SECRET", &p.Secret)
	loadEnvUint("JWT_REFRESH_TOKEN_EXP_TIME", &p.RefreshExpTime)
	loadEnvUint("JWT_ACCESS_TOKEN_EXP_TIME", &p.AccessExpTime)
}

type config struct {
	Listen listenConfig `yaml:"listen" json:"listen"`
	DBCfg  pgConfig     `yaml:"db" json:"db"`
	JwtCfg jwtConfig    `yaml:"jwt" json:"jwt"`
}

func (c *config) loadFromEnv() {
	c.Listen.loadFromEnv()
	c.DBCfg.loadFromEnv()
	c.JwtCfg.loadFromEnv()
}

func defaultConfig() config {
	return config{
		Listen: defaultListenConfig(),
		DBCfg:  defaultPgConfig(),
		JwtCfg: defaultJwtConfig(),
	}
}

func loadConfigFromFile(fn string, c *config) error {
	_, err := os.Stat(fn)
	if err != nil {
		return err
	}

	fl, err := os.Open(filepath.Clean(fn))
	if err != nil {

		return err
	}

	defer fl.Close()
	return yaml.NewDecoder(fl).Decode(c)
}

func loadConfig(fn string) config {
	cfg := defaultConfig()
	cfg.loadFromEnv()
	if err := loadConfigFromFile(fn, &cfg); err != nil {
		if err != nil {
			log.Warn().Str("file", fn).Err(err).Msg("cannot load config file, use defaults")
		}
	}
	return cfg
}
