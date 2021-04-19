package main

import (
	"errors"
	"strconv"
)

// ComputeResult stores the result on palindrome calculation and the time it took
type ComputeResult struct {
	number int
	binary string
	time   int
}

// MinerSingle executes the palindrome computations one at a time
func MinerSingle(number int) (ComputeResult, error) {
	// output := make(map[int]string)
	var output ComputeResult

	if isPalindrome(number) {
		if isBinaryPalindrome(number) {
			output.number = number
			output.binary = convertToBinary(number)
			return output, nil
		}
	}
	return output, errors.New("Not a palindrome")
}

func isPalindrome(number int) bool {
	// converts int into string and then check if the string is a palindrome
	forwardString := strconv.Itoa(number)
	reversedString := reverse(forwardString)

	return forwardString == reversedString
}

func isBinaryPalindrome(number int) bool {
	// converts int into a string of binary and then check if the string is a palindrome
	forwardBinaryString := convertToBinary(number)
	reversedBinaryString := reverse(forwardBinaryString)

	return forwardBinaryString == reversedBinaryString
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func convertToBinary(i int) string {
	i64 := int64(i)
	return strconv.FormatInt(i64, 2) // base 2 for binary
}