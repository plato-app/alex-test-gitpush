package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

const (
	platoConfigFileName = ".plato/config.json"
	testFile            = "test.txt"
)

var platoConfig = platoConfigStruct{
	GitAuthorName:  "Plato",
	GitAuthorEmail: "ops@platoteam.com",
}

type platoConfigStruct struct {
	GitAuthorEmail string `json:"GitAuthorEmail"`
	GitAuthorName  string `json:"GitAuthorName"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Debugf("fetch git author info from plato config (if present)")
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Debug(err)
	}
	platoConfigFile := filepath.Join(dirname, platoConfigFileName)
	data, err := ioutil.ReadFile(platoConfigFile)
	if err != nil {
		log.Debug(err)
	}
	if err = json.Unmarshal(data, &platoConfig); err != nil {
		log.Debug(err)
		log.Debug("fall back to default Git author info")
	}
	log.Debugf("using following Git author config: %v", platoConfig)
	log.Debug("open Git repo in current working directory")
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	if err = worktree.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
		log.Errorf("error while pulling: %v", err)
	}
	log.Debugf("write test file: %s", testFile)
	ioutil.WriteFile(testFile, []byte(time.Now().String()), 0644)
	_, err = worktree.Add(testFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("create test commit")
	hash, err := worktree.Commit("Test commit", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  platoConfig.GitAuthorName,
			Email: platoConfig.GitAuthorEmail,
		},
	})
	log.Debugf("commit hash: %s", hash)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("push commit to origin")
	if err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
	}); err != nil {
		log.Fatal(err)
	}
}
