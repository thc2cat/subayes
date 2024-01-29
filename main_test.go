package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/jbrukh/bayesian"
)

var (
	Spam bayesian.Class = "Spam"
	Ham  bayesian.Class = "Ham"

	tests = []struct {
		name       string
		s          string
		want       []string
		class      bayesian.Class
		wordscount int
	}{
		{"empty", "", nil, Ham, 0},
		{"A test !", "A test !", []string{"A", "test"}, Ham, 1},
		{"puctuation", "a lot of # numbers 1,2,5 ? ", []string{"a", "lot", "of", "numbers"}, Ham, 2},
		{"real1", "Re: [3616304] MODIFICATION // Fwd:  Senegal", []string{"Re", "MODIFICATION", "Fwd", "Senegal"}, Ham, 3},
		{"real2", "Chèque MAIF 1 000€", []string{"Chèque", "MAIF"}, Ham, 2},
		{"real3", "Re: Rattrapages (examens de seconde chance) 12-16 juin et 19-23  juin",
			[]string{"Re", "Rattrapages", "examens", "de", "seconde", "chance", "juin", "et", "juin"}, Ham, 5},
		{"real4", "Vaše vlastné dievča lokálne", []string{"Vaše", "vlastné", "dievča", "lokálne"}, Spam, 4},
		{"real5", "Compañeras własne Zdobądź nära pobliżu",
			[]string{"Compañeras", "własne", "Zdobądź", "nära", "pobliżu"}, Spam, 5},
		{"apostrophes", "Parlez avec des meufs pour s’envoyer en l’air dès maintenant",
			[]string{"Parlez", "avec", "des", "meufs", "pour", "s", "envoyer", "en", "l", "air", "dès", "maintenant"}, Spam, 9},
		{"with numbers", "ceci contient 1.999 test et ceci 10000 excuses",
			[]string{"ceci", "contient", "test", "et", "ceci", "excuses"}, Ham, 4},
	}

	K = bayesian.NewClassifier(Ham, Spam)
)

func Test_split(t *testing.T) {
	verbose = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_most(t *testing.T) {
	verbose = true
	if _, err := os.Stat("db"); !os.IsNotExist(err) {
		// path/to/whatever does not exist
		errcheck(K.ReadClassFromFile(Spam, "db"))
		errcheck(K.ReadClassFromFile(Ham, "db"))
		showClassesCount(K)
	}
}

func Test_removeDuplicate(t *testing.T) {
	verbose = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicate(split(tt.s), 3); len(got) != tt.wordscount {
				t.Errorf("removeDuplicate() '%s' is '%s' gave '%v', we want '%v'", tt.name, tt.s, len(got), tt.wordscount)
			}
		})
	}
}

func Test_classify(t *testing.T) {
	verbose = true
	if _, err := os.Stat("db"); !os.IsNotExist(err) {
		// path/to/whatever does not exist
		errcheck(K.ReadClassFromFile(Spam, "db"))
		errcheck(K.ReadClassFromFile(Ham, "db"))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classify(K, tt.want, Ham); got != tt.class {
				t.Errorf("For %v test : %v got %v but want %v\n",
					tt.name, classify(K, tt.want, Ham), got, tt.want)
			}
		})
	}
}
