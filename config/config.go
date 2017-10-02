package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type SConfigTest struct {
	TestKey     string `yaml:"consumer_key"`
	TestSecret  string `yaml: consumer_secret"`
	ATestKey    string `yaml:"access_token"`
	ATestSecret string `yaml:"access_token_secret"`
}
type SConfigToken struct {
	ConsumerKey       string `yaml:"consumer_key"`
	ConsumerSecret    string `yaml:"consumer_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`
}

type SConfigApp struct {
	DB    string       `yaml: "db"`
	Test  string       `yaml: "test"`
	Token SConfigToken `yaml: "token"`
}
type SConfig struct {
	App SConfigApp `yaml: "app"`
}

var Config = SConfig{}
var DB = ""
var Token = SConfigToken{}

func init() {
	yamlBytes, err := ioutil.ReadFile(`config.yml`)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(yamlBytes, &Config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	DB = Config.App.DB
	Token = Config.App.Token
}
