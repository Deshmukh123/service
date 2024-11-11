package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// Define a constant for the maximum number of concurrent uploads
const maxConcurrentUploads = 3

// Buffered channel to limit concurrent uploads
var uploadLimiter = make(chan struct{}, maxConcurrentUploads)

// Track upload progress
var uploadProgress = make(map[string]int)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Restrict the request method to POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Get the file from the request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Limit the number of concurrent uploads
	uploadLimiter <- struct{}{}
	defer func() { <-uploadLimiter }()

	// Track the file name for progress tracking
	fileName := fileHeader.Filename
	uploadProgress[fileName] = 0

	// Create the file on the server
	out, err := os.Create("./uploads/" + fileName)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Copy the file data to the server's file and track the progress
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		buf := make([]byte, 1024) // Buffer for file reading
		totalBytes := 0
		fileSize := fileHeader.Size

		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Error reading file:", err)
				return
			}

			totalBytes += n
			out.Write(buf[:n])

			// Update progress (as a percentage)
			progress := int(float64(totalBytes) / float64(fileSize) * 100)
			uploadProgress[fileName] = progress
			log.Printf("Uploading %s: %d%% complete\n", fileName, progress)
		}
	}()

	wg.Wait()
	fmt.Fprintf(w, "File upload complete: %s\n", fileName)
}
