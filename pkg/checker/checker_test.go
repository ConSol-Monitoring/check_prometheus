package checker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/consol-monitoring/check_prometheus/internal/helper"
	"github.com/consol-monitoring/check_x"
)

func TestCheckMainWritesQueryOutputToStdout(t *testing.T) {
	// Capture stdout so we can verify the CLI prints the check result.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = w

	oldTimestampFreshness := helper.TimestampFreshness
	oldInsecureSkipVerify := helper.InsecureSkipVerify
	oldCookies := helper.Cookies
	oldVerbose := helper.Verbose

	t.Cleanup(func() {
		os.Stdout = oldStdout
		r.Close()
		helper.TimestampFreshness = oldTimestampFreshness
		helper.InsecureSkipVerify = oldInsecureSkipVerify
		helper.Cookies = oldCookies
		helper.Verbose = oldVerbose
		address = nil
		timeout = 0
		warning = ""
		critical = ""
		query = ""
		queryDecoded = ""
		queryEncoding = Raw
		alias = ""
		search = ""
		replace = ""
		label = ""
		emptyQueryMessage = ""
		emptyQueryStatusArg = ""
		emptyQueryStatus = check_x.State{}
	})

	// Mock Prometheus' query API with a fixed vector result.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/query" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus"},"value":[%d,"1"]}]}}`, time.Now().Unix())
	}))
	t.Cleanup(server.Close)

	// Run the same mode the binary uses: `./check_prometheus m q`.
	code := CheckMain([]string{"check_prometheus", "m", "q", "--address", server.URL, "-q", "up"})

	if err := w.Close(); err != nil {
		t.Fatalf("close stdout pipe: %v", err)
	}

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}

	if code != 0 {
		t.Fatalf("CheckMain returned %d, want 0", code)
	}

	// The important part: the printed stdout should include the query output.
	got := string(output)
	if !strings.Contains(got, "OK - Query: 'up'|") {
		t.Fatalf("stdout %q does not contain query output", got)
	}
	if !strings.Contains(got, "'{__name__=\"up\", job=\"prometheus\"}'=1") {
		t.Fatalf("stdout %q does not contain perfdata output", got)
	}
}
