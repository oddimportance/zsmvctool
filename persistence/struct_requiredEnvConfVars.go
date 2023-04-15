package persistence

type RequiredEnvConfVars []struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}
