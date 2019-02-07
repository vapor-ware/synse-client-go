package scheme

// Config describes a response for `/config` endpoint.
type Config struct {
	Locale     string           `json:"locale"`
	Logging    string           `json:"logging"`
	Plugin     PluginOptions    `json:"plugin"`
	Cache      CacheOptions     `json:"cache"`
	GRPC       GRPCOptions      `json:"grpc"`
	Transport  TransportOptions `json:"transport"`
	Metrics    MetricsOptions   `json:"metrics"`
	PrettyJSON bool             `json:"pretty_json"`
}

// PluginOptions is the config options for plugin.
type PluginOptions struct {
	TCP      []string         `json:"tcp"`
	Unix     []string         `json:"unix"`
	Discover DiscoveryOptions `json:"discover"`
}

// DiscoveryOptions is the config options for service discovery.
type DiscoveryOptions struct {
	Kubernetes KubernetesOptions `json:"kubernetes"`
}

// KubernetesOptions is the config options for kubernetes.
type KubernetesOptions struct {
	Namespace string           `json:"namespace"`
	Endpoints EndpointsOptions `json:"endpoints"`
}

// EndpointsOptions is the config options for kubernetes's endpoint.
type EndpointsOptions struct {
	Labels map[string]string `json:"labels"`
}

// CacheOptions is the config options for cache.
type CacheOptions struct {
	Device      DeviceOptions      `json:"device"`
	Transaction TransactionOptions `json:"transaction"`
}

// DeviceOptions is the config options for device cache.
type DeviceOptions struct {
	TTL int `json:"ttl"`
}

// TransactionOptions is the config options for transaction cache.
type TransactionOptions struct {
	TTL int `json:"ttl"`
}

// GRPCOptions is the config options for grpc.
type GRPCOptions struct {
	Timeout int        `json:"timeout"`
	TLS     TLSOptions `json:"tls"`
}

// TLSOptions is the config options for tls communication.
type TLSOptions struct {
	Cert string `json:"cert"`
}

// MetricsOptions is the config options for metrics.
type MetricsOptions struct {
	Enabled bool `json:"enabled"`
}

// TransportOptions is the config options for transport communication layer.
type TransportOptions struct {
	HTTP      bool `json:"http"`
	WebSocket bool `json:"websocket"`
}
