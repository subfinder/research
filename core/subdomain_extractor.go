package core

import "errors"

// reverseBytes mutates the given slice of bytes, reversing
// the order of the slice
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

//var validURLChars = []byte(".abcdefghijklmnopqrstuvwxyz1234567890-ABCDEFGHIJKLMNOPQRSTUVWXYZ*")
var validURLChars = []byte(".abcdefghijklmnopqrstuvwxyz1234567890-*")
var zeroStr string
var dotChar = []byte(".")[0]
var starChar = []byte("*")[0]
var percChar = []byte("%")[0]

// NewSingleSubdomainExtractor creates a new extractor that looks for
// only one subdomain.
func NewSingleSubdomainExtractor(domain string) func([]byte) string {
	domain = "." + domain
	domainLen := len(domain)
	domainLenMinusOne := domainLen - 1
	lastByteInDomain := domain[domainLenMinusOne]

	return func(input []byte) string {
		if len(input) <= domainLen {
			return zeroStr
		}

		// scoped variables
		foundSomethingInteresting := false
		indexValue := domainLenMinusOne
		foundBuffer := []byte{}

		reader := readReverseBytes(input)

		for {
			nextByte, err := reader()
			if err != nil {
				if len(foundBuffer) > domainLen {
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

			if len(foundBuffer) < domainLen && nextByte == domain[indexValue] {
				foundBuffer = append(foundBuffer, nextByte)
				indexValue--
				continue
			}

			if len(foundBuffer) >= domainLen {
				foundNext := false
				for _, v := range validURLChars {
					if v == nextByte {
						if v == dotChar && v == foundBuffer[len(foundBuffer)-1] {
							foundNext = false
							// remove the last dot, since it was garbage
							foundBuffer = foundBuffer[:len(foundBuffer)-1]
							// remove shortfind
							if len(foundBuffer) == domainLen-1 {
								foundBuffer = []byte{}
							}
							break
						} else if v == dotChar && foundBuffer[len(foundBuffer)-1] == starChar {
							foundNext = false
							break
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
		if len(foundBuffer) == domainLen {
			return zeroStr
		}
		return string(reverseBytes(foundBuffer))
	}
}

// NewMultiSubdomainExtractor creates a new extractor that looks for
// as many subdomains as it can find.
func NewMultiSubdomainExtractor(domain string) func([]byte) []string {
	domain = "." + domain
	domainLen := len(domain)
	domainLenMinusOne := domainLen - 1
	lastByteInDomain := domain[domainLenMinusOne]

	return func(input []byte) (results []string) {
		if len(input) <= domainLen {
			return nil
		}

		// scoped variables
		foundSomethingInteresting := false
		indexValue := domainLenMinusOne
		foundBuffer := []byte{}

		reader := readReverseBytes(input)

		for {
			nextByte, err := reader()
			if err != nil {
				if len(foundBuffer) > domainLen {
					results = append(results, string(reverseBytes(foundBuffer)))
				}
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

			if len(foundBuffer) < domainLen && nextByte == domain[indexValue] {
				foundBuffer = append(foundBuffer, nextByte)
				indexValue--
				continue
			}

			if len(foundBuffer) >= domainLen {
				foundNext := false
				for _, v := range validURLChars {
					if v == nextByte {
						if v == dotChar && v == foundBuffer[len(foundBuffer)-1] {
							foundNext = false
							// remove the last dot, since it was garbage
							foundBuffer = foundBuffer[:len(foundBuffer)-1]
							// remove shortfind
							if len(foundBuffer) == domainLen-1 {
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
					if len(foundBuffer) > domainLen {
						results = append(results, string(reverseBytes(foundBuffer)))
					}
					foundSomethingInteresting = false
					foundBuffer = []byte{}
					indexValue = domainLenMinusOne
					continue
				}
			}
			continue
		}
	}
}
