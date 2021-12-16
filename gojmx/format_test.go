package gojmx

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"text/template"
)

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

	// Nil Config.
	assert.Equal(t, "", FormatConfig(nil, false))
}

func Test_FormatJMXAttributes(t *testing.T) {
	// Nil attributes
	assert.Equal(t, "", FormatJMXAttributes(nil))

	expected := `
-------------------------------------------------------
  - domain: abc
    beans:
      - query: def
        attributes:
          - xyz # Value[INT]: 3
          - xyz # Value[DOUBLE]: 3.2
      - query: ghi
        attributes:
          - xyz # Value[DOUBLE]: 3.2
-------------------------------------------------------
  - domain: jlk
    beans:
      - query: mno
        attributes:
          - xyz # Value[BOOL]: true`

	wrongAttributeFormat := []*JMXAttribute{
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

	assert.Equal(t, expected, FormatJMXAttributes(wrongAttributeFormat))
}

func Test_FormatJMXAttributes_WrongFormat(t *testing.T) {
	// Nil attributes
	assert.Equal(t, "", FormatJMXAttributes(nil))

	wrongAttributeFormat := []*JMXAttribute{
		{
			Attribute: "abc",
			ValueType: ValueTypeInt,
			IntValue:  3,
		},
	}

	assert.Equal(t, "", FormatJMXAttributes(wrongAttributeFormat))
}

func Test_CanParseTemplate(t *testing.T) {
	tpl, err := template.New("nrjmx output").Parse(outputTpl)
	assert.NoError(t, err)
	assert.NotNil(t, tpl)

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, nil)
	assert.NoError(t, err)
	assert.Equal(t, "", buf.String())
}
