package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudbuild/v1"
)

func main() {
	p := flag.String("p", "8080", "port to listen on")
	pr := flag.String("pr", "com-seankhliao", "GCP project to query")
	flag.Parse()

	// Setup / get client
	client, err := google.DefaultClient(oauth2.NoContext, cloudbuild.CloudPlatformScope)
	if err != nil {
		log.Fatal("get default client: ", err)
	}
	svc, err := cloudbuild.New(client)
	if err != nil {
		log.Fatal("get cloudbuild service: ", err)
	}

	// handle reuquests
	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NewShieldFromBuild("SUCCESS"))
	})
	http.HandleFunc("/failure", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NewShieldFromBuild("failure"))
	})
	http.HandleFunc("/status_unknown", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NewShieldFromBuild("STATUS_UNKNOWN"))
	})
	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		repoName := strings.Split(r.URL.Path, "/")[2]

		status, err := status(svc, *pr, repoName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(NewShieldFromBuild(status))
	})
	http.ListenAndServe(":"+*p, nil)
}

func status(svc *cloudbuild.Service, p, r string) (string, error) {
	res, err := svc.Projects.Builds.
		List(p).
		Filter(`source.repo_source.repo_name = "` + r + `"`).
		Fields("builds.status").
		Do()
	if err != nil {
		return "", err
	}

	// filter qorking / queued / cancelled
	status := "STATUS_UNKNOWN"
	for _, b := range res.Builds {
		if b.Status == "WORKING" || b.Status == "QUEUED" || b.Status == "CANCELLED" {
			continue
		}
		status = b.Status
		break
	}
	return status, nil

}

type Shield struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Color         string `json:"color,omitempty"`
	LabelColor    string `json:"labelColor,omitempty"`
	IsError       bool   `json:"isError,omitempty"`
	NamedLogo     string `json:"namedLogo,omitempty"`
	LogoSVG       string `json:"logoSvg,omitempty"`
	LogoWidth     string `json:"logoWidth,omitempty"`
	LogoPosition  string `json:"logoPosition,omitempty"`
	Style         string `json:"style,omitempty"`
	CacheSeconds  string `json:"CacheSeconds,omitempty"`
}

func NewShieldFromBuild(s string) Shield {
	color := "important"
	isErr := true
	switch s {
	case "SUCCESS":
		color = "success"
		isErr = false
	case "STATUS_UNKNOWN":
		color = "inactive"
		isErr = false
	}
	return Shield{
		SchemaVersion: 1,
		Label:         "build",
		Message:       s,
		Color:         color,
		IsError:       isErr,
		Style:         "for-the-badge",
	}
}

//   "STATUS_UNKNOWN" - Status of the build is unknown.
//   "QUEUED" - Build or step is queued; work has not yet begun.
//   "WORKING" - Build or step is being executed.
//   "SUCCESS" - Build or step finished successfully.
//   "FAILURE" - Build or step failed to complete successfully.
//   "INTERNAL_ERROR" - Build or step failed due to an internal cause.
//   "TIMEOUT" - Build or step took longer than was allowed.
//   "CANCELLED" - Build or step was canceled by a user.
