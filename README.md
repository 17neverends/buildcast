# BuildCast 🚀

**Frontend Deployment Automation Tool**

**Main problem:** `front-dev should manual build and send buildp path on other dev servers`

![Go](https://img.shields.io/badge/Go-1.18%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![SFTP](https://img.shields.io/badge/Protocol-SFTP-orange)

## Table of Contents
- [Features](#-features)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Technical Details](#-technical-details)

## ✨ Features
- 🔄 Environment variable substitution in `.env` files
- 🏗️ Build automation with custom commands
- 🚀 Multi-server SFTP deployment
- 🔒 Secure credential management
- 📊 Detailed logging for debugging

## 📦 Installation

### Prerequisites
- Go 1.18+
- Git

### Steps
```bash
git clone https://github.com/17neverends/buildcast.git
cd buildcast
go mod download
go build -o buildcast main.go
```

## ⚙️ Configuration

- config.json

```json
{
  // main settings
  "main_cmd": "npm run build", // cmd for build frontend app
  "build_output": "build", // path name for build files output
  "frontend_env_path": ".env", // path to env file
  "env_host": "REACT_APP_API_URL=", // field in env file for change
  
  // individual remote server settings
  "servers": [
    {
      "ip": "xx.xx.xx.xxx", // IP for connect
      "password": "...", // password for connect
      "user": "root", // system username
      "host": "https://test.io", // this value will be override in .env
       "sftp_port": 22, // sftp connection port
       "path": "/home/project" // path for download files on remore server
    }
  ]
}
```

## 🚀 Usage

- Basic commands

   `./buildcast --config=config.json --service=admin_dashboard`


- Command Line Options

   | Flag 	| Description 	|
   |------	|-------------	|
   | --config     	| Path to configuration file            	|
   |--service      	| Service name suffix for deployment path            	|

## 🛠 Technical Details

### Workflow
1. Reads configuration file

2. Backs up original .env file

3. For each server:
   - Updates environment variables

   - Executes build command

   - Cleans target directory

   - Deploys via SFTP

   - Restores original .env

### Dependencies
- github.com/pkg/sftp

- golang.org/x/crypto/ssh