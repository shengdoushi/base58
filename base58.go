package base58

import (
	"errors"
	"fmt"
)

var (
	Error_InvalidBase58String error = errors.New("invalid base58 string")

	// Alphabet: copy from https://en.wikipedia.org/wiki/Base58
	BitcoinAlphabet = NewAlphabet("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	IPFSAlphabet = NewAlphabet("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	FlickrAlphabet = NewAlphabet("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
	RippleAlphabet = NewAlphabet("rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz")
)

type Alphabet struct {
	encodeTable [58]rune
}

func NewAlphabet(alphabet string)*Alphabet{
	alphabetRunes := []rune(alphabet)
	if len(alphabetRunes) != 58 {
		panic(fmt.Sprintf("Base58 Alphabet length must 58, but %d", len(alphabetRunes)))
	}

	ret := new(Alphabet)
	for idx := 0; idx < 58; idx++ {
		ret.encodeTable[idx] = alphabetRunes[idx]
	}
	return ret
}

// encode with custom alphabet
func Encode(input []byte, alphabet *Alphabet)string{
	// prefix 0
	inputLength := len(input)
	prefixZeroes := 0
	for prefixZeroes < inputLength && input[prefixZeroes] == 0 {
		prefixZeroes++
	}

	capacity := len(input) * 138 / 100 + 1 // log256 / log58
	output := make([]byte, capacity)
	outputReverseEnd := capacity-1

	for inputPos := prefixZeroes; inputPos < inputLength; inputPos++{
		carry := uint32(input[inputPos])

		outputIdx := capacity-1
		for ; carry != 0 || outputIdx > outputReverseEnd ; outputIdx-- {
			carry += 256 * uint32(output[outputIdx])
			output[outputIdx] = byte(carry % 58)
			carry /= 58
		}
		outputReverseEnd= outputIdx
	}

	retStrRunes := make([]rune, prefixZeroes + (capacity-1-outputReverseEnd))
	for i := 0; i < prefixZeroes; i++ {
		retStrRunes[i] = alphabet.encodeTable[0]
	}
	for i := outputReverseEnd+1; i < capacity; i++{
		retStrRunes[prefixZeroes+i-(outputReverseEnd+1)] = alphabet.encodeTable[output[i]]
	}
	return string(retStrRunes)	
}

// Decode with custom alphabet
func Decode(input string, alphabet *Alphabet)([]byte, error){
	inputBytes := []rune(input)
	// pass first char
	capacity := len(inputBytes) * 733 /1000 + 1; // log(58) / log(256)
	output := make([]byte, capacity)
	outputReverseEnd := capacity-1

	// prefix 0
	zero58Byte := alphabet.encodeTable[0]
	inputLength := len(inputBytes)
	prefixZeroes := 0
	for prefixZeroes < inputLength && inputBytes[prefixZeroes] == zero58Byte {
		prefixZeroes++
	}

	for inputPos := prefixZeroes; inputPos < inputLength; inputPos++{
		carry := -1
		for i := 0; i < 58; i++ {
			if alphabet.encodeTable[i] == inputBytes[inputPos] {
				carry = i
				break
			}
		}
		if carry == -1 {
			return nil, Error_InvalidBase58String
		}

		outputIdx := capacity-1
		for ; carry != 0 || outputIdx > outputReverseEnd ; outputIdx-- {
			carry += 58 * int(output[outputIdx])
			output[outputIdx] = byte(carry % 256)
			carry /= 256
		}
		outputReverseEnd= outputIdx
	}

	retBytes := make([]byte, prefixZeroes + (capacity-1-outputReverseEnd))
	for i := outputReverseEnd+1; i < capacity; i++{
		retBytes[prefixZeroes+i-(outputReverseEnd+1)] = output[i]
	}
	return retBytes, nil
}
