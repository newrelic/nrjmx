package gojmx

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
	"text/template"
)

var noAttributesFormat = `
---
collect:`

func Test_FormatConfig(t *testing.T) {
	testURIPath := "test_URI_PATH"

	config := &JMXConfig{
		ConnectionURL:         "",
		Hostname:              "test_hostname",
		Port:                  int32(123),
		UriPath:               &testURIPath,
		Username:              "test_username",
		Password:              "test_password",
		KeyStore:              "test_keystore",
		KeyStorePassword:      "test_keystore_password",
		TrustStore:            "test_truststore",
		TrustStorePassword:    "test_truststore_password",
		IsRemote:              true,
		IsJBossStandaloneMode: true,
		UseSSL:                true,
		RequestTimoutMs:       4567,
	}

	// Exposed credentials.
	hideSecrets := false
	expected := "Hostname: 'test_hostname', Port: '123', IsJBossStandaloneMode: 'true', IsRemote: 'true', UseSSL: 'true', Username: 'test_username', Password: 'test_password', Truststore: 'test_truststore', TruststorePassword: 'test_truststore_password', Keystore: 'test_keystore', KeystorePassword: 'test_keystore_password', RequestTimeoutMs: '4567', URIPath: 'test_URI_PATH'"
	actual := FormatConfig(config, hideSecrets)
	assert.Equal(t, expected, actual)

	// Hidden credentials.
	hideSecrets = true
	expected = "Hostname: 'test_hostname', Port: '123', IsJBossStandaloneMode: 'true', IsRemote: 'true', UseSSL: 'true', Username: '<HIDDEN>', Password: '<HIDDEN>', Truststore: 'test_truststore', TruststorePassword: '<HIDDEN>', Keystore: 'test_keystore', KeystorePassword: '<HIDDEN>', RequestTimeoutMs: '4567', URIPath: 'test_URI_PATH'"
	actual = FormatConfig(config, hideSecrets)
	assert.Equal(t, expected, actual)

	// Empty credentials.
	config.Username = ""
	config.Password = ""
	config.KeyStorePassword = ""
	config.TrustStorePassword = ""

	expected = "Hostname: 'test_hostname', Port: '123', IsJBossStandaloneMode: 'true', IsRemote: 'true', UseSSL: 'true', Username: '<EMPTY>', Password: '<EMPTY>', Truststore: 'test_truststore', TruststorePassword: '<EMPTY>', Keystore: 'test_keystore', KeystorePassword: '<EMPTY>', RequestTimeoutMs: '4567', URIPath: 'test_URI_PATH'"
	actual = FormatConfig(config, hideSecrets)
	assert.Equal(t, expected, actual)

	config.UriPath = nil
	expected = "Hostname: 'test_hostname', Port: '123', IsJBossStandaloneMode: 'true', IsRemote: 'true', UseSSL: 'true', Username: '<EMPTY>', Password: '<EMPTY>', Truststore: 'test_truststore', TruststorePassword: '<EMPTY>', Keystore: 'test_keystore', KeystorePassword: '<EMPTY>', RequestTimeoutMs: '4567'"
	actual = FormatConfig(config, hideSecrets)
	assert.Equal(t, expected, actual)

	// Nil Config.
	assert.Equal(t, "", FormatConfig(nil, hideSecrets))
}

func Test_FormatConfig_ConnectionURL(t *testing.T) {
	testURIPath := "test_URI_PATH"

	config := &JMXConfig{
		ConnectionURL:         "service:jmx:rmi:///jndi/rmi://localhost:123/jmxrmi",
		Hostname:              "test_hostname",
		Port:                  int32(123),
		UriPath:               &testURIPath,
		Username:              "test_username",
		Password:              "test_password",
		KeyStore:              "test_keystore",
		KeyStorePassword:      "test_keystore_password",
		TrustStore:            "test_truststore",
		TrustStorePassword:    "test_truststore_password",
		IsRemote:              true,
		IsJBossStandaloneMode: true,
		UseSSL:                true,
		RequestTimoutMs:       4567,
	}

	hideSecrets := true
	expected := "ConnectionURL: 'service:jmx:rmi:///jndi/rmi://localhost:123/jmxrmi', Username: '<HIDDEN>', Password: '<HIDDEN>', Truststore: 'test_truststore', TruststorePassword: '<HIDDEN>', Keystore: 'test_keystore', KeystorePassword: '<HIDDEN>', RequestTimeoutMs: '4567', URIPath: 'test_URI_PATH'"
	actual := FormatConfig(config, hideSecrets)
	assert.Equal(t, expected, actual)
}

func Test_FormatJMXAttributes(t *testing.T) {
	// Nil attributes
	assert.Equal(t, noAttributesFormat, FormatJMXAttributes(nil))

	expected := `
---
collect:
##############################################
### BEGIN Beans for domain: "abc"
##############################################
  - domain: abc
    beans:
      - query: def
        attributes:
          # Attribute: "xyz", Sample[INT]: 3
          - xyz
          # Attribute: "xyz", Sample[DOUBLE]: 3.2
          - xyz
      
      - query: ghi
        attributes:
          # Attribute: "xyz", Sample[DOUBLE]: 3.2
          - xyz
      
##############################################
### END Beans for domain: "abc"
##############################################

##############################################
### BEGIN Beans for domain: "jlk"
##############################################
  - domain: jlk
    beans:
      - query: mno
        attributes:
          # Attribute: "xyz", Sample[BOOL]: true
          - xyz
      
##############################################
### END Beans for domain: "jlk"
##############################################
`

	attributes := []*JMXAttribute{
		{
			Attribute: "abc:def,attr=xyz",
			ValueType: ValueTypeInt,
			IntValue:  3,
		},
		{
			Attribute:   "abc:def,attr=xyz",
			ValueType:   ValueTypeDouble,
			DoubleValue: 3.2,
		},
		{
			Attribute:   "abc:ghi,attr=xyz",
			ValueType:   ValueTypeDouble,
			DoubleValue: 3.2,
		},
		{
			Attribute: "jlk:mno,attr=xyz",
			ValueType: ValueTypeBool,
			BoolValue: true,
		},
	}

	// WHEN FormatJMXAttributes
	formatted := FormatJMXAttributes(attributes)

	// THEN output is YAML format
	out := map[string]interface{}{}
	assert.NoError(t, yaml.Unmarshal([]byte(formatted), out))

	// AND has the expected format.
	assert.Equal(t, expected, formatted)
}

func Test_FormatJMXAttributes_WrongFormat(t *testing.T) {
	// Nil attributes
	assert.Equal(t, noAttributesFormat, FormatJMXAttributes(nil))

	wrongAttributeFormat := []*JMXAttribute{
		{
			Attribute: "abc",
			ValueType: ValueTypeInt,
			IntValue:  3,
		},
	}

	assert.Equal(t, noAttributesFormat, FormatJMXAttributes(wrongAttributeFormat))
}

func Test_CanParseTemplate(t *testing.T) {
	tpl, err := template.New("nrjmx output").Parse(outputTpl)
	assert.NoError(t, err)
	assert.NotNil(t, tpl)

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, nil)
	assert.NoError(t, err)
	assert.Equal(t, noAttributesFormat, buf.String())
}
