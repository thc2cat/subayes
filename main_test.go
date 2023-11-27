package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/jbrukh/bayesian"
)

func Test_split(t *testing.T) {

	tests := []struct {
		name string
		s    string
		want []string
	}{
		// {"empty", "", []string{""}},
		{"empty", "", nil},
		{"A test !", "A test !", []string{"A", "test"}},
		{"puctuation", "a lot of # numbers 1,2,5 ? ", []string{"a", "lot", "of", "numbers"}},
		{"real1", "Re: [3616304] MODIFICATION // Fwd:  Senegal", []string{"Re", "MODIFICATION", "Fwd", "Senegal"}},
		{"real2", "Chèque MAIF 1 000€", []string{"Chèque", "MAIF"}},
		{"real3", "Re: Rattrapages (examens de seconde chance) 12-16 juin et 19-23  juin",
			[]string{"Re", "Rattrapages", "examens", "de", "seconde", "chance", "juin", "et", "juin"}},
		{"real4", "Vaše vlastné dievča lokálne", []string{"Vaše", "vlastné", "dievča", "lokálne"}},
		{"real5", "Compañeras własne Zdobądź nära pobliżu",
			[]string{"Compañeras", "własne", "Zdobądź", "nära", "pobliżu"}},
		{"apostrophes", "Parlez avec des meufs pour s’envoyer en l’air dès maintenant",
			[]string{"Parlez", "avec", "des", "meufs", "pour", "s", "envoyer", "en", "l", "air", "dès", "maintenant"}},
	}
	// TODO: Add test cases.

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_most(t *testing.T) {
	var (
		Spam bayesian.Class = "Spam"
		Ham  bayesian.Class = "Ham"
	)
	verbose = true
	K := bayesian.NewClassifier(Ham, Spam)

	if _, err := os.Stat("db"); !os.IsNotExist(err) {
		// path/to/whatever does not exist
		errcheck(K.ReadClassFromFile(Spam, "db"))
		errcheck(K.ReadClassFromFile(Ham, "db"))
		showClassesCount(K)
	}

}

func Test_removeDuplicate(t *testing.T) {
	tests := []struct {
		name   string
		array  []string
		want   []string
		length int
	}{
		// WARNING : lowercase since V1.2
		{"basic", []string{"AAAA", "AAAA", "BBB"}, []string{"aaaa", "bbb"}, 3},
		{"length3", []string{"A", "BBB", "CCC", "BBB"}, []string{"bbb", "ccc"}, 3},
		{"length2", []string{"AA", "BBB", "CCC", "BBB"}, []string{"aa", "bbb", "ccc"}, 2},
		{"with numbers", []string{"ceci", "contient", "1.999", "test", "10000", "excuses"}, []string{"ceci", "contient", "test", "excuses"}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicate(tt.array, tt.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicate() gave %v, we want %v", got, tt.want)
			}
		})
	}
}
