package controllers

import (
	// "context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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

func getGitHubClient(githubProvider *lib.GitHubProvider) (*ghclient.Client, error) {
	githubClient, err := githubProvider.GitHubAuthenticator()
	if err != nil {
		log.Println(errors.New("GitHub Provider Authentication Failed"))
		return nil, err
	}
	log.Println("GitHub Provider Authentication Succeeded")
	return githubClient, nil
}

func targetDiscovery(githubClient *ghclient.Client, c *gin.Context, owner string, repo string, sha string) ([]string, error) {
	archiveLink, _, err := githubClient.Repositories.GetArchiveLink(c, owner, repo, "zipball", &ghclient.RepositoryContentGetOptions{Ref: sha}, true)
	if err != nil {
		log.Println(err)
	}
	archiveURL := archiveLink.String()
	fs := lib.OSFS{}
	domiID, err := lib.DownloadFile(fs, archiveURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	unzipErr := lib.UnZip(fmt.Sprintf("/tmp/%s.zip", domiID), fmt.Sprintf("/tmp/%s", domiID))
	if unzipErr != nil {
		log.Println(unzipErr)
		return nil, unzipErr
	}
	foundFiles, e := lib.FindFiles(fs, fmt.Sprintf("/tmp/%s", domiID), ".*\\.(tf|yaml|yml)")
	if e != nil {
		log.Println(e)
		return nil, e
	}
	return foundFiles, nil
}

func updateCheckRun(githubClient *ghclient.Client, c *gin.Context, owner string, repo string, checkRunID int64, status string, conclusion string, completedAt *ghclient.Timestamp, title string, summary string, text string) error {
	if conclusion != "" {
		_, _, checkError := githubClient.Checks.UpdateCheckRun(c, owner, repo, checkRunID, ghclient.UpdateCheckRunOptions{
			Name:        "domi - Policy-as-Code Enforcer",
			Status:      &status,
			Conclusion:  &conclusion,
			CompletedAt: completedAt,
			Output: &ghclient.CheckRunOutput{
				Title:   &title,
				Summary: &summary,
				Text:    &text,
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
				Text:    &text,
			},
		})
		if checkError != nil {
			return checkError
		}
	}
	return nil
}

// ReceiveGitHubWebHook - Receives and processes GitHub WebHook Events
func ReceiveGitHubWebHook(c *gin.Context) {
	// ctx := context.Background()
	githubProvider, err := lib.NewGitHubProvider()
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
	if err != nil {
		if err == github.ErrEventNotFound {
			c.String(http.StatusNotImplemented, "This event has not been implemented.")
		}
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
		targets, targetsError := targetDiscovery(githubClient, c, owner, repo, sha)
		if targetsError != nil {
			log.Println(targetsError)
		}
		status := "queued"
		title := "domi - Policy-as-Code Enforcer"
		summary := "Please stand by while we scan your repository... :thumbs-up:"
		text := "Something can go here"
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
					Text:    &text,
				},
			})
			if checkError != nil {
				log.Println(checkError)
			}
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
			summary := "Thanks for playing... :thumbs-up:"
			text := "Something can go here"
			inProgressCheckError := updateCheckRun(githubClient, c, owner, repo, checkRunID, "in_progress", "", nil, title, summary, text)
			if inProgressCheckError != nil {
				log.Println(inProgressCheckError)
			}
			_, targetsError := targetDiscovery(githubClient, c, owner, repo, sha)
			if targetsError != nil {
				log.Println(targetsError)
			}
			// *** Conftest Goes Here ***
			completedCheckError := updateCheckRun(githubClient, c, owner, repo, checkRunID, "completed", "neutral", &ghclient.Timestamp{Time: time.Now()}, title, summary, text)
			if completedCheckError != nil {
				log.Println(completedCheckError)
			}
		}
		c.String(http.StatusOK, "Check Run Payload Received")
	default:
		c.String(http.StatusNotImplemented, "Event Not Implemented")
	}
}
