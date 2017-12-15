package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func CheckIfError(err error, message string) {
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}

func main() {
	CheckRepo("https://github.com/src-d/go-siva")
}

func CheckRepo(repourl string) {
	fmt.Println("checking %s", repourl)

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repourl,
	})
	CheckIfError(err, "memory clone failed")

	ref, err := repo.Head()
	CheckIfError(err, "retreiving repo head failed")

	// ... retrieves the commit history
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err, "retreiving commit history failed")

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	CheckIfError(err, "checking commits failed")
}
