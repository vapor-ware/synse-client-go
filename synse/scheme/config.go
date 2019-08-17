package scheme

// Config describes a response for `/config` endpoint.
type Config struct {
	Locale     string           `json:"locale" yaml:"locale" mapstructure:"locale"`
	Logging    string           `json:"logging" yaml:"logging" mapstructure:"logging"`
	Plugin     PluginOptions    `json:"plugin" yaml:"plugin" mapstructure:"plugin"`
	Cache      CacheOptions     `json:"cache" yaml:"cache" mapstructure:"cache"`
	GRPC       GRPCOptions      `json:"grpc" yaml:"grpc" mapstructure:"grpc"`
	Transport  TransportOptions `json:"transport" yaml:"transport" mapstructure:"transport"`
	Metrics    MetricsOptions   `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
	PrettyJSON bool             `json:"pretty_json" yaml:"pretty_json" mapstructure:"pretty_json"`
}

// PluginOptions is the config options for plugin.
type PluginOptions struct {
	TCP      []string         `json:"tcp" yaml:"tcp" mapstructure:"tcp"`
	Unix     []string         `json:"unix" yaml:"unix" mapstructure:"unix"`
	Discover DiscoveryOptions `json:"discover" yaml:"discover" mapstructure:"discover"`
}

// DiscoveryOptions is the config options for service discovery.
type DiscoveryOptions struct {
	Kubernetes KubernetesOptions `json:"kubernetes" yaml:"kubernetes" mapstructure:"kubernetes"`
}

// KubernetesOptions is the config options for kubernetes.
type KubernetesOptions struct {
	Namespace string           `json:"namespace" yaml:"namespace" mapstructure:"namespace"`
	Endpoints EndpointsOptions `json:"endpoints" yaml:"endpoints" mapstructure:"endpoints"`
}

// EndpointsOptions is the config options for kubernetes's endpoint.
type EndpointsOptions struct {
	Labels map[string]string `json:"labels" yaml:"labels" mapstructure:"labels"`
}

// CacheOptions is the config options for cache.
type CacheOptions struct {
	Device      DeviceOptions      `json:"device" yaml:"device" mapstructure:"device"`
	Transaction TransactionOptions `json:"transaction" yaml:"transaction" mapstructure:"transaction"`
}

// DeviceOptions is the config options for device cache.
type DeviceOptions struct {
	TTL          int `json:"ttl" yaml:"ttl" mapstructure:"ttl"`
	RebuildEvery int `json:"rebuild_every" yaml:"rebuild_every" mapstructure:"rebuild_every"`
}

// TransactionOptions is the config options for transaction cache.
type TransactionOptions struct {
	TTL int `json:"ttl" yaml:"ttl" mapstructure:"ttl"`
}

// GRPCOptions is the config options for grpc.
type GRPCOptions struct {
	Timeout int        `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	TLS     TLSOptions `json:"tls" yaml:"tls" mapstructure:"tls"`
}

// TLSOptions is the config options for tls communication.
type TLSOptions struct {
	Cert string `json:"cert" yaml:"cert" mapstructure:"cert"`
}

// MetricsOptions is the config options for metrics.
type MetricsOptions struct {
	Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
}

// TransportOptions is the config options for transport communication layer.
type TransportOptions struct {
	HTTP      bool `json:"http" yaml:"http" mapstructure:"http"`
	WebSocket bool `json:"websocket" yaml:"websocket" mapstructure:"websocket"`
}
