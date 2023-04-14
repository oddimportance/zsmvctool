package persistence

import ()

type DockerRequiredEnvVars []struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}
