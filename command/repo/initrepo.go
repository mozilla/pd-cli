package repo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"

	"gopkg.in/urfave/cli.v1"
)

func initRepo(c *cli.Context) error {
	if err := initLabels(c); err != nil {
		return err
	}

	if err := createProject(c); err != nil {
		return err
	}

	return nil
}

// initRepo ensures the necessary labels, Milestone + Project exists
func initLabels(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	outInfo, outError := getWriters(c)

	fmt.Fprintf(outInfo, "Initializing labels on %s/%s\n", owner, repo)

	// make Issue labels
	labels := map[string]string{
		"P1":              "ffa32c",
		"P2":              "ffa32c",
		"P3":              "ffa32c",
		"P5":              "ffa32c",
		"bug":             "b60205",
		"security":        "b60205",
		"improvement":     "0e8a16",
		"documentation":   "0e8a16",
		"fix":             "0e8a16",
		"new-feature":     "0e8a16",
		"question":        "1d76db",
		"proposal":        "1d76db",
		"support-request": "1d76db",
	}

	ctx := context.Background()
	for name, color := range labels {
		// check if the label exists
		_ = color
		label, resp, err := ghClient.Issues.GetLabel(ctx, owner, repo, name)
		if err != nil && resp.StatusCode == http.StatusNotFound {
			label := &github.Label{Name: &name, Color: &color}
			_, _, err := ghClient.Issues.CreateLabel(ctx, owner, repo, label)
			if err != nil {
				fmt.Fprintf(outError, " - Error: Creating label %s\n", name)
				continue
			} else {
				fmt.Fprintf(outError, " - Created label %s\n", name)
			}
		} else if *label.Color != color {
			_, _, err := ghClient.Issues.EditLabel(ctx, owner, repo, name, &github.Label{Color: &color})
			if err != nil {
				fmt.Fprintf(outError, " - Error: changing color for %s\n", name)
			} else {
				fmt.Fprintf(outInfo, " - Changed color for label %s to %s\n", name, color)
			}
		}
	}

	return nil
}

func createMilestone(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	outInfo, outError := getWriters(c)

	milestoneTitle := c.String("milestone")
	fmt.Fprintf(outInfo, "Creating Milestone: %s\n", milestoneTitle)
	ctx := context.Background()
	milestones, _, err := ghClient.Issues.ListMilestones(ctx, owner, repo, nil)
	if err != nil {
		fmt.Fprintf(outError, " - Error: Could not list milestones\n")
		return err
	}

	// short circuit if milestone already exists
	for _, m := range milestones {
		if *m.Title == milestoneTitle {
			fmt.Fprintf(outError, " - Error: Milestone [%s] already exists\n", milestoneTitle)
			return errors.New("Duplicate Milestone Name")
		}
	}

	milestone := &github.Milestone{Title: &milestoneTitle}
	if _, _, err := ghClient.Issues.CreateMilestone(ctx, owner, repo, milestone); err != nil {
		fmt.Fprintf(outError, " - Error: Creating milestone %s: %s\n", milestoneTitle, err.Error())
		return err
	}

	return nil
}

func createProject(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	projectTitle := c.String("project")
	outInfo, outError := getWriters(c)
	ctx := context.Background()

	projects, _, err := ghClient.Repositories.ListProjects(ctx, owner, repo, nil)
	if err != nil {
		fmt.Fprintf(outError, " - Error: Fetching projects, %s\n", err.Error())
		return err
	}

	for _, project := range projects {
		if project.GetName() == projectTitle {
			fmt.Fprintf(outError, " - Error: Project [%s] already exists\n", projectTitle)
			return errors.New("Duplicate Project Name")
		}
	}

	fmt.Fprintf(outInfo, "Creating Project %s\n", projectTitle)
	project := &github.ProjectOptions{Name: projectTitle}
	if proj, _, err := ghClient.Repositories.CreateProject(ctx, owner, repo, project); err != nil {
		fmt.Fprintf(outError, " - Error: Creating Project %s: %s\n", projectTitle, err.Error())
		return err
	} else {
		fmt.Fprintf(outInfo, " - Created project %s\n", projectTitle)
		for _, colName := range []string{"Backlog", "In Progress", "Blocked", "Completed"} {
			ops := &github.ProjectColumnOptions{Name: colName}
			if _, _, err := ghClient.Projects.CreateProjectColumn(ctx, *proj.ID, ops); err != nil {
				fmt.Fprintf(outError, " - Error: Creating Project Column [%s]: %s\n", colName, err.Error())
			}
		}
	}

	return nil
}
