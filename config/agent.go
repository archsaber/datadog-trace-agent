package config

import (
	"bytes"
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/DataDog/datadog-trace-agent/flags"
	"github.com/DataDog/datadog-trace-agent/osutil"
	writerconfig "github.com/DataDog/datadog-trace-agent/writer/config"
	log "github.com/cihub/seelog"
)

var (
	// ErrMissingAPIKey is returned when the config could not be validated due to missing API key.
	ErrMissingAPIKey = errors.New("you must specify an API Key, either via a configuration file or the DD_API_KEY env var")

	// ErrMissingHostname is returned when the config could not be validated due to missing hostname.
	ErrMissingHostname = errors.New("failed to automatically set the hostname, you must specify it via configuration for or the DD_HOSTNAME env var")
)

// AgentConfig handles the interpretation of the configuration (with default
// behaviors) in one place. It is also a simple structure to share across all
// the Agent components, with 100% safe and reliable values.
// It is exposed with expvar, so make sure to exclude any sensible field
// from JSON encoding. Use New() to create an instance.
type AgentConfig struct {
	Enabled bool

	// Global
	Hostname   string
	DefaultEnv string // the traces will default to this environment
	ConfigPath string // the source of this config, if any

	// API
	APIEndpoint string
	APIKey      string `json:"-"` // never publish this
	APIEnabled  bool

	// Concentrator
	BucketInterval   time.Duration // the size of our pre-aggregation per bucket
	ExtraAggregators []string

	// Sampler configuration
	ExtraSampleRate float64
	PreSampleRate   float64
	MaxTPS          float64

	// Receiver
	ReceiverHost    string
	ReceiverPort    int
	ConnectionLimit int // for rate-limiting, how many unique connections to allow in a lease period (30s)
	ReceiverTimeout int

	// Writers
	ServiceWriterConfig writerconfig.ServiceWriterConfig
	StatsWriterConfig   writerconfig.StatsWriterConfig
	TraceWriterConfig   writerconfig.TraceWriterConfig

	// internal telemetry
	StatsdHost string
	StatsdPort int

	// logging
	LogLevel             string
	LogFilePath          string
	LogThrottlingEnabled bool

	// watchdog
	MaxMemory        float64       // MaxMemory is the threshold (bytes allocated) above which program panics and exits, to be restarted
	MaxCPU           float64       // MaxCPU is the max UserAvg CPU the program should consume
	MaxConnections   int           // MaxConnections is the threshold (opened TCP connections) above which program panics and exits, to be restarted
	WatchdogInterval time.Duration // WatchdogInterval is the delay between 2 watchdog checks

	// http/s proxying
	ProxyURL          *url.URL
	SkipSSLValidation bool

	// filtering
	Ignore map[string][]string

	// ReplaceTags is used to filter out sensitive information from tag values.
	// It maps tag keys to a set of replacements. Only supported in A6.
	ReplaceTags []*ReplaceRule

	// transaction analytics
	AnalyzedRateByServiceLegacy map[string]float64
	AnalyzedSpansByService      map[string]map[string]float64

	// infrastructure agent binary
	DDAgentBin string // DDAgentBin will be "" for Agent5 scenarios
}

// New returns a configuration with the default values.
func New() *AgentConfig {
	return &AgentConfig{
		Enabled:     true,
		DefaultEnv:  "none",
		APIEndpoint: "https://traces.archsaber.com",
		APIKey:      "",
		APIEnabled:  true,

		BucketInterval:   time.Duration(10) * time.Second,
		ExtraAggregators: []string{"http.status_code"},

		ExtraSampleRate: 1.0,
		PreSampleRate:   1.0,
		MaxTPS:          10,

		ReceiverHost:    "localhost",
		ReceiverPort:    8126,
		ConnectionLimit: 2000,

		ServiceWriterConfig: writerconfig.DefaultServiceWriterConfig(),
		StatsWriterConfig:   writerconfig.DefaultStatsWriterConfig(),
		TraceWriterConfig:   writerconfig.DefaultTraceWriterConfig(),

		StatsdHost: "localhost",
		StatsdPort: 8125,

		LogLevel:             "INFO",
		LogFilePath:          DefaultLogFilePath,
		LogThrottlingEnabled: true,

		MaxMemory:        5e8, // 500 Mb, should rarely go above 50 Mb
		MaxCPU:           0.5, // 50%, well behaving agents keep below 5%
		MaxConnections:   200, // in practice, rarely goes over 20
		WatchdogInterval: time.Minute,

		Ignore: make(map[string][]string),
		AnalyzedRateByServiceLegacy: make(map[string]float64),
		AnalyzedSpansByService:      make(map[string]map[string]float64),
	}
}

// LoadIni reads the contents of the given INI file into the config.
func (c *AgentConfig) LoadIni(path string) error {
	conf, err := NewIni(path)
	if err != nil {
		return err
	}
	c.loadIniConfig(conf)
	return nil
}

// LoadYaml reads the contents of the given YAML file into the config.
func (c *AgentConfig) LoadYaml(path string) error {
	conf, err := NewYaml(path)
	if err != nil {
		return err
	}
	c.loadYamlConfig(conf)
	return nil
}

// LoadEnv reads environment variable values into the config.
func (c *AgentConfig) LoadEnv() { c.loadEnv() }

// Validate validates if the current configuration is good for the agent to start with.
func (c *AgentConfig) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}
	if c.Hostname == "" {
		if err := c.acquireHostname(); err != nil {
			return ErrMissingHostname
		}
	}
	return nil
}

// acquireHostname attempts to acquire a hostname for this configuration. It
// tries to shell out to the infrastructure agent for this, if DD_AGENT_BIN is
// set, otherwise falling back to os.Hostname.
func (c *AgentConfig) acquireHostname() error {
	var err error
	if c.DDAgentBin == "" {
		c.Hostname, err = os.Hostname()
		return err
	}
	var out bytes.Buffer
	cmd := exec.Command(c.DDAgentBin, "hostname")
	cmd.Stdout = &out
	cmd.Env = append(os.Environ(), cmd.Env...) // needed for Windows
	err = cmd.Run()
	if err != nil {
		c.Hostname, err = os.Hostname()
		return err
	}
	if strings.TrimSpace(out.String()) == "" {
		c.Hostname, err = os.Hostname()
	}
	return err
}

// Load attempts to load the configuration from the given path. If it's not found
// it returns an error and a default configuration.
func Load(path string) (*AgentConfig, error) {
	cfgPath := path
	if cfgPath == flags.DefaultConfigPath && !osutil.Exists(cfgPath) && osutil.Exists(agent5Config) {
		// attempting to load inexistent default path, but found existing Agent 5
		// legacy config - try using it
		log.Warnf("Attempting to use Agent 5 configuration: %s", agent5Config)
		cfgPath = agent5Config
	}
	cfg := New()
	switch filepath.Ext(cfgPath) {
	case ".ini", ".conf":
		if err := cfg.LoadIni(cfgPath); err != nil {
			return cfg, err
		}
	case ".yaml":
		if err := cfg.LoadYaml(cfgPath); err != nil {
			return cfg, err
		}
	default:
		return cfg, errors.New("unrecognised file extension (need .yaml, .ini or .conf)")
	}
	cfg.ConfigPath = cfgPath
	return cfg, nil
}
