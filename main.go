package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
)


const version = "github-webhook v2"

var requiredGitVersion = "git version 2.11.0"

func main() {
	log.Println("Starting GitHub webhook...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on HTTP port: %s", port)

	cmd := exec.Command("git", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		log.Fatal(err)
	}

	gitVersion := strings.TrimSpace(string(output))

	if gitVersion != requiredGitVersion {
		log.Fatalf("git version mismatch. Got %s want %s", gitVersion, requiredGitVersion)
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version + "\n")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		githubRepo, err := githubRepositoryFromRequest(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		dir, err := ioutil.TempDir("", "repo")
		if err != nil {
			log.Printf("Unable to clone the %s repo: %s", *githubRepo.CloneURL, err)
			w.WriteHeader(500)
			return
		}
		defer os.RemoveAll(dir)

		log.Printf("Cloning %s", *githubRepo.CloneURL)
		cmd := exec.Command("git", "clone", *githubRepo.CloneURL, dir)
		output, err := cmd.CombinedOutput()
		log.Println(string(output))
		if err != nil {
			log.Printf("Unable to clone the %s repo: %s", *githubRepo.CloneURL, err)
			w.WriteHeader(500)
			return
		}
	})

	log.Fatal(http.ListenAndServe(net.JoinHostPort("", port), nil))
}

func githubRepositoryFromRequest(r *http.Request) (*github.PushEventRepository, error) {
	payload, err := github.ValidatePayload(r, []byte("pipeline"))
	if err != nil {
		return nil, fmt.Errorf("Unable to validate webhook payload: %s", err)
	}

	webHookType := github.WebHookType(r)
	if webHookType != "push" {
		return nil, fmt.Errorf("The %s event type is not supported", webHookType)
	}

	event, err := github.ParseWebHook(webHookType, payload)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse webhook payload: %s", err)
	}

	return event.(*github.PushEvent).GetRepo(), nil
}
