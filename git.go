package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitAuth "github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func (req *funcUpdateReq) update(repo *funcRepo) error {
	err := os.RemoveAll(repoPath)
	if err != nil {
		return fmt.Errorf("fail to remove local repo : %w", err)
	}
	auth := &gitAuth.BasicAuth{
		Username: repo.user,
		Password: repo.password,
	}
	// clone iac repo
	r, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth:     auth,
		URL:      repo.url,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("fail to clone repo : %w", err)
	}
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("fail to get worktree : %w", err)
	}

	// read service manifest
	svcFile := path.Join(repoPath, "func", req.name+`.yaml`)
	fileBytes, err := os.ReadFile(svcFile)
	if err != nil {
		return fmt.Errorf("fail to read svc file : %w", err)
	}

	// find and change image version
	var re = regexp.MustCompile(req.name + `:.*`)
	oldVersion := re.Find(fileBytes)
	newVersion := req.name + ":" + req.version
	log.Printf("update [%s] => [%s] \n", oldVersion, newVersion)
	newFileBytes := re.ReplaceAll(fileBytes, []byte(newVersion))
	// write back
	err = os.WriteFile(svcFile, newFileBytes, 0644)
	if err != nil {
		return fmt.Errorf("fail to write svc file : %w", err)
	}

	// commit to repo
	_, err = w.Add(strings.ReplaceAll(svcFile, repoPath, ""))
	if err != nil {
		return fmt.Errorf("fail to stage svc file : %s, %w", svcFile, err)
	}
	_, err = w.Commit("Auto Update Version By Webhook", &git.CommitOptions{
		Author: &object.Signature{
			Name:  repo.user,
			Email: repo.email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("fail to commit svc file : %w", err)
	}

	// push changes
	err = r.Push(&git.PushOptions{
		Auth: auth,
	})
	if err != nil {
		return fmt.Errorf("fail to push commit : %w", err)
	}

	return nil
}
