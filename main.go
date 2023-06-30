package main

// 2023/06 : cat : subayes : mail subject classification using bayesian filter
//
// v0.1 : working draft.
// v0.2 : minlen words, better split func, default bayes class, +main_test.go

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/navossoc/bayesian"
)

var (
	db, data                     string
	learnSpam, learnHam, verbose bool
	Spam                         bayesian.Class = "Spam"
	Ham                          bayesian.Class = "Ham"
	words                                       = regexp.MustCompile("[\\p{L}]+")
	minlength                                   = 4
)

func main() {

	// db is classes data store path
	flag.StringVar(&db, "db", "db", " db path")
	// data is the file to be read when learning
	flag.StringVar(&data, "d", "data", "data filename")
	// choosing between learning Spam or Ham (write db/classes files)
	flag.BoolVar(&learnSpam, "learnSpam", false, "Learn Spam subjects")
	flag.BoolVar(&learnHam, "learnHam", false, "Learn Ham subjects")
	flag.IntVar(&minlength, "m", 4, "word min length")

	flag.BoolVar(&verbose, "v", false, "verbose")

	// Default is to read stdin line per line for classification
	flag.Parse()

	K := bayesian.NewClassifier(Ham, Spam)

	switch {

	case learnHam && learnSpam:
		errcheck(errors.New("Please choose learn Ham or Spam, not Both !"))

	case learnHam:
		errcheck(learn(K, db, data, Ham))
		showClassesCount(K)

	case learnSpam:
		errcheck(learn(K, db, data, Spam))
		showClassesCount(K)

	case !learnHam && !learnSpam:
		errcheck(K.ReadClassFromFile(Spam, db))
		errcheck(K.ReadClassFromFile(Ham, db))
		// Is it needed ? TF-IDF
		K.ConvertTermsFreqToTfIdf()

		showClassesCount(K)

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() { // read line per line
			text := scanner.Text()
			if len(text) > minlength { // Minimum 'Re: '
				spl := removeDuplicate(split(text), minlength)
				fmt.Printf("%v: %s\n",
					classify(K, spl, Ham),
					text)
			}
		}
		if err := scanner.Err(); err != nil {
			errcheck(err)
		}
	}
}

func learn(c *bayesian.Classifier, xdb string, input string, class bayesian.Class) (err error) {

	err = c.ReadClassFromFile(class, xdb)
	// if db/class don't exist, we will create it, so any err is acceptable
	// errcheck(err)
	// Better error handling should test error for acceptables ones
	showClassesCount(c)

	in, err := os.ReadFile(input) // in  type is []byte
	errcheck(err)

	ins := string(in) // ins type is string

	indata := split(ins) // indata is []string
	indedup := removeDuplicate(indata, minlength)
	c.Learn(indedup, class)

	err = c.WriteClassToFile(class, xdb)
	errcheck(err)

	return nil
}
func classify(c *bayesian.Classifier, pattern []string, d bayesian.Class) bayesian.Class {
	//  ProbScores return scores ([]float64), indexofclass, strict(?)
	if len(pattern) == 0 { // return default class
		return d
	}
	_, likelyb, _ := c.ProbScores(pattern)
	return c.Classes[likelyb]
}

func showClassesCount(c *bayesian.Classifier) {
	if !verbose {
		return
	}
	fmt.Fprintf(os.Stderr, "INFO classifier corpus : ")
	for i := 0; i < len(c.Classes); i++ {
		if c.WordCount()[i] > 0 {
			fmt.Fprintf(os.Stderr,
				" [ %v -> %d items ]",
				c.Classes[i],
				c.WordCount()[i])
		}
	}
	fmt.Fprintln(os.Stderr)
}

func errcheck(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v", e)
		os.Exit(-1)
	}
}

func split(s string) []string {
	// 	return rxp.Split(s, -1)
	// words := regexp.MustCompile("\\w+")
	// words := regexp.MustCompile("\\P{M}+")
	// words := regexp.MustCompile("[\\p{L}]+")
	return words.FindAllString(s, -1)
}

// Generic version (string or int )
// func removeDuplicate[T string | int](sliceList []T) []T {
// 	allKeys := make(map[T]bool)
// 	list := []T{}
// 	for _, item := range sliceList {
// 		if _, value := allKeys[item]; !value {
// 			allKeys[item] = true
// 			list = append(list, item)
// 		}
// 	}
// 	return list
// }

// removeduplicate
func removeDuplicate(sliceList []string, length int) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range sliceList {
		if len(item) < length {
			continue
		}
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
