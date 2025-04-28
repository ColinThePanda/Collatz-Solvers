package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// CollatzResult holds the starting number and its sequence
type CollatzResult struct {
	startingNum uint64
	sequence    []uint64
	steps       int
	maxValue    uint64
}

// CollatzStats holds global statistics about the sequences
type CollatzStats struct {
	mutex              sync.Mutex
	longestSteps       int
	longestStepsNumber uint64
	largestValue       uint64
	largestValueNumber uint64
}

var stats CollatzStats

// Calculate the next number in the Collatz sequence
func collatz(num uint64) uint64 {
	if num%2 == 0 {
		// Even number: divide by 2
		return num / 2
	} else {
		// Odd number: 3n + 1
		return 3*num + 1
	}
}

// Process a single number through the Collatz sequence
func processCollatz(startingNum uint64, globalStats *CollatzStats, wg *sync.WaitGroup) {
	defer wg.Done()

	// Make a copy of the starting number to avoid modifications
	num := startingNum
	sequence := []uint64{num}
	maxValue := num

	// Generate the sequence until we reach 1
	for num != 1 {
		num = collatz(num)
		sequence = append(sequence, num)

		// Track the maximum value in this sequence
		if num > maxValue {
			maxValue = num
		}

		// Print progress periodically (every 10000 steps to reduce output)
		if len(sequence)%10000 == 0 {
			fmt.Printf("Number %d: %d steps so far\n", startingNum, len(sequence))
		}
	}

	steps := len(sequence) - 1 // Subtract 1 because we count transitions, not elements
	fmt.Printf("Completed sequence for %d with %d steps\n", startingNum, steps)

	// Update global statistics
	globalStats.mutex.Lock()
	if steps > globalStats.longestSteps {
		globalStats.longestSteps = steps
		globalStats.longestStepsNumber = startingNum
	}
	if maxValue > globalStats.largestValue {
		globalStats.largestValue = maxValue
		globalStats.largestValueNumber = startingNum
	}
	globalStats.mutex.Unlock()

	// Write the sequence to a file
	filename := fmt.Sprintf("collatz_%d.txt", startingNum)
	err := writeSequenceToFile(filename, sequence)
	if err != nil {
		fmt.Printf("Error writing to file for %d: %v\n", startingNum, err)
	}
}

// Write a sequence to a file
func writeSequenceToFile(filename string, sequence []uint64) error {
	// Create the output directory if it doesn't exist
	err := os.MkdirAll("collatz_results", 0755)
	if err != nil {
		return err
	}

	// Create the file
	filePath := filepath.Join("collatz_results", filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each number in the sequence to the file
	for _, num := range sequence {
		_, err := file.WriteString(strconv.FormatUint(num, 10) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Process numbers in a range
func processRange(start, end uint64, workerID int, globalStats *CollatzStats, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d processing range %d to %d\n", workerID, start, end)

	for i := start; i <= end; i++ {
		// Create a new wait group for each number
		var numWg sync.WaitGroup
		numWg.Add(1)
		processCollatz(i, globalStats, &numWg)
		numWg.Wait()

		// Print progress every 1000 numbers
		if i%1000 == 0 {
			fmt.Printf("Worker %d completed up to %d\n", workerID, i)
		}
	}

	fmt.Printf("Worker %d completed range %d to %d\n", workerID, start, end)
}

func main() {
	startTime := time.Now()

	// Define the range of numbers to process (1 to 1,000,000)
	const maxNumber uint64 = 100000

	// Initialize statistics
	stats = CollatzStats{
		longestSteps:       0,
		longestStepsNumber: 0,
		largestValue:       0,
		largestValueNumber: 0,
	}

	// Set GOMAXPROCS to use all available CPU cores
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Printf("Running with %d CPU cores\n", numCPU)

	// Divide the work among workers based on CPU cores
	var wg sync.WaitGroup
	numWorkers := numCPU
	numbersPerWorker := maxNumber / uint64(numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		// Calculate the range for this worker
		start := uint64(i)*numbersPerWorker + 1
		end := uint64(i+1) * numbersPerWorker

		// Adjust the last worker to include any remaining numbers
		if i == numWorkers-1 {
			end = maxNumber
		}

		// Start the worker
		go processRange(start, end, i+1, &stats, &wg)
	}

	// Wait for all workers to complete
	fmt.Println("Waiting for all sequences to complete...")
	wg.Wait()

	elapsed := time.Since(startTime)
	fmt.Printf("All sequences completed in %s\n", elapsed)

	// Print final statistics
	fmt.Println("\n=== COLLATZ CONJECTURE STATISTICS ===")
	fmt.Printf("Number with the longest sequence: %d (took %d steps)\n",
		stats.longestStepsNumber, stats.longestSteps)
	fmt.Printf("Largest value reached: %d (from starting number %d)\n",
		stats.largestValue, stats.largestValueNumber)
}
