package repo

// init.go initializes package level variables and other things shared by the packages
// libs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/urfave/cli.v1"
)

var (
	errMissingFlag = errors.New("Missing flag")

	// these are set in preflight()
	ghtoken  string
	ghClient *github.Client
)

// preflight ensures necessary flags and sets the package vars: ghClient, owner and repo
func preflight(c *cli.Context) (err error) {
	// set package vars
	ghtoken = c.String("ghtoken")

	if ghtoken == "" {
		fmt.Fprintf(c.App.Writer, "Error: github access token required\n")
		err = errMissingFlag
	}

	if err != nil {
		return
	}

	ctx := context.Background()

	// init ghClient if it's not already there
	// useful for checkAll where it calls all the sub-commands so they
	// don't have to recreate the ghClient
	if ghClient == nil {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghtoken})
		tc := oauth2.NewClient(ctx, ts)
		ghClient = github.NewClient(tc)
	}
	return
}

// extractOwnerRepo will find the owner and the repo information in the github
// url and return it as two strings
func extractOwnerRepo(arg string) (string, string, error) {
	parts := strings.Split(arg, "/")
	l := len(parts)
	if l < 2 {
		return "", "", errors.New("Invalid repo path, expect: owner/reponame")
	}

	return parts[l-2], parts[l-1], nil
}
