package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudbuild/v1"
)

const (
	ConsoleLink = "https://console.cloud.google.com/cloud-build/builds"
)

var (
	GCPProject = "com-seankhliao"
)

func main() {
	p := flag.String("p", "8080", "port to listen on")
	pr := flag.String("pr", GCPProject, "GCP project to query")
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
	http.HandleFunc("/i/", func(w http.ResponseWriter, r *http.Request) {
		repoName := strings.Split(r.URL.Path, "/")[2]
		u := "https://img.shields.io/badge/endpoint.svg?url=https://badger.seankhliao.com/r/" + repoName
		http.Redirect(w, r, u, http.StatusMovedPermanently)
	})
	http.HandleFunc("/l/", func(w http.ResponseWriter, r *http.Request) {
		repoName := strings.Split(r.URL.Path, "/")[2]
		vals := url.Values{
			"project": []string{
				*pr,
			},
			"query": []string{
				fmt.Sprintf(`source.repo_source.repo_name = "%s"`, repoName),
			},
		}
		u := ConsoleLink + "?" + vals.Encode()
		http.Redirect(w, r, u, http.StatusFound)

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		w.Write([]byte(`
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1" />

    <title>Badger</title>
    <meta name="description" content="Markdown badges for Cloud Build results" />

    <style>
	* {
	  box-sizing: border-box;
	}
	body, section {
	  display: flex;
	  flex-flow: column nowrap;
	  align-items: center;
	  margin: 0;
	  width: 100%;
	}
	label {
	  margin-top: 1em;
	}
	section input {
	  margin-bottom: 1em;
	  width: 90%;
	}
    </style>
  </head>
  <body>
    <h1>Badger</h1>
    <p>Badges for Cloud Build*</p>
    <p>
      <a href="https://github.com/seankhliao/badger">Source on Github</a>
      <a href="https://badger.seankhliao.com/l/github_seankhliao_badger"
	<img src="https://badger.seankhliao.com/i/github_seankhliao_badger" />
      </a>
    </p>
    <p>* <em>Only works with projects it had access to</em></p>

    <h2>Generate</h2>
    <label for"repo">Repo name (github format: <code>github_$user_$repo</code>)</label>
    <input type="text" id="repo" name="repo"/>

    <section>
	    <label for"out1">img url: </label>
	    <input type="text" id="out1" name="out1">

	    <label for"out2">markdown img: </label>
	    <input type="text" id="out2" name="out2">

	    <label for"out3">markdown with link: </label>
	    <input type="text" id="out3" name="out3">
    </section>


    <script>
      const repo = document.querySelector('#repo');
      const out1 = document.querySelector('#out1');
      const out2 = document.querySelector('#out2');
      const out3 = document.querySelector('#out3');
      const imgLinkUrl = 'https://badger.seankhliao.com/i/';
      repo.addEventListener('input', function(e){
	imgUrl = imgLinkUrl + repo.value;
	linkUrl = 'https://badger.seankhliao.com/l/' + repo.value;
	out1.value = imgUrl;
	out2.value = '![Build](' + imgUrl + ')';
	out3.value = '[![Build](' + imgUrl + ')](' + linkUrl + ')';
      });
    </script>
  </body>
</html>
		`))

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
