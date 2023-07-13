package main

import (
	"fmt"

	"github.com/jbrukh/bayesian"
)

const (
	Good bayesian.Class = "Good"
	Bad  bayesian.Class = "Bad"
	// Grey bayesian.Class = "Grey"
)

func main() {

	classifier := bayesian.NewClassifier(Good, Bad)

	goodStuff := []string{"tall", "rich", "handsome", "tennis", "chicken", "theatre", "reading", "books", "money", "suits", "golf", "eating", "wine"}
	badStuff := []string{"poor", "smelly", "ugly", "playing", "games", "burgers", "tv", "magazines", "debts", "jogging", "foot", "drinking", "beer", "Zdob<C4><85>d<C5><BA>pobli\xc5\xbcu"}
	greyStuff := []string{"man", "girl", "young", "citizen"}

	classifier.Learn(goodStuff, Good)
	classifier.Learn(badStuff, Bad)
	classifier.Learn(greyStuff, Bad)

	// classifier.Learn(greyStuff, Grey)

	//classifier.WriteToFile("data")

	fmt.Printf("classifier learned :\n")
	for i := 0; i < len(classifier.Classes); i++ {
		fmt.Printf("  %d words for classes %v \n",
			classifier.WordCount()[i],
			classifier.Classes[i])
	}
	fmt.Println()

	tests := [][]string{
		{"tall", "girl", "watching", "tv", "eating", "chicken", "wearing", "jogging"},
		{"small", "handsome", "man", "playing", "games", "before", "tv"},
		{"unknows", "words"},
	}

	for _, pattern := range tests {

		fmt.Printf("test : %v\n", pattern)

		probs, likelyb, _ := classifier.ProbScores(pattern)
		fmt.Printf("probs ")
		for i := 0; i < len(classifier.Classes); i++ {
			fmt.Printf("%s(%.2f) ",
				classifier.Classes[i], probs[i])
		}
		fmt.Printf("=> Class : %s\n\n",
			classifier.Classes[likelyb])
	}
}

// func round(num float64) int {
// 	return int(num + math.Copysign(0.5, num))
// }

// func toFixed(num float64, precision int) float64 {
// 	output := math.Pow(10, float64(precision))
// 	return float64(round(num*output)) / output
// }
