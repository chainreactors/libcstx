package cstx

type BatchResult struct {
	NewNodes     int `json:"new_nodes"`
	UpdatedNodes int `json:"updated_nodes"`
	Count        int `json:"count"`
}

type IngestResult struct {
	RecordsParsed int `json:"records_parsed"`
	NewNodes      int `json:"new_nodes"`
	UpdatedNodes  int `json:"updated_nodes"`
	NewEdges      int `json:"new_edges"`
	NodeCount     int `json:"node_count"`
	EdgeCount     int `json:"edge_count"`
}
