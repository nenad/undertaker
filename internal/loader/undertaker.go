package loader

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/yookoala/gofast"
)

type Undertaker struct {
	FPMAddr      string
	TombsAddress string
	PreloadFile  string
}

func (u *Undertaker) Preload() error {
	connFactory := gofast.SimpleConnFactory("tcp", u.FPMAddr)
	// The HTTP request is a dummy request just to satisfy the library requirement for a valid URL.
	req, _ := http.NewRequest("GET", "0.0.0.0", nil)
	resp := httptest.NewRecorder()
	h := gofast.NewHandler(
		gofast.NewFileEndpoint(u.PreloadFile)(gofast.BasicSession),
		gofast.SimpleClientFactory(connFactory, 0),
	)
	h.ServeHTTP(resp, req)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read body: %w", err)
	}

	if len(b) != 0 {
		return fmt.Errorf("error while executing preloader: %s", b)
	}

	return nil
}

func (u *Undertaker) Collect() ([]string, error) {
	conn, err := net.Dial("tcp", u.TombsAddress)
	if err != nil {
		return nil, fmt.Errorf("could not connect to tombs socket: %w", err)
	}
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("could not read bytes from tombs: %w", err)
	}

	return strings.Split(strings.TrimSpace(string(b)), "\n"), nil
}
