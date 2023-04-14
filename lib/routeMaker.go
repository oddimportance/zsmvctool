package lib

import (
	"fmt"

	//packages to read json
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	persistence "github.com/oddimportance/zsmvctool/persistence"
)

type RouteMaker struct {
	getVars map[string]string

	_handleFilePath HandleFilePath

	controllerRouteList persistence.ControllerRouteList
}

func (r *RouteMaker) MakeRoutes(_getVars map[string]string, handleFilePath HandleFilePath) {

	r.getVars = _getVars

	r._handleFilePath = handleFilePath

	r.readJson()

	//	fmt.Printf("Ya Rab %v", r.controllerRouteList)

	//	fmt.Printf("Ya Rab %v", r.getVars)

}

/**
 *
 *
 */
func (r *RouteMaker) checkIfRouteIsValid() bool {

	var isRouteValid = true

	for i := range r.controllerRouteList.ConrollerRouteList {
		fmt.Printf("Controller : %s\n", r.controllerRouteList.ConrollerRouteList[i].Name)
		if r.controllerRouteList.ConrollerRouteList[i].UrlKey == r.getVars["controller"] {
			for j := range r.controllerRouteList.ConrollerRouteList[i].Actionset {
				fmt.Printf("Action key %s\n", r.controllerRouteList.ConrollerRouteList[i].Actionset[j].ActionUrlKey)
			}
		}
	}

	return isRouteValid

}

func (r *RouteMaker) readJson() {

	// Open our jsonFile
	jsonFile, err := os.Open(strings.Join([]string{r._handleFilePath.GetPrivatePath(), "module_settings.json"}, "/"))
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	//	fmt.Println("Successfully Opened controllers.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	//var controllerRouteList persistence.ControllerRouteList

	json.Unmarshal([]byte(byteValue), &r.controllerRouteList)

	//	fmt.Printf("%v\n", controllerRouteList)

}
