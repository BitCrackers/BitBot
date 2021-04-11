package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Artifacts(owner, repo string) ([]Artifact, error) {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%v/%v/actions/artifacts", owner, repo)
	res, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("error while sending http request: %v", err)
	}
	var resData struct {
		Artifacts []Artifact `json:"artifacts"`
	}
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %v", err)
	}
	return resData.Artifacts, nil
}

func WorkflowRuns(owner, repo string) ([]WorkflowRun, error) {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%v/%v/actions/runs", owner, repo)
	res, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("error while sending http request: %v", err)
	}
	var resData struct {
		WorkflowRuns []WorkflowRun `json:"workflow_runs"`
	}
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return nil, fmt.Errorf("error while decoding response body: %v", err)
	}
	return resData.WorkflowRuns, nil
}
