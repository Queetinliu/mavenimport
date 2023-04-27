package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	//"time"
)

var wg sync.WaitGroup
var client *http.Client
func init() {
	     tr := &http.Transport{
	     	DisableKeepAlives:     false,
	     	DisableCompression:    false,
	     	MaxIdleConns:          0,
	     	MaxIdleConnsPerHost:   1024,
	     	MaxConnsPerHost:       0,
	     	IdleConnTimeout:       0,
	     	ResponseHeaderTimeout: 0,
	     	ExpectContinueTimeout: 0,
	     	MaxResponseHeaderBytes: 0,
	     	WriteBufferSize:        0,
	     	ReadBufferSize:         0,
	     	ForceAttemptHTTP2:      false,
	     }
		 client = &http.Client{Transport: tr}
		}
func main() {
	username := flag.String("u", "admin", "input your username")
	password := flag.String("p", "admin", "input your password")
	repositoryurl := flag.String("r", "http://nexus.com", "input your repository url")
	flag.Parse()
    const maxConcurrent = 600
	filech := make(chan string,1)
	wg.Add(1)
	go uploadfile(*username, *password, *repositoryurl, filech,maxConcurrent)
	wg.Add(1)
	go findfile(filech)
	wg.Wait()

}

func findfile(ch chan string) {
	defer wg.Done()
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		patters := "(.|/)+/\\.(.)*|(.|/)+/\\^archetype-catalog\\.xml(.)*|(.|/)+/\\^maven-metadata-local\\.xml|(.|/)+/\\^maven-metadata-deployment\\.xml|(.|/)*\\.sh"
		matched, err := regexp.Match(patters, []byte(path))
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() && !matched {
			ch <- path
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	close(ch)
}

func uploadfile(username, password, repositoryurl string, ch chan string,maxConcurrent int) {
	defer wg.Done()
	var innerwg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)
	for path := range ch {
		sem <- struct{}{}
		innerwg.Add(1)
		go func(f string) {
			defer innerwg.Done()
			defer func() { <-sem }()
			form := new(bytes.Buffer)
			writer := multipart.NewWriter(form)
			fw, err := writer.CreateFormFile("fileUploadName", f)
			if err != nil {
				fmt.Println(err)
			}
			fd, err := os.Open(f)
			if err != nil {
				fmt.Println(err)
			}

			_, err = io.Copy(fw, fd)
			if err != nil {
				fmt.Println(err)
			}

			writer.Close()
			/*
			tr := &http.Transport{
				MaxIdleConns: 0,
				IdleConnTimeout:       180 * time.Second,
				ResponseHeaderTimeout: 60 * time.Second,
				DisableKeepAlives:     false,
			}
			client := &http.Client{Transport: tr}
			*/
			url := repositoryurl + f
			req, err := http.NewRequest("PUT", url, form)
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.SetBasicAuth(username, password)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
            io.Copy(io.Discard, resp.Body);
			defer resp.Body.Close()
		}(path)
	}
	innerwg.Wait()
}
