package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/hive/hivesim"
	"github.com/gorilla/websocket"
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
		httpClient = &http.Client{Timeout: 5 * time.Second}
		httpURL    = fmt.Sprintf("http://%s", net.JoinHostPort(c.IP.String(), "8545"))
		wsURL      = url.URL{Scheme: "ws", Host: net.JoinHostPort(c.IP.String(), "8546")}
		wsConn     *websocket.Conn
		vars       = make(map[string]string)
		err        error
		respBytes  []byte
	)

	for _, msg := range test.messages {
		if msg.send {
			data := substituteVars(msg.data, vars)
			t.Log(">>", data)
			if msg.ws {
				if wsConn == nil {
					dialer := websocket.Dialer{
						HandshakeTimeout: 5 * time.Second,
					}
					// Match the Origin scheme to the WS URL so the underlying
					// golang.org/x/net/websocket same-origin check passes.
					headers := http.Header{}
					headers.Set("Origin", fmt.Sprintf("ws://%s", wsURL.Host))
					wsConn, _, err = dialer.Dial(wsURL.String(), headers)
					if err != nil {
						return fmt.Errorf("websocket dial error: %v", err)
					}
					defer wsConn.Close()
				}
				respBytes, err = postWS(wsConn, strings.NewReader(data))
			} else {
				respBytes, err = postHTTP(httpClient, httpURL, strings.NewReader(data))
			}
			if err != nil {
				return err
			}
		} else {
			if respBytes == nil {
				return fmt.Errorf("invalid test, response before request")
			}
			expectedData := substituteVars(msg.data, vars)
			resp := string(bytes.TrimSpace(respBytes))
			t.Log("<<", resp)
			if !gjson.Valid(resp) {
				return fmt.Errorf("invalid JSON response")
			}

			expectedData, err = applyCaptures(expectedData, resp, vars)
			if err != nil {
				return fmt.Errorf("failed to apply captures: %v", err)
			}

			// Remove error message strings from comparison when both sides error.
			var errorRedacted bool
			resp, expectedData, errorRedacted = redactErrorMessages(resp, expectedData)

			// Ignore response fields not present in the expected payload.
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

func postWS(conn *websocket.Conn, d io.Reader) ([]byte, error) {
	payload, err := io.ReadAll(d)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return nil, fmt.Errorf("websocket write error: %v", err)
	}
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("websocket read error: %v", err)
	}
	return data, nil
}

var varPlaceholder = regexp.MustCompile(`\{\{var:([a-zA-Z0-9_]+)\}\}`)
var capturePlaceholder = regexp.MustCompile(`^\{\{capture:([a-zA-Z0-9_]+)\}\}$`)

func substituteVars(data string, vars map[string]string) string {
	return varPlaceholder.ReplaceAllStringFunc(data, func(match string) string {
		name := varPlaceholder.FindStringSubmatch(match)[1]
		if v, ok := vars[name]; ok {
			return v
		}
		return match
	})
}

func applyCaptures(expected, actual string, vars map[string]string) (string, error) {
	return applyCaptureValue(expected, actual, vars)
}

func applyCaptureValue(expected, actual string, vars map[string]string) (string, error) {
	e := gjson.Parse(expected)
	r := gjson.Parse(actual)

	switch {
	case e.IsObject() && r.IsObject():
		out := "{}"
		var err error
		e.ForEach(func(key, expVal gjson.Result) bool {
			k := key.String()
			respChild := r.Get(k)
			if !respChild.Exists() {
				return true
			}
			filtered, fErr := applyCaptureValue(expVal.Raw, respChild.Raw, vars)
			if fErr != nil {
				err = fErr
				return false
			}
			out, fErr = sjson.SetRaw(out, k, filtered)
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
			filtered, fErr := applyCaptureValue(expVal.Raw, respItem.Raw, vars)
			if fErr != nil {
				err = fErr
				return false
			}
			items = append(items, filtered)
			return true
		})
		return "[" + strings.Join(items, ",") + "]", err

	default:
		if e.Type == gjson.String {
			if m := capturePlaceholder.FindStringSubmatch(e.Str); m != nil {
				vars[m[1]] = r.String()
				return r.Raw, nil
			}
		}
		return expected, nil
	}
}
