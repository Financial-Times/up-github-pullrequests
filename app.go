package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func sPtr(s string) *string { return &s }
func main() {

	if len(os.Args) < 2 {
		fmt.Print("No key specified")
		return
	}
	accessToken := os.Args[1]

	gistOutput := ""
	prCount := 0

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	opt := &github.ListOptions{PerPage: 10}
	var allRepos []*github.Repository

	for {
		repos, resp, _ := client.Organizations.ListTeamRepos(1544114, opt)

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, r := range allRepos {
		opt := &github.PullRequestListOptions{
			ListOptions: github.ListOptions{PerPage: 10},
		}
		for {
			prs, resp, _ := client.PullRequests.List("Financial-Times", *r.Name, opt)

			for _, p := range prs {
				gistOutput += *r.Name + " - " + *p.Title + " (Created: " + p.CreatedAt.Format("2006-01-02") + ")\n"
				gistOutput += "\t" + *p.HTMLURL + "\n"
				prCount++
			}

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}
	ft, _, _ := client.Users.Get("Financial-Times")
	gistFile := &github.GistFile{
		Type:    sPtr("test"),
		Content: &gistOutput,
	}
	files := make(map[github.GistFilename]github.GistFile)
	files["Open Pull Requests"] = *gistFile

	gist := &github.Gist{
		Owner:       ft,
		Description: sPtr("List of all open pull requests for " + time.Now().Format("2 January 2006")),
		Files:       files,
	}

	gistOut, _, err := client.Gists.Create(gist)
	if err != nil {
		panic(err)
	}

	fmt.Printf("There are <"+*gistOut.HTMLURL+"|%d pull requests> open in Universal Publishing repositories", prCount)

}
