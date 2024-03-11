package template

import (
	"errors"
	"github.com/GSH-LAN/Unwindia_common/src/go/matchservice"
	"github.com/rs/zerolog/log"
	"net"
	"strings"
	"text/template"
)

func ParseTemplateForMatch(tpl string, matchinfo *matchservice.MatchInfo) (string, error) {
	if matchinfo == nil {
		return "", errors.New("empty matchinfo")
	}

	funcs := map[string]any{
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"splitHostPort": func(s string) (map[string]string, error) {
			host, port, err := net.SplitHostPort(s)
			if err != nil {
				return nil, err
			}
			return map[string]string{"host": host, "port": port}, nil
		},
	}

	tmpl, err := template.New("match").Option("missingkey=error").Funcs(funcs).Parse(tpl)
	if err != nil {
		log.Err(err).Msg("Error parsing template")
		return "", err
	}

	parsedTemplate := strings.Builder{}
	err = tmpl.Execute(&parsedTemplate, matchinfo)
	if err != nil {
		log.Err(err).Msg("Error parsing matchinfo into template")
		return "", err
	}

	return parsedTemplate.String(), nil
}
