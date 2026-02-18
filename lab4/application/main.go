package main

import (
	"application-for-kubernetes/logger"
	"fmt"
	"net/http"
	"sync/atomic"
)

var counter int32

func main() {
	const headerPassword = "iad1i9f1yj0ne5zo"
	log := logger.SetupLogger("K8S LAB4")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		fmt.Fprintf(w, "Hello\n")
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK\n")
		log.Info("Readiness check OK")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		if r.Header.Get("Authorization") != headerPassword {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "UNAUTH\n")
			log.Warnf("Unauthorized health check attempt from %s", r.RemoteAddr)
			return
		}

		c := atomic.AddInt32(&counter, 1)

		if c > 5 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "UNHEALTHY\n")
			log.Errorf("Health check failed: unhealthy (counter=%d)", c)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK\n")
		log.Infof("Health check OK (counter=%d)", c)
	})

	log.Infof("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}