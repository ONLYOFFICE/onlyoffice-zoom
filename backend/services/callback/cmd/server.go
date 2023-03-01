package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/service/repl"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/shared"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/config"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/web/service"
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
			&cli.IntFlag{
				Name:  "circuit_breaker_timeout",
				Usage: "sets hystrix timeout",
			},
			&cli.IntFlag{
				Name:  "circuit_breaker_max_concurrent",
				Usage: "sets hystrix max concurrency level",
			},
			&cli.IntFlag{
				Name:  "circuit_breaker_volume_threshold",
				Usage: "sets hystrix volume threshold",
			},
			&cli.IntFlag{
				Name:  "circuit_breaker_sleep_window",
				Usage: "sets hystrix sleep window",
			},
			&cli.IntFlag{
				Name:  "circuit_breaker_error_percent",
				Usage: "sets hystrix error percent threshold",
			},
			&cli.StringFlag{
				Name:  "onlyoffice_doc_secret",
				Usage: "sets onlyoffice document server jwt secret",
			},
			&cli.StringFlag{
				Name:  "callback_max_size",
				Usage: "sets callback handler max file size",
			},
			&cli.StringFlag{
				Name:  "callback_upload_timeout",
				Usage: "sets callback file upload timeout",
			},
			&cli.IntFlag{
				Name:  "broker_type",
				Usage: "sets http server's broker type (1 - RabbitMQ, 2 - Memory)",
			},
			&cli.StringSliceFlag{
				Name:  "broker_addresses",
				Usage: "sets http server's broker addresses",
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
			&cli.UintFlag{
				Name:  "worker_max_concurrency",
				Usage: "sets max worker's concurrency",
			},
			&cli.StringSliceFlag{
				Name:  "worker_addresses",
				Usage: "sets worker's broker address",
			},
			&cli.StringFlag{
				Name:  "worker_username",
				Usage: "sets worker's broker username",
			},
			&cli.StringFlag{
				Name:  "worker_password",
				Usage: "sets worker's broker password",
			},
			&cli.IntFlag{
				Name:  "worker_database",
				Usage: "sets worker's broker database",
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

				HTTP_VERSION = c.Int("http_version")
				HTTP_ADDRESS = c.String("http_address")

				HTTP_LIMITS                   = c.Uint64("http_limits")
				HTTP_LIMITS_IP                = c.Uint64("http_iplimits")
				HTTP_CIRCUIT_TIMEOUT          = c.Int("circuit_breaker_timeout")
				HTTP_CIRCUIT_MAX_CONCURRENT   = c.Int("circuit_breaker_max_concurrent")
				HTTP_CIRCUIT_VOLUME_THRESHOLD = c.Int("circuit_breaker_volume_threshold")
				HTTP_CIRCUIT_SLEEP_WINDOW     = c.Int("circuit_breaker_sleep_window")
				HTTP_CIRCUIT_ERROR_PERCENT    = c.Int("circuit_breaker_error_percent")

				ONLYOFFICE_DOC_SECRET = c.String("onlyoffice_doc_secret")

				CALLBACK_MAX_SIZE       = c.String("callback_max_size")
				CALLBACK_UPLOAD_TIMEOUT = c.String("callback_upload_timeout")

				REGISTRY_ADDRESSES = c.StringSlice("registry_addresses")
				REGISTRY_TYPE      = c.Int("registry_type")
				REGISTRY_TTL       = c.Duration("registry_ttl")

				BROKER_TYPE      = c.Int("broker_type")
				BROKER_ADDRESSES = c.StringSlice("broker_addresses")

				TRACER_ADDRESS = c.String("tracer_address")
				TRACER_TYPE    = c.String("tracer_type")
				TRACER_RATIO   = c.Float64("tracer_ratio")

				WORKER_MAX_CONCURRENCY = c.Int("worker_max_concurrency")
				WORKER_ADDRESS         = c.StringSlice("worker_addresses")
				WORKER_USERNAME        = c.String("worker_username")
				WORKER_PASSWORD        = c.String("worker_password")
				WORKER_DATABASE        = c.Int("worker_database")

				REPL_NAME    = c.String("repl_name")
				REPL_ADDRESS = c.String("repl_address")
				REPL_DEBUG   = c.Bool("repl_debug")
			)

			config.Server.Namespace = "onlyoffice"
			config.REPL.Namespace = "onlyoffice"
			config.Server.Name = "callback"
			config.REPL.Name = "callback.repl"

			config.Callback.MaxSize = 2100000
			config.Callback.UploadTimeout = 10
			config.Worker.RedisDatabase = WORKER_DATABASE

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

			if err := envconfig.Process(context.Background(), config); err != nil {
				return err
			}

			if _, ok := shared.SUPPORTED_ENVIRONMENTS[config.Environment]; !ok {
				config.Environment = shared.SUPPORTED_ENVIRONMENTS["development"]
			}

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
				config.Server.Resilience.RateLimiter.Limit = HTTP_LIMITS
			}

			if HTTP_LIMITS_IP > 0 {
				config.Server.Resilience.RateLimiter.IPLimit = HTTP_LIMITS_IP
			}

			if HTTP_CIRCUIT_TIMEOUT > 0 {
				config.Server.Resilience.CircuitBreaker.Timeout = HTTP_CIRCUIT_TIMEOUT
			}

			if HTTP_CIRCUIT_MAX_CONCURRENT > 0 {
				config.Server.Resilience.CircuitBreaker.MaxConcurrent = HTTP_CIRCUIT_MAX_CONCURRENT
			}

			if HTTP_CIRCUIT_VOLUME_THRESHOLD > 0 {
				config.Server.Resilience.CircuitBreaker.VolumeThreshold = HTTP_CIRCUIT_VOLUME_THRESHOLD
			}

			if HTTP_CIRCUIT_SLEEP_WINDOW > 0 {
				config.Server.Resilience.CircuitBreaker.SleepWindow = HTTP_CIRCUIT_SLEEP_WINDOW
			}

			if HTTP_CIRCUIT_ERROR_PERCENT > 0 {
				config.Server.Resilience.CircuitBreaker.ErrorPercentThreshold = HTTP_CIRCUIT_ERROR_PERCENT
			}

			if ONLYOFFICE_DOC_SECRET != "" {
				config.Onlyoffice.DocSecret = ONLYOFFICE_DOC_SECRET
			}

			if CALLBACK_MAX_SIZE != "" {
				v, err := strconv.ParseInt(CALLBACK_MAX_SIZE, 10, 0)
				if err != nil {
					config.Callback.MaxSize = v
				}
			}

			if CALLBACK_UPLOAD_TIMEOUT != "" {
				v, err := strconv.Atoi(CALLBACK_UPLOAD_TIMEOUT)
				if err != nil {
					config.Callback.UploadTimeout = v
				}
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

			if WORKER_MAX_CONCURRENCY > 0 {
				config.Worker.MaxConcurrency = WORKER_MAX_CONCURRENCY
			}

			if len(WORKER_ADDRESS) > 0 {
				config.Worker.RedisAddresses = WORKER_ADDRESS
			}

			if WORKER_USERNAME != "" {
				config.Worker.RedisUsername = WORKER_USERNAME
			}

			if WORKER_PASSWORD != "" {
				config.Worker.RedisPassword = WORKER_PASSWORD
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
			service.WithCallback(config.Callback),
			service.WithOnlyoffice(config.Onlyoffice),
			service.WithTracer(config.Tracer),
			service.WithBroker(config.Broker),
			service.WithLogger(logger),
			service.WithRegistry(config.Registry),
			service.WithCache(config.Cache),
			service.WithWorker(config.Worker),
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
