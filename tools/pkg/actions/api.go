package actions

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v69/github"
)

type DefaultClient struct {
	repo   string
	owner  string
	client *github.Client
}

func NewDefaultClient() *DefaultClient {
	dc := &DefaultClient{
		client: github.NewClient(nil).WithAuthToken(Token()),
	}
	dc.owner, dc.repo = Repository()
	return dc
}

func (d *DefaultClient) HasBranch(branchName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	branch, resp, err := d.client.Repositories.GetBranch(
		ctx, d.owner, d.repo, branchName, 0,
	)
	exists := err == nil &&
		branch != nil &&
		resp.StatusCode == http.StatusOK
	return exists, err
}

func (d *DefaultClient) IsAssociatedWithPullRequest(sha string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	pulls, _, err := d.client.PullRequests.ListPullRequestsWithCommit(
		ctx, d.owner, d.repo, sha, &github.ListOptions{},
	)

	return len(pulls) > 0 &&
		pulls[0].GetMerged(), err
}
