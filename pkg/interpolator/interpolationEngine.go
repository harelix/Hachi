package interpolator

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rills-ai/Hachi/pkg/helper"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
)

var InterpolationRegex = regexp.MustCompile("{{\\.((local|remote|route|resolver)::(.*?))}}")
var stanzaVars map[string]string
var envars map[string]string

func InitInterpolationValues(vars map[string]string) {
	if envars != nil {
		return
	}
	//stanza vars override envars values
	stanzaVars = helper.MapKeys[string, string](vars, strings.ToLower)
	//index envars
	envars = make(map[string]string)
	for _, e := range os.Environ() {
		before, after, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		envars[strings.ToLower(before)] = after
	}

}

// InterpolateStrings  we currently support interpolation from envars and Hachi stanza vars
func InterpolateStrings(content string) (string, error) {

	matches := InterpolationRegex.FindAllString(content, -1)
	for _, v := range matches {
		interpolatedPlaceholder := v
		instructions := InterpolationRegex.FindStringSubmatch(v)
		instruct := instructions[2]
		key := instructions[3]
		//todo: add resolver implementation in the future
		if instruct == "local" {
			interpolatedValue := envars[key]
			if val, ok := stanzaVars[key]; ok {
				content = strings.Replace(content, interpolatedPlaceholder, val, -1)
			} else {
				content = strings.Replace(content, interpolatedPlaceholder, interpolatedValue, -1)
				if interpolatedValue == "" {
					return "", fmt.Errorf("ERROR: key " + interpolatedPlaceholder + " value is missing, check your configuration file and machine ENVARS")
				}
			}
		}
	}
	return content, nil
}

func InterpolateCapsuleValues(c echo.Context, InterpolationValues map[string]string, content string) string {

	for name, pattern := range InterpolationValues {
		value := c.Param(name)
		content = strings.Replace(content, pattern, value, -1)
	}

	content, err := InterpolateStrings(content)
	if err != nil {
		log.Error("capsule message failed to interpolate, err: %w", err)
	}
	return content
}
