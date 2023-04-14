package persistence

import (

)



// Controllers struct which contains
// an array of controllers
type ControllerRouteList struct {
	ConrollerRouteList []ControllerRouteDetails `json:"controllerList"`
}


// Controller struct which contains a name
// a type and a list of social links
type ControllerRouteDetails struct {
	Name   string `json:"name"'`
	UrlKey   string `json:"urlKey"`
	Description	string `json:"description"`
	Actionset []ActionRouteDetails `json:"actionSet"`
}

type ActionRouteDetails struct {
	ActionName string `json:"actionName"`
	ActionUrlKey string `json:"actionUrlKey"`
}


