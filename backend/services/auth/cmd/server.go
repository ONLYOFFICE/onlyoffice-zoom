package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/service/repl"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/shared"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/config"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/service"
	"github.com/oklog/run"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func Server(config *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    "starts a new rpc server instance",
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
				Name:  "rpc_version",
				Usage: "sets rpc server's version",
			},
			&cli.StringFlag{
				Name:  "rpc_address",
				Usage: "sets rpc server's address",
			},
			&cli.Uint64Flag{
				Name:  "rpc_limits",
				Usage: "sets rpc server's limits",
			},
			&cli.StringFlag{
				Name:  "zoom_client_id",
				Usage: "sets zoom oauth clientID",
			},
			&cli.StringFlag{
				Name:  "zoom_client_secret",
				Usage: "sets zoom oauth clientSecret",
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
			&cli.StringFlag{
				Name:  "registry_secure",
				Usage: "sets http server's registry secure flag",
			},
			&cli.IntFlag{
				Name:  "registry_type",
				Usage: "sets http server's registry type (1 - Kubernetes, 2 - Consul, 3 - ETCD)",
			},
			&cli.DurationFlag{
				Name:  "registry_ttl",
				Usage: "sets http server's registry cache ttl",
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
				Name:  "repl_address",
				Usage: "sets repl server's address",
			},
			&cli.BoolFlag{
				Name:  "repl_debug",
				Usage: "sets repl server's profiler flag",
			},
			&cli.StringFlag{
				Name:  "persistence_url",
				Usage: "sets mongodb url",
			},
		},
		Action: func(c *cli.Context) error {
			var (
				CONFIG_PATH = c.String("config_path")
				ENVIRONMENT = c.String("environment")

				RPC_VERSION = c.Int("rpc_version")
				RPC_ADDRESS = c.String("rpc_address")
				RPC_LIMITS  = c.Uint64("rpc_limits")

				ZOOM_CLIENT_ID     = c.String("zoom_client_id")
				ZOOM_CLIENT_SECRET = c.String("zoom_client_secret")

				REGISTRY_ADDRESSES = c.StringSlice("registry_addresses")
				REGISTRY_SECURE    = c.String("registry_secure")
				REGISTRY_TYPE      = c.Int("registry_type")
				REGISTRY_TTL       = c.Duration("registry_ttl")

				BROKER_TYPE      = c.Int("broker_type")
				BROKER_ADDRESSES = c.StringSlice("broker_addresses")
				BROKER_SECURE    = c.String("broker_secure")

				TRACER_ADDRESS = c.String("tracer_address")
				TRACER_TYPE    = c.String("tracer_type")
				TRACER_RATIO   = c.Float64("tracer_ratio")

				REPL_ADDRESS = c.String("repl_address")
				REPL_DEBUG   = c.Bool("repl_debug")

				PERSISTENCE_URL = c.String("persistence_url")
			)

			if err := envconfig.Process(context.Background(), config); err != nil {
				return err
			}

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
			config.Server.Name = "auth"
			config.REPL.Namespace = "onlyoffice"
			config.REPL.Name = "auth.repl"

			if env, ok := shared.SUPPORTED_ENVIRONMENTS[ENVIRONMENT]; ok {
				config.Environment = env
			}

			if RPC_VERSION > 0 {
				config.Server.Version = RPC_VERSION
				config.REPL.Version = RPC_VERSION
			}

			if RPC_ADDRESS != "" {
				config.Server.Address = RPC_ADDRESS
			}

			if RPC_LIMITS > 0 {
				config.Server.RateLimiter.Limit = RPC_LIMITS
			}

			if ZOOM_CLIENT_ID != "" {
				config.Zoom.ClientID = ZOOM_CLIENT_ID
			}

			if ZOOM_CLIENT_SECRET != "" {
				config.Zoom.ClientSecret = ZOOM_CLIENT_SECRET
			}

			if len(REGISTRY_ADDRESSES) > 0 {
				config.Registry.Addresses = REGISTRY_ADDRESSES
			}

			if REGISTRY_SECURE != "" {
				flag, err := strconv.ParseBool(REGISTRY_SECURE)
				if err == nil {
					config.Registry.Secure = flag
				}
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

			if REPL_ADDRESS != "" {
				config.REPL.Address = REPL_ADDRESS
			}

			if !config.REPL.Debug {
				config.REPL.Debug = REPL_DEBUG
			}

			if PERSISTENCE_URL != "" {
				config.Persistence.Url = PERSISTENCE_URL
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
			service.WithPersistence(config.Persistence),
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
