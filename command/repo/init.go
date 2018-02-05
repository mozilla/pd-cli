package repo

// init.go initializes package level variables and other things shared by the packages
// libs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	_, outError := getWriters(c)

	// set package vars
	ghtoken = c.String("ghtoken")

	if ghtoken == "" {
		fmt.Fprintf(outError, "Error: github access token required\n")
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

	// check that we got the repo value, which should be the last thing in
	// the arg list
	args := c.Args()
	_, _, err = extractOwnerRepo(args[len(args)-1])
	if err != nil {
		fmt.Fprintf(outError, "ERROR: %s\n", err.Error())
		return err
	}
	return
}

// extractOwnerRepo will find the owner and the repo information in the github
// url and return it as two strings
func extractOwnerRepo(arg string) (string, string, error) {
	parts := strings.Split(arg, "/")
	l := len(parts)
	if l < 2 {
		return "", "", errors.New("Invalid or missing repo path, expect: owner/reponame")
	}

	return parts[l-2], parts[l-1], nil
}

func getWriters(c *cli.Context) (info, err io.Writer) {
	if c.Bool("quiet") {
		return ioutil.Discard, c.App.Writer
	} else {
		return c.App.Writer, c.App.Writer
	}
}
