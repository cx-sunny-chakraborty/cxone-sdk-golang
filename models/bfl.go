package models

// BFLResponseModel is the response of GET /api/bfl. BFL is "Best Fix
// Location" — given a SAST query and a scan, the platform returns the
// recommended single point in the data flow that, if patched, would
// remediate the largest number of related findings.
type BFLResponseModel struct {
	ID         string         `json:"id"`
	Trees      []BFLTreeModel `json:"trees"`
	TotalCount int            `json:"totalCount"`
}

// BFLTreeModel is one BFL tree returned in [BFLResponseModel.Trees].
type BFLTreeModel struct {
	ID                  string                     `json:"id"`
	BFL                 *ScanResultNode            `json:"bfl"`
	Results             []*ScanResultData          `json:"results"`
	Nodes               map[string]*ScanResultNode `json:"nodes,omitempty"`
	NodesAdjacencyPairs [][]string                 `json:"nodesAdjacencyPairs,omitempty"`
}
