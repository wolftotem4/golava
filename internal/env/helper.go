package env

import (
	"fmt"
	"os"
	"regexp"
)

const DotEnvFile = ".env"

func SetKeyInEnvironmentFile(key, value string) error {
	content, err := os.ReadFile(DotEnvFile)
	if err != nil {
		return err
	}

	content, appends := SetEnvVar(content, key, value)
	if appends {
		f, err := os.OpenFile(DotEnvFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(content)
		return err
	} else {
		return os.WriteFile(DotEnvFile, content, 0644)
	}
}

func SetEnvVar(content []byte, key, value string) (newcontent []byte, appends bool) {
	env := string(content)
	old := os.Getenv(key)

	r := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*%s\s*=\s*%s.*`, regexp.QuoteMeta(key), regexp.QuoteMeta(old)))

	if r.MatchString(env) {
		newEnv := r.ReplaceAllString(env, fmt.Sprintf("%s=%s", key, value))
		return []byte(newEnv), false
	} else {
		return []byte(fmt.Sprintf("\n%s=%s\n", key, value)), true
	}
}
