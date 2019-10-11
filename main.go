package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"google.golang.org/api/cloudbuild/v1"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	initLog()

	ctx := context.Background()
	cb, err := cloudbuild.NewService(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cloudbuild.NewService")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	project := os.Getenv("PROJECT")
	if project == "" {
		log.Fatal().Msg("PROJECT must be set")
	}

	http.Handle("/badger/i/", http.StripPrefix("/badger/i/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tag string
		if tag = strings.Split(r.URL.Path, "/")[0]; tag == "" {
			log.Error().Str("url", r.URL.String()).Msg("no tag specified")
			linkImg(w, r, http.StatusText(http.StatusBadRequest), "red")
			return
		}

		res, err := cb.Projects.Builds.List(project).Filter(fmt.Sprintf(`tags="%s"`, tag)).Fields("builds.status").Do()
		if err != nil {
			log.Error().Str("tag", tag).Err(err).Msg("list builds")
			linkImg(w, r, http.StatusText(http.StatusInternalServerError), "red")
			return
		}

		status, color := "no builds", "red"
		for _, b := range res.Builds {
			if b.Status == "CANCELLED" {
				continue
			}
			status = strings.ToLower(b.Status)
			color = "orange"
			if status == "success" {
				color = "brightgreen"
			}
			break
		}

		log.Info().Str("tag", tag).Str("status", status).Msg("served")
		linkImg(w, r, status, color)
	})))
	http.Handle("/badger/i/", http.StripPrefix("/badger/i/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tag string
		if tag = strings.Split(r.URL.Path, "/")[0]; tag == "" {
			log.Error().Str("url", r.URL.String()).Msg("no tag specified")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		log.Info().Str("tag", tag).Msg("redirected")
		linkLog(w, r, project, tag)
	})))

	log.Info().Str("port", port).Msg("serving")
	http.ListenAndServe(":"+port, nil)
}

func initLog() {
	logfmt := os.Getenv("LOGFMT")
	if logfmt != "json" {
		logfmt = "text"
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: !terminal.IsTerminal(int(os.Stdout.Fd()))})
	}

	level, _ := zerolog.ParseLevel(os.Getenv("LOGLVL"))
	if level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	log.Info().Str("FMT", logfmt).Str("LVL", level.String()).Msg("log initialized")
	zerolog.SetGlobalLevel(level)
}

func linkImg(w http.ResponseWriter, r *http.Request, status, color string) {
	u := url.URL{
		Scheme:   "https:",
		Host:     "img.shields.io",
		Path:     "/badge/build-" + status + "-" + color + ".svg",
		RawQuery: r.URL.RawQuery,
	}
	http.Redirect(w, r, u.String(), http.StatusFound)
}
func linkLog(w http.ResponseWriter, r *http.Request, project, tag string) {
	u := url.URL{
		Scheme:   "https",
		Host:     "console.cloud.google.com",
		Path:     "/cloud-build/builds",
		RawQuery: url.Values{"project": []string{project}, "query": []string{`tags="` + tag + `"`}}.Encode(),
	}
	http.Redirect(w, r, u.String(), http.StatusFound)
}
