package main

import (
	"github.com/xanzy/go-gitlab"
	"log"
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	gitClone()
}

func gitClone() {
	CheckArgs("<url>", "<directory>", "<github_username>", "<github_password>")
	url, directory, username, password := os.Args[1], os.Args[2], os.Args[3], os.Args[4]

	// Clone the given repository to the given directory
	Info("git clone %s %s", url, directory)

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
		URL:      url,
		Progress: os.Stdout,
	})
	CheckIfError(err)

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}



func repositoryFileExample() {
	git, err := gitlab.NewClient("glpat-J2Vy55PwsvXac7xkhs6B")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new repository file
	cf := &gitlab.CreateFileOptions{
		Branch:        gitlab.String("master"),
		Content:       gitlab.String("test file content "),
		CommitMessage: gitlab.String("Adding a test file"),
	}
	_, _, err = git.RepositoryFiles.CreateFile("fulcrum29/argoapps", "test/testfile.txt", cf)
	if err != nil {
		log.Fatal(err)
	}

	
}
