package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	//useEnv = flag.Bool("env", true, "是否开启读取env")
	//configFileName = flag.String("c", "./configs.yaml", "设置配置文件地址")
	// 远程配置文件参数 http://127.0.0.1:4001/config/hugo.json
	//remoteConfigEndPoint = flag.String("rcep", "", "远程配置文件url")
	defaultTag = "json"
	defaultEnv = true
)

type (
	// Option defines the method to customize the config options.
	Option func(opt *options)

	options struct {
		env                  bool
		openRemoteConfig     bool
		configName           string
		configPath           []string
		remoteConfigEndPoint []string
		defaultTag           string
		openWatch            bool
	}
)

func newOption() options {
	return options{env: defaultEnv}
}

func (o *options) ResetTag(m *mapstructure.DecoderConfig) {
	m.TagName = o.defaultTag
	if o.defaultTag == "" {
		m.TagName = defaultTag
	}
}

func UseCloseEnv() Option {
	return func(opt *options) {
		opt.env = false
	}
}

// AddConfigPath 添加额外的配置文件所包含的目录
func AddConfigPath(configPath ...string) Option {
	return func(opt *options) {
		opt.configPath = configPath
	}
}

func SetConfigName(configName string) Option {
	return func(opt *options) {
		opt.configName = configName
	}
}

func UseRemoteConfig(endPoint ...string) Option {
	return func(opt *options) {
		if len(endPoint) > 0 {
			opt.openRemoteConfig = true
			opt.remoteConfigEndPoint = endPoint
		}
	}
}

func ResetTag(tag string) Option {
	return func(opt *options) {
		opt.defaultTag = tag
	}
}

func OpenWatch() Option {
	return func(opt *options) {
		opt.openWatch = true
	}
}

type cfgLock interface {
	// TryLock 尝试加锁
	TryLock() error
	// UnLock 需要会自动释放
	UnLock() error
}

// Setup 载入配置文件
func Setup(defaultPath string, cfg cfgLock, ops ...Option) error {
	var urls *url.URL
	var err error
	opt := newOption()
	for _, o := range ops {
		o(&opt)
	}
	v := viper.NewWithOptions()

	v.SetConfigType("yaml")
	v.AddConfigPath(defaultPath)

	// 配置文件名称 如果文件名是aaa.yaml,这里就要设置configName=aaa(viper默认的文件名为config)
	if opt.configName != "" {
		v.SetConfigName(opt.configName)
	}

	//配置额外的文件位置
	for i := range opt.configPath {
		v.AddConfigPath(opt.configPath[i])
	}
	//自动获取全部的env加入到viper中。（如果环境变量多就全部加进来）默认别名和环境变量名一致
	if opt.env {
		v.AutomaticEnv()
	}
	//读文件到viper配置中
	err = v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}
	// 获取远程配置文件
	if opt.openRemoteConfig {
		for i := range opt.remoteConfigEndPoint {
			urls, err = url.Parse(opt.remoteConfigEndPoint[i])
			if err != nil {
				return err
			}
			if err = v.AddRemoteProvider("remote", fmt.Sprintf("%s://%s", urls.Scheme, urls.Host), urls.Path); err != nil {
				return err
			}
		}
		if err = v.ReadRemoteConfig(); err != nil {
			return err
		}
	}
	// 是否开启监听
	if opt.openWatch {
		v.OnConfigChange(func(in fsnotify.Event) {
			// 缓一下确认更改完文件
			time.Sleep(10 * time.Second)
			if err := cfg.TryLock(); err != nil {
				return
			}
			defer cfg.UnLock()
			if err = v.Unmarshal(cfg, opt.ResetTag); err != nil {
				return
			}
		})
		v.WatchConfig()
	}
	// 系列化成config对象
	if err = v.Unmarshal(cfg, opt.ResetTag); err != nil {
		return err
	}
	return nil
}

func setEnvToViper(v *viper.Viper) {
	v.AutomaticEnv()
	keys := os.Environ()
	for i := range keys {
		cache := strings.Split(keys[i], "=")
		if strings.Contains(cache[0], "PATH") {
			continue
		}
		if len(cache) > 1 {
			v.Set(strings.ToLower(cache[0]), cache[1])
		}
	}
}
