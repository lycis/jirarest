package jirarest

import (
	"fmt"
	"testing"
)

// test cases

func TestGetIssue(t *testing.T) {
	client := getTestClient()

	issue, err := client.GetIssue(TEST_EXISTING_ISSUE)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if issue.Key == "" {
		t.Error("empty issue key")
		t.Fail()
	} else {
		t.Logf("found issue %s\n", issue.Key)
	}
}

func TestSearchIssue(t *testing.T) {
	client := getTestClient()

	issues, err := client.SearchIssue(fmt.Sprintf("key = %s", TEST_EXISTING_ISSUE))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if issues.Total == 0 {
		t.Error("existing issue not found")
		t.Fail()
	} else {
		t.Logf("found %d issues", issues.Total)
	}
}

func TestCreateIssue(t *testing.T) {
	client := getTestClient()

	// create test issue
	var issue Issue
	issue.Fields.Description = "This is a test issue. Please delete."
	issue.Fields.Summary = "Test issue"

	project := new(Project)
	project.Key = TEST_PROJECT_KEY
	issue.Fields.Project = project

	itype := new(IssueType)
	itype.Name = "Bug"
	issue.Fields.IssueType = itype

	priority := new(Priority)
	priority.Name = "Major"
	issue.Fields.Priority = priority

	newIssue, err := client.CreateIssue(issue)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if newIssue.Key == "" {
		t.Error("created issue has no key")
		t.Fail()
	}

	t.Logf("created issue %s", issue.Key)
}

// test helper
func getTestClient() Client {
	var c Client
	c.Uri = TEST_URI
	c.User = TEST_USER
	c.Password = TEST_PASSWORD
	return c
}
