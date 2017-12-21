package repo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"

	"gopkg.in/urfave/cli.v1"
)

func initRepo(c *cli.Context) error {
	if err := initLabels(c); err != nil {
		return err
	}

	if err := createMilestone(c); err != nil {
		return err
	}

	return nil
}

// initRepo ensures the necessary labels, Milestone + Project exists
func initLabels(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, "Initializing labels on %s/%s\n", owner, repo)

	// make Issue labels
	labels := map[string]string{
		"P1":              "ffa32c",
		"P2":              "ffa32c",
		"P3":              "ffa32c",
		"P5":              "ffa32c",
		"bug":             "b60205",
		"security":        "b60205",
		"documentation":   "0e8a16",
		"fix":             "0e8a16",
		"new-feature":     "0e8a16",
		"question":        "1d76db",
		"propsal":         "1d76db",
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
				fmt.Fprintf(c.App.Writer, " - Error: Creating label %s\n", name)
				continue
			} else {
				fmt.Fprintf(c.App.Writer, " - Created label %s\n", name)
			}
		} else if *label.Color != color {
			_, _, err := ghClient.Issues.EditLabel(ctx, owner, repo, name, &github.Label{Color: &color})
			if err != nil {
				fmt.Fprintf(c.App.Writer, " - Error: changing color for %s\n", name)
			} else {
				fmt.Fprintf(c.App.Writer, " - Changed color for label %s to %s\n", name, color)
			}
		}
	}

	return nil
}

func createMilestone(c *cli.Context) error {
	milestoneTitle := c.String("milestone")
	fmt.Fprintf(c.App.Writer, "Creating Milestone/Project: %s\n", milestoneTitle)
	ctx := context.Background()
	milestones, _, err := ghClient.Issues.ListMilestones(ctx, owner, repo, nil)
	if err != nil {
		fmt.Fprintf(c.App.Writer, " - Error: Could not list milestones\n")
		return err
	}

	// short circuit if milestone already exists
	for _, m := range milestones {
		if *m.Title == milestoneTitle {
			fmt.Fprintf(c.App.Writer, " - Error: Milestone %s already exists\n", milestoneTitle)
			return nil
		}
	}

	milestone := &github.Milestone{Title: &milestoneTitle}
	if _, _, err := ghClient.Issues.CreateMilestone(ctx, owner, repo, milestone); err != nil {
		fmt.Fprintf(c.App.Writer, " - Error: Creating milestone %s: %s\n", milestoneTitle, err.Error())
		return err
	}

	// Make Project w/ Backlog, In Progress, Blocked and Completed
	// ... note: it will always create a new Project 1.0. Your repo will just have
	// multiple ones if you run this several times..
	fmt.Fprintf(c.App.Writer, " - Creating milestone %s\n", milestoneTitle)
	project := &github.ProjectOptions{Name: milestoneTitle}
	if proj, _, err := ghClient.Repositories.CreateProject(ctx, owner, repo, project); err != nil {
		fmt.Fprintf(c.App.Writer, " - Error: Creating Project %s: %s\n", milestoneTitle, err.Error())
		return err
	} else {
		fmt.Fprintf(c.App.Writer, " - Created project %s\n", milestoneTitle)
		for _, colName := range []string{"Backlog", "In Progress", "Blocked", "Complete"} {
			ops := &github.ProjectColumnOptions{Name: colName}
			if _, _, err := ghClient.Projects.CreateProjectColumn(ctx, *proj.ID, ops); err != nil {
				fmt.Fprintf(c.App.Writer, " - Error: creating Project Column [%s]: %s\n", colName, err.Error())
			}
		}
	}

	return nil
}
