package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudbuild/v1"
)

const (
	// ConsoleLink is a link to build history
	ConsoleLink = "https://console.cloud.google.com/cloud-build/builds?%s"
	// ShieldsIOLink is a link to the shields image
	ShieldsIOLink = "https://img.shields.io/badge/BUILD-%s-%s.svg?style=for-the-badge&maxAge=31536000"
)

var (
	Port       = os.Getenv("PORT")
	GCPProject = os.Getenv("PROJECT")

	// Possible values:
	//   "STATUS_UNKNOWN" - Status of the build is unknown.
	//   "QUEUED" - Build or step is queued; work has not yet begun.
	//   "WORKING" - Build or step is being executed.
	//   "SUCCESS" - Build or step finished successfully.
	//   "FAILURE" - Build or step failed to complete successfully.
	//   "INTERNAL_ERROR" - Build or step failed due to an internal cause.
	//   "TIMEOUT" - Build or step took longer than was allowed.
	//   "CANCELLED" - Build or step was canceled by a user.
	Color = map[string]string{
		"SUCCESS":        "brightgreen",
		"FAILURE":        "orange",
		"STATUS_UNKNOWN": "blue",
		"NOT_FOUND":      "lightgrey",
		"INTERNAL_ERROR": "red",
	}
)

func init() {
	if Port == "" {
		Port = ":8080"
	}
	if Port[0] != ':' {
		Port = ":" + Port
	}
}

func main() {
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
	http.HandleFunc("/i/", func(w http.ResponseWriter, r *http.Request) {
		repoName := strings.Split(r.URL.Path, "/")[2]
		status := "NOT_FOUND"
		res, err := svc.Projects.Builds.
			List(GCPProject).
			Filter(`source.repo_source.repo_name = "` + repoName + `"`).
			Fields("builds.status").
			Do()
		if err != nil {
			status = "INTERNAL_ERROR"
			log.Println("get builds status", err)
		}

		for _, b := range res.Builds {
			if b.Status == "WORKING" || b.Status == "QUEUED" || b.Status == "CANCELLED" {
				continue
			}
			status = b.Status
			break
		}
		http.Redirect(w, r, fmt.Sprintf(ShieldsIOLink, status, Color[status]), http.StatusFound)
	})

	http.HandleFunc("/l/", func(w http.ResponseWriter, r *http.Request) {
		repoName := strings.Split(r.URL.Path, "/")[2]
		vals := url.Values{
			"project": []string{
				GCPProject,
			},
			"query": []string{
				fmt.Sprintf(`source.repo_source.repo_name = "%s"`, repoName),
			},
		}
		http.Redirect(w, r, fmt.Sprintf(ConsoleLink, vals.Encode()), http.StatusFound)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		w.Write(indexHTML)
	})

	log.Fatal(http.ListenAndServe(Port, nil))
}

var indexHTML = []byte(`
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
`)
