package main

import (
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

	verbose = true
	K := bayesian.NewClassifier(Ham, Spam)

	errcheck(K.ReadClassFromFile(Spam, "db"))
	errcheck(K.ReadClassFromFile(Ham, "db"))
	showClassesCount(K)
}

func Test_removeDuplicate(t *testing.T) {
	tests := []struct {
		name   string
		array  []string
		want   []string
		length int
	}{
		{"basic", []string{"AAAA", "AAAA", "BBB"}, []string{"AAAA", "BBB"}, 3},
		{"length3", []string{"A", "BBB", "CCC", "BBB"}, []string{"BBB", "CCC"}, 3},
		{"length2", []string{"AA", "BBB", "CCC", "BBB"}, []string{"AA", "BBB", "CCC"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicate(tt.array, tt.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
