package interpolator

import (
	"github.com/pkg/errors"
	"github.com/rills-ai/Hachi/pkg/helper"
	"os"
	"regexp"
	"strings"
)

var InterpolationRegex = regexp.MustCompile("{{\\.((local|remote|route|resolver)::(.*?))}}")

// InterpolateStrings  we currently support interpolation from envars and Hachi stanza vars
func InterpolateFromValues(stanzaVars map[string]string, content string) (string, error) {
	//stanza vars override envars values
	lstanzaVars := helper.MapKeys[string, string](stanzaVars, strings.ToLower)

	//index envars
	envars := make(map[string]string)
	for _, e := range os.Environ() {
		before, after, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		envars[strings.ToLower(before)] = after
	}

	matches := InterpolationRegex.FindAllString(content, -1)
	for _, v := range matches {
		interpolatedPlaceholder := v
		instructions := InterpolationRegex.FindStringSubmatch(v)
		instruct := instructions[2]
		key := instructions[3]
		//todo: add resolver implementation in the future
		if instruct != "local" {
			continue
		}
		interpolatedValue := envars[key]
		if val, ok := lstanzaVars[key]; ok {
			content = strings.Replace(content, interpolatedPlaceholder, val, -1)
		} else {
			content = strings.Replace(content, interpolatedPlaceholder, interpolatedValue, -1)
			if interpolatedValue == "" {
				return "", errors.New("ERROR: key " + interpolatedPlaceholder + " value is missing, check your configuration file and machine ENVARS")
			}
		}
	}
	return content, nil
}
