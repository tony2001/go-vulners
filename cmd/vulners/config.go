package main

type Config struct {
	ApiKey  string `yaml:"apiKey" env:"API_KEY" env-description:"Vulners API key" env-required:"true"`
	Server  string `yaml:"server" env:"SERVER" env-description:"Vulners API server address" env-default:"https://vulners.com"`
	Verbose bool   `yaml:"verbose" env:"VERBOSE" env-description:"Enable verbose output"`
}
