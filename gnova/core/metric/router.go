package metric

import (
	zlog "Advanced_Shop/pkg/log"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"net/http"
)

func StartMetricsServer(port int) {
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf(":%d", port)
		zlog.Infof("metrics server listen on %s", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			zlog.Errorf("metrics server error: %s", err)
		}
	}()
}
