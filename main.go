// Subayes main package is a bayesian cli build around github.com/jbrukh/bayesian
package main

// 2023/06 : cat : subayes : mail subject classification using bayesian filter
//
// v0.1 : working draft.
// v0.2 : minlen words, better split func, default bayes class, +main_test.go
// v0.3 : -E options for explaining and showing scores
// v1.0 tag for go doc
// V1.1 : ignore numbers
//
// TODO :
// - how to remove item from db ?

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	// Credits to "github.com/jbrukh/bayesian"
	"github.com/jbrukh/bayesian"
)

var (
	db, data            string
	learnSpam, learnHam bool
	explain, verbose    bool
)

func main() {
	var (
		minlength                = 4
		Spam      bayesian.Class = "Spam"
		Ham       bayesian.Class = "Ham"
	)

	// db is classes data store path
	flag.StringVar(&db, "db", "db", " db path")
	// data is the file to be read when learning
	flag.StringVar(&data, "d", "subayes.spam", "data filename")
	// choosing between learning Spam or Ham (write db/classes files)
	flag.BoolVar(&learnSpam, "learnSpam", false, "learn Spam subjects")
	flag.BoolVar(&learnHam, "learnHam", false, "learn Ham subjects")
	flag.IntVar(&minlength, "m", 4, "word min length")

	flag.BoolVar(&explain, "E", false, "explain words scores")

	flag.BoolVar(&verbose, "v", false, "verbose")

	// Default is to read stdin line per line for classification
	flag.Parse()

	K := bayesian.NewClassifier(Ham, Spam)

	switch {

	case learnHam && learnSpam:
		errcheck(errors.New("please choose learn Ham or Spam, not Both"))

	case learnHam:
		errcheck(learn(K, db, data, Ham, minlength))
		showClassesCount(K)

	case learnSpam:
		errcheck(learn(K, db, data, Spam, minlength))
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
			} else {
				if verbose {
					fmt.Fprintf(os.Stderr, "Warning short string : \"%s\"", text)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			errcheck(err)
		}
	}
}

// learn ingest data file into bayesian class and save to classifier db
func learn(c *bayesian.Classifier, xdb string, input string, class bayesian.Class, minilength int) (err error) {

	c.ReadClassFromFile(class, xdb)
	// if db/class don't exist, we will create it, so any err is acceptable
	// errcheck(err)
	// Better error handling should test error for acceptables ones
	showClassesCount(c)

	in, err := os.ReadFile(input) // in  type is []byte
	errcheck(err)

	ins := string(in) // ins type is string

	indata := split(ins) // indata is []string
	indedup := removeDuplicate(indata, minilength)
	c.Learn(indedup, class)

	err = c.WriteClassToFile(class, xdb)
	errcheck(err)

	return nil
}

// classify return bayesian Class of []string from a classifier
// if explain option is given, it print out class score on stderr
func classify(c *bayesian.Classifier, pattern []string, d bayesian.Class) bayesian.Class {
	if len(pattern) == 0 { // return default class
		if verbose {
			fmt.Fprintf(os.Stderr, "Warning Empty pattern\n")
		}
		return d
	}

	if explain {
		for _, word := range pattern {
			var wordarr []string
			wordarr = append(wordarr, word)
			scores, likelyb, _ := c.ProbScores(wordarr)
			fmt.Fprintf(os.Stderr, "[ %s = %s ] : ", word, c.Classes[likelyb])
			for i := 0; i < len(c.Classes); i++ {
				fmt.Fprintf(os.Stderr, "[%v]{ %.4f } ", c.Classes[i], scores[i])
			}
			fmt.Fprintf(os.Stderr, "\n")
		}
	}

	//  ProbScores return scores ([]float64), indexofclass, strict(?)
	_, likelyb, _ := c.ProbScores(pattern)
	// Would testing strict should be done ?
	// _, likelyb, strict := c.ProbScores(pattern)
	// if false returning default class d ?
	return c.Classes[likelyb]
}

// showClassesCount display classes and item counts
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

// errcheck func perform basic error check
func errcheck(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v", e)
		os.Exit(-1)
	}
}

// split function return []string of words from string
func split(s string) []string {
	var words = regexp.MustCompile(`[\p{L}]+`)
	// See http://www.unicode.org/reports/tr44/#General_Category_Values

	// 	return rxp.Split(s, -1)
	// words := regexp.MustCompile("\\w+")
	// words := regexp.MustCompile("\\P{M}+")
	// words := regexp.MustCompile("[\\p{L}]+")
	return words.FindAllString(s, -1)
}

// removeDuplicate function remove duplicate entries from []string
// and entries length must be > length parameter
func removeDuplicate(sliceList []string, length int) []string {
	var digits = regexp.MustCompile(`^[0-9\.]+$`)

	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range sliceList {
		if digits.MatchString(item) {
			continue
		}

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
