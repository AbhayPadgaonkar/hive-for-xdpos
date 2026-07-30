package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var bom = []byte{0xEF, 0xBB, 0xBF}

func readJSONFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.TrimPrefix(data, bom), nil
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// --- Run handlers ---

func listRuns(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT r.id, r.client_id, c.name, r.simulator_id, s.name,
		       r.date_run, r.total, r.passed, r.failed, COALESCE(r.version,''), r.created_at
		FROM runs r
		JOIN clients c ON c.id = r.client_id
		JOIN simulators s ON s.id = r.simulator_id
		ORDER BY r.created_at DESC
	`)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var runs []Run
	for rows.Next() {
		var ru Run
		if err := rows.Scan(&ru.ID, &ru.ClientID, &ru.ClientName, &ru.SimulatorID, &ru.SimName,
			&ru.DateRun, &ru.Total, &ru.Passed, &ru.Failed, &ru.Version, &ru.CreatedAt); err != nil {
			continue
		}
		runs = append(runs, ru)
	}
	if runs == nil {
		runs = []Run{}
	}
	writeJSON(w, 200, runs)
}

func getRun(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	var ru Run
	err = db.QueryRow(`
		SELECT r.id, r.client_id, c.name, r.simulator_id, s.name,
		       r.date_run, r.total, r.passed, r.failed, COALESCE(r.version,''), COALESCE(r.raw_log_path,''), r.created_at
		FROM runs r
		JOIN clients c ON c.id = r.client_id
		JOIN simulators s ON s.id = r.simulator_id
		WHERE r.id = ?
	`, id).Scan(&ru.ID, &ru.ClientID, &ru.ClientName, &ru.SimulatorID, &ru.SimName,
		&ru.DateRun, &ru.Total, &ru.Passed, &ru.Failed, &ru.Version, &ru.RawLogPath, &ru.CreatedAt)
	if err != nil {
		writeError(w, 404, "run not found")
		return
	}

	trows, err := db.Query("SELECT id, run_id, test_name, passed FROM tests WHERE run_id = ?", id)
	if err == nil {
		defer trows.Close()
		for trows.Next() {
			var tr TestResult
			if trows.Scan(&tr.ID, &tr.RunID, &tr.TestName, &tr.Passed) == nil {
				// included inline
			}
		}
	}

	writeJSON(w, 200, ru)
}

func createRun(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Client      string `json:"client"`
		Simulator   string `json:"simulator"`
		DateRun     string `json:"date_run"`
		Total       int    `json:"total"`
		Passed      int    `json:"passed"`
		Failed      int    `json:"failed"`
		Version     string `json:"version"`
		RawLogPath  string `json:"raw_log_path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid JSON")
		return
	}
	clientID, err := ensureClient(input.Client)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	simID, err := ensureSimulator(input.Simulator)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	res, err := db.Exec(
		"INSERT INTO runs (client_id, simulator_id, date_run, total, passed, failed, version, raw_log_path) VALUES (?,?,?,?,?,?,?,?)",
		clientID, simID, input.DateRun, input.Total, input.Passed, input.Failed, input.Version, input.RawLogPath,
	)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	writeJSON(w, 201, map[string]int64{"id": id})
}

// --- Probe handlers ---

func listProbes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT p.id, p.client_id, c.name, COALESCE(p.version,''),
		       p.supported, p.unsupported, p.total, p.date_run, p.created_at
		FROM probes p
		JOIN clients c ON c.id = p.client_id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var probes []Probe
	for rows.Next() {
		var p Probe
		if err := rows.Scan(&p.ID, &p.ClientID, &p.ClientName, &p.Version,
			&p.Supported, &p.Unsupported, &p.Total, &p.DateRun, &p.CreatedAt); err != nil {
			continue
		}
		probes = append(probes, p)
	}
	if probes == nil {
		probes = []Probe{}
	}
	writeJSON(w, 200, probes)
}

func createProbe(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Client      string   `json:"client"`
		Version     string   `json:"version"`
		Supported   int      `json:"supported"`
		Unsupported int      `json:"unsupported"`
		Total       int      `json:"total"`
		Modules     string   `json:"modules"`
		DateRun     string   `json:"date_run"`
		Methods     []struct {
			Method      string `json:"method"`
			Supported   bool   `json:"supported"`
			SampleValue string `json:"sample_value,omitempty"`
			Error       string `json:"error,omitempty"`
		} `json:"methods"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid JSON: "+err.Error())
		return
	}
	clientID, err := ensureClient(input.Client)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	tx, err := db.Begin()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO probes (client_id, version, supported, unsupported, total, modules, date_run) VALUES (?,?,?,?,?,?,?)",
		clientID, input.Version, input.Supported, input.Unsupported, input.Total, input.Modules, input.DateRun,
	)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	probeID, _ := res.LastInsertId()

	for _, m := range input.Methods {
		sv := m.SampleValue
		errStr := m.Error
		_, err := tx.Exec(
			"INSERT INTO probe_methods (probe_id, method, supported, sample_value, error) VALUES (?,?,?,?,?)",
			probeID, m.Method, boolToInt(m.Supported), sv, errStr,
		)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]int64{"id": probeID})
}

// --- Gap matrix handlers ---

func listGapMatrices(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT g.id, g.client_a_id, ca.name, g.client_b_id, cb.name,
		       COALESCE(g.version_a,''), COALESCE(g.version_b,''),
		       g.total_methods, g.both_supported, g.both_unsupported, g.in_a_not_b, g.in_b_not_a,
		       g.date_created
		FROM gap_matrices g
		JOIN clients ca ON ca.id = g.client_a_id
		JOIN clients cb ON cb.id = g.client_b_id
		ORDER BY g.date_created DESC
	`)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var matrices []GapMatrix
	for rows.Next() {
		var m GapMatrix
		if err := rows.Scan(&m.ID, &m.ClientAID, &m.ClientAName, &m.ClientBID, &m.ClientBName,
			&m.VersionA, &m.VersionB, &m.TotalMethods, &m.BothSupported, &m.BothUnsupp,
			&m.InANotB, &m.InBNotA, &m.DateCreated); err != nil {
			continue
		}
		matrices = append(matrices, m)
	}
	if matrices == nil {
		matrices = []GapMatrix{}
	}
	writeJSON(w, 200, matrices)
}

func getGapMatrix(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	var m GapMatrix
	err = db.QueryRow(`
		SELECT g.id, g.client_a_id, ca.name, g.client_b_id, cb.name,
		       COALESCE(g.version_a,''), COALESCE(g.version_b,''),
		       COALESCE(g.modules_a,''), COALESCE(g.modules_b,''),
		       g.total_methods, g.both_supported, g.both_unsupported, g.in_a_not_b, g.in_b_not_a,
		       g.date_created
		FROM gap_matrices g
		JOIN clients ca ON ca.id = g.client_a_id
		JOIN clients cb ON cb.id = g.client_b_id
		WHERE g.id = ?
	`, id).Scan(&m.ID, &m.ClientAID, &m.ClientAName, &m.ClientBID, &m.ClientBName,
		&m.VersionA, &m.VersionB, &m.ModulesA, &m.ModulesB,
		&m.TotalMethods, &m.BothSupported, &m.BothUnsupp, &m.InANotB, &m.InBNotA, &m.DateCreated)
	if err != nil {
		writeError(w, 404, "gap matrix not found")
		return
	}

	mrows, err := db.Query(
		"SELECT method, a_supported, b_supported, COALESCE(a_error,''), COALESCE(b_error,'') FROM gap_matrix_methods WHERE matrix_id = ? ORDER BY method",
		id,
	)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer mrows.Close()

	type gapMethod struct {
		Method     string `json:"method"`
		ASupported bool   `json:"a_supported"`
		BSupported bool   `json:"b_supported"`
		AError     string `json:"a_error,omitempty"`
		BError     string `json:"b_error,omitempty"`
	}
	var methods []gapMethod
	for mrows.Next() {
		var gm gapMethod
		var aSup, bSup int
		if err := mrows.Scan(&gm.Method, &aSup, &bSup, &gm.AError, &gm.BError); err != nil {
			continue
		}
		gm.ASupported = aSup == 1
		gm.BSupported = bSup == 1
		methods = append(methods, gm)
	}
	if methods == nil {
		methods = []gapMethod{}
	}

	writeJSON(w, 200, map[string]any{
		"id":               m.ID,
		"client_a_name":    m.ClientAName,
		"client_b_name":    m.ClientBName,
		"version_a":        m.VersionA,
		"version_b":        m.VersionB,
		"modules_a":        m.ModulesA,
		"modules_b":        m.ModulesB,
		"total_methods":    m.TotalMethods,
		"both_supported":   m.BothSupported,
		"both_unsupported": m.BothUnsupp,
		"in_a_not_b":       m.InANotB,
		"in_b_not_a":       m.InBNotA,
		"date_created":     m.DateCreated,
		"methods":          methods,
	})
}

func createGapMatrix(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ClientA    string `json:"client_a"`
		ClientB    string `json:"client_b"`
		VersionA   string `json:"version_a"`
		VersionB   string `json:"version_b"`
		ModulesA   string `json:"modules_a"`
		ModulesB   string `json:"modules_b"`
		TotalM     int    `json:"total_methods"`
		BothSupp   int    `json:"both_supported"`
		BothUnsupp int    `json:"both_unsupported"`
		InANotB    int    `json:"in_a_not_b"`
		InBNotA    int    `json:"in_b_not_a"`
		DateCreate string `json:"date_created"`
		Methods    []struct {
			Method     string `json:"method"`
			ASupported bool   `json:"a_supported"`
			BSupported bool   `json:"b_supported"`
			AError     string `json:"a_error,omitempty"`
			BError     string `json:"b_error,omitempty"`
		} `json:"methods"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid JSON")
		return
	}
	caID, err := ensureClient(input.ClientA)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	cbID, err := ensureClient(input.ClientB)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	tx, err := db.Begin()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO gap_matrices (client_a_id, client_b_id, version_a, version_b, modules_a, modules_b,
			total_methods, both_supported, both_unsupported, in_a_not_b, in_b_not_a, date_created)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
	`, caID, cbID, input.VersionA, input.VersionB, input.ModulesA, input.ModulesB,
		input.TotalM, input.BothSupp, input.BothUnsupp, input.InANotB, input.InBNotA, input.DateCreate)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	matrixID, _ := res.LastInsertId()

	for _, m := range input.Methods {
		_, err := tx.Exec(
			"INSERT INTO gap_matrix_methods (matrix_id, method, a_supported, b_supported, a_error, b_error) VALUES (?,?,?,?,?,?)",
			matrixID, m.Method, boolToInt(m.ASupported), boolToInt(m.BSupported), m.AError, m.BError,
		)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]int64{"id": matrixID})
}

// --- Comparison handlers ---

func listComparisons(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, simulator, client_a_name, client_b_name, both_pass, a_only, b_only, both_fail, total_matched, date_compared, created_at
		FROM comparisons ORDER BY created_at DESC
	`)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var comps []Comparison
	for rows.Next() {
		var c Comparison
		if err := rows.Scan(&c.ID, &c.Simulator, &c.ClientAName, &c.ClientBName,
			&c.BothPass, &c.AOnly, &c.BOnly, &c.BothFail, &c.TotalMatched, &c.DateCompared, &c.CreatedAt); err != nil {
			continue
		}
		comps = append(comps, c)
	}
	if comps == nil {
		comps = []Comparison{}
	}
	writeJSON(w, 200, comps)
}

func createComparison(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Simulator    string `json:"simulator"`
		ClientA      string `json:"client_a_name"`
		ClientB      string `json:"client_b_name"`
		BothPass     int    `json:"both_pass"`
		AOnly        int    `json:"a_only"`
		BOnly        int    `json:"b_only"`
		BothFail     int    `json:"both_fail"`
		TotalMatched int    `json:"total_matched"`
		DateCompared string `json:"date_compared"`
		Tests        []struct {
			TestName string `json:"test_name"`
			APass    *bool  `json:"a_pass"`
			BPass    *bool  `json:"b_pass"`
			Status   string `json:"status"`
		} `json:"tests"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid JSON")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO comparisons (simulator, client_a_name, client_b_name, both_pass, a_only, b_only, both_fail, total_matched, date_compared)
		VALUES (?,?,?,?,?,?,?,?,?)
	`, input.Simulator, input.ClientA, input.ClientB, input.BothPass, input.AOnly, input.BOnly,
		input.BothFail, input.TotalMatched, input.DateCompared)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	compID, _ := res.LastInsertId()

	for _, t := range input.Tests {
		_, err := tx.Exec(
			"INSERT INTO comparison_tests (comparison_id, test_name, a_pass, b_pass, status) VALUES (?,?,?,?,?)",
			compID, t.TestName, nullableBool(t.APass), nullableBool(t.BPass), t.Status,
		)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]int64{"id": compID})
}

// --- Stats handler ---

func getStats(w http.ResponseWriter, r *http.Request) {
	var s Stats
	db.QueryRow("SELECT COUNT(*) FROM runs").Scan(&s.TotalRuns)
	db.QueryRow("SELECT COUNT(*) FROM probes").Scan(&s.TotalProbes)
	db.QueryRow("SELECT COUNT(*) FROM comparisons").Scan(&s.TotalComparisons)
	db.QueryRow("SELECT COUNT(*) FROM gap_matrices").Scan(&s.TotalGapMatrices)
	db.QueryRow("SELECT COUNT(*) FROM tests").Scan(&s.TotalTests)

	// Pass rates for xdc-geth-audit and go-ethereum
	var xdcTotal, xdcPassed, goTotal, goPassed int
	db.QueryRow("SELECT COALESCE(SUM(total),0), COALESCE(SUM(passed),0) FROM runs WHERE client_id = (SELECT id FROM clients WHERE name = 'xdc-geth-audit')").Scan(&xdcTotal, &xdcPassed)
	db.QueryRow("SELECT COALESCE(SUM(total),0), COALESCE(SUM(passed),0) FROM runs WHERE client_id = (SELECT id FROM clients WHERE name = 'go-ethereum')").Scan(&goTotal, &goPassed)
	if xdcTotal > 0 {
		s.XdcGethPassRate = float64(xdcPassed) / float64(xdcTotal) * 100
	}
	if goTotal > 0 {
		s.GoEthPassRate = float64(goPassed) / float64(goTotal) * 100
	}

	writeJSON(w, 200, s)
}

// --- Import from error_ledger JSON files ---

func importRunFromFile(path string) error {
	data, err := readJSONFile(path)
	if err != nil {
		return err
	}
	var input struct {
		Client    string `json:"client"`
		Simulator string `json:"simulator"`
		Date      string `json:"date"`
		Timestamp string `json:"timestamp"`
		Command   string `json:"command"`
		Version   string `json:"version"`
		Suites    int    `json:"suites"`
		Total     int    `json:"total"`
		Passed    int    `json:"passed"`
		Failed    int    `json:"failed"`
		RawLog    string `json:"raw_log"`
		Tests     []struct {
			Name string `json:"name"`
			Pass *bool  `json:"pass"`
		} `json:"tests"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	clientID, err := ensureClient(input.Client)
	if err != nil {
		return err
	}
	simName := input.Simulator
	if simName == "" {
		simName = "xdc/rpc-compat"
	}
	simID, err := ensureSimulator(simName)
	if err != nil {
		return err
	}

	dateRun := input.Date
	version := input.Version
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO runs (client_id, simulator_id, date_run, total, passed, failed, version, raw_log_path) VALUES (?,?,?,?,?,?,?,?)",
		clientID, simID, dateRun, input.Total, input.Passed, input.Failed, version, input.RawLog,
	)
	if err != nil {
		return err
	}
	runID, _ := res.LastInsertId()

	for _, t := range input.Tests {
		if t.Pass != nil {
			_, err := tx.Exec("INSERT INTO tests (run_id, test_name, passed) VALUES (?,?,?)",
				runID, t.Name, boolToInt(*t.Pass))
			if err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("Imported run: %s/%s %d/%d (%d tests)", input.Client, simName, input.Passed, input.Total, len(input.Tests))
	return nil
}

func importProbeFromFile(path string) error {
	data, err := readJSONFile(path)
	if err != nil {
		return err
	}
	var input struct {
		Client      string `json:"client"`
		Version     string `json:"version"`
		Supported   int    `json:"supported"`
		Unsupported int    `json:"unsupported"`
		Total       int    `json:"total"`
		Modules     map[string]string `json:"modules"`
		MethodList []struct {
			Method      string `json:"method"`
			Supported   bool   `json:"supported"`
			SampleValue string `json:"sample_value,omitempty"`
			Error       string `json:"error,omitempty"`
		} `json:"methods"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("parse probe %s: %w", path, err)
	}
	clientID, err := ensureClient(input.Client)
	if err != nil {
		return err
	}
	modJSON, _ := json.Marshal(input.Modules)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO probes (client_id, version, supported, unsupported, total, modules, date_run) VALUES (?,?,?,?,?,?,?)",
		clientID, input.Version, input.Supported, input.Unsupported, input.Total, string(modJSON), "2026-07-30",
	)
	if err != nil {
		return err
	}
	probeID, _ := res.LastInsertId()

	for _, m := range input.MethodList {
		_, err := tx.Exec(
			"INSERT INTO probe_methods (probe_id, method, supported, sample_value, error) VALUES (?,?,?,?,?)",
			probeID, m.Method, boolToInt(m.Supported), m.SampleValue, m.Error,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("Imported probe: %s (%d/%d methods)", input.Client, input.Supported, input.Total)
	return nil
}

func importGapMatrixFromFile(path string) error {
	data, err := readJSONFile(path)
	if err != nil {
		return err
	}
	var input struct {
		ClientA      string `json:"client_a"`
		ClientB      string `json:"client_b"`
		VersionA     string `json:"version_a"`
		VersionB     string `json:"version_b"`
		ModulesA     map[string]string `json:"modules_a"`
		ModulesB     map[string]string `json:"modules_b"`
		TotalM       int    `json:"total_methods"`
		BothSupp     int    `json:"both_supported"`
		BothUnsupp   int    `json:"both_unsupported"`
		InANotB      int    `json:"in_a_not_b"`
		InBNotA      int    `json:"in_b_not_a"`
		DateCreate   string `json:"date_created"`
		Methods      []struct {
			Method     string `json:"method"`
			ASupported bool   `json:"a_supported"`
			BSupported bool   `json:"b_supported"`
			AError     string `json:"a_error,omitempty"`
			BError     string `json:"b_error,omitempty"`
		} `json:"matrix"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("parse gap matrix %s: %w", path, err)
	}
	caID, err := ensureClient(input.ClientA)
	if err != nil {
		return err
	}
	cbID, err := ensureClient(input.ClientB)
	if err != nil {
		return err
	}
	modAJSON, _ := json.Marshal(input.ModulesA)
	modBJSON, _ := json.Marshal(input.ModulesB)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO gap_matrices (client_a_id, client_b_id, version_a, version_b, modules_a, modules_b,
			total_methods, both_supported, both_unsupported, in_a_not_b, in_b_not_a, date_created)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
	`, caID, cbID, input.VersionA, input.VersionB, string(modAJSON), string(modBJSON),
		input.TotalM, input.BothSupp, input.BothUnsupp, input.InANotB, input.InBNotA, input.DateCreate)
	if err != nil {
		return err
	}
	matrixID, _ := res.LastInsertId()

	for _, m := range input.Methods {
		_, err := tx.Exec(
			"INSERT INTO gap_matrix_methods (matrix_id, method, a_supported, b_supported, a_error, b_error) VALUES (?,?,?,?,?,?)",
			matrixID, m.Method, boolToInt(m.ASupported), boolToInt(m.BSupported), m.AError, m.BError,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("Imported gap matrix: %s vs %s", input.ClientA, input.ClientB)
	return nil
}

func handleImport(w http.ResponseWriter, r *http.Request) {
	errorLedger := filepath.Join("..", "error_ledger")
	entries, err := os.ReadDir(errorLedger)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	var imported int
	var errors []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(errorLedger, e.Name())
		switch {
		case strings.HasPrefix(e.Name(), "gapmatrix-"):
			if err := importGapMatrixFromFile(path); err != nil {
				msg := fmt.Sprintf("gapmatrix %s: %v", e.Name(), err)
				log.Print(msg)
				errors = append(errors, msg)
			} else {
				imported++
			}
		case strings.HasPrefix(e.Name(), "probe-"):
			if err := importProbeFromFile(path); err != nil {
				msg := fmt.Sprintf("probe %s: %v", e.Name(), err)
				log.Print(msg)
				errors = append(errors, msg)
			} else {
				imported++
			}
		case strings.HasPrefix(e.Name(), "compare-"):
			if err := importComparisonFromFile(path); err != nil {
				msg := fmt.Sprintf("comparison %s: %v", e.Name(), err)
				log.Print(msg)
				errors = append(errors, msg)
			} else {
				imported++
			}
		default:
			if err := importRunFromFile(path); err != nil {
				msg := fmt.Sprintf("run %s: %v", e.Name(), err)
				log.Print(msg)
				errors = append(errors, msg)
			} else {
				imported++
			}
		}
	}
	writeJSON(w, 200, map[string]any{"imported": imported, "errors": errors})
}

func importComparisonFromFile(path string) error {
	data, err := readJSONFile(path)
	if err != nil {
		return err
	}
	var input struct {
		Comparison   string `json:"comparison"`
		Simulator    string `json:"simulator"`
		Date         string `json:"date"`
		BothPass     int    `json:"both_pass"`
		AOnly        int    `json:"a_only"`
		BOnly        int    `json:"b_only"`
		BothFail     int    `json:"both_fail"`
		TotalMatched int    `json:"total_matched"`
		ClientA      struct {
			Name   string `json:"name"`
			Total  int    `json:"total"`
			Passed int    `json:"passed"`
			Failed int    `json:"failed"`
		} `json:"client_a"`
		ClientB struct {
			Name   string `json:"name"`
			Total  int    `json:"total"`
			Passed int    `json:"passed"`
			Failed int    `json:"failed"`
		} `json:"client_b"`
		Tests []struct {
			Name   string `json:"name"`
			APass  *bool  `json:"a_pass"`
			BPass  *bool  `json:"b_pass"`
			Status string `json:"status"`
		} `json:"tests"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO comparisons (simulator, client_a_name, client_b_name, both_pass, a_only, b_only, both_fail, total_matched, date_compared)
		VALUES (?,?,?,?,?,?,?,?,?)
	`, input.Simulator, input.ClientA.Name, input.ClientB.Name, input.BothPass, input.AOnly, input.BOnly,
		input.BothFail, input.TotalMatched, input.Date)
	if err != nil {
		return err
	}
	compID, _ := res.LastInsertId()

	for _, t := range input.Tests {
		_, err := tx.Exec(
			"INSERT INTO comparison_tests (comparison_id, test_name, a_pass, b_pass, status) VALUES (?,?,?,?,?)",
			compID, t.Name, nullableBool(t.APass), nullableBool(t.BPass), t.Status,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("Imported comparison: %s vs %s", input.ClientA.Name, input.ClientB.Name)
	return nil
}

// --- Helpers ---

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullableBool(b *bool) interface{} {
	if b == nil {
		return nil
	}
	if *b {
		return 1
	}
	return 0
}
