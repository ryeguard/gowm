package internal

import "os"

func LoadEnvVar() (string, bool) {
	for _, key := range []string{"OWM_APP_ID", "OWM_API_KEY"} {
		appId, ok := os.LookupEnv(key)
		if ok {
			return appId, true
		}
	}
	return "", false
}
