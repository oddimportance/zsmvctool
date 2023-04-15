package lib

import (
	"encoding/json"
	"fmt"

	"github.com/oddimportance/zsmvctool/persistence"
)

type GetEnvConfigVars struct {
	envVars                    persistence.EnvConfigVars
	sensibleConfigSecurityVars []string
}

func (g *GetEnvConfigVars) Initiate(envConfigFilePath string) persistence.EnvConfigVars {

	var _readFile = new(ReadFile)
	var fileContent = _readFile.Initiate(envConfigFilePath)

	g.unmarshalToJson(fileContent)

	// fmt.Println("App Version", g.envVars.EnvConfigVarList["version"])

	//fmt.Println()

	return g.envVars

}

func (g *GetEnvConfigVars) unmarshalToJson(fileContent []byte) {

	json.Unmarshal([]byte(fileContent), &g.envVars)

	//fmt.Println(env)

}

// Set a list of sensible vars which may leak
// important information, like passwords, port
// and so on
func (g *GetEnvConfigVars) SetSensibleConfigSecurityVars(sensibleList []string) {
	g.sensibleConfigSecurityVars = sensibleList
}

func (g *GetEnvConfigVars) GetConfVar(key string) string {
	return g.envVars.EnvConfigVarList[key]
}

func (g *GetEnvConfigVars) SetConfVar(key, value string) {
	g.envVars.EnvConfigVarList[key] = value
}

func (g *GetEnvConfigVars) GetAllConfVarsForView() map[string]string {
	return g.filterSensibleVars()
}

func (g *GetEnvConfigVars) filterSensibleVars() map[string]string {
	var filteredConfVars map[string]string = g.envVars.EnvConfigVarList

	for _, sensibleVar := range g.sensibleConfigSecurityVars {
		delete(filteredConfVars, sensibleVar)
	}

	return filteredConfVars
}

func (g *GetEnvConfigVars) PrintAllConfigVars() {

	fmt.Println()
	fmt.Println()
	fmt.Println("+++ Available General Configuration Vars +++")
	fmt.Println("===================================================")

	for key, configVar := range g.envVars.EnvConfigVarList {
		fmt.Printf("- %s : %s\n", key, configVar)
	}

	fmt.Println("===================================================")
	fmt.Println()

}
