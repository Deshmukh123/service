package main

import (
	"fmt"
	"time"
)

func uploadFile(filename string) {
	fmt.Printf("Uploading %s\n", filename)
	time.Sleep(500 * time.Millisecond)
}

// main uploads files sequentially and measures the total time taken.
func main() {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}

	start := time.Now()

	for _, file := range files {
		uploadFile(file)
	}

	elapsed := time.Since(start)
	fmt.Printf("All files uploaded sequentially in %s.\n", elapsed)
}
