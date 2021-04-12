package github

type WorkflowRun struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	NodeID           string `json:"node_id"`
	HeadBranch       string `json:"head_branch"`
	HeadSha          string `json:"head_sha"`
	RunNumber        int    `json:"run_number"`
	Event            string `json:"event"`
	Status           string `json:"status"`
	Conclusion       string `json:"conclusion"`
	WorkflowID       int    `json:"workflow_id"`
	CheckSuiteID     int64  `json:"check_suite_id"`
	CheckSuiteNodeID string `json:"check_suite_node_id"`
}
