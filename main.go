package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/cloudbuild/v1"
)

func main() {
	for i, e := range os.Environ() {
		fmt.Println(i, e)
	}

	svc, err := cloudbuild.New(http.DefaultClient)
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
