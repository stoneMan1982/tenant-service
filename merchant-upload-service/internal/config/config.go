package config

import (
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server   ServerConfig `yaml:"server"`
	Database DBConfig     `yaml:"database"`
}

type ServerConfig struct {
	Host         string   `yaml:"host" default:"localhost"`
	Port         int      `yaml:"port" default:"8080"`
	ReadTimeout  Duration `yaml:"read_timeout" default:"30s"`
	WriteTimeout Duration `yaml:"write_timeout" default:"30s"`
}

type DBConfig struct {
	// for sqlite, only support in memory mode
	Dialect  string `yaml:"dialect" validate:"required,oneof=postgres sqlite mysql" comment:"valid value: postgres,sqlite,mysql"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`

	// 连接池配置
	MaxOpenConns    int      `yaml:"max_open_conns" default:"25" comment:"maximum number of open connections"`
	MaxIdleConns    int      `yaml:"max_idle_conns" default:"10" comment:"maximum number of idle connections"`
	ConnMaxLifetime Duration `yaml:"conn_max_lifetime" default:"1h" comment:"maximum lifetime of a connection"`
	ConnMaxIdleTime Duration `yaml:"conn_max_idle_time" default:"30m" comment:"maximum idle time of a connection"`
}

// Duration 是 time.Duration 的自定义类型，用于实现 YAML 序列化/反序列化
type Duration time.Duration

// UnmarshalYAML 实现 yaml.Unmarshaler 接口
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var v interface{}
	if err := value.Decode(&v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		// 解析字符串格式的 duration (如 "5s", "1m")
		dur, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(dur)
	case int:
		// 处理整数形式的毫秒数
		*d = Duration(time.Duration(value) * time.Millisecond)
	case float64:
		// 处理浮点数形式的毫秒数
		*d = Duration(time.Duration(value) * time.Millisecond)
	default:
		// 默认值为 0
		*d = Duration(time.Duration(0))
	}
	return nil
}

// MarshalYAML 实现 yaml.Marshaler 接口
func (d Duration) MarshalYAML() (interface{}, error) {
	// 序列化为字符串形式 (如 "5s")
	return time.Duration(d).String(), nil
}

// Duration 返回 time.Duration 类型的值
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}
