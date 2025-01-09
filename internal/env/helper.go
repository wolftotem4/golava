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

	content, appends := SetEnvVar(content, key, value, os.Getenv(key))
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

func SetEnvVar(content []byte, key, value, old string) (newcontent []byte, appends bool) {
	const space = `[\t\f\v ]*`
	const unit = space + "%s" + space
	const regEx = `(?m)^` + unit + `=` + unit + `$`

	env := string(content)

	r := regexp.MustCompile(fmt.Sprintf(regEx, regexp.QuoteMeta(key), regexp.QuoteMeta(old)))

	if r.MatchString(env) {
		newEnv := r.ReplaceAllString(env, fmt.Sprintf("%s=%s", key, value))
		return []byte(newEnv), false
	} else {
		return []byte(fmt.Sprintf("\n%s=%s\n", key, value)), true
	}
}
