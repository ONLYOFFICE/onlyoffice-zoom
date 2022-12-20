package log

import (
	"os"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log/hook"
	"github.com/natefinch/lumberjack"
	elastic "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

var levels = map[LogLevel]logrus.Level{
	LEVEL_TRACE:   logrus.TraceLevel,
	LEVEL_DEBUG:   logrus.DebugLevel,
	LEVEL_INFO:    logrus.InfoLevel,
	LEVEL_WARNING: logrus.WarnLevel,
	LEVEL_ERROR:   logrus.ErrorLevel,
	LEVEL_FATAL:   logrus.FatalLevel,
}

// LogrusLogger is a logrus logger wrapper.
type LogrusLogger struct {
	logger  *logrus.Logger
	options Options
}

// createElasticHook opens a new elastic client and generates an elastic hook.
func createElasticHook(options Options) (*hook.ElasticHook, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(options.Elastic.Address),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(options.Elastic.BasicAuthUsername, options.Elastic.BasicAuthPassword),
		elastic.SetHealthcheck(options.Elastic.HealthcheckEnabled),
		elastic.SetGzip(options.Elastic.GzipEnabled),
	)

	if err != nil {
		return nil, &LogElasticInitializationError{
			Address: options.Elastic.Address,
			Cause:   err,
		}
	}

	if options.Elastic.Bulk {
		return hook.NewBulkProcessorElasticHook(client, options.Elastic.Address, levels[options.Level], options.Elastic.Index)
	}

	if options.Elastic.Async {
		return hook.NewAsyncElasticHook(client, options.Elastic.Address, levels[options.Level], options.Elastic.Index)
	}

	return hook.NewElasticHook(client, options.Elastic.Address, levels[options.Level], options.Elastic.Index)
}

// NewLogrusLogger creates a new logger compliant with the Logger interface.
func NewLogrusLogger(opts ...Option) (Logger, error) {
	options := NewOptions(opts...)

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: !options.Color,
		FullTimestamp: true,
	})

	if lvl, ok := levels[options.Level]; ok {
		log.SetLevel(lvl)
	}

	log.SetReportCaller(options.ReportCaller)
	log.SetOutput(os.Stdout)

	if options.File.Filename != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   options.File.Filename,
			MaxSize:    options.File.MaxSize,
			MaxBackups: options.File.MaxBackups,
			MaxAge:     options.File.MaxAge,
			LocalTime:  options.File.LocalTime,
			Compress:   options.File.Compress,
		})
	}

	if options.File.Filename == "" && options.Elastic.Address != "" && options.Elastic.Index != "" {
		hook, err := createElasticHook(options)

		if err != nil {
			return nil, &LogElasticInitializationError{
				Address: options.Elastic.Address,
				Cause:   err,
			}
		}

		log.AddHook(hook)
	}

	return LogrusLogger{
		logger:  log,
		options: options,
	}, nil
}

func (l LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Debugf(format, args...)
}

func (l LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Infof(format, args...)
}

func (l LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Warnf(format, args...)
}

func (l LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Errorf(format, args...)
}

func (l LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Fatalf(format, args...)
}

func (l LogrusLogger) Debug(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Debug(args...)
}

func (l LogrusLogger) Info(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Info(args...)
}

func (l LogrusLogger) Warn(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Warn(args...)
}

func (l LogrusLogger) Error(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Error(args...)
}

func (l LogrusLogger) Fatal(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"name":    l.options.Name,
		"machine": l.options.Machine,
		"env":     l.options.Environment,
	}).Fatal(args...)
}
