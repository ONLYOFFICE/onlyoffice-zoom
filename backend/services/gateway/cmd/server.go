package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/service/repl"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/shared"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/config"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/service"
	"github.com/oklog/run"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func Server(config *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    "starts a new http server instance",
		Category: "server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config_path",
				Usage:   "sets custom configuration path",
				Aliases: []string{"config", "conf", "c"},
			},
			&cli.StringFlag{
				Name:    "environment",
				Usage:   "sets servers environment (development, testing, production)",
				Aliases: []string{"env", "e"},
			},
			&cli.IntFlag{
				Name:  "http_version",
				Usage: "sets http server's version",
			},
			&cli.StringFlag{
				Name:  "http_address",
				Usage: "sets http server's address",
			},
			&cli.Uint64Flag{
				Name:  "http_limits",
				Usage: "sets http server's limits",
			},
			&cli.Uint64Flag{
				Name:  "http_iplimits",
				Usage: "sets http server's IP limits",
			},
			&cli.StringFlag{
				Name:  "zoom_client_id",
				Usage: "sets zoom oauth clientID",
			},
			&cli.StringFlag{
				Name:  "zoom_client_secret",
				Usage: "sets zoom oauth clientSecret",
			},
			&cli.StringFlag{
				Name:  "zoom_redirect_uri",
				Usage: "sets zoom oauth redirect uri",
			},
			&cli.IntFlag{
				Name:  "broker_type",
				Usage: "sets http server's broker type (1 - RabbitMQ, 2 - Memory)",
			},
			&cli.StringSliceFlag{
				Name:  "broker_addresses",
				Usage: "sets http server's broker addresses",
			},
			&cli.StringFlag{
				Name:  "broker_secure",
				Usage: "sets http server's broker secure flag",
			},
			&cli.StringSliceFlag{
				Name:  "registry_addresses",
				Usage: "sets http server's registry addresses",
			},
			&cli.IntFlag{
				Name:  "registry_type",
				Usage: "sets http server's registry type (1 - Kubernetes, 2 - Consul, 3 - ETCD)",
			},
			&cli.DurationFlag{
				Name:  "registry_ttl",
				Usage: "sets http server's registry cache ttl",
			},
			&cli.BoolFlag{
				Name:  "registry_secure",
				Usage: "sets http server's registry secure flag",
			},
			&cli.StringFlag{
				Name:  "tracer_address",
				Usage: "sets distributed tracing address",
			},
			&cli.StringFlag{
				Name:  "tracer_type",
				Usage: "sets distributed tracing provider (0 - Console, 1 - Zipkin)",
			},
			&cli.Float64Flag{
				Name:  "tracer_ratio",
				Usage: "sets distributed tracing ratio",
			},
			&cli.StringFlag{
				Name:  "repl_name",
				Usage: "sets repl server's name",
			},
			&cli.StringFlag{
				Name:  "repl_address",
				Usage: "sets repl server's address",
			},
			&cli.BoolFlag{
				Name:  "repl_debug",
				Usage: "sets repl server's profiler flag",
			},
		},
		Action: func(c *cli.Context) error {
			var (
				ENVIRONMENT = c.String("environment")

				CONFIG_PATH = c.String("config_path")

				HTTP_VERSION   = c.Int("http_version")
				HTTP_ADDRESS   = c.String("http_address")
				HTTP_LIMITS    = c.Uint64("http_limits")
				HTTP_LIMITS_IP = c.Uint64("http_iplimits")

				ZOOM_CLIENT_ID      = c.String("zoom_client_id")
				ZOOM_CLIENT_SECRET  = c.String("zoom_client_secret")
				ZOOM_WEBHOOK_SECRET = c.String("zoom_webhook_secret")
				ZOOM_REDIRECT_URI   = c.String("zoom_redirect_uri")

				REGISTRY_ADDRESSES = c.StringSlice("registry_addresses")
				REGISTRY_TYPE      = c.Int("registry_type")
				REGISTRY_TTL       = c.Duration("registry_ttl")
				REGISTRY_SECURE    = c.Bool("registry_secure")

				BROKER_TYPE      = c.Int("broker_type")
				BROKER_ADDRESSES = c.StringSlice("broker_addresses")
				BROKER_SECURE    = c.String("broker_secure")

				TRACER_ADDRESS = c.String("tracer_address")
				TRACER_TYPE    = c.String("tracer_type")
				TRACER_RATIO   = c.Float64("tracer_ratio")

				REPL_NAME    = c.String("repl_name")
				REPL_ADDRESS = c.String("repl_address")
				REPL_DEBUG   = c.Bool("repl_debug")
			)

			if err := envconfig.Process(context.Background(), config); err != nil {
				return err
			}

			config.Registry.Secure = REGISTRY_SECURE
			if CONFIG_PATH != "" {
				file, err := os.Open(CONFIG_PATH)
				if err != nil {
					return err
				}
				defer file.Close()

				decoder := yaml.NewDecoder(file)

				if err := decoder.Decode(&config); err != nil {
					return err
				}
			}

			if _, ok := shared.SUPPORTED_ENVIRONMENTS[config.Environment]; !ok {
				config.Environment = shared.SUPPORTED_ENVIRONMENTS["development"]
			}

			config.Server.Namespace = "onlyoffice"
			config.REPL.Namespace = "onlyoffice"
			config.Server.Name = "gateway"
			config.REPL.Name = "gateway.repl"

			if env, ok := shared.SUPPORTED_ENVIRONMENTS[ENVIRONMENT]; ok {
				config.Environment = env
			}

			if HTTP_VERSION > 0 {
				config.Server.Version = HTTP_VERSION
				config.REPL.Version = HTTP_VERSION
			}

			if HTTP_ADDRESS != "" {
				config.Server.Address = HTTP_ADDRESS
			}

			if HTTP_LIMITS > 0 {
				config.Server.RateLimiter.Limit = HTTP_LIMITS
			}

			if HTTP_LIMITS_IP > 0 {
				config.Server.RateLimiter.IPLimit = HTTP_LIMITS_IP
			}

			if ZOOM_CLIENT_ID != "" {
				config.Zoom.ClientID = ZOOM_CLIENT_ID
			}

			if ZOOM_CLIENT_SECRET != "" {
				config.Zoom.ClientSecret = ZOOM_CLIENT_SECRET
			}

			if ZOOM_WEBHOOK_SECRET != "" {
				config.Zoom.WebhookSecret = ZOOM_WEBHOOK_SECRET
			}

			if ZOOM_REDIRECT_URI != "" {
				config.Zoom.RedirectURI = ZOOM_REDIRECT_URI
			}

			if len(REGISTRY_ADDRESSES) > 0 {
				config.Registry.Addresses = REGISTRY_ADDRESSES
			}

			if BROKER_TYPE > 0 {
				config.Broker.Type = BROKER_TYPE
			}

			if len(BROKER_ADDRESSES) > 0 {
				config.Broker.Addrs = BROKER_ADDRESSES
			}

			if BROKER_SECURE != "" {
				flag, err := strconv.ParseBool(BROKER_SECURE)
				if err == nil {
					config.Broker.Secure = flag
				}
			}

			if REGISTRY_TYPE > 0 {
				config.Registry.RegistryType = REGISTRY_TYPE
			}

			if REGISTRY_TTL > 0 {
				config.Registry.CacheTTL = REGISTRY_TTL
			}

			if TRACER_ADDRESS != "" {
				config.Tracer.Address = TRACER_ADDRESS
			}

			if TRACER_TYPE != "" {
				t, err := strconv.ParseInt(TRACER_TYPE, 10, 0)
				if err == nil {
					config.Tracer.TracerType = int(t)
				}
			}

			if TRACER_RATIO > 0 {
				config.Tracer.FractionRatio = TRACER_RATIO
			}

			if REPL_NAME != "" {
				config.REPL.Name = REPL_NAME
			}

			if REPL_ADDRESS != "" {
				config.REPL.Address = REPL_ADDRESS
			}

			if !config.REPL.Debug {
				config.REPL.Debug = REPL_DEBUG
			}

			if err := config.Validate(); err != nil {
				return err
			}

			return startGroup(config)
		},
	}
}

func startGroup(config *config.Config) error {
	runGroup := run.Group{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := log.NewLogrusLogger(
		log.WithName(config.Server.Name),
		log.WithEnvironment(config.Environment),
		log.WithLevel(log.LogLevel(config.Logger.Level)),
		log.WithPretty(config.Logger.Pretty),
		log.WithColor(config.Logger.Color),
		log.WithReportCaller(false),
		log.WithFile(log.LumberjackOption{
			Filename:   config.Logger.File.Filename,
			MaxSize:    config.Logger.File.MaxSize,
			MaxAge:     config.Logger.File.MaxAge,
			MaxBackups: config.Logger.File.MaxBackups,
			LocalTime:  config.Logger.File.LocalTime,
			Compress:   config.Logger.File.Compress,
		}),
		log.WithElastic(log.ElasticOption{
			Address:            config.Logger.Elastic.Address,
			Index:              config.Logger.Elastic.Index,
			Bulk:               config.Logger.Elastic.Bulk,
			Async:              config.Logger.Elastic.Async,
			HealthcheckEnabled: config.Logger.Elastic.HealthcheckEnabled,
			BasicAuthUsername:  config.Logger.Elastic.BasicAuthUsername,
			BasicAuthPassword:  config.Logger.Elastic.BasicAuthPassword,
			GzipEnabled:        config.Logger.Elastic.GzipEnabled,
		}),
	)

	if err != nil {
		return err
	}

	{
		server, err := service.NewService(
			service.WithConfig(config.Server),
			service.WithZoomConfig(config.Zoom),
			service.WithTracer(config.Tracer),
			service.WithBroker(config.Broker),
			service.WithLogger(logger),
			service.WithRegistry(config.Registry),
			service.WithContext(ctx),
		)

		if err != nil {
			logger.Errorf("failed to initialize %s server. ", config.Server.Name, err)
			return err
		}

		runGroup.Add(func() error {
			return server.Run()
		}, func(e error) {
			logger.Warnf("shutting down %s server", config.Server.Name)
		})
	}

	{
		repl := repl.NewService(
			repl.WithNamespace(config.REPL.Namespace),
			repl.WithName(config.REPL.Name),
			repl.WithVersion(config.REPL.Version),
			repl.WithAddress(config.REPL.Address),
			repl.WithDebug(config.REPL.Debug || config.Environment == "development"),
		)

		runGroup.Add(func() error {
			logger.Warnf("starting a repl server %s", config.REPL.Address)
			return repl.ListenAndServe()
		}, func(e error) {
			logger.Warnf("shutting down %s server", config.REPL.Name)
			repl.Shutdown(context.Background())
		})
	}

	return runGroup.Run()
}
