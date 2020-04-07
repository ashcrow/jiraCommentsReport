package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift/gmarkley-VI/jiraSosRepot/functions"
	"log"
	"strings"
)

func main() {
	jiraURL := "https://issues.redhat.com"
	username, password := functions.ReadCredentials()

	var jiraJQL [1][2]string
	jiraJQL[0][0] = "project = WINC AND status in (\"In Progress\", \"Code Review\")AND(sprint in openSprints())"

	//Create the client
	client, _ := functions.CreatTheClient(username, password, jiraURL)

	//Loop over the jiraJQL array and Request the issue objects
	for z := 0; z < len(jiraJQL); z++ {

		var issues []jira.Issue

		// append the jira issues to []jira.Issue
		appendFunc := func(i jira.Issue) (err error) {
			issues = append(issues, i)
			return err
		}

		// SearchPages will page through results and pass each issue to appendFunc taken from the Jira Example implementation
		// In this example, we'll search for all the issues with the provided JQL filter and Print the header that goes with it.
		err := client.Issue.SearchPages(fmt.Sprintf(`%s`, jiraJQL[z][0]), nil, appendFunc)
		if err != nil {
			log.Fatal(err)
		}

		for _, i := range issues {
			options := &jira.GetQueryOptions{Expand: "renderedFields"}
			u, _, err := client.Issue.Get(i.Key, options)

			if err != nil {
				fmt.Printf("\n==> error: %v\n", err)
				return
			}

			if len(u.RenderedFields.Comments.Comments) >= 1 {
				c := u.RenderedFields.Comments.Comments[len(u.RenderedFields.Comments.Comments)-1]
				if strings.Contains(c.Updated, "days ago") {
					commentString := fmt.Sprintf("%s Please comment/update - Last update was %+v", i.Fields.Assignee.DisplayName, c.Updated)
					com := jira.Comment{
						ID:           i.ID,
						Self:         "",
						Name:         "",
						Author:       jira.User{},
						Body:         commentString,
						UpdateAuthor: jira.User{},
						Updated:      "",
						Created:      "",
						Visibility:   jira.CommentVisibility{},
					}
					commentOUT, _, err := client.Issue.AddComment(i.Key, &com)
					if err != nil {
						panic(err)
					}
					fmt.Printf("ID - %s \n Body - %+v\n", i.Key, commentOUT.Body)
				}
			}
		}
	}
}
