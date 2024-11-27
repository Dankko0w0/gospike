package confManager

import (
	"fmt"
	"sync"

	"github.com/Dankko0w0/gospike/models"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	once            sync.Once
	v               *viper.Viper
	ConfInitialized bool = false
)

// InitConfig 初始化配置管理器
func InitConfig(configPath string, configName string, configType string) error {
	var err error
	once.Do(func() {
		v = viper.New()

		// 设置配置文件信息
		v.AddConfigPath(".")        // 添加当前目录作为查找路径
		v.AddConfigPath(configPath) // 配置文件路径
		if configName != "" {
			v.SetConfigName(configName) // 配置文件名称(无扩展名)
		}
		if configType != "" {
			v.SetConfigType(configType) // 配置文件类型
		}

		// 读取环境变量
		v.AutomaticEnv()
		v.SetEnvPrefix("GOSPIKE") // 环境变量前缀

		// 读取配置文件
		if err = v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err = fmt.Errorf("no config file found: %w", err)
			} else {
				err = fmt.Errorf("error reading config file: %w", err)
			}
		}
		ConfInitialized = true

		// 监听配置文件变化
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("Config file changed: %s\n", e.Name)
			// 重新读取配置文件
			if err := v.ReadInConfig(); err != nil {
				fmt.Printf("Error reloading config: %v\n", err)
			}
		})
	})

	return err
}

// SetDefault 设置单个默认配置项
func SetDefault(key string, value interface{}) {
	v.SetDefault(key, value)
}

// SetDefaults 批量设置默认配置项
func SetDefaults(defaults []models.DefaultKV) {
	for _, d := range defaults {
		v.SetDefault(d.Key, d.Value)
	}
}

// Set 设置配置项
func Set(key string, value interface{}) {
	v.Set(key, value)
}

// GetString 获取字符串配置
func GetString(key string) string {
	return v.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return v.GetInt(key)
}

// GetBool 获取布尔配置
func GetBool(key string) bool {
	return v.GetBool(key)
}

// GetFloat64 获取浮点数配置
func GetFloat64(key string) float64 {
	return v.GetFloat64(key)
}

// GetStringSlice 获取字符串切片配置
func GetStringSlice(key string) []string {
	return v.GetStringSlice(key)
}

// GetStringMap 获取字符串映射配置
func GetStringMap(key string) map[string]interface{} {
	return v.GetStringMap(key)
}

// GetAll 获取所有配置
func GetAll() map[string]interface{} {
	return v.AllSettings()
}
