package synse

// scheme.go defines the scheme for the server's responses.

// Error describes an error response.
type Error struct {
	HTTPCode    int    `json:"http_code"`
	ErrorID     int    `json:"error_id"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
	Context     string `json:"context"`
}

// Status describes a response for `/test` endpoint.
type Status struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Version describes a response for `version` endpoint.
type Version struct {
	Version    string `json:"version"`
	APIVersion string `json:"api_version"`
}

// Config describes a response for `/config` endpoint.
type Config struct {
	Logging    string           `json:"logging"`
	Locale     string           `json:"locale"`
	PrettyJSON bool             `json:"pretty_json"`
	Plugin     PluginOptions    `json:"plugin"`
	Cache      CacheOptions     `json:"cache"`
	GRPC       GRPCOptions      `json:"grpc"`
	Metrics    MetricsOptions   `json:"metrics"`
	Transport  TransportOptions `json:"transport"`
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
	Label LabelOptions `json:"label"`
}

// LabelOptions is the config options for kubernetes's label.
type LabelOptions struct {
	App       string `json:"app"`
	Component string `json:"server"`
}

// CacheOptions is the config options for cache.
type CacheOptions struct {
	Meta        MetaOptions        `json:"meta"` // FIXME: not in v3
	Device      DeviceOptions      `json:"device"`
	Transaction TransactionOptions `json:"transaction"`
}

// MetaOptions is the config options for meta cache.
type MetaOptions struct {
	TTL int `json:"ttl"`
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
	Protocol string `json:"protocol"`
}
