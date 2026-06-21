package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/hive/hivesim"
	"github.com/nsf/jsondiff"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var files = map[string]string{
	"genesis.json": "./tests/genesis.json",
}

func main() {
	runTestSuiteMain()
}

func runTestSuiteMain() {
	suite := hivesim.Suite{
		Name: "xdc-rpc-compat",
		Description: `
The XDC RPC-compatibility test suite runs a curated set of JSON-RPC tests
against an XDC-compatible client. It validates standard and XDC-specific RPC
endpoints using XDC testnet genesis and expected responses.`,
	}
	suite.Add(&hivesim.ClientTestSpec{
		Role:        "eth1",
		Name:        "client launch",
		Description: `This test launches the XDC client and runs all RPC tests.`,
		Parameters: map[string]string{
			"HIVE_MINER": "0x746249c61f5832c5eed53172776b460491bdcd5c",
		},
		Files: files,
		Run: func(t *hivesim.T, c *hivesim.Client) {
			runAllTests(t, c, c.Type)
		},
		AlwaysRun: true,
	})
	sim := hivesim.New()
	hivesim.MustRunSuite(sim, suite)
}

func runAllTests(t *hivesim.T, c *hivesim.Client, clientName string) {
	_, testPattern := t.Sim.TestPattern()
	re := regexp.MustCompile(testPattern)
	tests := loadTests(t, "tests", re)
	for _, test := range tests {
		test := test
		t.Run(hivesim.TestSpec{
			Name:        fmt.Sprintf("%s (%s)", test.name, clientName),
			Description: test.comment,
			Run: func(t *hivesim.T) {
				if err := runTest(t, c, &test); err != nil {
					t.Fatal(err)
				}
			},
		})
	}
}

func runTest(t *hivesim.T, c *hivesim.Client, test *rpcTest) error {
	var (
		client    = &http.Client{Timeout: 5 * time.Second}
		url       = fmt.Sprintf("http://%s", net.JoinHostPort(c.IP.String(), "8545"))
		err       error
		respBytes []byte
	)

	for _, msg := range test.messages {
		if msg.send {
			t.Log(">>", msg.data)
			respBytes, err = postHTTP(client, url, strings.NewReader(msg.data))
			if err != nil {
				return err
			}
		} else {
			if respBytes == nil {
				return fmt.Errorf("invalid test, response before request")
			}
			expectedData := msg.data
			resp := string(bytes.TrimSpace(respBytes))
			t.Log("<<", resp)
			if !gjson.Valid(resp) {
				return fmt.Errorf("invalid JSON response")
			}

			// Remove error message strings from comparison when both sides error.
			var errorRedacted bool
			resp, expectedData, errorRedacted = redactErrorMessages(resp, expectedData)

			// Ignore response fields not present in the expected payload, allowing
			// fixtures to assert only the subset of fields they care about.
			resp, err = filterToExpected(resp, expectedData)
			if err != nil {
				return fmt.Errorf("failed to filter response: %v", err)
			}

			opts := &jsondiff.Options{
				Added:            jsondiff.Tag{Begin: "++ "},
				Removed:          jsondiff.Tag{Begin: "-- "},
				Changed:          jsondiff.Tag{Begin: "-- "},
				ChangedSeparator: " ++ ",
				Indent:           "  ",
				CompareNumbers:   numbersEqual,
			}
			diffStatus, diffText := jsondiff.Compare([]byte(resp), []byte(expectedData), opts)
			if diffStatus != jsondiff.FullMatch {
				if errorRedacted {
					t.Log("note: error messages removed from comparison")
				}
				return fmt.Errorf("response differs from expected (-- client, ++ test):\n%s", diffText)
			}
			respBytes = nil
		}
	}

	if respBytes != nil {
		t.Fatalf("unhandled response in test case")
	}
	return nil
}

func redactErrorMessages(resp, expected string) (string, string, bool) {
	paths := collectErrorMessagePaths(nil, "", gjson.Parse(resp), gjson.Parse(expected))
	if len(paths) == 0 {
		return resp, expected, false
	}
	for _, p := range paths {
		resp, _ = sjson.Delete(resp, p)
		expected, _ = sjson.Delete(expected, p)
	}
	return resp, expected, true
}

func collectErrorMessagePaths(paths []string, path string, respVal, expectedVal gjson.Result) []string {
	switch {
	case expectedVal.IsObject():
		expectedVal.ForEach(func(key, val gjson.Result) bool {
			k := key.String()
			respChild := respVal.Get(k)
			if !respChild.Exists() {
				return true
			}
			childPath := joinPath(path, k)
			if k == "error" {
				if val.Get("message").Exists() && respChild.Get("message").Exists() {
					paths = append(paths, childPath+".message")
				}
				return true
			}
			paths = collectErrorMessagePaths(paths, childPath, respChild, val)
			return true
		})
	case expectedVal.IsArray():
		i := 0
		expectedVal.ForEach(func(_, val gjson.Result) bool {
			idx := strconv.Itoa(i)
			i++
			respChild := respVal.Get(idx)
			if !respChild.Exists() {
				return true
			}
			paths = collectErrorMessagePaths(paths, joinPath(path, idx), respChild, val)
			return true
		})
	}
	return paths
}

func joinPath(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}

func numbersEqual(a, b json.Number) bool {
	af, err1 := a.Float64()
	bf, err2 := b.Float64()
	if err1 == nil && err2 == nil {
		return af == bf || math.IsNaN(af) && math.IsNaN(bf)
	}
	return a == b
}

func filterToExpected(resp, expected string) (string, error) {
	return filterValue(resp, expected)
}

func filterValue(resp, expected string) (string, error) {
	r := gjson.Parse(resp)
	e := gjson.Parse(expected)

	switch {
	case e.IsObject() && r.IsObject():
		out := "{}"
		var err error
		e.ForEach(func(key, expVal gjson.Result) bool {
			if !r.Get(key.String()).Exists() {
				return true
			}
			filtered, fErr := filterValue(r.Get(key.String()).Raw, expVal.Raw)
			if fErr != nil {
				err = fErr
				return false
			}
			out, fErr = sjson.SetRaw(out, key.String(), filtered)
			if fErr != nil {
				err = fErr
				return false
			}
			return true
		})
		return out, err

	case e.IsArray() && r.IsArray():
		var items []string
		i := 0
		var err error
		e.ForEach(func(_, expVal gjson.Result) bool {
			idx := strconv.Itoa(i)
			i++
			respItem := r.Get(idx)
			if !respItem.Exists() {
				return true
			}
			filtered, fErr := filterValue(respItem.Raw, expVal.Raw)
			if fErr != nil {
				err = fErr
				return false
			}
			items = append(items, filtered)
			return true
		})
		return "[" + strings.Join(items, ",") + "]", err

	default:
		return resp, nil
	}
}

func postHTTP(c *http.Client, url string, d io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", url, d)
	if err != nil {
		return nil, fmt.Errorf("error building request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("write error: %v", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
