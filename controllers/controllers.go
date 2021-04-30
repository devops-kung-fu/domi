package controllers

import (
	// "context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/devops-kung-fu/domi/integrations"
	"github.com/devops-kung-fu/domi/lib"
	"github.com/gin-gonic/gin"
	ghclient "github.com/google/go-github/v33/github"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// InstallationID is a temporary fix until https://github.com/go-playground/webhooks/pull/128 is merged.
var InstallationID int64 = 0

// CanYouHearMeNow - Responds to liveliness checks
func CanYouHearMeNow(c *gin.Context) {
	c.Status(200)
}

func getGitHubClient(githubProvider *integrations.GitHubProvider) (*ghclient.Client, error) {
	githubClient, err := githubProvider.GitHubAuthenticator("/domi")
	if err != nil {
		log.Println(errors.New("GitHub Provider Authentication Failed"))
		return nil, err
	}
	log.Println("GitHub Provider Authentication Succeeded")
	return githubClient, nil
}

func downloadRepo(fs lib.FileSystem, githubClient *ghclient.Client, c *gin.Context, owner string, repo string, sha string) (string, error) {
	archiveLink, _, err := githubClient.Repositories.GetArchiveLink(c, owner, repo, "zipball", &ghclient.RepositoryContentGetOptions{Ref: sha}, true)
	if err != nil {
		log.Println(err)
	}
	archiveURL := archiveLink.String()
	client := lib.NewHTTPClient(http.DefaultClient, archiveURL)
	domiID, err := client.DownloadFile(fs)
	if err != nil {
		log.Println(err)
		return "", err
	}
	unzipErr := lib.UnZip(fmt.Sprintf("/domi/%s.zip", domiID), fmt.Sprintf("/domi/%s", domiID))
	if unzipErr != nil {
		log.Println(unzipErr)
		return "", unzipErr
	}
	return domiID, nil
}

// Move this to lib/filesystem.go
func targetDiscovery(fs lib.FileSystem, domiID string) ([]string, error) {
	foundFiles, e := lib.FindFiles(fs, fmt.Sprintf("/domi/%s", domiID), ".*\\.(tf|yaml|yml)")
	if e != nil {
		log.Println(e)
		return nil, e
	}
	return foundFiles, nil
}

func updateCheckRun(githubClient *ghclient.Client, c *gin.Context, owner string, repo string, checkRunID int64, status string, conclusion string, completedAt *ghclient.Timestamp, title string, summary string) error {
	if conclusion != "" {
		_, _, checkError := githubClient.Checks.UpdateCheckRun(c, owner, repo, checkRunID, ghclient.UpdateCheckRunOptions{
			Name:        "domi - Policy-as-Code Enforcer",
			Status:      &status,
			Conclusion:  &conclusion,
			CompletedAt: completedAt,
			Output: &ghclient.CheckRunOutput{
				Title:   &title,
				Summary: &summary,
			},
		})
		if checkError != nil {
			return checkError
		}
	} else {
		_, _, checkError := githubClient.Checks.UpdateCheckRun(c, owner, repo, checkRunID, ghclient.UpdateCheckRunOptions{
			Name:        "domi - Policy-as-Code Enforcer",
			Status:      &status,
			CompletedAt: completedAt,
			Output: &ghclient.CheckRunOutput{
				Title:   &title,
				Summary: &summary,
			},
		})
		if checkError != nil {
			return checkError
		}
	}
	return nil
}

func downloadPolicyRepo(fs lib.FileSystem, githubClient *ghclient.Client, c *gin.Context) (string, error) {
	var policyRepo string
	if os.Getenv("POLICY_REPO") != "" {
		policyRepo = os.Getenv("POLICY_REPO")
	} else {
		policyRepo = "https://github.com/devops-kung-fu/domi-policies"
	}
	policyRepoRegex := regexp.MustCompile(`https://github.com/(?P<owner>[-_a-zA-Z0-9]*)/(?P<repo>[-_a-zA-Z0-9]*)`)
	policyRepoMatch := policyRepoRegex.FindAllStringSubmatch(policyRepo, -1)
	policyRepoOwner := policyRepoMatch[0][1]
	policyRepoRepo := policyRepoMatch[0][2]
	policyRepoID, policyRepoIDErr := downloadRepo(fs, githubClient, c, policyRepoOwner, policyRepoRepo, "")
	if policyRepoIDErr != nil {
		return "", policyRepoIDErr
	}
	return policyRepoID, nil
}

func cleanUp(fs lib.FileSystem, id string) error {
	idZipError := fs.Remove(fmt.Sprintf("/domi/%s.zip", id))
	if idZipError != nil {
		return idZipError
	}
	idDirError := fs.RemoveAll(fmt.Sprintf("/domi/%s", id))
	if idDirError != nil {
		return idDirError
	}
	return nil
}

func downloadAndDiscover(fs lib.FileSystem, githubClient *ghclient.Client, c *gin.Context, owner string, repo string, sha string) (string, []string) {
	domiID, downloadRepoErr := downloadRepo(fs, githubClient, c, owner, repo, sha)
	if downloadRepoErr != downloadRepoErr {
		log.Println(downloadRepoErr)
	}
	targets, targetsError := targetDiscovery(fs, domiID)
	if targetsError != nil {
		log.Println(targetsError)
	}
	return domiID, targets
}

// ReceiveGitHubWebHook - Receives and processes GitHub WebHook Events
func ReceiveGitHubWebHook(c *gin.Context) {
	fs := lib.OSFS{}
	// ctx := context.Background()
	githubProvider, err := integrations.NewGitHubProvider()
	if err != nil {
		http.Error(c.Writer, "Could not get a provider.", 500)
		return
	}
	var hook *github.Webhook
	if githubProvider.GithubWebhookSecret != "" {
		hook, _ = github.New(github.Options.Secret(githubProvider.GithubWebhookSecret))
	} else {
		hook, _ = github.New()
	}
	payload, err := hook.Parse(c.Request, github.PushEvent, github.CheckRunEvent)
	if err != nil && err == github.ErrEventNotFound {
		c.String(http.StatusNotImplemented, "This event has not been implemented.")
	}

	switch payload.(type) {
	case github.PushPayload:
		push := payload.(github.PushPayload)
		githubProvider.InstallationID = int64(push.Installation.ID)
		InstallationID = int64(push.Installation.ID)
		owner := push.Repository.Owner.Login
		repo := push.Repository.Name
		sha := push.After
		githubClient, githubClientError := getGitHubClient(githubProvider)
		if githubClientError != nil {
			log.Println(githubClientError)
		}
		domiID, targets := downloadAndDiscover(fs, githubClient, c, owner, repo, sha)
		status := "queued"
		title := "domi - Policy-as-Code Enforcer"
		summary := "**Status**: Queued"
		if len(targets) > 0 {
			_, _, checkError := githubClient.Checks.CreateCheckRun(c, owner, repo, ghclient.CreateCheckRunOptions{
				Name:    "domi - Policy-as-Code Enforcer",
				HeadSHA: sha,
				Status:  &status,
				StartedAt: &ghclient.Timestamp{
					Time: time.Now(),
				},
				Output: &ghclient.CheckRunOutput{
					Title:   &title,
					Summary: &summary,
				},
			})
			if checkError != nil {
				log.Println(checkError)
			}
		}
		cleanUpError := cleanUp(fs, domiID)
		if cleanUpError != nil {
			log.Println(cleanUpError)
		}
		c.String(http.StatusOK, "Push Payload Received")
	case github.CheckRunPayload:
		check := payload.(github.CheckRunPayload)
		if check.Action == "created" && check.CheckRun.App.ID == int64(githubProvider.AppID) {
			githubProvider.InstallationID = InstallationID
			owner := check.Repository.Owner.Login
			repo := check.Repository.Name
			checkRunID := check.CheckRun.ID
			sha := check.CheckRun.HeadSHA
			githubClient, githubClientError := getGitHubClient(githubProvider)
			if githubClientError != nil {
				log.Println(githubClientError)
			}
			title := "domi - Policy-as-Code Enforcer"
			summary := "**Status**: Scanning"
			inProgressCheckError := updateCheckRun(githubClient, c, owner, repo, checkRunID, "in_progress", "", nil, title, summary)
			if inProgressCheckError != nil {
				log.Println(inProgressCheckError)
			}
			domiID, targets := downloadAndDiscover(fs, githubClient, c, owner, repo, sha)
			policyRepoID, policyRepoIDErr := downloadPolicyRepo(fs, githubClient, c)
			if policyRepoIDErr != nil {
				log.Println(policyRepoIDErr)
			}
			scanResults := lib.Scan(fs, "/domi", policyRepoID, targets)
			scanSummary, scanConclusion := lib.SummaryBuilder(scanResults)
			completedCheckError := updateCheckRun(githubClient, c, owner, repo, checkRunID, "completed", scanConclusion, &ghclient.Timestamp{Time: time.Now()}, title, scanSummary)
			if completedCheckError != nil {
				log.Println(completedCheckError)
			}
			cleanUpDomiIDError := cleanUp(fs, domiID)
			if cleanUpDomiIDError != nil {
				log.Println(cleanUpDomiIDError)
			}
			cleanUpPolicyIDError := cleanUp(fs, policyRepoID)
			if cleanUpPolicyIDError != nil {
				log.Println(cleanUpPolicyIDError)
			}
		}
		c.String(http.StatusOK, "Check Run Payload Received")
	default:
		c.String(http.StatusNotImplemented, "Event Not Implemented")
	}
}
