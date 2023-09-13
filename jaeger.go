package tracer

import (
	"fmt"
	"github.com/go-tron/config"
	"github.com/opentracing/opentracing-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"io"
	"log"
)

func NewWithConfig(c *config.Config) (opentracing.Tracer, io.Closer) {
	host := c.GetString("jaeger.host")
	if host == "cluster.nodeIP" {
		host = c.GetString(host)
	}
	if host == "" {
		panic("host 必须设置")
	}
	port := c.GetString("jaeger.port")
	if port == "" {
		panic("port 必须设置")
	}
	return New(c.GetString("application.name")+"."+c.GetString("cluster.namespace"), host+":"+port)
}

func New(serviceName string, hostPort string) (opentracing.Tracer, io.Closer) {
	log.Println("jaeger hostPort", hostPort)
	cfg := &jaegerconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: hostPort,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
