package main

import (
	"errors"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type issue struct {
	reason string
	commit *object.Commit
}

func CheckIfError(err error, message string) {
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}

func main() {
	issues, err := CheckRepo("https://github.com/src-d/go-siva")
	CheckIfError(err, "CheckRepo failed...")
	printIssues(issues)
}

func printIssues(issues []issue) error {
	fmt.Println(issues)
	return nil
}
func CheckRepo(repourl string) ([]issue, error) {
	fmt.Println("checking ", repourl)

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repourl,
	})
	CheckIfError(err, "memory clone failed")

	ref, err := repo.Head()
	CheckIfError(err, "retreiving repo head failed")

	// ... retrieves the commit history
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err, "retreiving commit history failed")

	issues := make([]issue, 0)
	err = cIter.ForEach(func(c *object.Commit) error {
		return checkCommit(c, issues)
	})
	fmt.Println(issues)
	//err = cIter.ForEach(func(c *object.Commit) error {
	//	fmt.Println(c)
	//	return nil
	//})
	CheckIfError(err, "checking commits failed")
	return nil, nil
}

func checkCommit(c *object.Commit, issues []issue) error {
	commitissues := make([]issue, 0)
	newissue := issue{"commit has issue: " + c.String(), c}
	commitissues = append(commitissues, newissue)
	issues = append(commitissues)
	return errors.New("not implemented")
}


func checkEntropy(issues []issue) error {
	return errors.New("not implemented")
}
func checkRegexes(issues []issue) error {
	return errors.New("not implemented")
}
