package integrations

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v33/github"
)

// GitHubProvider - Contains the components necessary for GitHub
type GitHubProvider struct {
	GithubWebhookSecret   string
	githubToken           string
	githubPrivateKey      string
	AppID                 int
	InstallationID        int64
	oauthClient           *http.Client
	appClient             *http.Client
	githubClient          *github.Client
	installationTransport *ghinstallation.Transport
}

// GitHub interface
type GitHub interface {
	GitHubAuthenticator() (*github.Client, error)
}

// NewGitHubProvider - Starts a new GitHub Provider
func NewGitHubProvider() (*GitHubProvider, error) {
	gitHubProvider := GitHubProvider{}
	ghPrivateKey := os.Getenv("GITHUB_PRIVATE_KEY")
	ghAppID, _ := strconv.Atoi(os.Getenv("GITHUB_APP_ID"))
	gitHubProvider.GithubWebhookSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")
	if ghPrivateKey != "" {
		gitHubProvider.githubPrivateKey = ghPrivateKey
		gitHubProvider.AppID = ghAppID
	} else {
		return nil, errors.New("domi has no way to authenticate into GitHub")
	}
	return &gitHubProvider, nil
}

// GitHubAuthenticator - Authenticates Personal Access Token
func (githubProvider *GitHubProvider) GitHubAuthenticator(rootPath string) (*github.Client, error) {
	transport := http.DefaultTransport
	pemFile, err := os.Create(fmt.Sprintf("%s/private-key.pem", rootPath))
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error creating PEM file")
	}
	defer func() {
		err = pemFile.Close()
	}()
	bytesWritten, err := pemFile.WriteString(githubProvider.githubPrivateKey)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error writing to PEM file")
	}
	fmt.Println(bytesWritten, "bytes written successfully to PEM File")
	err = pemFile.Close()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Error closing PEM file")
	}
	itr, err := ghinstallation.NewKeyFromFile(transport, int64(githubProvider.AppID), int64(githubProvider.InstallationID), "/domi/private-key.pem")
	if err != nil {
		log.Println(err)
	}
	githubProvider.installationTransport = itr
	githubProvider.appClient = &http.Client{Transport: itr}
	githubProvider.githubClient = github.NewClient(githubProvider.appClient)
	return githubProvider.githubClient, nil
}
