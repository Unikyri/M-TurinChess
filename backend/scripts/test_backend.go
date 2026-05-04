package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	fileDir := "../docs/PGN/Example01.pgn"
	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, _ := writer.CreateFormFile("pgn", "Example01.pgn")
	file, _ := os.Open(fileDir)
	io.Copy(part, file)
	file.Close()

	writer.WriteField("elo", "1500")
	writer.Close()

	req, _ := http.NewRequest("POST", "http://localhost:8080/api/analyze", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Println("Body length:", len(respBody))
}
