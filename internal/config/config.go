package config

import (
	"encoding/json"
	"os"
)

type Server struct {
	IP       string `json:"ip"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	SFTPPort int    `json:"sftp_port"`
	Path     string `json:"path"`
}

type Config struct {
	MainCmd         string   `json:"main_cmd"`
	BuildOutput     string   `json:"build_output"`
	FrontendEnvPath string   `json:"frontend_env_path"`
	EnvHost         string   `json:"env_host"`
	Servers         []Server `json:"servers"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	return &config, err
}
