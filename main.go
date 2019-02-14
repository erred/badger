package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudbuild/v1"
)

func main() {
	client, err := google.DefaultClient(oauth2.NoContext, cloudbuild.CloudPlatformScope)
	if err != nil {
		log.Fatal("get default client: ", err)
	}
	svc, err := cloudbuild.New(client)
	if err != nil {
		log.Fatal("get cloudbuild service: ", err)
	}

	res, err := svc.Projects.Builds.List("com-seankhliao").Do()
	if err != nil {
		log.Fatal("project list: ", err)
	}
	for i, b := range res.Builds {
		fmt.Printf("%d: %s %s %s %s\n", i, b.ProjectId, b.Source.RepoSource.RepoName, b.Status, b.LogUrl)
	}

	time.Sleep(10 * time.Minute)
}
