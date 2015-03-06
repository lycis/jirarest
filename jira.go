// This library provides access to Atlassian Jira from within Go code.
// It relies on using the REST API to do so.
package jirarest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// This type will get you acces to the REST API.
type Client struct {
	Uri      string
	User     string
	Password string
	client   http.Client
}

// Representation object for Jira issues
type Issue struct {
	Id        string      `json:"id,omitempty"`
	Key       string      `json:"key,omitempty"`
	Self      string      `json:"self,omitempty"`
	Expand    string      `json:"expand,omitempty"`
	CreatedAt time.Time   `json:"createdat,omitempty"`
	Fields    IssueFields `json:"fields,omitempty"`
}

// Fields that are available in an issue
type IssueFields struct {
	IssueType   *IssueType `json:"issuetype,omitempty"`
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Reporter    *User      `json:"reporter,omitempty"`
	Assignee    *User      `json:"assignee,omitempty"`
	Project     *Project   `json:"project,omitempty"`
	Created     string     `json:"created,omitempty"`
	Creator     *User      `json:"creator,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
	Priority    *Priority  `json:"priority,omitempty"`
}

// type of an issue
type IssueType struct {
	Self        string `json:"self,omitempty"`
	Id          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	IconUrl     string `json:"iconurl,omitempty"`
	Name        string `json:"name,omitempty"`
	Subtask     bool   `json:"subtask,omitempty"`
}

// representation of a user
type User struct {
	Self         string            `json:"self,omitempty"`
	Name         string            `json:"name,omitempty"`
	EmailAddress string            `json:"emailaddress,omitempty"`
	AvatarUrls   map[string]string `json:"avatarurls,omitempty"`
	DisplayName  string            `json:"displayname,omitempty"`
	Active       bool              `json:"active,omitempty"`
}

// project within Jira
type Project struct {
	Self       string            `json:"self,omitempty"`
	Id         string            `json:"id,omitempty"`
	Key        string            `json:"key,omitempty"`
	Name       string            `json:"name,omitempty"`
	AvatarUrls map[string]string `json:"avatarurls,omitempty"`
}

// a list of issues as returned by e.g. search
type IssueList struct {
	Expand    string   `json:"expand,omitempty"`
	StartAt   int      `json:"startat,omitempty"`
	MaxReslts int      `json:"maxresults,omitempty"`
	Total     int      `json:"total,omitempty"`
	Issues    []*Issue `json:"issues,omitempty"`
}

type Priority struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	IconUrl string `json:"iconurl,omitempty"`
	Self    string `json:"self,omitempty"`
}

// Gives a specific issue identified by Jira-Key (e.g. PROJECT-1).
func (jc Client) GetIssue(key string) (Issue, error) {
	request, err := jc.buildRequest("GET", fmt.Sprintf("issue/%s", key), nil)
	if err != nil {
		return Issue{}, err
	}

	response, err := jc.client.Do(request)
	if err != nil {
		return Issue{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if ok, jErr := toJiraError(body); ok {
		return Issue{}, jErr
	}

	if err != nil {
		return Issue{}, err
	}

	var issue Issue
	err = json.Unmarshal(body, &issue)
	if err != nil {
		return Issue{}, err
	}

	return issue, nil
}

// Search issues by using the given JQL
func (jc Client) SearchIssue(jql string) (IssueList, error) {
	request, err := jc.buildRequest("GET",
		fmt.Sprintf("search?jql=%s", url.QueryEscape(jql)), nil)

	if err != nil {
		return IssueList{}, err
	}

	response, err := jc.client.Do(request)
	if err != nil {
		return IssueList{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if ok, jErr := toJiraError(body); ok {
		return IssueList{}, jErr
	}

	if err != nil {
		return IssueList{}, err
	}

	var issues IssueList
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return IssueList{}, err
	}

	return issues, nil
}

// Creates a Jira issue and returns rudimentary information about it.
// The returned Issue will only contain ID, key and self reference
func (jc Client) CreateIssue(issue Issue) (Issue, error) {
	data, err := json.Marshal(issue)
	if err != nil {
		return Issue{}, err
	}

	request, err := jc.buildRequest("POST", "issue", bytes.NewBuffer(data))
	if err != nil {
		return Issue{}, err
	}

	response, err := jc.client.Do(request)

	body, err := ioutil.ReadAll(response.Body)
	if ok, jErr := toJiraError(body); ok {
		return Issue{}, jErr
	}

	if err != nil {
		return Issue{}, err
	}

	var newIssue Issue
	err = json.Unmarshal(body, &newIssue)
	if err != nil {
		return Issue{}, err
	}

	return newIssue, nil
}

// --- INTERNAL FUNCTIONS ---

// internal: build request URI
func (jc Client) buildUri(path string) string {
	var fmtString string
	if !strings.HasSuffix(jc.Uri, "/") {
		fmtString = "%s/rest/api/2/%s"
	} else {
		fmtString = "%srest/api/2/%s"
	}

	return fmt.Sprintf(fmtString, jc.Uri, path)
}

// internal: build request with authentication
func (jc Client) buildRequest(method, path string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, jc.buildUri(path), body)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(jc.User, jc.Password)

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	return request, nil
}
