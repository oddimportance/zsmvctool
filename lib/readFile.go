package lib

import (
	"fmt"

	//packages to read json
	//	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ReadFile struct {
}

func (r *ReadFile) Initiate(filePath string) []byte {

	var fileToOpen, _ = filepath.Abs(filePath)

	// Open our jsonFile
	file, err := os.Open(fileToOpen)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	//	fmt.Println("Successfully Opened controllers.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer file.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(file)

	//fmt.Println(string(byteValue))

	//return string(byteValue)
	return byteValue

}

// Check if the give file exists
// @param pathToFile string
// @return bool
func (r *ReadFile) FileExists(pathToFile string) bool {
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		// file does not exist
		return false
	}
	return true
}
