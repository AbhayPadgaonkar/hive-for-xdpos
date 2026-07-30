package main

type ClientInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SimulatorInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Run struct {
	ID          int64  `json:"id"`
	ClientID    int64  `json:"client_id"`
	ClientName  string `json:"client_name,omitempty"`
	SimulatorID int64  `json:"simulator_id"`
	SimName     string `json:"sim_name,omitempty"`
	DateRun     string `json:"date_run"`
	Total       int    `json:"total"`
	Passed      int    `json:"passed"`
	Failed      int    `json:"failed"`
	Version     string `json:"version,omitempty"`
	RawLogPath  string `json:"raw_log_path,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type TestResult struct {
	ID       int64  `json:"id"`
	RunID    int64  `json:"run_id"`
	TestName string `json:"test_name"`
	Passed   bool   `json:"passed"`
}

type Probe struct {
	ID          int64  `json:"id"`
	ClientID    int64  `json:"client_id"`
	ClientName  string `json:"client_name,omitempty"`
	Version     string `json:"version,omitempty"`
	Supported   int    `json:"supported"`
	Unsupported int    `json:"unsupported"`
	Total       int    `json:"total"`
	Modules     string `json:"modules,omitempty"`
	DateRun     string `json:"date_run"`
	CreatedAt   string `json:"created_at"`
}

type ProbeMethod struct {
	ID          int64  `json:"id"`
	ProbeID     int64  `json:"probe_id"`
	Method      string `json:"method"`
	Supported   bool   `json:"supported"`
	SampleValue string `json:"sample_value,omitempty"`
	Error       string `json:"error,omitempty"`
}

type GapMatrix struct {
	ID             int64  `json:"id"`
	ClientAID      int64  `json:"client_a_id"`
	ClientAName    string `json:"client_a_name,omitempty"`
	ClientBID      int64  `json:"client_b_id"`
	ClientBName    string `json:"client_b_name,omitempty"`
	VersionA       string `json:"version_a,omitempty"`
	VersionB       string `json:"version_b,omitempty"`
	ModulesA       string `json:"modules_a,omitempty"`
	ModulesB       string `json:"modules_b,omitempty"`
	TotalMethods   int    `json:"total_methods"`
	BothSupported  int    `json:"both_supported"`
	BothUnsupp     int    `json:"both_unsupported"`
	InANotB        int    `json:"in_a_not_b"`
	InBNotA        int    `json:"in_b_not_a"`
	DateCreated    string `json:"date_created"`
}

type GapMatrixMethod struct {
	ID          int64  `json:"id"`
	MatrixID    int64  `json:"matrix_id"`
	Method      string `json:"method"`
	ASupported  bool   `json:"a_supported"`
	BSupported  bool   `json:"b_supported"`
	AError      string `json:"a_error,omitempty"`
	BError      string `json:"b_error,omitempty"`
}

type Comparison struct {
	ID           int64  `json:"id"`
	Simulator    string `json:"simulator"`
	ClientAName  string `json:"client_a_name"`
	ClientBName  string `json:"client_b_name"`
	BothPass     int    `json:"both_pass"`
	AOnly        int    `json:"a_only"`
	BOnly        int    `json:"b_only"`
	BothFail     int    `json:"both_fail"`
	TotalMatched int    `json:"total_matched"`
	DateCompared string `json:"date_compared"`
	CreatedAt    string `json:"created_at"`
}

type ComparisonTest struct {
	ID           int64  `json:"id"`
	ComparisonID int64  `json:"comparison_id"`
	TestName     string `json:"test_name"`
	APass        *bool  `json:"a_pass"`
	BPass        *bool  `json:"b_pass"`
	Status       string `json:"status"`
}

type Stats struct {
	TotalRuns         int `json:"total_runs"`
	TotalProbes       int `json:"total_probes"`
	TotalComparisons  int `json:"total_comparisons"`
	TotalGapMatrices  int `json:"total_gap_matrices"`
	TotalTests        int `json:"total_tests"`
	XdcGethPassRate   float64 `json:"xdc_geth_pass_rate"`
	GoEthPassRate     float64 `json:"go_eth_pass_rate"`
}
