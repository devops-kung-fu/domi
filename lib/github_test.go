package lib

import (
	"os"
	"testing"
)

func TestNewGitHubProvider(t *testing.T) {
	os.Setenv("GITHUB_ACCESS_TOKEN", "okay")
	os.Setenv("GITHUB_APP_ID", "1")
	os.Setenv("GITHUB_WEBHOOK_SECRET", "okay")
	_, err1 := NewGitHubProvider()
	if err1 != nil {
		t.Fail()
	}
	os.Setenv("GITHUB_ACCESS_TOKEN", "")
	os.Setenv("GITHUB_PRIVATE_KEY", "okay")
	_, err2 := NewGitHubProvider()
	if err2 != nil {
		t.Fail()
	}
	os.Setenv("GITHUB_ACCESS_TOKEN", "")
	os.Setenv("GITHUB_PRIVATE_KEY", "")
}

func TestNewGitHubProviderEmptyEnvVar(t *testing.T) {
	_, err := NewGitHubProvider()
	if err != nil {
		t.Log()
	}
}