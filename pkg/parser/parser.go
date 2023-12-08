package parser

import (
	"errors"
	"fmt"
	"os"

	"github.com/bl4ko/netbox-ssot/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  *LoggerConfig  `yaml:"logger"`
	Netbox  *NetboxConfig  `yaml:"netbox"`
	Sources []SourceConfig `yaml:"source"`
}

type LoggerConfig struct {
	Level int    `yaml:"level"`
	Dest  string `yaml:"dest"`
}

func (l LoggerConfig) String() string {
	if l.Dest == "" {
		return fmt.Sprintf("LoggerConfig{Level: %d, Dest: stdout}", l.Level)
	}
	return fmt.Sprintf("LoggerConfig{Level: %d, Dest: %s}", l.Level, l.Dest)
}

type HTTPScheme string

const (
	HTTP  HTTPScheme = "http"
	HTTPS HTTPScheme = "https"
)

type NetboxConfig struct {
	ApiToken string `yaml:"apiToken"`
	Hostname string `yaml:"hostname"`
	// Can be http or https (default)
	HTTPScheme   HTTPScheme `yaml:"httpScheme"`
	Port         int        `yaml:"port"`
	ValidateCert bool       `yaml:"validateCert"`
	Timeout      int        `yaml:"timeout"`
	MaxRetries   int        `yaml:"maxRetries"`
	Tag          string     `yaml:"tag"`
	TagColor     string     `yaml:"tagColor"`
}

func (n NetboxConfig) String() string {
	return fmt.Sprintf("NetboxConfig{ApiToken: %s, Hostname: %s, Port: %d, HTTPScheme: %s, ValidateCert: %t, Timeout: %d, MaxRetries: %d, Tag: %s, TagColor: %s}", n.ApiToken, n.Hostname, n.Port, n.HTTPScheme, n.ValidateCert, n.Timeout, n.MaxRetries, n.Tag, n.TagColor)
}

type SourceType string

const (
	Ovirt  SourceType = "ovirt"
	Vmware SourceType = "vmware"
)

type SourceConfig struct {
	Name                   string     `yaml:"name"`
	Type                   SourceType `yaml:"type"`
	HTTPScheme             HTTPScheme `yaml:"httpScheme"`
	Hostname               string     `yaml:"hostname"`
	Port                   int        `yaml:"port"`
	Username               string     `yaml:"username"`
	Password               string     `yaml:"password"`
	PermittedSubnets       []string   `yaml:"permittedSubnets"`
	ValidateCert           bool       `yaml:"validateCert"`
	Tag                    string     `yaml:"tag"`
	TagColor               string     `yaml:"tagColor"`
	HostSiteRelations      []string   `yaml:"hostSiteRelations"`
	ClusterSiteRelations   []string   `yaml:"clusterSiteRelations"`
	ClusterTenantRelations []string   `yaml:"clusterTenantRelations"`
	HostTenantRelations    []string   `yaml:"hostTenantRelations"`
	VmTenantRelations      []string   `yaml:"vmTenantRelations"`
}

func (s SourceConfig) String() string {
	return fmt.Sprintf("SourceConfig{Name: %s, Type: %s, HTTPScheme: %s, Hostname: %s, Port: %d, Username: %s, Password: %s, PermittedSubnets: %v, ValidateCert: %t, Tag: %s, TagColor: %s, HostSiteRelations: %v, ClusterSiteRelations: %v, clusterTenantRelations: %v, HostTenantRelations: %v, VmTenantRelations %v}", s.Name, s.Type, s.HTTPScheme, s.Hostname, s.Port, s.Username, s.Password, s.PermittedSubnets, s.ValidateCert, s.Tag, s.TagColor, s.HostSiteRelations, s.ClusterSiteRelations, s.ClusterTenantRelations, s.HostTenantRelations, s.VmTenantRelations)
}

// Validates the user's config for limits and required fields
func validateConfig(config *Config) error {
	err := validateLoggerConfig(config)
	if err != nil {
		return err
	}

	err = validateNetboxConfig(config)
	if err != nil {
		return err
	}

	err = validateSourceConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func validateLoggerConfig(config *Config) error {
	if config.Logger.Level < 0 || config.Logger.Level > 4 {
		return errors.New("logger.level: must be between 0 and 4")
	}
	return nil
}

// Function that validates NetboxConfig
func validateNetboxConfig(config *Config) error {
	// Validate Netbox config
	if config.Netbox.ApiToken == "" {
		return errors.New("netbox.apiToken: cannot be empty")
	}
	if config.Netbox.HTTPScheme != HTTP && config.Netbox.HTTPScheme != HTTPS {
		return errors.New("netbox.httpScheme: must be either http or https. Is " + string(config.Netbox.HTTPScheme))
	}
	if config.Netbox.Hostname == "" {
		return errors.New("netbox.hostname: cannot be empty")
	}
	if config.Netbox.Port < 0 || config.Netbox.Port > 65535 {
		return errors.New("netbox.port: must be between 0 and 65535. Is " + fmt.Sprintf("%d", config.Netbox.Port))
	}
	if config.Netbox.Timeout < 0 {
		return errors.New("netbox.timeout: cannot be negative")
	}
	if config.Netbox.MaxRetries < 0 {
		return errors.New("netbox.maxRetries: cannot be negative")
	}
	if config.Netbox.Tag == "" {
		config.Netbox.Tag = "netbox-ssot"
	}
	if config.Netbox.TagColor == "" {
		config.Netbox.TagColor = "00add8"
	} else {
		// Ensure that TagColor is a string of 6 hexadecimal characters
		if len(config.Netbox.TagColor) != 6 {
			return errors.New("netbox.tagColor: must be a string of 6 hexadecimal characters")
		}
		for _, c := range config.Netbox.TagColor {
			if c < '0' || c > '9' && c < 'a' || c > 'f' {
				return errors.New("netbox.tagColor: must be a string of 6 lowercase hexadecimal characters")
			}
		}
	}
	return nil
}

func validateSourceConfig(config *Config) error {
	// Valicate Sources
	for i := range config.Sources {
		externalSource := &config.Sources[i]
		externalSourceStr := "source[" + externalSource.Name + "]."
		if externalSource.Name == "" {
			return errors.New(externalSourceStr + "name cannot be empty")
		}
		if externalSource.HTTPScheme == "" {
			externalSource.HTTPScheme = "https"
		} else if externalSource.HTTPScheme != HTTP && externalSource.HTTPScheme != HTTPS {
			return errors.New(externalSourceStr + "httpScheme must be either http or https. Is " + string(externalSource.HTTPScheme))
		}
		if externalSource.Hostname == "" {
			return errors.New(externalSourceStr + "hostname cannot be empty")
		}
		if externalSource.Port == 0 {
			externalSource.Port = 443
		} else if externalSource.Port < 0 || externalSource.Port > 65535 {
			return errors.New(externalSourceStr + "port must be between 0 and 65535. Is " + fmt.Sprintf("%d", externalSource.Port))
		}
		if externalSource.Username == "" {
			return errors.New(externalSourceStr + "username cannot be empty")
		}
		if externalSource.Password == "" {
			return errors.New(externalSourceStr + "password cannot be empty")
		}
		if externalSource.Tag == "" {
			externalSource.Tag = fmt.Sprintf("Source: %s", externalSource.Name)
		}
		if externalSource.TagColor == "" {
			source2color := map[SourceType]string{
				Ovirt: "07426b",
			}
			externalSource.TagColor = source2color[externalSource.Type]
		}
		switch externalSource.Type {
		case Ovirt:
			// Do nothing
		default:
			return errors.New(externalSourceStr + "type is not valid")
		}
		if len(externalSource.HostSiteRelations) > 0 {
			err := utils.ValidateRegexRelations(externalSource.HostSiteRelations)
			if err != nil {
				return errors.New(externalSourceStr + "hostSiteRelations: " + err.Error())
			}
		}
		if len(externalSource.ClusterSiteRelations) > 0 {
			err := utils.ValidateRegexRelations(externalSource.ClusterSiteRelations)
			if err != nil {
				return errors.New(externalSourceStr + "clusterSiteRelations: " + err.Error())
			}
		}
		if len(externalSource.ClusterTenantRelations) > 0 {
			err := utils.ValidateRegexRelations(externalSource.ClusterTenantRelations)
			if err != nil {
				return errors.New(externalSourceStr + "clusterTenantRelations: " + err.Error())
			}
		}
		if len(externalSource.HostTenantRelations) > 0 {
			err := utils.ValidateRegexRelations(externalSource.HostTenantRelations)
			if err != nil {
				return errors.New(externalSourceStr + "hostTenantRelations: " + err.Error())
			}
		}
		if len(externalSource.VmTenantRelations) > 0 {
			err := utils.ValidateRegexRelations(externalSource.VmTenantRelations)
			if err != nil {
				return errors.New(externalSourceStr + "vmTenantRelations: " + err.Error())
			}
		}
	}
	return nil
}

func ParseConfig(filename string) (*Config, error) {
	// First we read the config file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Define Config with default values
	config := &Config{
		Logger: &LoggerConfig{
			Level: 1,
			Dest:  "",
		},
		Netbox: &NetboxConfig{
			HTTPScheme: "https",
			Port:       443,
			Timeout:    30,
			MaxRetries: 3,
		},
		Sources: []SourceConfig{},
	}

	// Parse the config file into a Config struct
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	// Validate the config for limits and required fields
	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
