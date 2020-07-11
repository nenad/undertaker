package main

import (
	"net/http"
	"strings"
	"time"
)

func NewServer(undertaker Undertaker, port string) http.Server {
	preloadHandler := &preloadHandler{undertaker: undertaker}
	collectHandler := &collectHandler{undertaker: undertaker}

	mux := http.NewServeMux()
	mux.Handle("/preload", preloadHandler)
	mux.Handle("/collect", collectHandler)

	return http.Server{
		Handler:           mux,
		Addr:              ":" + port,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 10,
		IdleTimeout:       time.Second * 10,
	}
}

type preloadHandler struct {
	undertaker Undertaker
}

func (h *preloadHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if err := h.undertaker.Preload(); err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write([]byte("Successfully preloaded classes"))
}

type collectHandler struct {
	undertaker Undertaker
}

func (h *collectHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	funcs, err := h.undertaker.Collect()
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, _ = w.Write([]byte(strings.Join(funcs, "\n")))
}
