package main

import (
	"errors"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"strings"
)

type issue struct {
	commit *object.Commit // the commit that failed
	source string         // the objectionable string from the source
	file   string         // the file containing the issue
	line   int            // the line in the file with the issue
	reason string         // the semantic reason (eg, entropy high, which regex rule, etc)
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

	//ref, err := repo.Head()
	//CheckIfError(err, "retreiving repo head failed")

	// ... retrieves the commit history
	//cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	//CheckIfError(err, "retreiving commit history failed")

	// -- get all commits in repo?
	cIter, err := repo.CommitObjects()

	issues := make([]issue, 0)
	err = cIter.ForEach(func(c *object.Commit) error {
		return checkCommit(c, &issues)
	})
	CheckIfError(err, "checking commits failed")
	return issues, nil
}

func checkCommit(c *object.Commit, issues *[]issue) error {
	//	fmt.Println(fmt.Sprintf("type(c)=%T", c))
	//	fmt.Println(c)

	fileIter, err := c.Files() // possible to iterate by diff/patch instead of file?
	CheckIfError(err, "issues getting files from commit")

	fileIter.ForEach(func(file *object.File) error {

		fisb, err := file.IsBinary()
		CheckIfError(err, "can't determine if file is binary")
		if fisb {
			return nil
		}
		lines, err := file.Lines()
		//TODO filter lines to remove false positives here?
		CheckIfError(err, "can't get lines in file")
		entropyissues, err := checkEntropy(lines, file.Name, c)
		CheckIfError(err, "checking entropy failed")
		*issues = append(*issues, entropyissues...)

		//TODO check regexes here
		return nil
	})
	return nil
}

func checkEntropy(lines []string, filename string, c *object.Commit) ([]issue, error) {
	entropyissues := make([]issue, 0)
	fmt.Println(lines)
	for linenum, lineval := range lines {
		for _, word := range strings.Fields(lineval) {
			highentropy, err := wordEntropy(word)
			CheckIfError(err, "couldn't calculate word entropy?")
			if highentropy {
				newissue := issue{c, word, filename, linenum, "high entropy word"}
				entropyissues = append(entropyissues, newissue)
			}
		}
	}
	return entropyissues, nil
}

func wordEntropy(word string) (bool, error) {
	return true, errors.New("not implemented")
}
