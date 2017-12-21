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
	ghtoken, repo, owner string
	ghClient             *github.Client
)

// preflight ensures necessary flags and sets the package vars: ghClient, owner and repo
func preflight(c *cli.Context) (err error) {
	// set package vars
	ghtoken = c.String("ghtoken")
	owner = c.String("owner")
	repo = c.String("repo")

	if ghtoken == "" {
		fmt.Fprintf(c.App.Writer, "Error: github access token required\n")
		err = errMissingFlag
	}

	if owner == "" {
		fmt.Fprintf(c.App.Writer, "Error: owner required\n")
		err = errMissingFlag
	}

	if repo == "" {
		fmt.Fprintf(c.App.Writer, "Error: repo name required\n")
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

	r, _, err := ghClient.Repositories.Get(ctx, owner, repo)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Preflight Error: %s\n", err.Error())
		if strings.Contains(err.Error(), "404 Not Found") {
			fmt.Fprintf(c.App.Writer, "  !!! Is this a private repo? Make sure your GH Token has the full Repo scope !!! \n")
		}
		return err
	} else if r.FullName == nil {
		fmt.Fprintf(c.App.Writer, "Preflight Error: [%s] does not have a FullName entry\n", repo)
		return errors.New("No repo fullname")
	}

	return
}
