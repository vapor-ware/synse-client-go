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
	PrettyJSON bool             `json:"pretty_json"`
	Locale     string           `json:"locale"`
	Plugin     PluginOptions    `json:"plugin"`
	Cache      CacheOptions     `json:"cache"`
	GRPC       GRPCOptions      `json:"grpc"`
	Metrics    MetricsOptions   `json:"metrics"`
	Transport  TransportOptions `json:"transport"`
}

// PluginsOptions
type PluginOptions struct {
	TCP      []string         `json:"tcp"`
	Unix     []string         `json:"unix"`
	Discover DiscoveryOptions `json:"discover"`
}

// DiscoveryOptions
type DiscoveryOptions struct {
	Kubernetes KubernetesOptions `json:"kubernetes"`
}

// KubernetesOptions
type KubernetesOptions struct {
	Namespace string           `json:"namespace"`
	Endpoints EndpointsOptions `json:"endpoints"`
}

// EndpointsOptions
type EndpointsOptions struct {
	Label LabelOptions `json:"label"`
}

// LabelOptions
type LabelOptions struct {
	App       string `json:"app"`
	Component string `json:"server"`
}

// CacheOptions
type CacheOptions struct {
	Meta        MetaOptions        `json:"meta"` // FIXME: not in v3
	Device      DeviceOptions      `json:"device"`
	Transaction TransactionOptions `json:"transaction"`
}

// MetaOptions
type MetaOptions struct {
	TTL int `json:"ttl"`
}

// DeviceOptions
type DeviceOptions struct {
	TTL int `json:"ttl"`
}

// TransactionOptions
type TransactionOptions struct {
	TTL int `json:"ttl"`
}

// GRPCOptions
type GRPCOptions struct {
	Timeout int        `json:"timeout"`
	TLS     TLSOptions `json:"tls"`
}

// TLSOptions
type TLSOptions struct {
	Cert string `json:"cert"`
}

// MetricsOptions
type MetricsOptions struct {
	Enabled bool `json:"enabled"`
}

// TransportOptions
type TransportOptions struct {
	Protocol string `json:"protocol"`
}
