package main

import (
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
	ConsoleLink   = "https://console.cloud.google.com/cloud-build/builds"
	ShieldsIOLink = "https://img.shields.io/badge/BUILD-%s-%s.svg?style=for-the-badge&maxAge=31536000"
)

var (
	GCPProject = "com-seankhliao"
	ColorMap   = map[string]string{
		"SUCCESS":        "success",
		"FAILURE":        "warning",
		"STATUS_UNKNOWN": "informational",
		"NOT_FOUND":      "inactive",
		"INTERNAL_ERROR": "critical",
	}
)

func ShieldLink(status string) string {
	return fmt.Sprintf(ShieldsIOLink, status, ColorMap[status])
}

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
	http.HandleFunc("/i/", func(w http.ResponseWriter, r *http.Request) {
		repoName := strings.Split(r.URL.Path, "/")[2]
		s, err := status(svc, *pr, repoName)
		if err != nil {
			s = "INTERNAL_ERROR"
		}
		http.Redirect(w, r, ShieldLink(s), http.StatusFound)

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		w.Write(indexHTML)
	})
	http.ListenAndServe(":"+*p, nil)
}

// Possible values:
//   "STATUS_UNKNOWN" - Status of the build is unknown.
//   "QUEUED" - Build or step is queued; work has not yet begun.
//   "WORKING" - Build or step is being executed.
//   "SUCCESS" - Build or step finished successfully.
//   "FAILURE" - Build or step failed to complete successfully.
//   "INTERNAL_ERROR" - Build or step failed due to an internal cause.
//   "TIMEOUT" - Build or step took longer than was allowed.
//   "CANCELLED" - Build or step was canceled by a user.
func status(svc *cloudbuild.Service, p, r string) (string, error) {
	res, err := svc.Projects.Builds.
		List(p).
		Filter(`source.repo_source.repo_name = "` + r + `"`).
		Fields("builds.status").
		Do()
	if err != nil {
		return "INTERNAL_ERROR", err
	}

	// filter qorking / queued / cancelled
	status := "NOT_FOUND"
	for _, b := range res.Builds {
		if b.Status == "WORKING" || b.Status == "QUEUED" || b.Status == "CANCELLED" {
			continue
		}
		status = b.Status
		break
	}
	return status, nil
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
