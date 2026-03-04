package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/utils"
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

func (l *LoggerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rawMarshal := make(map[string]interface{})
	if err := unmarshal(&rawMarshal); err != nil {
		return fmt.Errorf("logger: %s", err)
	}

	switch item := rawMarshal["dest"].(type) {
	case string:
		l.Dest = item
	default:
		return fmt.Errorf("logger.dest: %v is not a valid type", rawMarshal["dest"])
	}

	if rawMarshal["level"] == nil || rawMarshal["level"] == "" {
		l.Level = 1
		return nil
	}

	// Check the type of level
	switch item := rawMarshal["level"].(type) {
	case int:
		l.Level = item
	case string:
		switch item {
		case "debug", "DEBUG", "Debug":
			l.Level = 0
		case "info", "INFO", "Info":
			l.Level = 1
		case "warn", "WARN", "Warn", "warning", "WARNING", "Warning":
			l.Level = 2
		case "error", "ERROR", "Error":
			l.Level = 3
		default:
			return fmt.Errorf("logger.level: %s is not a valid level", rawMarshal["level"])
		}
	}
	return nil
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

// Configuration that can be used for Netbox.
// In netbox block.
type NetboxConfig struct {
	APIToken string `yaml:"apiToken"`
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	// Can be http or https (default)
	HTTPScheme             HTTPScheme `yaml:"httpScheme"`
	ValidateCert           bool       `yaml:"validateCert"`
	Timeout                int        `yaml:"timeout"`
	Tag                    string     `yaml:"tag"`
	TagColor               string     `yaml:"tagColor"`
	RemoveOrphans          bool       `yaml:"removeOrphans"`
	RemoveOrphansAfterDays int        `yaml:"removeOrphansAfterDays"`
	SourcePriority         []string   `yaml:"sourcePriority"`
	CAFile                 string     `yaml:"caFile"`
}

func (n NetboxConfig) String() string {
	return fmt.Sprintf(
		"NetboxConfig{ApiToken: %s, Hostname: %s, Port: %d, "+
			"HTTPScheme: %s, ValidateCert: %t, Timeout: %d, "+
			"Tag: %s, TagColor: %s, RemoveOrphans: %t, RemoveOrphansAfterDays: %d}",
		n.APIToken,
		n.Hostname,
		n.Port,
		n.HTTPScheme,
		n.ValidateCert,
		n.Timeout,
		n.Tag,
		n.TagColor,
		n.RemoveOrphans,
		n.RemoveOrphansAfterDays,
	)
}

// Configuration that can be used for each of the sources.
type SourceConfig struct {
	Name                string               `yaml:"name"`
	Type                constants.SourceType `yaml:"type"`
	HTTPScheme          HTTPScheme           `yaml:"httpScheme"`
	Hostname            string               `yaml:"hostname"`
	Port                int                  `yaml:"port"`
	Username            string               `yaml:"username"`
	Password            string               `yaml:"password"`
	APIToken            string               `yaml:"apiToken"`
	ValidateCert        bool                 `yaml:"validateCert"`
	Tag                 string               `yaml:"tag"`
	TagColor            string               `yaml:"tagColor"`
	IgnoredSubnets      []string             `yaml:"ignoredSubnets"`
	PermittedSubnets    []string             `yaml:"permittedSubnets"`
	InterfaceFilter     string               `yaml:"interfaceFilter"`
	CollectArpData      bool                 `yaml:"collectArpData"`
	CAFile              string               `yaml:"caFile"`
	IgnoreAssetTags     bool                 `yaml:"ignoreAssetTags"`
	IgnoreSerialNumbers bool                 `yaml:"ignoreSerialNumbers"`
	IgnoreVMTemplates   bool                 `yaml:"ignoreVMTemplates"`
	AssignDomainName    string               `yaml:"assignDomainName"`
	ContinueOnError     bool                 `yaml:"continueOnError"`
	VlanPrefix          string               `yaml:"vlanPrefix"`

	// Relations
	DatacenterClusterGroupRelations map[string]string `yaml:"datacenterClusterGroupRelations"`
	HostSiteRelations               map[string]string `yaml:"hostSiteRelations"`
	HostRoleRelations               map[string]string `yaml:"hostRoleRelations"`
	ClusterSiteRelations            map[string]string `yaml:"clusterSiteRelations"`
	ClusterTenantRelations          map[string]string `yaml:"clusterTenantRelations"`
	HostTenantRelations             map[string]string `yaml:"hostTenantRelations"`
	VMTenantRelations               map[string]string `yaml:"vmTenantRelations"`
	VMRoleRelations                 map[string]string `yaml:"vmRoleRelations"`
	VlanGroupRelations              map[string]string `yaml:"vlanGroupRelations"`
	VlanGroupSiteRelations          map[string]string `yaml:"vlanGroupSiteRelations"`
	VlanTenantRelations             map[string]string `yaml:"vlanTenantRelations"`
	VlanSiteRelations               map[string]string `yaml:"vlanSiteRelations"`
	IPVrfRelations                  map[string]string `yaml:"ipVrfRelations"`
	WlanTenantRelations             map[string]string `yaml:"wlanTenantRelations"`
	CustomFieldMappings             map[string]string `yaml:"customFieldMappings"`
}

// UnmarshalYAML is a custom unmarshal function for SourceConfig.
// This is needed because we map relations to the map[string]string.
func (sc *SourceConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type realSourceConfig struct {
		Name                            string               `yaml:"name"`
		Type                            constants.SourceType `yaml:"type"`
		HTTPScheme                      HTTPScheme           `yaml:"httpScheme"`
		Hostname                        string               `yaml:"hostname"`
		Port                            int                  `yaml:"port"`
		Username                        string               `yaml:"username"`
		Password                        string               `yaml:"password"`
		APIToken                        string               `yaml:"apiToken"`
		ValidateCert                    bool                 `yaml:"validateCert"`
		Tag                             string               `yaml:"tag"`
		TagColor                        string               `yaml:"tagColor"`
		AssignDomainName                string               `yaml:"assignDomainName"`
		VlanPrefix                      string               `yaml:"vlanPrefix"`
		IgnoredSubnets                  []string             `yaml:"ignoredSubnets"`
		PermittedSubnets                []string             `yaml:"permittedSubnets"`
		InterfaceFilter                 string               `yaml:"interfaceFilter"`
		CollectArpData                  bool                 `yaml:"collectArpData"`
		CAFile                          string               `yaml:"caFile"`
		IgnoreSerialNumbers             bool                 `yaml:"ignoreSerialNumbers"`
		IgnoreAssetTags                 bool                 `yaml:"ignoreAssetTags"`
		IgnoreVMTemplates               bool                 `yaml:"ignoreVMTemplates"`
		ContinueOnError                 bool                 `yaml:"continueOnError"`
		DatacenterClusterGroupRelations []string             `yaml:"datacenterClusterGroupRelations"`
		HostSiteRelations               []string             `yaml:"hostSiteRelations"`
		HostRoleRelations               []string             `yaml:"hostRoleRelations"`
		ClusterSiteRelations            []string             `yaml:"clusterSiteRelations"`
		ClusterTenantRelations          []string             `yaml:"clusterTenantRelations"`
		HostTenantRelations             []string             `yaml:"hostTenantRelations"`
		VMTenantRelations               []string             `yaml:"vmTenantRelations"`
		VMRoleRelations                 []string             `yaml:"vmRoleRelations"`
		VlanGroupRelations              []string             `yaml:"vlanGroupRelations"`
		VlanGroupSiteRelations          []string             `yaml:"vlanGroupSiteRelations"`
		VlanTenantRelations             []string             `yaml:"vlanTenantRelations"`
		VlanSiteRelations               []string             `yaml:"vlanSiteRelations"`
		IPVrfRelations                  []string             `yaml:"ipVrfRelations"`
		WlanTenantRelations             []string             `yaml:"wlanTenantRelations"`
		CustomFieldMappings             []string             `yaml:"customFieldMappings"`
	}
	rawMarshal := realSourceConfig{}
	if err := unmarshal(&rawMarshal); err != nil {
		return err
	}
	sc.Name = rawMarshal.Name
	sc.Type = rawMarshal.Type
	sc.HTTPScheme = rawMarshal.HTTPScheme
	sc.Hostname = rawMarshal.Hostname
	sc.Port = rawMarshal.Port
	sc.Username = rawMarshal.Username
	sc.Password = rawMarshal.Password
	sc.APIToken = rawMarshal.APIToken
	sc.ValidateCert = rawMarshal.ValidateCert
	sc.Tag = rawMarshal.Tag
	sc.TagColor = rawMarshal.TagColor
	sc.AssignDomainName = rawMarshal.AssignDomainName
	sc.VlanPrefix = rawMarshal.VlanPrefix
	sc.IgnoredSubnets = rawMarshal.IgnoredSubnets
	sc.PermittedSubnets = rawMarshal.PermittedSubnets
	sc.InterfaceFilter = rawMarshal.InterfaceFilter
	sc.CollectArpData = rawMarshal.CollectArpData
	sc.CAFile = rawMarshal.CAFile
	sc.IgnoreSerialNumbers = rawMarshal.IgnoreSerialNumbers
	sc.IgnoreAssetTags = rawMarshal.IgnoreAssetTags
	sc.IgnoreVMTemplates = rawMarshal.IgnoreVMTemplates
	sc.ContinueOnError = rawMarshal.ContinueOnError

	if len(rawMarshal.DatacenterClusterGroupRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.DatacenterClusterGroupRelations)
		if err != nil {
			return fmt.Errorf("%s.datacenterClusterGroupRelations: %s", rawMarshal.Name, err)
		}
		sc.DatacenterClusterGroupRelations = utils.ConvertStringsToRegexPairs(
			rawMarshal.DatacenterClusterGroupRelations,
		)
	}
	if len(rawMarshal.HostSiteRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("%s.hostSiteRelations: %s", rawMarshal.Name, err)
		}
		sc.HostSiteRelations = utils.ConvertStringsToRegexPairs(rawMarshal.HostSiteRelations)
	}
	if len(rawMarshal.HostRoleRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.HostRoleRelations)
		if err != nil {
			return fmt.Errorf("%s.hostRoleRelations: %s", rawMarshal.Name, err)
		}
		sc.HostRoleRelations = utils.ConvertStringsToRegexPairs(rawMarshal.HostRoleRelations)
	}
	if len(rawMarshal.ClusterSiteRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.ClusterSiteRelations)
		if err != nil {
			return fmt.Errorf("%s.clusterSiteRelations: %s", rawMarshal.Name, err)
		}
		sc.ClusterSiteRelations = utils.ConvertStringsToRegexPairs(rawMarshal.ClusterSiteRelations)
	}
	if len(rawMarshal.ClusterTenantRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.ClusterTenantRelations)
		if err != nil {
			return fmt.Errorf("%s.clusterTenantRelations: %s", rawMarshal.Name, err)
		}
		sc.ClusterTenantRelations = utils.ConvertStringsToRegexPairs(
			rawMarshal.ClusterTenantRelations,
		)
	}
	if len(rawMarshal.HostTenantRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("%s.hostTenantRelations: %s", rawMarshal.Name, err)
		}
		sc.HostTenantRelations = utils.ConvertStringsToRegexPairs(rawMarshal.HostTenantRelations)
	}
	if len(rawMarshal.VMTenantRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.VMTenantRelations)
		if err != nil {
			return fmt.Errorf("%s.vmTenantRelations: %s", rawMarshal.Name, err)
		}
		sc.VMTenantRelations = utils.ConvertStringsToRegexPairs(rawMarshal.VMTenantRelations)
	}
	if len(rawMarshal.VMRoleRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.VMRoleRelations)
		if err != nil {
			return fmt.Errorf("%s.vmRoleRelations: %s", rawMarshal.Name, err)
		}
		sc.VMRoleRelations = utils.ConvertStringsToRegexPairs(rawMarshal.VMRoleRelations)
	}
	if len(rawMarshal.VlanGroupRelations) > 0 {
		err := utils.ValidateRegexRelations(rawMarshal.VlanGroupRelations)
		if err != nil {
			return fmt.Errorf("%s.vlanGroupRelations: %v", rawMarshal.Name, err)
		}
		sc.VlanGroupRelations = utils.ConvertStringsToRegexPairs(rawMarshal.VlanGroupRelations)
	}
	if len(rawMarshal.VlanTenantRelations) > 0 {
		err := utils.ValidateRegexRelations((rawMarshal.VlanTenantRelations))
		if err != nil {
			return fmt.Errorf("%s.vlanTenantRelations: %v", rawMarshal.Name, err)
		}
		sc.VlanTenantRelations = utils.ConvertStringsToRegexPairs(rawMarshal.VlanTenantRelations)
	}
	if len(rawMarshal.VlanSiteRelations) > 0 {
		err := utils.ValidateRegexRelations((rawMarshal.VlanSiteRelations))
		if err != nil {
			return fmt.Errorf("%s.vlanSiteRelations: %v", rawMarshal.Name, err)
		}
		sc.VlanSiteRelations = utils.ConvertStringsToRegexPairs(rawMarshal.VlanSiteRelations)
	}
	if len(rawMarshal.IPVrfRelations) > 0 {
    	err := utils.ValidateRegexRelations(rawMarshal.IPVrfRelations)
    	if err != nil {
        	return fmt.Errorf("%s.ipVrfRelations: %v", rawMarshal.Name, err)
    }
    sc.IPVrfRelations = utils.ConvertStringsToRegexPairs(rawMarshal.IPVrfRelations)
	}
	if len(rawMarshal.VlanGroupSiteRelations) > 0 {
		err := utils.ValidateRegexRelations((rawMarshal.VlanGroupSiteRelations))
		if err != nil {
			return fmt.Errorf("%s.vlanGroupSiteRelations: %v", rawMarshal.Name, err)
		}
		sc.VlanGroupSiteRelations = utils.ConvertStringsToRegexPairs(
			rawMarshal.VlanGroupSiteRelations,
		)
	}
	if len(rawMarshal.WlanTenantRelations) > 0 {
		err := utils.ValidateRegexRelations((rawMarshal.WlanTenantRelations))
		if err != nil {
			return fmt.Errorf("%s.wlanTenantRelations: %v", rawMarshal.Name, err)
		}
		sc.WlanTenantRelations = utils.ConvertStringsToRegexPairs(rawMarshal.WlanTenantRelations)
	}
	if len(rawMarshal.CustomFieldMappings) > 0 {
		err := utils.ValidateRegexRelations((rawMarshal.CustomFieldMappings))
		if err != nil {
			return fmt.Errorf("%s.customFieldMappings: %v", rawMarshal.Name, err)
		}
		sc.CustomFieldMappings = utils.ConvertStringsToRegexPairs(rawMarshal.CustomFieldMappings)
	}
	return nil
}

func (sc SourceConfig) String() string {
	return fmt.Sprintf(
		"SourceConfig{Name: %s, Type: %s, HTTPScheme: %s, Hostname: %s, Port: %d, "+
			"Username: %s, Password: %s, PermittedSubnets: %v, ValidateCert: %t, "+
			"Tag: %s, TagColor: %s, AssignDomainName: %s, VlanPrefix: %s, DatacenterClusterGroupRelations: %s, "+
			"HostSiteRelations: %v, ClusterSiteRelations: %v, ClusterTenantRelations: %v, "+
			"HostTenantRelations: %v, VmTenantRelations: %v, VlanGroupRelations: %v, "+
			"VlanTenantRelations: %v, WlanTenantRelations: %v}",
		sc.Name,
		sc.Type,
		sc.HTTPScheme,
		sc.Hostname,
		sc.Port,
		sc.Username,
		sc.Password,
		sc.IgnoredSubnets,
		sc.ValidateCert,
		sc.Tag,
		sc.TagColor,
		sc.AssignDomainName,
		sc.VlanPrefix,
		sc.DatacenterClusterGroupRelations,
		sc.HostSiteRelations,
		sc.ClusterSiteRelations,
		sc.ClusterTenantRelations,
		sc.HostTenantRelations,
		sc.VMTenantRelations,
		sc.VlanGroupRelations,
		sc.VlanTenantRelations,
		sc.WlanTenantRelations,
	)
}

// Validates the user's config for limits and required fields.
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
	if config.Logger.Level < 0 || config.Logger.Level > 3 {
		return errors.New("logger.level: must be between 0 and 3")
	}
	return nil
}

// Function that validates NetboxConfig.
func validateNetboxConfig(config *Config) error {
	// Validate Netbox config
	if config.Netbox.APIToken == "" {
		return errors.New("netbox.apiToken: cannot be empty")
	}
	if config.Netbox.HTTPScheme != HTTP && config.Netbox.HTTPScheme != HTTPS {
		return errors.New(
			"netbox.httpScheme: must be either http or https. Is " + string(
				config.Netbox.HTTPScheme,
			),
		)
	}
	if config.Netbox.Hostname == "" {
		return errors.New("netbox.hostname: cannot be empty")
	}
	if config.Netbox.Port < 0 || config.Netbox.Port > 65535 {
		return errors.New(
			"netbox.port: must be between 0 and 65535. Is " + fmt.Sprintf("%d", config.Netbox.Port),
		)
	}
	if config.Netbox.Timeout < 0 {
		return errors.New("netbox.timeout: cannot be negative")
	}
	if config.Netbox.Tag == "" {
		config.Netbox.Tag = constants.SsotTagName
	}
	if !config.Netbox.RemoveOrphans {
		if config.Netbox.RemoveOrphansAfterDays < 0 {
			return fmt.Errorf("netbox.RemoveOrphansAfterDays: must be positive integer")
		}
		if config.Netbox.RemoveOrphansAfterDays == 0 {
			config.Netbox.RemoveOrphansAfterDays = constants.CustomFieldOrphanLastSeenDefaultValue
		}
	} else if config.Netbox.RemoveOrphansAfterDays != 0 {
		return fmt.Errorf("netbox.removeOrphansAfterDays has no effect when netbox.removeOrphans is set to true")
	}
	if config.Netbox.TagColor == "" {
		config.Netbox.TagColor = constants.SsotTagColor
	} else {
		// Ensure that TagColor is a string of 6 hexadecimal characters
		if len(config.Netbox.TagColor) != len("ffffff") {
			return errors.New("netbox.tagColor: must be a string of 6 hexadecimal characters")
		}
		for _, c := range config.Netbox.TagColor {
			if c < '0' || c > '9' && c < 'a' || c > 'f' {
				return errors.New("netbox.tagColor: must be a string of 6 lowercase hexadecimal characters")
			}
		}
	}
	if len(config.Netbox.SourcePriority) > 0 {
		if len(config.Netbox.SourcePriority) != len(config.Sources) {
			return fmt.Errorf(
				"netbox.sourcePriority: len(config.Netbox.SourcePriority) != len(config.Sources)",
			)
		}
		for _, sourceName := range config.Netbox.SourcePriority {
			contains := false
			for _, source := range config.Sources {
				if source.Name == sourceName {
					contains = true
					break
				}
			}
			if !contains {
				return fmt.Errorf(
					"netbox.sourcePriority: %s doesn't exist in the sources array",
					sourceName,
				)
			}
		}
	}
	if config.Netbox.CAFile != "" {
		_, err := os.ReadFile(config.Netbox.CAFile)
		if err != nil {
			return fmt.Errorf("netbox.caFile: %s", err)
		}
	}
	return nil
}

//nolint:gocyclo
func validateSourceConfig(config *Config) error {
	// Validate Sources
	for i := range config.Sources {
		externalSource := &config.Sources[i]
		externalSourceStr := externalSource.Name
		if externalSource.Name == "" {
			return fmt.Errorf("source name: cannot be empty")
		}
		switch externalSource.Type {
		case constants.Ovirt:
		case constants.Vmware:
		case constants.Dnac:
		case constants.Proxmox:
		case constants.PaloAlto:
		case constants.Fortigate:
		case constants.FMC:
		case constants.IOSXE:
		default:
			return fmt.Errorf("%s.type is not valid", externalSourceStr)
		}
		if externalSource.HTTPScheme == "" {
			externalSource.HTTPScheme = "https"
		} else if externalSource.HTTPScheme != HTTP && externalSource.HTTPScheme != HTTPS {
			return fmt.Errorf(
				"%s.httpScheme: must be either http or https. Is %s",
				externalSourceStr,
				string(externalSource.HTTPScheme),
			)
		}
		if externalSource.Hostname == "" {
			return fmt.Errorf("%s.hostname: cannot be empty", externalSourceStr)
		}
		if externalSource.Port == 0 {
			externalSource.Port = 443
		} else if externalSource.Port < 0 || externalSource.Port > 65535 {
			return fmt.Errorf("%s.port: must be between 0 and 65535. Is %d", externalSourceStr, externalSource.Port)
		}
		if externalSource.APIToken == "" && externalSource.Type == constants.Fortigate {
			return fmt.Errorf(
				"%s.apiToken is required for %s",
				externalSourceStr,
				constants.Fortigate,
			)
		}
		if externalSource.Username == "" && externalSource.Type != constants.Fortigate {
			return fmt.Errorf("%s.username: cannot be empty", externalSourceStr)
		}
		if externalSource.Password == "" && externalSource.Type != constants.Fortigate {
			return fmt.Errorf("%s.password: cannot be empty", externalSourceStr)
		}
		if externalSource.Tag == "" {
			externalSource.Tag = fmt.Sprintf("Source: %s", externalSource.Name)
		}
		if externalSource.TagColor == "" {
			externalSource.TagColor = constants.SourceTagColorMap[externalSource.Type]
		}
		if externalSource.CAFile != "" {
			if _, err := os.ReadFile(externalSource.CAFile); err != nil {
				return fmt.Errorf("%s.caFile: %s", externalSourceStr, err)
			}
		}
		if len(externalSource.IgnoredSubnets) > 0 {
			for _, ignoredSubnet := range externalSource.IgnoredSubnets {
				if !utils.VerifySubnet(ignoredSubnet) {
					return fmt.Errorf(
						"%s.ignoredSubnets: wrong format: %s",
						externalSourceStr,
						ignoredSubnet,
					)
				}
			}
		}
		if len(externalSource.PermittedSubnets) > 0 {
			for _, permittedSubnet := range externalSource.PermittedSubnets {
				if !utils.VerifySubnet(permittedSubnet) {
					return fmt.Errorf(
						"%s.permittedSubnets: wrong format: %s",
						externalSourceStr,
						permittedSubnet,
					)
				}
			}
		}

		// Try to compile interfaceFilter
		_, err := regexp.Compile(externalSource.InterfaceFilter)
		if err != nil {
			return fmt.Errorf("%s.interfaceFilter: wrong format: %s", externalSourceStr, err)
		}
	}
	return nil
}

func ParseConfig(configFilename string) (*Config, error) {
	// First we read the config file
	file, err := os.Open(configFilename)
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
			HTTPScheme:    "https",
			Port:          constants.HTTPSDefaultPort,
			Timeout:       constants.DefaultAPITimeout,
			RemoveOrphans: true,
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
