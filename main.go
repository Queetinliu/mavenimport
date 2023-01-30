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
	"time"
)
	func main() {
        username := flag.String("u","admin","input your username")
		password := flag.String("p","admin","input your password")
		repositoryurl := flag.String("r","http://nexus.z-bank.com","input your repository url")
		flag.Parse()
		var wg sync.WaitGroup
		wg.Wait()
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
			    return err
			}
			patters := "(.|/)+/\\.(.)*|(.|/)+/\\^archetype-catalog\\.xml(.)*|(.|/)+/\\^maven-metadata-local\\.xml|(.|/)+/\\^maven-metadata-deployment\\.xml|(.|/)*\\.sh"
            matched,err := regexp.Match(patters,[]byte(path))
            if err != nil {
              fmt.Println(err)
			  return err
			}
			if ! info.IsDir() && ! matched {
				wg.Add(1)
			go func(file string)  {
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
			IdleConnTimeout:       15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			DisableKeepAlives:     false,
		}
		client := &http.Client{Transport: tr,}
		url := *repositoryurl+path
		req, err := http.NewRequest("PUT", url, form)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		//fmt.Println(req.Header)
		//fmt.Println(req.Method)
		req.SetBasicAuth(*username, *password)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
	    wg.Done()
		resp.Body.Close()
		
		fmt.Println(resp.StatusCode)	
			}(path)
	}
	return nil
})
if err != nil {
	fmt.Println(err)
}
	}
/*	
func uploadfile(ch chan string,username string,password string,repourl string) (errch chan error) {

    form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	path := <- ch 
	fw, err := writer.CreateFormFile("fileUploadName", path)
		if err != nil {
			fmt.Println(err)
			errch <- err
			return errch
		}
		fd, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			errch <- err
			return errch
		}
		defer fd.Close()
		_, err = io.Copy(fw, fd)
		if err != nil {
			fmt.Println(err)
			errch <- err
			return errch
		}
	
		writer.Close()
	
		client := &http.Client{}
		url := repourl+path
		req, err := http.NewRequest("PUT", url, form)
		if err != nil {
			fmt.Println(err)
			errch <- err
			return errch
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		//fmt.Println(req.Header)
		//fmt.Println(req.Method)
		req.SetBasicAuth(username, password)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			errch <- err
			return errch
		}
		defer resp.Body.Close()
		
		fmt.Println(resp.StatusCode)	
	return nil
}
*/