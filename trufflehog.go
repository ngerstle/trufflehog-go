package main

import (
	"errors"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"math"
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
	// TODO get rid of helper function
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}

func main() {
	// run code
	// TODO arguments/configuration file
	// TODO curl | jq for all repositories in org/GHE

	repo := "https://github.com/comoyo/terraform-modules"
	//	repo := "https://github.com/src-d/go-siva"
	issues, err := CheckRepo(repo)
	CheckIfError(err, "CheckRepo failed...")
	printIssues(issues)
}

func printIssues(issues []issue) error {
	// pretty print discovered issues. i
	// TODO output json

	fmt.Println("printing issues:", len(issues))
	fmt.Println(issues)
	return nil
}
func CheckRepo(repourl string) ([]issue, error) {
	// checks a repository for any issues, in particular, secrets.

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
	//Checks a commit for any issues

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
	// checks each line in a file for high entropy words, creating issues as discovered

	entropyissues := make([]issue, 0)
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
	// calculates shannon entropy of a word. if the word surpasses thresholds, returns true.
	//TODO optimize character counting loops...

	if len(word) < 1 {
		return false, nil
	}
	wordlen := float64(len(word))
	entropy := 0.0
	charset := "1234567890abcdefABCDEF"
	for _, char := range charset {
		px := float64(strings.Count(word, string(char))) / wordlen
		if px > 0.0 {
			entropy += (-px * math.Log2(px))
		}
	}
	exentropy := 0.0
	excharset := "ghijklmnopqrstuvwyzGHIJKLMNOPQRSTUVWXYZ+/="
	for _, char := range excharset {
		px := float64(strings.Count(word, string(char))) / wordlen
		if px > 0.0 {
			exentropy += (-px * math.Log2(px))
		}
	}
	if exentropy > 0.0 { //must be base64
		entropy += exentropy
		return (entropy > 4.5), nil
	}
	return (entropy > 3), nil
}

func checkRegexes(lines []string, filename string, c *object.Commit) ([]issue, error) {
	return nil, errors.New("not implemented")
}
