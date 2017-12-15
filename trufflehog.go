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
	fmt.Println("printing issues:", len(issues))
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
		return checkCommit(c, &issues)
	})
	CheckIfError(err, "checking commits failed")
	return issues, nil
}

func checkCommit(c *object.Commit, issues *[]issue) error {
	entropyissues, err := checkEntropy(c)
	CheckIfError(err, "checking entropy failed")
	regexissues, err := checkRegexes(c)
	CheckIfError(err, "checking regex failed")
	*issues = append(*issues, entropyissues...)
	*issues = append(*issues, regexissues...)
	return nil
}

func checkEntropy(c *object.Commit) ([]issue, error) {
	//	entropyissues := make([]issue, 0)
	//	newissue := issue{"commit has issue: " + c.String(), c}
	return nil, errors.New("not implemented")
}
func checkRegexes(c *object.Commit) ([]issue, error) {
	//	regexissue := make([]issue, 0)
	return nil, errors.New("not implemented")
}
