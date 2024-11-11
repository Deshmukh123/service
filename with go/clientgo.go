package main

import (
	"fmt"
	"sync"
	"time"
)

func uploadFile(filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Uploading %s\n", filename)
	time.Sleep(500 * time.Millisecond)
}

func main() {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}

	start := time.Now()

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go uploadFile(file, &wg)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("All files uploaded concurrently in %s.\n", elapsed)
}
