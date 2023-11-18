package main

import (
	"log"
	"testing"
)

func TestUpdate(t *testing.T) {
	var req = &funcUpdateReq{
		name:    "func-iptv",
		version: "23.11.3",
	}

	repo := funcRepo{
		url:      envOrThrow("GIT_REPO"),
		user:     envOrThrow("GIT_USER"),
		email:    envOrThrow("GIT_EMAIL"),
		password: envOrThrow("GIT_PASSWORD"),
	}

	err := req.update(&repo)
	if err != nil {
		log.Panicln(err.Error())
	}

	out, err := req.deploy(envOrThrow("KUBE_CONFIG"))
	if err != nil {
		log.Panicln(err.Error())
	} else {
		log.Println(out)
	}
}
