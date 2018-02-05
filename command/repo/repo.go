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
				Usage: "initializes labels and a project for a repo",
				Flags: []cli.Flag{cli.StringFlag{
					Name:  "project, p",
					Value: "Version 1.0",
					Usage: "initial project to create",
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
				Name:  "create-project",
				Usage: "creates a project with standard columns",
				Flags: []cli.Flag{cli.StringFlag{
					Name: "p, project",
				}},
				Action: createProject,
			},
			cli.Command{
				Name:  "check",
				Usage: "checks to verify standards conformity",
				Subcommands: cli.Commands{
					cli.Command{
						Name:   "all",
						Usage:  "Runs all checks on a repository",
						Action: checkAll,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "quiet, q",
								Usage: "only show errors",
							},
						},
					},
					cli.Command{
						Name:   "topic",
						Usage:  "Checks `product-delivery` topic is set",
						Action: checkTopic,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "quiet, q",
								Usage: "only show errors",
							},
						},
					},
					cli.Command{
						Name:   "labels",
						Usage:  "Checks Product Delivery standard labels are set",
						Action: checkLabels,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "quiet, q",
								Usage: "only show errors",
							},
						},
					},
					cli.Command{
						Name:   "unassigned",
						Usage:  "verify P1 issues are assigned to somebody",
						Action: checkUnassigned,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "quiet, q",
								Usage: "only show errors",
							},
						},
					},
					cli.Command{
						Name:   "unlabled",
						Usage:  "finds issues that do not have a label",
						Action: checkUnlabled,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "quiet, q",
								Usage: "only show errors",
							},
						},
					},
					cli.Command{
						Name:   "projects",
						Usage:  "verify projects fit standards",
						Action: checkProjects,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "quiet, q",
								Usage: "only show errors",
							},
						},
					},
				},
			},
		},
	}

}
