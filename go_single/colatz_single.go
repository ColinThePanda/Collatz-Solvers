package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	// Precomputed constants
	one := big.NewInt(1)
	two := big.NewInt(2)
	three := big.NewInt(3)
	zero := big.NewInt(0)

	// Starting number: 2^10000 - 1
	startingNum := new(big.Int).Sub(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(10000), nil),
		big.NewInt(1),
	)

	fmt.Println("Calculating 2^10,000 - 1...")

	num := new(big.Int).Set(startingNum)
	temp := new(big.Int) // temporary result reuse

	// Setup buffered file writing
	filename := fmt.Sprintf("collatz_conjecture_%s.txt", firstDigits(startingNum, 187))
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 1<<20) // 1MB buffer

	// Output starting number (if you still want to log it)
	// fmt.Println("Starting with: " + num.String())
	step := 0

	// Main Collatz Loop
	for num.Cmp(one) != 0 {
		// Check even
		if temp.Mod(num, two).Cmp(zero) == 0 {
			num.Div(num, two)
		} else {
			num.Mul(num, three)
			num.Add(num, one)
		}

		// Write to file
		_, err := writer.WriteString(num.String() + "\n")
		if err != nil {
			panic(err)
		}

		step++
		// Optional: Flush every 1000 steps for real-time file output
		if step%1000 == 0 {
			writer.Flush()
		}
	}

	// Final flush
	writer.Flush()

	// Print completion message
	fmt.Println("Done! Total steps:", step)
}

// firstDigits returns a string of first N digits of a big.Int
func firstDigits(n *big.Int, length int) string {
	s := n.String()
	if len(s) <= length {
		return s
	}
	return s[:length]
}
