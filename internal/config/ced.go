package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type CEDConfig struct {
	CED            CED `mapstructure:"ced"`
	ConfigFilePath string
	_v             *viper.Viper
}

// Currently only supports rabbitmq, so type has to be
// only rabbitmq
type CED struct {
	Broker   Broker   `mapstructure:"broker"`
	Simulate Simulate `mapstructure:"simulate"`
	DryRun   bool     `mapstructure:"dryRun"`
}

type Broker struct {
	Type       string     `mapstructure:"type"`
	Connection Connection `mapstructure:"connection"`
	Exchange   Exchange   `mapstructure:"exchange"`
}

type Connection struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Exchange struct {
	Name        string `mapstructure:"name"`
	Key         string `mapstructure:"key"`
	Type        string `mapstructure:"type"`
	ContentType string `mapstructure:"contentType"`
	Durable     bool   `mapstructure:"durable"`
	AutoDeleted bool   `mapstructure:"autoDeleted"`
	Internal    bool   `mapstructure:"internal"`
	NoWait      bool   `mapstructure:"noWait"`
}

type Simulate struct {
	Nodes       Nodes      `mapstructure:"nodes"`
	Namespaces  Namespaces `mapstructure:"namespaces"`
	Pods        Pods       `mapstructure:"pods"`
	Workloads   Workloads  `mapstructure:"workloads"`
	Load        Load       `mapstructure:"load"`
	ClusterID   int64      `mapstructure:"clusterID"`
	WorkspaceID int64      `mapstructure:"workspaceID"`
	Log         bool       `mapstructure:"log"`
}

type Nodes struct {
	Base    int  `mapstructure:"base"`
	Enabled bool `mapstructure:"enabled"`
}

type Namespaces struct {
	Base    int  `mapstructure:"base"`
	Enabled bool `mapstructure:"enabled"`
}

type Pods struct {
	Base    int  `mapstructure:"base"`
	Enabled bool `mapstructure:"enabled"`
}

type Workloads struct {
	Base    int  `mapstructure:"base"`
	Enabled bool `mapstructure:"enabled"`
}

type Load struct {
	Function string `mapstructure:"function"`
	Enabled  bool   `mapstructure:"enabled"`
}

func NewCEDConfig(configFilePath string) *CEDConfig {
	v := viper.New()
	return &CEDConfig{
		_v:             v,
		ConfigFilePath: configFilePath,
	}
}

func (cc *CEDConfig) LoadConfig() error {
	cc._v.SetConfigType("yaml")
	cc._v.SetConfigFile(cc.ConfigFilePath)

	if err := cc._v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

func (cc *CEDConfig) Unmarshal() error {
	if err := cc._v.Unmarshal(cc); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}
