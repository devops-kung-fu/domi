package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// ConftestResults - Holds conftest results
type ConftestResults []struct {
	Filename  string `json:"filename"`
	Successes int    `json:"successes,omitempty"`
	Failures  []struct {
		Msg      string `json:"msg"`
		Metadata struct {
			details struct{}
		} `json:"metadata,omitempty"`
	} `json:"failures,omitempty"`
	Warnings []struct {
		Msg      string `json:"msg"`
		Metadata struct {
			details struct{}
		} `json:"metadata,omitempty"`
	} `json:"warnings,omitempty"`
}

// Scan - Use conftest to Scan discovered files.
func Scan(fs FileSystem, rootPath string, policyID string, files []string) ConftestResults {
	policyIDPath := fmt.Sprintf("%s/%s", rootPath, policyID)
	policyPaths, policyPathsError := FindFiles(fs, policyIDPath, "policy$")
	if policyPathsError != nil {
		return nil
	}
	var policyPath string = ""
	for _, policyPathCandidate := range policyPaths {
		if len(strings.Split(policyPathCandidate, "/")) == 5 {
			policyPath = policyPathCandidate
			break
		}
	}
	cmd := exec.Command("conftest", "test", "--all-namespaces", "--fail-on-warn", "-o", "json", "-p", policyPath)
	cmd.Args = append(cmd.Args, files...)
	output, _ := cmd.Output()
	log.Println(string(output))
	conftestResults := ConftestResults{}
	IfErrorLog(json.Unmarshal(output, &conftestResults), "[ERROR]")
	return conftestResults
}
