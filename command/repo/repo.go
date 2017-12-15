package repo

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"

	"gopkg.in/urfave/cli.v1"
)

var (
	errMissingFlag = errors.New("Missing flag")

	// these are set in preflight()
	ghtoken, repo, owner string
	ghClient             *github.Client
)

// NewRepoCommand creates the `repo` command tree
func NewCommand() cli.Command {
	return cli.Command{
		Name:  "repo",
		Usage: "tools for managing our github repos",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "ghtoken, g",
				Usage:  "github access token",
				EnvVar: "GH_ACCESS_TOKEN",
			},
			cli.StringFlag{
				Name:  "owner,o",
				Usage: "Owner of the repo",
			},
			cli.StringFlag{
				Name:  "repo, r",
				Usage: "Name of the repo",
			},
		},
		Before: preflight,
		Subcommands: cli.Commands{
			cli.Command{
				Name:  "init",
				Usage: "initializes labels and a milestone for a repo",
				Flags: []cli.Flag{cli.StringFlag{
					Name:  "milestone, m",
					Value: "Version 1.0",
					Usage: "initial milestone",
				}},
				Action: initRepo,
			},
			cli.Command{
				Name:   "init-labels",
				Usage:  "creates or updates standard labels",
				Action: initLabels,
			},
			cli.Command{
				Name:  "create-milestone",
				Usage: "creates a new milestone",
				Flags: []cli.Flag{cli.StringFlag{
					Name: "milestone, m",
				}},
				Action: createMilestone,
			},

			cli.Command{
				Name:  "check",
				Usage: "checks to verify standards conformity",
				Subcommands: cli.Commands{
					cli.Command{
						Name:   "all",
						Usage:  "Runs all checks on a repository",
						Action: checkAll,
					},
					cli.Command{
						Name:   "topic",
						Usage:  "Checks `product-delivery` topic is set",
						Action: checkTopic,
					},
					cli.Command{
						Name:   "labels",
						Usage:  "Checks `product-delivery` topic is set",
						Action: checkLabels,
					},
					cli.Command{
						Name:   "unassigned",
						Usage:  "verify P1 issues are assigned to somebody",
						Action: checkUnassigned,
					},
					cli.Command{
						Name:   "unlabled",
						Usage:  "finds issues that do not have a label",
						Action: checkUnlabled,
					},
					cli.Command{
						Name:   "milestones",
						Usage:  "verify milestones have a project to track them",
						Action: checkMilestones,
					},
				},
			},
		},
	}

}

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
		fmt.Fprintf(c.App.Writer, "Error: %s\n", err.Error())
		return err
	} else if r.FullName == nil {
		fmt.Fprintf(c.App.Writer, "Error: [%s] does not have a FullName entry\n", repo)
		return errors.New("No repo fullname")
	}

	return
}
