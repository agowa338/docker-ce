package config

import (
	"strings"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/discovery"
	"github.com/docker/libkv/store"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/netlabel"
)

// Config encapsulates configurations of various Libnetwork components
type Config struct {
	Daemon  DaemonCfg
	Cluster ClusterCfg
	Scopes  map[string]*datastore.ScopeCfg
}

// DaemonCfg represents libnetwork core configuration
type DaemonCfg struct {
	Debug          bool
	DataDir        string
	DefaultNetwork string
	DefaultDriver  string
	Labels         []string
	DriverCfg      map[string]interface{}
}

// ClusterCfg represents cluster configuration
type ClusterCfg struct {
	Watcher   discovery.Watcher
	Address   string
	Discovery string
	Heartbeat uint64
}

// LoadDefaultScopes loads default scope configs for scopes which
// doesn't have explicit user specified configs.
func (c *Config) LoadDefaultScopes(dataDir string) {
	for k, v := range datastore.DefaultScopes(dataDir) {
		if _, ok := c.Scopes[k]; !ok {
			c.Scopes[k] = v
		}
	}
}

// ParseConfig parses the libnetwork configuration file
func ParseConfig(tomlCfgFile string) (*Config, error) {
	cfg := &Config{
		Scopes: map[string]*datastore.ScopeCfg{},
	}

	if _, err := toml.DecodeFile(tomlCfgFile, cfg); err != nil {
		return nil, err
	}

	cfg.LoadDefaultScopes(cfg.Daemon.DataDir)
	return cfg, nil
}

// Option is a option setter function type used to pass varios configurations
// to the controller
type Option func(c *Config)

// OptionDefaultNetwork function returns an option setter for a default network
func OptionDefaultNetwork(dn string) Option {
	return func(c *Config) {
		log.Infof("Option DefaultNetwork: %s", dn)
		c.Daemon.DefaultNetwork = strings.TrimSpace(dn)
	}
}

// OptionDefaultDriver function returns an option setter for default driver
func OptionDefaultDriver(dd string) Option {
	return func(c *Config) {
		log.Infof("Option DefaultDriver: %s", dd)
		c.Daemon.DefaultDriver = strings.TrimSpace(dd)
	}
}

// OptionDriverConfig returns an option setter for driver configuration.
func OptionDriverConfig(networkType string, config map[string]interface{}) Option {
	return func(c *Config) {
		c.Daemon.DriverCfg[networkType] = config
	}
}

// OptionLabels function returns an option setter for labels
func OptionLabels(labels []string) Option {
	return func(c *Config) {
		for _, label := range labels {
			if strings.HasPrefix(label, netlabel.Prefix) {
				c.Daemon.Labels = append(c.Daemon.Labels, label)
			}
		}
	}
}

// OptionKVProvider function returns an option setter for kvstore provider
func OptionKVProvider(provider string) Option {
	return func(c *Config) {
		log.Infof("Option OptionKVProvider: %s", provider)
		if _, ok := c.Scopes[datastore.GlobalScope]; !ok {
			c.Scopes[datastore.GlobalScope] = &datastore.ScopeCfg{}
		}
		c.Scopes[datastore.GlobalScope].Client.Provider = strings.TrimSpace(provider)
	}
}

// OptionKVProviderURL function returns an option setter for kvstore url
func OptionKVProviderURL(url string) Option {
	return func(c *Config) {
		log.Infof("Option OptionKVProviderURL: %s", url)
		if _, ok := c.Scopes[datastore.GlobalScope]; !ok {
			c.Scopes[datastore.GlobalScope] = &datastore.ScopeCfg{}
		}
		c.Scopes[datastore.GlobalScope].Client.Address = strings.TrimSpace(url)
	}
}

// OptionDiscoveryWatcher function returns an option setter for discovery watcher
func OptionDiscoveryWatcher(watcher discovery.Watcher) Option {
	return func(c *Config) {
		c.Cluster.Watcher = watcher
	}
}

// OptionDiscoveryAddress function returns an option setter for self discovery address
func OptionDiscoveryAddress(address string) Option {
	return func(c *Config) {
		c.Cluster.Address = address
	}
}

// OptionDataDir function returns an option setter for data folder
func OptionDataDir(dataDir string) Option {
	return func(c *Config) {
		c.Daemon.DataDir = dataDir
	}
}

// ProcessOptions processes options and stores it in config
func (c *Config) ProcessOptions(options ...Option) {
	for _, opt := range options {
		if opt != nil {
			opt(c)
		}
	}
}

// IsValidName validates configuration objects supported by libnetwork
func IsValidName(name string) bool {
	if strings.TrimSpace(name) == "" || strings.Contains(name, ".") {
		return false
	}
	return true
}

// OptionLocalKVProvider function returns an option setter for kvstore provider
func OptionLocalKVProvider(provider string) Option {
	return func(c *Config) {
		log.Infof("Option OptionLocalKVProvider: %s", provider)
		if _, ok := c.Scopes[datastore.LocalScope]; !ok {
			c.Scopes[datastore.LocalScope] = &datastore.ScopeCfg{}
		}
		c.Scopes[datastore.LocalScope].Client.Provider = strings.TrimSpace(provider)
	}
}

// OptionLocalKVProviderURL function returns an option setter for kvstore url
func OptionLocalKVProviderURL(url string) Option {
	return func(c *Config) {
		log.Infof("Option OptionLocalKVProviderURL: %s", url)
		if _, ok := c.Scopes[datastore.LocalScope]; !ok {
			c.Scopes[datastore.LocalScope] = &datastore.ScopeCfg{}
		}
		c.Scopes[datastore.LocalScope].Client.Address = strings.TrimSpace(url)
	}
}

// OptionLocalKVProviderConfig function returns an option setter for kvstore config
func OptionLocalKVProviderConfig(config *store.Config) Option {
	return func(c *Config) {
		log.Infof("Option OptionLocalKVProviderConfig: %v", config)
		if _, ok := c.Scopes[datastore.LocalScope]; !ok {
			c.Scopes[datastore.LocalScope] = &datastore.ScopeCfg{}
		}
		c.Scopes[datastore.LocalScope].Client.Config = config
	}
}
