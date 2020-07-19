package loader_test

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/nenad/undertaker/internal/loader"
	"github.com/nenad/undertaker/internal/storage"
	dbtest "github.com/nenad/undertaker/internal/testing"
)

func TestUndertaker_CollectIntegrationTest(t *testing.T) {
	dsn := os.Getenv("TEST_STORAGE_DSN")

	testDB, err := dbtest.NewSQLDatabase(dsn, "test.__undertaker_test")
	if err != nil {
		t.Fatalf("could not init test db: %s", err)
	}

	testCases := []struct {
		name    string
		fixture string
	}{
		{
			name:    "initial dump is exactly the same",
			fixture: "test_1",
		},
		{
			name:    "existing data will shrink when less functions are reported by tombs",
			fixture: "test_2",
		},
		{
			name:    "existing data will remain the same when initial dump is ran",
			fixture: "test_3",
		},
		{
			name:    "new functions will be added if they are not in the initial dump",
			fixture: "test_4",
		},
	}

	for _, subTest := range testCases {

		srv, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil {
			t.Fatalf("could not start listener: %s", err)
		}

		gravedigger, err := storage.NewPostgres(dsn, "test.__undertaker_test")
		if err != nil {
			t.Fatalf("could not connect to test database: %s", err)
		}

		t.Run(subTest.name, func(t *testing.T) {
			defer func() {
				if err := testDB.Reset(); err != nil {
					t.Fatalf("could not reset database: %s", err)
				}
			}()

			// Prepare database
			sqlFixture := path.Join("testdata", subTest.fixture+".sql")
			if err := testDB.LoadFixture(sqlFixture); err != nil {
				t.Fatalf("could not reset database: %s", err)
			}

			// Prepare tombs output and listener
			tombsOutput, err := ioutil.ReadFile(path.Join("testdata", subTest.fixture+".tombs"))
			if err != nil {
				t.Fatalf("could not open file: %s", err)
			}
			go func() {
				c, err := srv.Accept()
				if err != nil {
					t.Errorf("could not accept connection: %s", err)
				}
				_, _ = c.Write(tombsOutput)
				_ = c.Close()
			}()
			time.Sleep(time.Millisecond * 10)

			// Get expectation data
			wantRaw, err := ioutil.ReadFile(path.Join("testdata", subTest.fixture+".out"))
			if err != nil {
				t.Fatalf("could not open expectation file: %s", err)
			}
			want := strings.Split(strings.TrimSpace(string(wantRaw)), "\n")

			u := loader.Undertaker{
				Gravedigger:  gravedigger,
				TombsAddress: srv.Addr().String(),
			}

			got, err := u.Collect()
			if err != nil {
				t.Fatalf("did not expect Collect to fail: %s", err)
			}

			sort.Slice(want, func(i, j int) bool {
				return want[i] < want[j]
			})
			sort.Slice(got, func(i, j int) bool {
				return got[i] < got[j]
			})

			if !reflect.DeepEqual(want, got) {
				t.Errorf("difference found in tombs; want:\n%s\ngot:\n%s\n", strings.Join(want, "\n"), strings.Join(got, "\n"))
			}
		})
	}
}
