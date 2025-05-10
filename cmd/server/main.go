package main

import (
	"net/http"
	"os"

	shared "tracking-service/internal"
	component "tracking-service/internal/components"
	handler "tracking-service/internal/handlers"
	repository "tracking-service/internal/repositories"
	route "tracking-service/internal/routes"
	service "tracking-service/internal/services"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"github.com/urfave/cli/v2"
	metricssdk "go.opentelemetry.io/otel/sdk/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	config shared.Config
)

// @title       Tracking Service
// @version     1.0
// @description 行為追蹤服務

// @securityDefinitions.apikey Bearer
// @in            header
// @name          Authorization

// @securityDefinitions.apikey Bearer
// @in            header
// @name          Authorization

// @schemes      http https
// @termsOfService http://your-terms-of-service.url

// @contact.name ddjoanna
// @contact.url  https://github.com/ddjoanna/web-go-tracking-service
// @contact.email joann.chang0722@gmail.com

// @host     localhost:8080
// @BasePath /
func main() {
	app := &cli.App{
		Name:  "notify",
		Usage: "notify service server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "grpc-port",
				Usage:       "GRPC server port",
				Value:       50051,
				EnvVars:     []string{"GRPC_PORT"},
				Destination: &config.GrpcPort,
			},
			&cli.IntFlag{
				Name:        "http-port",
				Usage:       "HTTP server port",
				Value:       8080,
				EnvVars:     []string{"HTTP_PORT"},
				Destination: &config.HttpPort,
			},
			&cli.StringFlag{
				Name:        "gin-mode",
				Usage:       "gin mode",
				Value:       "release",
				EnvVars:     []string{"GIN_MODE"},
				Destination: &config.GinMode,
			},
			&cli.StringFlag{
				Name:        "postgres-host",
				Usage:       "PostgresSQL DB host address",
				EnvVars:     []string{"POSTGRES_HOST"},
				Destination: &config.PostgresHost,
			},
			&cli.IntFlag{
				Name:        "postgres-port",
				Usage:       "PostgresSQL DB port number",
				Value:       5432,
				EnvVars:     []string{"POSTGRES_PORT"},
				Destination: &config.PostgresPort,
			},
			&cli.StringFlag{
				Name:        "postgres-user",
				Usage:       "PostgresSQL DB user",
				EnvVars:     []string{"POSTGRES_USER"},
				Destination: &config.PostgresUser,
			},
			&cli.StringFlag{
				Name:        "postgres-password",
				Usage:       "PostgresSQL DB password",
				EnvVars:     []string{"POSTGRES_PASSWORD"},
				Destination: &config.PostgresPassword,
			},
			&cli.StringFlag{
				Name:        "postgres-db",
				Usage:       "PostgresSQL DB name",
				EnvVars:     []string{"POSTGRES_DB"},
				Destination: &config.PostgresDb,
			},
			&cli.StringFlag{
				Name:        "postgres-schema",
				Usage:       "PostgresSQL DB schema",
				EnvVars:     []string{"POSTGRES_SCHEMA"},
				Destination: &config.PostgresSchema,
			},
			&cli.IntFlag{
				Name:        "db-max-idle-conns",
				Usage:       "PostgresSQL DB max idle connections",
				EnvVars:     []string{"DB_MAX_IDLE_CONNS"},
				Value:       2,
				Destination: &config.DbMaxIdleConns,
			},
			&cli.IntFlag{
				Name:        "db-max-open-conns",
				Usage:       "PostgresSQL DB max open connections",
				EnvVars:     []string{"DB_MAX_OPEN_CONNS"},
				Value:       5,
				Destination: &config.DbMaxOpenConns,
			},
			&cli.StringFlag{
				Name:        "otlp-service-name",
				Usage:       "Service name for observability",
				EnvVars:     []string{"OTLP_SERVICE_NAME"},
				Destination: &config.OtlpServiceName,
			},
			&cli.StringFlag{
				Name:        "otlp-endpoint",
				Usage:       "The endpoint of the OTLP collector",
				EnvVars:     []string{"OTLP_ENDPOINT"},
				Destination: &config.OtlpEndpoint,
			},
			&cli.StringFlag{
				Name:        "log-format",
				Usage:       "Log format",
				EnvVars:     []string{"LOG_FORMAT"},
				Destination: &config.LogFormat,
			},
			&cli.StringFlag{
				Name:        "kafka-broker",
				Usage:       "Kafka broker",
				EnvVars:     []string{"KAFKA_BROKERS"},
				Destination: &config.KafkaBrokers,
			},
			&cli.StringFlag{
				Name:        "kafka-version",
				Usage:       "Kafka version",
				EnvVars:     []string{"KAFKA_VERSION"},
				Destination: &config.KafkaVersion,
			},
			&cli.StringFlag{
				Name:        "admin-api-key",
				Usage:       "Admin API key",
				EnvVars:     []string{"ADMIN_API_KEY"},
				Destination: &config.AdminApiKey,
			},
		},
		Action: execute,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func execute(cCtx *cli.Context) error {
	log.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
	)))
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	switch config.LogFormat {
	case shared.LOG_FORMAT_JSON:
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	log.Infof("Starting %s", config.OtlpServiceName)

	fx.New(
		fx.Supply(&config),
		fx.Provide(
			component.NewOtlpConn,
			component.NewTracerProvider,
			component.NewMeterProvider,
			component.NewSnowflake,
			component.NewDb,
			component.NewValidator,
			component.NewProducer,
			component.NewHttpServer,
			fx.Annotate(
				component.NewRouter,
				fx.ParamTags("", `group:"routeRegistrars"`),
			),
			AsRouteRegistrar(route.NewAdminRoutes),
			AsRouteRegistrar(route.NewTenantRoutes),
			handler.NewAdminHandler,
			handler.NewTenantHandler,
			service.NewTenantService,
			service.NewPlatformService,
			service.NewApplicationService,
			service.NewEventService,
			repository.NewTenantRepository,
			repository.NewPlatformRepository,
			repository.NewApplicationRepository,
			repository.NewEventRepository,
		),
		fx.Invoke(
			func(*tracesdk.TracerProvider) {},
			func(*metricssdk.MeterProvider) {},
			func(*gorm.DB) {},
			func(*http.Server) {},
			func(*validator.Validate) {},
		),
	).Run()
	return nil
}

func AsRouteRegistrar(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(component.RouteRegistrar)),
		fx.ResultTags(`group:"routeRegistrars"`),
	)
}
