package core

import "testing"
import "fmt"
import "regexp"
import "errors"

func TestNewSubdomainExtractor(t *testing.T) {
	_, err := NewSubdomainExtractor("google.com")
	if err != nil {
		t.Error(err)
	}
}

func ExampleNewSubdomainExtractor() {
	exampleText := `<a href="https://subdomain.google.com">`

	extractor, err := NewSubdomainExtractor("google.com")
	if err != nil {
		panic(err)
	}

	subdomain := extractor.FindString(exampleText)

	fmt.Println(subdomain)
	// Output: subdomain.google.com
}

func BenchmarkSubdomainExtractorOld(b *testing.B) {
	domain := "google.com"

	var result string

	extractor, err := regexp.Compile(`[a-zA-Z0-9\*_.-]+\.` + domain)
	if err != nil {
		b.Error(err)
	}

	for n := 0; n < b.N; n++ {
		exampleText := `<a href="https://subdomain.google.com">`

		result = extractor.FindString(exampleText)
	}

	b.Log(result)
}

func BenchmarkSubdomainExtractorNew(b *testing.B) {
	domain := "google.com"

	var result string

	extractor, err := regexp.Compile(`[\w-\*]+\.` + domain)
	if err != nil {
		b.Error(err)
	}

	exampleText := `<a href="https://subdomain.google.com">`

	for n := 0; n < b.N; n++ {
		result = extractor.FindString(exampleText)
	}

	b.Log(result)
}

// below is an expieriment in removing regex

var zeroStr string
var validURLChars = []byte("abcdefghijklmnopqrstuvwxyz1234567890.-*ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var dotChar = []byte(".")[0]
var starChar = []byte("*")[0]

func reverseBytes(input []byte) []byte {
	// i starts at index 0
	// j starts at the len of the input, minus 1
	for i, j := 0, len(input)-1;
	// while i is less than j
	i < j;
	// add to i, minus from j
	i, j = i+1, j-1 {
		//fmt.Println(input[i])
		input[i], input[j] = input[j], input[i]
	}
	return input
}

var zeroByte byte
var errNoMoreBytes = errors.New("no more bytes")

func readReverseBytes(input []byte) func() (byte, error) {
	i := len(input)

	return func() (byte, error) {
		i--
		if i >= 0 {
			return input[i], nil
		}
		return zeroByte, errNoMoreBytes
	}
}

func singleSubdomainExtractor(domain string) func([]byte) string {
	domain = "." + domain
	lastByteInDomain := domain[len(domain)-1]

	return func(input []byte) string {
		foundSomethingInteresting := false
		indexValue := len(domain) - 1
		foundBuffer := []byte{}

		reader := readReverseBytes(input)

		for {
			nextByte, err := reader()
			if err != nil {
				if len(foundBuffer) >= len(domain) {
					return string(reverseBytes(foundBuffer))
				}
				return zeroStr
			}

			// wait until we found the first interesting byte,
			// which is the last byte in the domain
			if !foundSomethingInteresting {
				if nextByte == lastByteInDomain {
					foundSomethingInteresting = true
					foundBuffer = append(foundBuffer, nextByte)
					indexValue--
					continue
				}
			}

			if len(foundBuffer) < len(domain) && nextByte == domain[indexValue] {
				foundBuffer = append(foundBuffer, nextByte)
				indexValue--
				continue
			}

			if len(foundBuffer) >= len(domain) {
				foundNext := false
				for _, v := range validURLChars {
					if v == nextByte {
						if v == dotChar {
							lastIn := foundBuffer[len(foundBuffer)-1]
							if v == lastIn {
								foundNext = false
								// remove the last dot, since it was garbage
								foundBuffer = foundBuffer[:len(foundBuffer)-1]
								// remove shortfind
								if len(foundBuffer) == len(domain)-1 {
									foundBuffer = []byte{}
								}
								break
							}
							if lastIn == starChar {
								foundNext = false
								break
							}
						} else {
							foundNext = true
							foundBuffer = append(foundBuffer, nextByte)
							indexValue++
							break
						}
					}
				}
				if foundNext {
					continue
				} else {
					break
				}
			}
		}
		return string(reverseBytes(foundBuffer))
	}
}

func multiSubdomainExtractor(domain string) func([]byte) []string {
	domain = "." + domain
	lastByteInDomain := domain[len(domain)-1]

	return func(input []byte) (results []string) {
		foundSomethingInteresting := false
		indexValue := len(domain) - 1
		foundBuffer := []byte{}

		reader := readReverseBytes(input)

		for {
			nextByte, err := reader()
			if err != nil {
				return results
			}

			// wait until we found the first interesting byte,
			// which is the last byte in the domain
			if !foundSomethingInteresting {
				if nextByte == lastByteInDomain {
					foundSomethingInteresting = true
					foundBuffer = append(foundBuffer, nextByte)
					indexValue--
					continue
				}
			}

			if len(foundBuffer) < len(domain) && nextByte == domain[indexValue] {
				foundBuffer = append(foundBuffer, nextByte)
				indexValue--
				continue
			}

			if len(foundBuffer) >= len(domain) {
				foundNext := false
				for _, v := range validURLChars {
					if v == nextByte {
						if v == dotChar && v == foundBuffer[len(foundBuffer)-1] {
							foundNext = false
							// remove the last dot, since it was garbage
							foundBuffer = foundBuffer[:len(foundBuffer)-1]
							// remove shortfind
							if len(foundBuffer) == len(domain)-1 {
								foundBuffer = []byte{}
							}
							break
						} else if v == dotChar && foundBuffer[len(foundBuffer)-1] == starChar {
							foundNext = false
							break
						} else {
							foundNext = true
							foundBuffer = append(foundBuffer, nextByte)
							break
						}
					}
				}
				if foundNext {
					continue
				} else {
					if len(foundBuffer) >= len(domain) {
						results = append(results, string(reverseBytes(foundBuffer)))
					}
					foundSomethingInteresting = false
					foundBuffer = []byte{}
					indexValue = len(domain) - 1
					continue
				}
			}
			continue
		}
	}
}

func BenchmarkSubdomainExtractorCustom(b *testing.B) {

	domain := "google.com"

	var result string

	exampleText := []byte(`<a href="https://subdomain.google.com">`)

	extractor := singleSubdomainExtractor(domain)

	for n := 0; n < b.N; n++ {
		result = extractor(exampleText)
	}

	b.Log(result)
}

func BenchmarkSubdomainExtractorCustomMulti(b *testing.B) {

	domain := "google.com"

	var results []string

	exampleText := []byte(`<a href="https://subdomain.google.com">`)

	extractor := multiSubdomainExtractor(domain)

	for n := 0; n < b.N; n++ {
		results = extractor(exampleText)
	}

	b.Log(results)
}
