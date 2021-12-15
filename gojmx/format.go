package gojmx

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var outputTpl = `
{{- range $mBean, $attrs := . }}
-------------------------------------------------------
    MBean:
      Domain: {{ $mBean.Domain }}
      Query: {{ $mBean.Query }}
	Attributes:
	{{- range $i, $attr := $attrs }}
        {{ $i }}: {{ $attr.Attribute }}
           Value [{{ $attr.ValueType }}]: {{ $attr.Value }}
    {{- end }}
{{- end }}`

type mBeanFormat struct {
	Domain string
	Query  string
}

type attributeFormat struct {
	Attribute string
	Value     interface{}
	ValueType string
}

// FormatJMXAttributes will prettify JMXAttributes.
func FormatJMXAttributes(attrs []*JMXAttribute) string {
	result := map[mBeanFormat][]attributeFormat{}

	separator := ",attr="

	for _, attr := range attrs {
		i := strings.LastIndex(attr.Attribute, separator)
		if i < 0 {
			continue
		}

		split := strings.SplitAfterN(attr.Attribute[:i], ":", 2)
		if len(split) != 2 {
			continue
		}

		key := mBeanFormat{
			Domain: split[0],
			Query:  split[1],
		}

		result[key] = append(result[key], attributeFormat{
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

	sb.WriteString(fmt.Sprintf(", RequestTimeoutMs: '%d', URIPath: '%v'", config.RequestTimoutMs, config.UriPath))

	return sb.String()
}
