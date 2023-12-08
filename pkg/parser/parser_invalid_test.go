package parser

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestInvalidConfig1(t *testing.T) {
	filename := filepath.Join("testdata", "invalid_config1.yaml")
	expectedErr := "netbox.hostname: cannot be empty"
	_, err := ParseConfig(filename)
	fmt.Printf("%v", err)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expecrted error: %v, got: %v", expectedErr, err)
		return
	}
}

func TestInvalidConfig2(t *testing.T) {
	filename := filepath.Join("testdata", "invalid_config2.yaml")
	expectedErr := "netbox.port: must be between 0 and 65535. Is 333333"
	_, err := ParseConfig(filename)
	fmt.Printf("%v", err)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
		return
	}
}

func TestInvalidConfig3(t *testing.T) {
	filename := filepath.Join("testdata", "invalid_config3.yaml")
	expectedErr := "source[testolvm].type is not valid"
	_, err := ParseConfig(filename)
	fmt.Printf("%v", err)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
		return
	}
}

func TestInvalidConfig4(t *testing.T) {
	filename := filepath.Join("testdata", "invalid_config4.yaml")
	expectedErr := "netbox.httpScheme: must be either http or https. Is httpd"
	_, err := ParseConfig(filename)
	fmt.Printf("%v", err)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
		return
	}
}

func TestInvalidConfig5(t *testing.T) {
	filename := filepath.Join("testdata", "invalid_config5.yaml")
	expectedErr := "source[prodovirt].httpScheme must be either http or https. Is httpd"
	_, err := ParseConfig(filename)
	fmt.Printf("%v", err)
	if err == nil {
		t.Errorf("%s", err)
		return
	} else if err.Error() != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
		return
	}
}

func TestInvalidConfig6(t *testing.T) {
	filename := filepath.Join("testdata", "invalid_config6.yaml")
	expectedErr := "source[testolvm].hostTenantRelations: invalid regex relation: This should not work. Should be of format: regex = value"
	_, err := ParseConfig(filename)
	fmt.Printf("%v", err)
	if err == nil {
		t.Errorf("%s", err)
		return
	} else if err.Error() != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
		return
	}
}

func TestInvalidConfig7(t *testing.T) {
	filename := filepath.Join("testdata", "invalid_config7.yaml")
	expectedErr := "source[prodolvm].hostTenantRelations: invalid regex: [a-z++, in relation: [a-z++ = Should not work"
	_, err := ParseConfig(filename)
	if err == nil {
		t.Errorf("%s", err)
		return
	} else if err.Error() != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
		return
	}
}
