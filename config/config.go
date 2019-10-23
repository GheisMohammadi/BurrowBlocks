package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const Config_File string = "config.toml"

type Config struct {
	GRPC     *GRPCConfig     `toml:"grpc"`
	DataBase *DataBaseConfig `toml:"database"`
	App      *AppConfig      `toml:"app"`
}

type GRPCConfig struct {
	Name string `toml:"name"`
	URL  string `toml:"host"`
	Port string `toml:"port"`
}

type DataBaseConfig struct {
	Type     string `toml:"type"`
	DBName   string `toml:"dbname"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type AppConfig struct {
	CheckingInterval int `toml:"checking interval"`
}

func DefaultGRPCConfig() *GRPCConfig {
	return &GRPCConfig{
		Name: "Hyperledger Burrow",
		URL:  "68.183.183.19",
		Port: "50052",
	}
}

func DefaultDataBaseConfig() *DataBaseConfig {
	return &DataBaseConfig{
		Type:     "Postgre",
		DBName:   "HubbleScan",
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "123456",
	}
}

func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		CheckingInterval: 1000,
	}
}

func LoadConfigFile(create bool) (*Config, error) {
	conf, err := LoadFromFile(Config_File)
	if err != nil {
		if create {
			conf = DefaultConfig()
			conf.SaveToFile(Config_File)
		} else {
			return nil, err
		}

	}
	return conf, nil
}

func (conf *Config) ToTOML() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(conf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (conf *Config) SaveToFile(file string) error {
	toml, err := conf.ToTOML()
	if err != nil {
		return err
	}
	if err := WriteFile(file, toml); err != nil {
		return err
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		GRPC:     DefaultGRPCConfig(),
		DataBase: DefaultDataBaseConfig(),
		App:      DefaultAppConfig(),
	}
}

func FromTOML(t string) (*Config, error) {
	conf := DefaultConfig()

	if _, err := toml.Decode(t, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func LoadFromFile(file string) (*Config, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return FromTOML(string(dat))
}

func WriteFile(filename string, data []byte) error {
	if err := Mkdir(filepath.Dir(filename)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0777); err != nil {
		return fmt.Errorf("config file (%s) writing failed. error: %s", filename, err.Error())
	}
	return nil
}

func Mkdir(dir string) error {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("Creating directory failed. error: %s", err.Error())
	}
	return nil
}
