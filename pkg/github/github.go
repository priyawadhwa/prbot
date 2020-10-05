package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	v1 "github.com/priyawadhwa/prbot/pkg/config/v1"
	"golang.org/x/oauth2"
)

// Comment comments on the PR
func Comment(ctx context.Context, cfg *v1.Config, contents []byte) error {
	c := newClient(ctx, cfg.Github)
	fmt.Println(c)
	return nil
}

// client provides the context and client with necessary auth
// for interacting with the Github API
type client struct {
	ctx context.Context
	*github.Client
	owner string
	repo  string
}

// newClient returns a github client with the necessary auth
func newClient(ctx context.Context, cfg v1.Github) *client {
	githubToken := os.Getenv(cfg.AccessTokenEnvVar)
	// Setup the token for github authentication
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	// Return a client instance from github
	c := github.NewClient(tc)
	return &client{
		ctx:    ctx,
		Client: c,
		owner:  cfg.Owner,
		repo:   cfg.Repo,
	}
}

// CommentOnPR comments message on the PR
func (g *client) CommentOnPR(pr int, message string) error {
	comment := &github.IssueComment{
		Body: &message,
	}

	log.Printf("Creating comment on PR %d: %s", pr, message)
	_, _, err := g.Client.Issues.CreateComment(g.ctx, g.owner, g.repo, pr, comment)
	if err != nil {
		return errors.Wrap(err, "creating github comment")
	}
	log.Printf("Successfully commented on PR %d.", pr)
	return nil
}

// ListOpenPRsWithLabel returns all open PRs with the specified label
func (g *client) ListOpenPRsWithLabel(label string) ([]int, error) {
	validPrs := []int{}
	prs, _, err := g.Client.PullRequests.List(g.ctx, g.owner, g.repo, &github.PullRequestListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "listing pull requests")
	}
	for _, pr := range prs {
		if prContainsLabel(pr.Labels, label) {
			validPrs = append(validPrs, pr.GetNumber())
		}
	}
	return validPrs, nil
}

func prContainsLabel(labels []*github.Label, label string) bool {
	for _, l := range labels {
		if l == nil {
			continue
		}
		if l.GetName() == label {
			return true
		}
	}
	return false
}

// NewCommitsExist checks if new commits exist since minikube-pr-bot
// commented on the PR. If so, return true.
func (g *client) NewCommitsExist(pr int, login string) (bool, error) {
	lastCommentTime, err := g.timeOfLastComment(pr, login)
	if err != nil {
		return false, errors.Wrapf(err, "getting time of last comment by %s on pr %d", login, pr)
	}
	lastCommitTime, err := g.timeOfLastCommit(pr)
	if err != nil {
		return false, errors.Wrapf(err, "getting time of last commit on pr %d", pr)
	}
	return lastCommentTime.Before(lastCommitTime), nil
}

func (g *client) timeOfLastCommit(pr int) (time.Time, error) {
	var commits []*github.RepositoryCommit

	page := 0
	resultsPerPage := 30
	for {
		c, _, err := g.Client.PullRequests.ListCommits(g.ctx, g.owner, g.repo, pr, &github.ListOptions{
			Page:    page,
			PerPage: resultsPerPage,
		})
		if err != nil {
			return time.Time{}, err
		}
		commits = append(commits, c...)
		if len(c) < resultsPerPage {
			break
		}
		page++
	}

	lastCommitTime := time.Time{}
	for _, c := range commits {
		if newCommitTime := c.GetCommit().GetAuthor().GetDate(); newCommitTime.After(lastCommitTime) {
			lastCommitTime = newCommitTime
		}
	}
	return lastCommitTime, nil
}

func (g *client) timeOfLastComment(pr int, login string) (time.Time, error) {
	var comments []*github.IssueComment

	page := 0
	resultsPerPage := 30
	for {
		c, _, err := g.Client.Issues.ListComments(g.ctx, g.owner, g.repo, pr, &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: resultsPerPage,
			},
		})
		if err != nil {
			return time.Time{}, err
		}
		comments = append(comments, c...)
		if len(c) < resultsPerPage {
			break
		}
		page++
	}

	// go through comments backwards to find the most recent
	lastCommentTime := time.Time{}

	for _, c := range comments {
		if u := c.GetUser(); u != nil {
			if u.GetLogin() == login {
				if c.GetCreatedAt().After(lastCommentTime) {
					lastCommentTime = c.GetCreatedAt()
				}
			}
		}
	}

	return lastCommentTime, nil
}
