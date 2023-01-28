package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	//"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)
	func main() {
        username := flag.String("u","admin","input your username")
		password := flag.String("p","admin","input your password")
		repositoryurl := flag.String("r","http://nexus.z-bank.com","input your repository url")
		flag.Parse()
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
			return err
			}
			patters := "(.|/)+/\\.(.)*|(.|/)+/\\^archetype-catalog\\.xml(.)*|(.|/)+/\\^maven-metadata-local\\.xml|(.|/)+/\\^maven-metadata-deployment\\.xml|(.|/)*\\.sh"
	        //patters := "(.|/)+/\\.(.)*"
            matched,err := regexp.Match(patters,[]byte(path))
            if err != nil {
              fmt.Println(err)
			  return err
			}
			
			if ! info.IsDir() && ! matched {
            //fmt.Println(path)
             err := uploadfile(path,*username,*password,*repositoryurl)
			 if err != nil {
				fmt.Println(err)
				return err
			 }
			}
			
		
			return nil
		})

		if err != nil {
			fmt.Println(err)
		}
	}
	
func uploadfile(path string,username string,password string,repourl string) error {

    form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
		fw, err := writer.CreateFormFile("fileUploadName", path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fd, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer fd.Close()
		_, err = io.Copy(fw, fd)
		if err != nil {
			fmt.Println(err)
			return err
		}
	
		writer.Close()
	
		client := &http.Client{}
		url := repourl+path
		req, err := http.NewRequest("PUT", url, form)
		if err != nil {
			fmt.Println(err)
			return err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		//fmt.Println(req.Header)
		//fmt.Println(req.Method)
		req.SetBasicAuth(username, password)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer resp.Body.Close()
		
		//fmt.Println(resp.StatusCode)	
	return nil
}
