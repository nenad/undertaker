package loader_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"reflect"
	"testing"
	"time"

	"github.com/nenad/undertaker/internal/loader"
)

func TestUndertaker_Preload(t *testing.T) {
	tests := []struct {
		name    string
		message string
		err     error
	}{
		{
			name:    "Successful preload without nothing echoed",
			message: "",
			err:     nil,
		},
		{
			name:    "Failed preload when something is echoed",
			message: "Should fail here",
			err:     fmt.Errorf("error while executing preloader: Should fail here"),
		},
	}

	preloadFile := "/server/undertaker.php"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()

			go func() {
				_ = fcgi.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fcgiHeaders := fcgi.ProcessEnv(r)
					actualFilename := fcgiHeaders["SCRIPT_FILENAME"]
					if actualFilename != preloadFile {
						t.Errorf("preload script filename mismatch; got %q want %q", actualFilename, preloadFile)
					}
					_, _ = w.Write([]byte(tt.message))
				}))
			}()

			u := loader.Undertaker{
				FPMAddr:      listener.Addr().String(),
				TombsAddress: listener.Addr().String(),
				PreloadFile:  preloadFile,
			}

			time.Sleep(time.Millisecond * 50)

			err = u.Preload()
			if err == nil && tt.err != nil {
				t.Errorf("error expectation failed; want %q got nil", tt.err)
			}

			if err != nil && tt.err == nil {
				t.Errorf("error expectation failed; want nil got %q", err)
			}

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("error expectation failed; want %q got %q", tt.err, err)
			}
		})
	}
}

func TestUndertaker_Collect(t *testing.T) {
	srv, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Error(err)
	}

	go func() {
		c, err := srv.Accept()
		if err != nil {
			t.Error(err)
		}
		_, _ = c.Write([]byte("hello\nworld"))
		_ = c.Close()
	}()

	u := loader.Undertaker{TombsAddress: srv.Addr().String()}

	str, err := u.Collect()
	if err != nil {
		t.Error(err)
	}

	expect := []string{
		"hello",
		"world",
	}

	if !reflect.DeepEqual(expect, str) {
		t.Errorf("collected functions are not equal; got %#v want %#v", str, expect)
	}
}
