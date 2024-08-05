package templating

import (
	"fmt"
	"github.com/catalystcommunity/app-utils-go/errorutils"
	"github.com/joomcode/errorx"
	"regexp"
	"strings"
)

func TemplateString(source string, replacements ...string) (string, error) {
	templated := source
	if len(replacements)%2 != 0 {
		return "", errorx.IllegalArgument.New("received invalid number of replacements. Should be an even number")
	}

	for i, key := range replacements {
		if i%2 != 0 {
			continue
		}
		// ensure proper format
		key = strings.ReplaceAll(key, "<<", "")
		key = strings.ReplaceAll(key, ">>", "")
		templated = strings.ReplaceAll(templated, fmt.Sprintf("<<%s>>", key), replacements[i+1])
	}
	return templated, nil
}

func MustTemplateString(source string, replacements ...string) string {
	templated, err := TemplateString(source, replacements...)
	errorutils.PanicOnErr(nil, "error templating string", err)
	panicIfNotTemplated(templated)
	return templated
}

func TemplateStringWithMap(source string, replacements map[string]string) string {
	templated := source
	for key, value := range replacements {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("<<%s>>", key), value)
	}
	return templated
}

func MustTemplateStringWithMap(source string, replacements map[string]string) string {
	templated := TemplateStringWithMap(source, replacements)
	panicIfNotTemplated(templated)
	return templated
}

func panicIfNotTemplated(templated string) {
	if strings.Contains(templated, "<<") || strings.Contains(templated, ">>") {
		errorutils.PanicOnErr(nil, "string is not fully templated", errorx.IllegalState.New("Templated string: %s", templated))
	}
}

func TemplateWithFunction(source string, replaceFunction func(key string) (string, error)) (string, error) {
	r := regexp.MustCompile(`<<.*>>`)
	keys := r.FindAllString(source, -1)
	templatedString := source
	for _, key := range keys {
		replacement, err := replaceFunction(key)
		if err != nil {
			return "", err
		}
		templatedString, err = TemplateString(templatedString, key, replacement)
		if err != nil {
			return "", err
		}
	}
	return templatedString, nil
}

func MustTemplateWithFunction(source string, replaceFunction func(key string) (string, error)) string {
	templatedString, err := TemplateWithFunction(source, replaceFunction)
	errorutils.PanicOnErr(nil, "error templating string with function", err)
	return templatedString
}
