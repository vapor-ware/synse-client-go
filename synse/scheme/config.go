package scheme

// Config describes a response for `/config` endpoint.
type Config struct {
	Locale     string           `json:"locale" mapstructure:"locale"`
	Logging    string           `json:"logging" mapstructure:"logging"`
	Plugin     PluginOptions    `json:"plugin" mapstructure:"plugin"`
	Cache      CacheOptions     `json:"cache" mapstructure:"cache"`
	GRPC       GRPCOptions      `json:"grpc" mapstructure:"grpc"`
	Transport  TransportOptions `json:"transport" mapstructure:"transport"`
	Metrics    MetricsOptions   `json:"metrics" mapstructure:"metrics"`
	PrettyJSON bool             `json:"pretty_json" mapstructure:"pretty_json"`
}

// PluginOptions is the config options for plugin.
type PluginOptions struct {
	TCP      []string         `json:"tcp" mapstructure:"tcp"`
	Unix     []string         `json:"unix" mapstructure:"unix"`
	Discover DiscoveryOptions `json:"discover" mapstructure:"discover"`
}

// DiscoveryOptions is the config options for service discovery.
type DiscoveryOptions struct {
	Kubernetes KubernetesOptions `json:"kubernetes" mapstructure:"kubernetes"`
}

// KubernetesOptions is the config options for kubernetes.
type KubernetesOptions struct {
	Namespace string           `json:"namespace" mapstructure:"namespace"`
	Endpoints EndpointsOptions `json:"endpoints" mapstructure:"endpoints"`
}

// EndpointsOptions is the config options for kubernetes's endpoint.
type EndpointsOptions struct {
	Labels map[string]string `json:"labels" mapstructure:"labels"`
}

// CacheOptions is the config options for cache.
type CacheOptions struct {
	Device      DeviceOptions      `json:"device" mapstructure:"device"`
	Transaction TransactionOptions `json:"transaction" mapstructure:"transaction"`
}

// DeviceOptions is the config options for device cache.
type DeviceOptions struct {
	TTL int `json:"ttl" mapstructure:"ttl"`
}

// TransactionOptions is the config options for transaction cache.
type TransactionOptions struct {
	TTL int `json:"ttl" mapstructure:"ttl"`
}

// GRPCOptions is the config options for grpc.
type GRPCOptions struct {
	Timeout int        `json:"timeout" mapstructure:"timeout"`
	TLS     TLSOptions `json:"tls" mapstructure:"tls"`
}

// TLSOptions is the config options for tls communication.
type TLSOptions struct {
	Cert string `json:"cert" mapstructure:"cert"`
}

// MetricsOptions is the config options for metrics.
type MetricsOptions struct {
	Enabled bool `json:"enabled" mapstructure:"enabled"`
}

// TransportOptions is the config options for transport communication layer.
type TransportOptions struct {
	HTTP      bool `json:"http" mapstructure:"http"`
	WebSocket bool `json:"websocket" mapstructure:"websocket"`
}
