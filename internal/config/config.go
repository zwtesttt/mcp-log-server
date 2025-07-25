package config

import "fmt"

// Environment 环境配置结构
type Environment struct {
	Name   string
	Host   string
	User   string
	Port   string
	SSHKey string
}

// LogFile 日志文件配置结构
type LogFile struct {
	Name        string
	Path        string
	Description string
	Aliases     []string
}

// Config 应用配置结构
type Config struct {
	Environments map[string]Environment
	LogFiles     map[string]LogFile
	Ollama       OllamaConfig
}

// OllamaConfig Ollama配置
type OllamaConfig struct {
	BaseURL      string
	DefaultModel string
	Timeout      int // 超时时间（秒）
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Environments: map[string]Environment{
			"dev": {
				Name: "开发环境",
				Host: "localhost",
				User: "local",
				Port: "22",
			},
			"test": {
				Name: "测试环境",
				Host: "test.example.com",
				User: "testuser",
				Port: "22",
			},
			"staging": {
				Name: "预发布环境",
				Host: "staging.example.com",
				User: "staginguser",
				Port: "22",
			},
			"prod": {
				Name: "生产环境",
				Host: "prod.example.com",
				User: "produser",
				Port: "22",
			},
		},
		LogFiles: map[string]LogFile{
			"blackhole": {
				Name:        "blackhole",
				Path:        "", // 动态路径，将在运行时根据环境设置
				Description: "黑洞服务日志",
				Aliases:     []string{"黑洞"},
			},
			"oms": {
				Name:        "oms",
				Path:        "", // 动态路径，将在运行时根据环境设置
				Description: "OMS系统日志",
				Aliases:     []string{},
			},
		},
		Ollama: OllamaConfig{
			BaseURL:      "http://192.168.194.90:11434",
			DefaultModel: "gemma3:27b",
			Timeout:      120, // 增加到120秒以支持长时间AI分析
		},
	}
}

// GetEnvironment 获取环境配置
func (c *Config) GetEnvironment(name string) (Environment, bool) {
	env, exists := c.Environments[name]
	return env, exists
}

// GetLogFile 获取日志文件配置，根据环境动态设置路径
func (c *Config) GetLogFile(name string) (LogFile, bool) {
	return c.GetLogFileForEnvironment(name, "dev") // 默认使用dev环境
}

// GetLogFileForEnvironment 根据环境获取日志文件配置
func (c *Config) GetLogFileForEnvironment(name, envName string) (LogFile, bool) {
	var logFile LogFile
	var exists bool

	// 直接查找
	if logFile, exists = c.LogFiles[name]; exists {
		// 动态设置路径
		logFile.Path = c.getLogPath(envName, name)
		return logFile, true
	}

	// 通过别名查找
	for logName, lf := range c.LogFiles {
		for _, alias := range lf.Aliases {
			if alias == name {
				logFile = lf
				logFile.Path = c.getLogPath(envName, logName)
				return logFile, true
			}
		}
	}

	return LogFile{}, false
}

// getLogPath 根据环境和日志类型生成日志文件路径
func (c *Config) getLogPath(envName, logType string) string {
	basePath := "/Users/lms/GolandProjects/mcp-log-server/logs"
	return fmt.Sprintf("%s/%s/%s.log", basePath, envName, logType)
}

// ListEnvironments 列出所有环境
func (c *Config) ListEnvironments() map[string]Environment {
	return c.Environments
}

// ListLogFiles 列出所有日志文件类型
func (c *Config) ListLogFiles() map[string]LogFile {
	return c.LogFiles
}
