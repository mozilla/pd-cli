package repo

import "gopkg.in/urfave/cli.v1"

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
