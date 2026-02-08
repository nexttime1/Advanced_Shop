package options

import (
	"Advanced_Shop/pkg/errors"
	"github.com/spf13/pflag"
)

type TelemetryOptions struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Sampler  float64 `json:"sampler"`
	Batcher  string  `json:"batcher"`
}

func NewTelemetryOptions() *TelemetryOptions {
	return &TelemetryOptions{
		Name:     "xshop",
		Endpoint: "http://127.0.0.1:14268/api/traces",
		Sampler:  1.0,
		Batcher:  "jaeger",
	}
}

func (t *TelemetryOptions) Validate() []error {
	errs := []error{}
	if t.Batcher != "jaeger" && t.Batcher != "zipkin" {
		errs = append(errs, errors.New("opentelemetry batcher only support jaeger or zipkin"))
	}
	return errs
}

// AddFlags adds flags related to open telemetry for a specific tracing to the specified FlagSet.
func (to *TelemetryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&to.Name, "telemetry.name", to.Name, "opentelemetry name")

	fs.StringVar(&to.Endpoint, "telemetry.endpoint", to.Endpoint, "opentelemetry endpoint")
	fs.Float64Var(&to.Sampler, "telemetry.sampler", to.Sampler, "telemetry sampler")
	fs.StringVar(&to.Batcher, "telemetry.batcher", to.Batcher, "telemetry batcher, only support jaeger and zipkin")
}
