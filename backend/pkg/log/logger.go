package log

import (
	"fmt"
	"log"
	"os"
)

// Logger is a generic logger interface.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// EmptyLogger is an empty Logger implementation.
type EmptyLogger struct{}

// NewEmptyLogger is an empty logger constructor.
func NewEmptyLogger() Logger {
	return EmptyLogger{}
}

func (l EmptyLogger) Debugf(format string, args ...interface{}) {}
func (l EmptyLogger) Infof(format string, args ...interface{})  {}
func (l EmptyLogger) Warnf(format string, args ...interface{})  {}
func (l EmptyLogger) Errorf(format string, args ...interface{}) {}
func (l EmptyLogger) Fatalf(format string, args ...interface{}) {}
func (l EmptyLogger) Debug(args ...interface{})                 {}
func (l EmptyLogger) Info(args ...interface{})                  {}
func (l EmptyLogger) Warn(args ...interface{})                  {}
func (l EmptyLogger) Error(args ...interface{})                 {}
func (l EmptyLogger) Fatal(args ...interface{})                 {}

// DefaultLogger is a golang log package Logger implementation.
type DefaultLogger struct {
	debugL  *log.Logger
	infoL   *log.Logger
	warnL   *log.Logger
	errorL  *log.Logger
	fatalL  *log.Logger
	options Options
}

// NewDefaultLogger is a golang log package Logger constructor.
func NewDefaultLogger(opts ...Option) Logger {
	o := NewOptions(opts...)
	return DefaultLogger{
		debugL:  log.New(os.Stdout, fmt.Sprintf("[DEBUG - Default %s]: ", o.Name), log.Ldate|log.Ltime|log.Llongfile),
		infoL:   log.New(os.Stdout, fmt.Sprintf("[INFO - Default %s]: ", o.Name), log.Ldate|log.Ltime|log.Lshortfile),
		warnL:   log.New(os.Stdout, fmt.Sprintf("[WARN - Default %s]: ", o.Name), log.Ldate|log.Ltime|log.Lshortfile),
		errorL:  log.New(os.Stdout, fmt.Sprintf("[ERROR - Default %s]: ", o.Name), log.Ldate|log.Ltime|log.Lshortfile),
		fatalL:  log.New(os.Stderr, fmt.Sprintf("[FATAL - Default %s]: ", o.Name), log.Ldate|log.Ltime|log.Llongfile),
		options: o,
	}
}

func (l DefaultLogger) Debugf(format string, args ...interface{}) {
	if l.options.Level <= 2 {
		l.debugL.Printf(format+"\n", args...)
	}
}

func (l DefaultLogger) Infof(format string, args ...interface{}) {
	if l.options.Level <= 3 {
		l.infoL.Printf(format+"\n", args...)
	}
}

func (l DefaultLogger) Warnf(format string, args ...interface{}) {
	if l.options.Level <= 4 {
		l.warnL.Printf(format+"\n", args...)
	}
}

func (l DefaultLogger) Errorf(format string, args ...interface{}) {
	if l.options.Level <= 5 {
		l.errorL.Printf(format+"\n", args...)
	}
}

func (l DefaultLogger) Fatalf(format string, args ...interface{}) {
	if l.options.Level <= 6 {
		l.fatalL.Fatalf(format+"\n", args...)
	}
}

func (l DefaultLogger) Debug(args ...interface{}) {
	if l.options.Level <= 2 {
		l.debugL.Println(args...)
	}
}

func (l DefaultLogger) Info(args ...interface{}) {
	if l.options.Level <= 3 {
		l.infoL.Println(args...)
	}
}

func (l DefaultLogger) Warn(args ...interface{}) {
	if l.options.Level <= 4 {
		l.warnL.Println(args...)
	}
}

func (l DefaultLogger) Error(args ...interface{}) {
	if l.options.Level <= 5 {
		l.errorL.Println(args...)
	}
}

func (l DefaultLogger) Fatal(args ...interface{}) {
	if l.options.Level <= 6 {
		l.fatalL.Fatalln(args...)
	}
}
