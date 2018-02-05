package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

// this file contains all the verification checking functions. Look in repo.go for
// where all these private functions are wired up. These functions depend on
// preflight() which sets ghClient, owner and repo

func checkAll(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	r, _, err := ghClient.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Error: %s\n", err.Error())
		if strings.Contains(err.Error(), "404 Not Found") {
			fmt.Fprintf(c.App.Writer, "  !!! Is this a private repo? Make sure your GH Token has the full Repo scope !!! \n")
		}
		return err
	} else if r.FullName == nil {
		fmt.Fprintf(c.App.Writer, "Error: [%s] does not have a FullName entry\n", repo)
		return errors.New("No repo fullname")
	}

	funcs := []cli.ActionFunc{
		checkTopic,
		checkLabels,
		checkUnassigned,
		checkUnlabled,
		checkProjects,
	}

	for _, f := range funcs {
		if err := f(c); err != nil {
			return err
		}
	}
	return nil
}

// checkTopic ensures "product-delivery" topic is assigned to the repo
func checkTopic(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	outInfo, outError := getWriters(c)

	fmt.Fprintf(outInfo, "Checking for [product-delivery] topic\n")
	ctx := context.Background()
	topics, _, err := ghClient.Repositories.ListAllTopics(ctx, owner, repo)
	if err != nil {
		fmt.Fprintf(outError, " - Error: %s\n", err.Error())
		return err
	}

	for _, name := range topics {
		if name == "product-delivery" {
			fmt.Fprintf(outInfo, " - OK. Found product-delivery topic\n")
			return nil
		}
	}

	fmt.Fprintf(outError, " - Error: product-delivery topic not set\n")
	return errors.New("product-delivery topic not set")
}

func checkLabels(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	outInfo, outError := getWriters(c)

	fmt.Fprintf(outInfo, "Checking Labels\n")
	labels, _, err := ghClient.Issues.ListLabels(context.Background(), owner, repo, nil)

	if err != nil {
		fmt.Fprintf(outError, " - Error: %s\n", err.Error())
		return err
	}

	// expected labels and their color
	standardLabels := map[string]string{
		"bug":             "b60205",
		"security":        "b60205",
		"improvement":     "0e8a16",
		"documentation":   "0e8a16",
		"fix":             "0e8a16",
		"new-feature":     "0e8a16",
		"P1":              "ffa32c",
		"P2":              "ffa32c",
		"P3":              "ffa32c",
		"P5":              "ffa32c",
		"proposal":        "1d76db",
		"question":        "1d76db",
		"support-request": "1d76db",
	}

	for _, label := range labels {
		if label == nil || label.Name == nil || label.Color == nil {
			continue
		}

		name := label.GetName()
		color := label.GetColor()

		if expectedColor, ok := standardLabels[name]; !ok {
			// not a standard label
			if color != "5319e7" {
				fmt.Fprintf(outError, " - Error: [%s] should have color #5319e7\n", name)
			} else {
				fmt.Fprintf(outInfo, " - OK. [%s] verified\n", name)
			}
		} else {
			// check standard label has correct color
			if color != expectedColor {
				fmt.Fprintf(outError, " - Error: standard label [%s] should have color #%s\n", name, expectedColor)
			} else {
				fmt.Fprintf(outInfo, " - OK. [%s] verified\n", name)
			}

			// delete it so we know how many are missing
			delete(standardLabels, name)
		}
	}

	// check for missing standard labels
	for missing, color := range standardLabels {
		fmt.Fprintf(outError, " - Error: missing %s (%s)\n", missing, color)
	}

	return nil
}

func checkUnassigned(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	outInfo, outError := getWriters(c)
	fmt.Fprintf(outInfo, "Checking Unassigned Issues\n")

	query := fmt.Sprintf("repo:%s/%s is:open no:assignee label:P1", owner, repo)
	results, _, err := ghClient.Search.Issues(context.Background(), query, nil)
	if err != nil {
		fmt.Fprintf(outError, " - Error: %s\n", err.Error())
		return err
	}

	count := results.GetTotal()
	if count > 0 {
		fmt.Fprintf(outError, " - Error: %d unassigned P1 issues\n", count)
		for _, issue := range results.Issues {
			fmt.Fprintf(outError, "  #%-4d %s", issue.GetNumber(), issue.GetTitle())
		}
	} else {
		fmt.Fprintf(outInfo, " - OK. All P1 issues assigned\n")
	}
	return nil
}
func checkUnlabled(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	outInfo, outError := getWriters(c)
	fmt.Fprintf(outInfo, "Checking Unlabled\n")
	query := fmt.Sprintf("repo:%s/%s is:open no:label is:issue", owner, repo)
	results, _, err := ghClient.Search.Issues(context.Background(), query, nil)
	if err != nil {
		fmt.Fprintf(outError, " - Error: %s\n", err.Error())
		return err
	}

	unassigned := results.GetTotal()
	if unassigned > 0 {
		fmt.Fprintf(outError, " - Error: %d issues unlabeled\n", unassigned)
		for _, issue := range results.Issues {
			fmt.Fprintf(outError, "   #%-4d %s\n", issue.GetNumber(), issue.GetTitle())
		}
	} else {
		fmt.Fprintf(outInfo, " - OK. All issues are labeled\n")
	}

	return nil
}

func checkProjects(c *cli.Context) error {
	owner, repo, err := extractOwnerRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	outInfo, outError := getWriters(c)
	fmt.Fprintf(outInfo, "Checking Projects\n")
	ctx := context.Background()

	projects, _, err := ghClient.Repositories.ListProjects(ctx, owner, repo, nil)
	if err != nil {
		fmt.Fprintf(outError, " - Error: Fetching projects, %s\n", err.Error())
		return err
	}

	projError := false
	for _, project := range projects {

		if project.GetBody() == "" {
			fmt.Fprintf(outError, " - WARNING: Project [%s] *SHOULD* have a description\n", project.GetName())
		}

		// check the project's columns
		pCols, _, err := ghClient.Projects.ListProjectColumns(ctx, project.GetID(), nil)
		if err != nil {
			fmt.Fprintf(outError, " - Error: Fetching project columns, %s\n", err.Error())
			return err
		}

		flags := 0
		for _, col := range pCols {
			switch col.GetName() {
			case "Backlog":
				flags |= 0x01
			case "In Progress":
				flags |= 0x02
			case "Blocked":
				flags |= 0x04
			case "Completed":
				flags |= 0x08
			default:
				fmt.Fprintf(outError, " - Error: Project [%s] has unexpected column %s\n",
					project.GetName(), col.GetName())
				projError = true
			}
		}

		if flags&0x01 == 0 {
			fmt.Fprintf(outError, " - Error: Project [%s] missing [Backlog] column\n", project.GetName())
			projError = true
		}
		if flags&0x02 == 0 {
			fmt.Fprintf(outError, " - Error: Project [%s] missing [In Progress] column\n", project.GetName())
			projError = true
		}
		if flags&0x04 == 0 {
			fmt.Fprintf(outError, " - Error: Project [%s] missing [Blocked] column\n", project.GetName())
			projError = true
		}
		if flags&0x08 == 0 {
			fmt.Fprintf(outError, " - Error: Project [%s] missing [Completed] column\n", project.GetName())
			projError = true
		}
	}

	if !projError {
		fmt.Fprintf(outInfo, " - OK. Projects verified\n")
	}

	return nil
}
