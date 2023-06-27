package main

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
	db, data            string
	learnSpam, learnHam bool
	Spam                bayesian.Class = "Spam"
	Ham                 bayesian.Class = "Ham"
	rxp                                = regexp.MustCompile(" |'|,|\t|\n")
)

func main() {

	// db is path for storing data classes
	flag.StringVar(&db, "db", "db", " db path")
	// data is the file to be read
	flag.StringVar(&data, "d", "data", "data filename")
	// choosing between learning Spam or Han (write db/classes files)
	flag.BoolVar(&learnSpam, "learnSpam", false, "Learn Spam subjects")
	flag.BoolVar(&learnHam, "learnHam", false, "Learn Ham subjects")
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
		showClassesCount(K)

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() { // read line per line
			text := scanner.Text()

			if len(text) > 3 { // Minimum 'Re:'

				spl := rxp.Split(text, -1)

				fmt.Printf("%v: %s\n",
					classify(K, spl),
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

	in, err := os.ReadFile(input) // in  type is []byte
	errcheck(err)

	ins := string(in) // ins type is string

	indata := rxp.Split(ins, -1) // indata is []string
	c.Learn(indata, class)

	err = c.WriteClassToFile(class, xdb)
	errcheck(err)

	return nil
}
func classify(c *bayesian.Classifier, pattern []string) bayesian.Class {
	//  ProbScores return scores ([]float64), indexofclass, strict(?)
	_, likelyb, _ := c.ProbScores(pattern)
	return c.Classes[likelyb]
}

func showClassesCount(c *bayesian.Classifier) {
	fmt.Printf("INFO classifier corpus : ")
	for i := 0; i < len(c.Classes); i++ {
		fmt.Printf(" [ %v -> %d items ]",
			c.Classes[i],
			c.WordCount()[i])

	}
	fmt.Println()
}

func errcheck(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(-1)
	}
}
