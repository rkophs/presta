package presta

import (
	"fmt"
	"io"
)

func Tokenize(reader io.Reader) bool {
	s := NewScanner(reader)
	for {
		tok, lit := s.Scan()
		fmt.Printf("%i\t%q\n", tok, lit)
		if tok == EOF {
			return false
		} else if tok == ILLEGAL {
			fmt.Printf("Illigal token: %q\n", lit)
			return true
		}
	}
}
