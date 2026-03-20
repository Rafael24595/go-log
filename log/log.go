package log

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/model/record"
)

var (
	log  Log
	once sync.Once
)

func init() {
	logg, err := newBootstrapLogger()
	if err != nil {
		panic(err.Error())
	}
	log = logg
}

type Provider interface {
	Build() (Log, error)
}

func New(provider Provider) (Log, error) {
	return provider.Build()
}

func DefaultFromProvider(provider Provider) error {
	target, err := provider.Build()
	if err != nil {
		return err
	}
	return DefaultFromLog(target)
}

func DefaultFromLog(target Log) error {
	if target == nil {
		return errors.New("nil logger")
	}

	var init bool
	var err error

	once.Do(func() {
		if b, ok := log.(Bootstrap); ok {
			err = b.Flush(target)
			if err != nil {
				return
			}
		}

		log = target
		log.Messagef("Logging is configured to use the %s instance.", target.Name())

		init = true
	})

	if !init {
		return errors.New("logger already initialized")
	}

	return err
}

func OnClose() error {
	_, err := log.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error closing logger: %v\n", err)
	}
	return err
}

type Bootstrap interface {
	Log
	Flush(Log) error
}

type Log interface {
	Name() logger.Logger
	Records() []record.Record
	Custom(string, string) record.Record
	Custome(string, error) record.Record
	Customf(string, string, ...any) record.Record
	Message(string) record.Record
	Messagef(string, ...any) record.Record
	Warning(string) record.Record
	Warningf(string, ...any) record.Record
	Error(error) record.Record
	Errors(string) record.Record
	Errorf(string, ...any) record.Record
	Record(...record.Record) []record.Record
	Close() ([]record.Record, error)
}

func Name() logger.Logger {
	return log.Name()
}

func Records() []record.Record {
	return log.Records()
}

func Custom(category string, message string) {
	log.Custom(category, message)
}

func Custome(category string, err error) {
	log.Custome(category, err)
}

func Customf(category string, format string, args ...any) {
	log.Customf(category, format, args...)
}

func Message(message string) {
	log.Message(message)
}

func Messagef(format string, args ...any) {
	log.Messagef(format, args...)
}

func Warning(message string) {
	log.Warning(message)
}

func Warningf(format string, args ...any) {
	log.Warningf(format, args...)
}

func Error(err error) {
	log.Error(err)
}

func Errors(message string) {
	log.Errors(message)
}

func Errorf(format string, args ...any) {
	log.Errorf(format, args...)
}
