package conf

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
	"easy-go/library"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     int
	Charset  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var logger = library.Logger

// config配置
var Config Configuration

// 数据配置
var DBConfigMap = make(map[string]DBConfig)

// appconfig
var AppConfig = make(map[string]string)

func init() {

	// 加载项目配置
	loadConfig()

	// 加载数据库配置
	loadDBConfig()

	// 加载app.config 文件
	loadAppConfig()
}

func loadConfigAll() {
	Config = Configuration{}
	_, _ = library.ReadFile("conf/files/config.json", Config)
}

func loadConfig() {
	file, err := os.Open("conf/files/config.json")
	if err != nil {
		logger.Fatalln("connot open config file", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)

	if err != nil {
		logger.Fatalln("connot get configuration from file", err)
	}
}

func loadDBConfig() {
	dbBytes, err := os.Open("conf/files/databases.json")
	defer dbBytes.Close()
	if err != nil {
		logger.Fatalln("connot open db config file", err)
	}
	decoder := json.NewDecoder(dbBytes)
	err = decoder.Decode(&DBConfigMap)
	if err != nil {
		logger.Fatalln("connot get configuration from file", err)
	}
}

// load app conf file
func loadAppConfig() {
	configFile, err := os.Open("conf/files/app.conf")
	defer configFile.Close()
	if err != nil {
		logger.Fatalln("connot open app config file", err)
	}

	r := bufio.NewReader(configFile)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Fatalln("read config file err", err)
		}

		lineStr := strings.TrimSpace(string(line))
		if strings.Index(lineStr, "#") == 0 {
			continue
		}

		equalIndex := strings.Index(lineStr, "=")
		if equalIndex < 0 {
			continue
		}
		configKey := strings.TrimSpace(lineStr[:equalIndex])
		if len(configKey) == 0 {
			continue
		}
		configValue := strings.TrimSpace(lineStr[equalIndex+1:])
		pos := strings.Index(configValue, "\t#")
		if pos > -1 {
			configValue = configValue[:pos]
		}
		pos = strings.Index(configValue, "#")
		if pos > -1 {
			configValue = configValue[:pos]
		}
		pos = strings.Index(configValue, "\t//")
		if pos > -1 {
			configValue = configValue[:pos]
		}
		pos = strings.Index(configValue, "//")
		if pos > -1 {
			configValue = configValue[:pos]
		}
		configValue = strings.TrimSpace(configValue)
		if len(configValue) == 0 {
			continue
		}

		AppConfig[configKey] = configValue
	}

}
