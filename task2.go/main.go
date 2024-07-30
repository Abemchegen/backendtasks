package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func frequency(s string) map[string]int {
	word := false
	start := 0
	wordmap := make(map[string]int)

	for i, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !word {
				word = true
				start = i
			}
		} else {
			if word {
				w := strings.ToLower(s[start:i])
				wordmap[w]++
				word = false
			}
		}
	}
	if word {
		w := strings.ToLower(s[start:])
		wordmap[w]++
	}
	return wordmap
}

func palindrome(word string) bool {
	i := 0
	j := len(word) - 1
	word = strings.ToLower(word)

	for i < j {
		for i < j && !unicode.IsDigit(rune(word[i])) && !unicode.IsLetter(rune(word[i])) {
			i++
		}
		for i < j && !unicode.IsDigit(rune(word[j])) && !unicode.IsLetter(rune(word[j])) {
			j--
		}
		if word[i] != word[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func main() {
	fmt.Println("Enter the string to count the frequency:")
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	str = strings.TrimSpace(str)

	f := frequency(str)
	for char, freq := range f {
		fmt.Printf("Char: %q, Frequency: %v\n", char, freq)
	}

	fmt.Println("Enter the word to check for palindrome:")
	pal, _ := reader.ReadString('\n')
	pal = strings.TrimSpace(pal)

	ans := palindrome(pal)
	if ans {
		fmt.Println(pal, "is a palindrome")
	} else {
		fmt.Println(pal, "is not a palindrome")
	}
	fmt.Println(f)
}
