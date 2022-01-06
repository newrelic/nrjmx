/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var outputTpl = `
---
collect:
{{- range $domain, $queryFormat := . }}
##############################################
### BEGIN Beans for domain: "{{ $domain }}"
##############################################
  - domain: {{ $domain }}
    beans:
      {{- range $query, $attrs := $queryFormat }}
      - query: {{ $query }}
        attributes:
	    {{- range $attr := $attrs }}
          # Attribute: "{{ $attr.Attribute }}", Sample[{{ $attr.ValueType }}]: {{ $attr.Value }}
          - {{ $attr.Attribute }}
        {{- end }}
      {{ end }}
##############################################
### END Beans for domain: "{{ $domain }}"
##############################################
{{ end -}}`

type queryFormat map[string][]attributeFormat

type attributeFormat struct {
	Attribute string
	Value     interface{}
	ValueType string
}

// FormatJMXAttributes will prettify JMXAttributes.
func FormatJMXAttributes(attrs []*JMXAttribute) string {
	result := map[string]queryFormat{}
	separator := ",attr="

	for _, attr := range attrs {
		i := strings.LastIndex(attr.Attribute, separator)
		if i < 0 {
			continue
		}

		split := strings.SplitN(attr.Attribute[:i], ":", 2)
		if len(split) != 2 {
			continue
		}

		domain, query := split[0], split[1]

		if _, exists := result[domain]; !exists {
			result[domain] = queryFormat{}
		}

		result[domain][query] = append(result[domain][query], attributeFormat{
			Attribute: attr.Attribute[i+len(separator):],
			Value:     attr.GetValue(),
			ValueType: fmt.Sprintf("%v", attr.ValueType),
		})
	}

	buf := new(bytes.Buffer)
	tpl, err := template.New("nrjmx output").Parse(outputTpl)
	if err != nil {
		return err.Error()
	}
	err = tpl.Execute(buf, result)
	if err != nil {
		return err.Error()
	}

	return buf.String()
}

// FormatConfig will convert the JMXConfig into a string.
func FormatConfig(config *JMXConfig, hideSecrets bool) string {
	sb := strings.Builder{}
	if config == nil {
		return sb.String()
	}

	if config.ConnectionURL != "" {
		sb.WriteString(fmt.Sprintf("ConnectionURL: '%s'", config.ConnectionURL))
	} else {
		sb.WriteString(fmt.Sprintf("Hostname: '%s', Port: '%d', IsJBossStandaloneMode: '%t', IsRemote: '%t', UseSSL: '%t'",
			config.Hostname,
			config.Port,
			config.IsJBossStandaloneMode,
			config.IsRemote,
			config.UseSSL,
		))
	}

	obfuscate := func(in string) string {
		if in == "" {
			return "<EMPTY>"
		}
		if hideSecrets {
			return "<HIDDEN>"
		}
		return in
	}

	sb.WriteString(fmt.Sprintf(", Username: '%s', Password: '%s'",
		obfuscate(config.Username),
		obfuscate(config.Password),
	))

	sb.WriteString(fmt.Sprintf(", Truststore: '%s', TruststorePassword: '%s', Keystore: '%s', KeystorePassword: '%s'",
		config.TrustStore,
		obfuscate(config.TrustStorePassword),
		config.KeyStore,
		obfuscate(config.KeyStorePassword),
	))

	sb.WriteString(fmt.Sprintf(", RequestTimeoutMs: '%d'", config.RequestTimoutMs))
	if config.UriPath != nil {
		sb.WriteString(fmt.Sprintf(", URIPath: '%v'", *config.UriPath))
	}

	return sb.String()
}
