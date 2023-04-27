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

func main() {
	username := flag.String("u", "admin", "input your username")
	password := flag.String("p", "admin", "input your password")
	repositoryurl := flag.String("r", "http://nexus.z-bank.com", "input your repository url")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Wait()
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
			wg.Add(1)
			go func(file string) {
				form := new(bytes.Buffer)
				writer := multipart.NewWriter(form)
				fw, err := writer.CreateFormFile("fileUploadName", path)
				if err != nil {
					fmt.Println(err)
				}
				fd, err := os.Open(path)
				if err != nil {
					fmt.Println(err)
				}
				defer fd.Close()
				_, err = io.Copy(fw, fd)
				if err != nil {
					fmt.Println(err)
				}

				writer.Close()
				tr := &http.Transport{
					MaxIdleConns:          10,
					//IdleConnTimeout:       60 * time.Second,
					IdleConnTimeout:  0,
					//ResponseHeaderTimeout: 60 * time.Second,
					ResponseHeaderTimeout: 0,
					DisableKeepAlives:     false,
				}
				client := &http.Client{Transport: tr}
				url := *repositoryurl + path
				req, err := http.NewRequest("PUT", url, form)
				if err != nil {
					fmt.Println(err)
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
				req.SetBasicAuth(*username, *password)
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
				}
				wg.Done()
				if _, err := io.Copy(io.Discard, resp.Body); err != nil {
					fmt.Println(err)
				}
				resp.Body.Close()
			}(path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
