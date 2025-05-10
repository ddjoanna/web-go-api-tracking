package shared

const (
	LOG_FORMAT_JSON = "json"
	LOG_FORMAT_TEXT = "text"
)

const (
	KafkaTopic   = "tracking"
	KafkaGroupId = "tracking_group"
)

type contextKey string

const (
	AdminApiKey            contextKey = "admin_key"
	TenantApplicationIDKey contextKey = "tenant_application_id"
	TenantIDKey            contextKey = "tenant_id"
)

type Config struct {
	Env              string
	GrpcPort         int
	HttpPort         int
	GinMode          string
	ExternalHost     string
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
	PostgresSchema   string
	DbMaxIdleConns   int
	DbMaxOpenConns   int
	OtlpEndpoint     string
	OtlpServiceName  string
	LogFormat        string
	KafkaBrokers     string
	KafkaVersion     string
	AdminApiKey      string
}
