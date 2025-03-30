package app

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/unbxd/go-base/kit/transport/http"
	"github.com/unbxd/go-base/utils/log"
	"github.com/unbxd/go-base/utils/metrics"
	"github.com/unbxd/go-base/utils/notifier"
)

// Option defines the way in which we can modify the app
// object for initialisation.
//
// This method along with `Newapp` will be used to
// initialise the project properties
type Option func(*App) (err error)

func WithLogger(
	level string,
	encoding string,
	output string,
) Option {
	return func(s *App) (err error) {
		logger, err := log.NewZapLogger(
			log.ZapWithLevel(level),
			log.ZapWithEncoding(encoding),
			log.ZapWithOutput([]string{output}),
		)
		if err != nil {
			return errors.Wrap(err, "failed to initialise logging")
		}

		s.logger = logger
		return err
	}
}

func WithCustomLogger(logger log.Logger) Option {
	return func(s *App) (err error) {
		s.logger = logger
		return
	}
}

func WithMetrics(
	enabled bool,
	conn string,
	namespace string,
	tags []string,
) Option {
	return func(s *App) (err error) {
		met, err := metrics.NewDatadogMetrics(
			metrics.WithDatadogEnabled(enabled),
			metrics.WithDatadogNamespace(namespace),
			metrics.WithDatadogLogger(s.logger),
			metrics.WithDatadogServerConnstr(conn),
			metrics.WithDatadogTags(tags),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create metrics")
		}

		s.metrics = met
		return
	}
}

// WithNotifier sets the notifier for the Overpass
func WithNotifier(
	enabled bool,
	hosts []string,
	name string,
	prefix string,
) Option {
	return func(s *App) error {
		if !enabled {
			return nil
		}

		not, err := notifier.NewNotifier(
			strings.Join(hosts, ","),
			notifier.WithSubjectPrefix(prefix),
			notifier.WithName(name),
			notifier.WithDefaultOptions(),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create notifier")
		}

		s.notifier = not
		return err
	}
}

// WithHTTPTransport supports custom http transport
func WithHTTPTransport(
	host, port string,
	monitor []string,
	opts ...http.TransportOption,
) Option {
	return func(s *App) (err error) {
		options := append(
			[]http.TransportOption{
				http.WithLogger(s.logger),
				http.WithFullDefaults(),
				http.WithMonitors(monitor),
			}, opts...)

		tr, err := http.NewTransport(
			host,
			port,
			options...,
		)

		if err != nil {
			return err
		}

		s.httpTransport = tr
		return
	}
}
