package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/oddimportance/zsmvctool/persistence"
)

type ManageUploadedFiles struct {
	WaitGroup *sync.WaitGroup
}

func (m *ManageUploadedFiles) MoveUploadedFile(uploadedFile multipart.File, handle *multipart.FileHeader, destination, renameTo string) bool {
	return m.saveFile(uploadedFile, handle, destination, renameTo)
}

func (m *ManageUploadedFiles) saveFile(file multipart.File, handle *multipart.FileHeader, destination, renameTo string) bool {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("%v", err)
		return false
	}

	fileName := renameTo

	if fileName == "" {
		fileName = handle.Filename
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", destination, fileName), data, 0666)
	if err != nil {
		fmt.Printf("%v", err)
		return false
	}

	m.WaitGroup.Done()

	return true

}

func (m *ManageUploadedFiles) DeleteFile(path string) {

	// make sure the file to delete exists
	// befor the remove func is exectued
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	} else {
		err := os.Remove(path)
		if err != nil {
			fmt.Println(err)
			//os.Exit(1)
			panic(err)
		}
	}

}

// move or rename file, basically
// executes the mv (linux move)
// command
func (m *ManageUploadedFiles) MoveFile(source, destination string) {
	// move file
	var err = os.Rename(source, destination)
	if err != nil {
		fmt.Printf("%v", err)
	}

	if m.WaitGroup != nil {
		m.WaitGroup.Done()
	}

}

func (m *ManageUploadedFiles) DeleteFileFromFileServeCloud(cloudDetails persistence.CloudToken) bool {
	fmt.Println(cloudDetails)
	return m.makeHttpRequestToCloud(cloudDetails, bytes.NewReader(nil))
}

func (m *ManageUploadedFiles) UploadFileToFileServeCloud(cloudDetails persistence.CloudToken, fileName, filePath string) bool {
	reader, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", filePath, fileName))
	if err != nil {
		log.Println(err)
		return false
	}
	// convert byte slice to io.Reader
	data := bytes.NewReader(reader)
	return m.makeHttpRequestToCloud(cloudDetails, data)
}

func (m *ManageUploadedFiles) makeHttpRequestToCloud(cloudDetails persistence.CloudToken, data *bytes.Reader) bool {
	// client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, cloudDetails.CloudUrl, data)
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Set("file", cloudDetails.FileName)
	req.Header.Set("file-size", cloudDetails.FileSize)
	// fmt.Println(req.Header)
	// resp, err := client.Do(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	// return resp.Status, content
	if resp.Status == "200 OK" {
		var responseText = map[string]string{}
		json.Unmarshal([]byte(string(content)), &responseText)
		fmt.Println(responseText)
		if responseText["success"] == "1" {
			return true
		}
	}
	return false
}
