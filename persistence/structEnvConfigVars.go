package persistence

import ()

// EnvConfigVars struct which contains
// an array of coniguration vars
type EnvConfigVars struct {
	EnvConfigVarList ConfDetails `json:"envConfigVars"`
}

// Controller struct which contains a name
// a type and a list of social links
type ConfDetails map[string]string
