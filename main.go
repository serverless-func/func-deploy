package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type GithubWebhook struct {
	Action      string `json:"action"` // must be "completed"
	WorkflowRun struct {
		HeadBranch string `json:"head_branch"`
	} `json:"workflow_run"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

type funcRepo struct {
	url      string
	user     string
	email    string
	password string
}

type funcUpdateReq struct {
	name    string
	version string
}

var repoPath = "/tmp/repo/"

func main() {
	repo := funcRepo{
		url:      envOrThrow("GIT_REPO"),
		user:     envOrThrow("GIT_USER"),
		email:    envOrThrow("GIT_EMAIL"),
		password: envOrThrow("GIT_PASSWORD"),
	}
	kubeconfig := envOrThrow("KUBE_CONFIG")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "func-deploy")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "pong")
	})
	http.HandleFunc("/github", func(w http.ResponseWriter, r *http.Request) {
		req, err := parse(r)
		if err != nil {
			handleError(w, err)
			return
		}
		if req.version == "latest" {
			_, _ = fmt.Fprintf(w, "skip deploy latest version")
			return
		}
		if err = req.update(&repo); err != nil {
			handleError(w, err)
			return
		}
		out, err := req.deploy(kubeconfig)
		if err != nil {
			handleError(w, err)
			return
		}
		_, _ = fmt.Fprintf(w, out)
	})

	port := envOrDefault("FC_SERVER_PORT", "9000")

	log.Println("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func envOrDefault(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func envOrThrow(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("env '%s' not set", name)
	}
	return value
}

func handleError(w http.ResponseWriter, err error) {
	log.Println("%w", err)
	_, _ = fmt.Fprintf(w, err.Error())
}

func parse(r *http.Request) (*funcUpdateReq, error) {
	decoder := json.NewDecoder(r.Body)
	var webhook GithubWebhook
	err := decoder.Decode(&webhook)
	if err != nil {
		return nil, err
	}
	if webhook.Action != "completed" {
		return nil, errors.New("workflow not completed")
	}
	return &funcUpdateReq{
		name:    webhook.Repository.Name,
		version: webhook.WorkflowRun.HeadBranch,
	}, nil
}
