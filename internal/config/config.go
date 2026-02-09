package config

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构体
type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
	Redis  RedisConfig  `yaml:"redis" json:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        string `yaml:"port" json:"port"`
	Host        string `yaml:"host" json:"host"`
	Mode        string `yaml:"mode" json:"mode"` // gin模式: debug, release, test
	ReadTimeout int    `yaml:"read_timeout" json:"read_timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `yaml:"addr" json:"addr"`
	Password     string `yaml:"password" json:"password"`
	DB           int    `yaml:"db" json:"db"`
	PoolSize     int    `yaml:"pool_size" json:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns" json:"min_idle_conns"`
}

// Load 加载配置，支持YAML文件和环境变量
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Host:        getEnv("SERVER_HOST", ""),
			Mode:        getEnv("GIN_MODE", "debug"),
			ReadTimeout: getEnvAsInt("SERVER_READ_TIMEOUT", 60),
		},
		Redis: RedisConfig{
			Addr:         getEnv("REDIS_ADDR", "localhost:6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		},
	}

	// 尝试从配置文件加载
	if err := loadFromFile(config); err != nil {
		log.Printf("配置文件加载失败，使用默认配置: %v", err)
	}

	// 环境变量覆盖配置文件
	overrideWithEnv(config)

	log.Printf("配置加载完成: Server Port=%s, Redis Addr=%s", config.Server.Port, config.Redis.Addr)
	return config
}

// loadFromFile 从YAML文件加载配置
func loadFromFile(config *Config) error {
	configPaths := []string{
		"./config.yaml",
		"./config/config.yaml",
		"./configs/config.yaml",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("找到配置文件: %s", path)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			return yaml.Unmarshal(data, config)
		}
	}

	return os.ErrNotExist
}

// overrideWithEnv 用环境变量覆盖配置
func overrideWithEnv(config *Config) {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		config.Server.Mode = mode
	}
	if timeout := getEnvAsInt("SERVER_READ_TIMEOUT", 0); timeout > 0 {
		config.Server.ReadTimeout = timeout
	}

	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		config.Redis.Addr = addr
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}
	if db := getEnvAsInt("REDIS_DB", -1); db >= 0 {
		config.Redis.DB = db
	}
	if poolSize := getEnvAsInt("REDIS_POOL_SIZE", 0); poolSize > 0 {
		config.Redis.PoolSize = poolSize
	}
	if minIdle := getEnvAsInt("REDIS_MIN_IDLE_CONNS", 0); minIdle > 0 {
		config.Redis.MinIdleConns = minIdle
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量作为整数，如果不存在则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
