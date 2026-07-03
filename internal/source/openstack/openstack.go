package openstack

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
)

type Server struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Status          string                   `json:"status"`
	VMState         string                   `json:"OS-EXT-STS:vm_state"`
	Flavor          any                      `json:"flavor"`
	Addresses       map[string]interface{}   `json:"addresses"`
	Metadata        any                      `json:"metadata"`
	Image           any                      `json:"image"`
	ImageMetadata   any                      `json:"image_metadata"`
	AttachedVolumes []servers.AttachedVolume `json:"os-extended-volumes:volumes_attached"`
}

type Source struct {
	common.Config

	// OpenStack API data
	Servers  []Server
	Flavors  []flavors.Flavor
	Networks []networks.Network
	Volumes  []volumes.Volume
	Images   []images.Image

	// Gophercloud clients
	ComputeClient      *gophercloud.ServiceClient
	NetworkClient      *gophercloud.ServiceClient
	BlockStorageClient *gophercloud.ServiceClient
	ImageClient        *gophercloud.ServiceClient
}

// domainConfig holds resolved domain configuration for OpenStack authentication.
type domainConfig struct {
	domainName        string
	domainID          string
	projectDomainName string
	projectDomainID   string
}

// resolveDomainConfig resolves domain configuration with proper fallbacks and precedence.
// It handles:
// - id| prefix in domainName
// - Default domain fallback
// - ID precedence over Name
// - ProjectDomain* fallback to Domain*
func resolveDomainConfig(cfg *parser.SourceConfig) domainConfig {
	domainName := cfg.DomainName
	domainID := cfg.DomainID

	// Handle id| prefix in domainName (same as MH WHMCS module)
	if strings.HasPrefix(domainName, "id|") {
		domainID = strings.TrimPrefix(domainName, "id|")
		domainName = ""
	}

	if domainName == "" && domainID == "" {
		domainName = "Default"
	}

	// Enforce ID precedence: if ID is set, clear Name to avoid ambiguity
	if domainID != "" {
		domainName = ""
	}

	projectDomainName := cfg.ProjectDomainName
	if projectDomainName == "" {
		projectDomainName = domainName
	}
	projectDomainID := cfg.ProjectDomainID
	if projectDomainID == "" && cfg.ProjectDomainName == "" {
		// Only inherit domainID when neither project domain field is explicitly set
		projectDomainID = domainID
	} else if projectDomainID != "" {
		// Enforce ID precedence: if ID is set, clear Name to avoid ambiguity
		projectDomainName = ""
	}

	return domainConfig{
		domainName:        domainName,
		domainID:          domainID,
		projectDomainName: projectDomainName,
		projectDomainID:   projectDomainID,
	}
}

func (oss *Source) Init() error {
	projectName := oss.SourceConfig.ProjectName
	if projectName == "" {
		projectName = oss.SourceConfig.TenantName
	}

	projectID := oss.SourceConfig.ProjectID
	if projectID == "" {
		projectID = oss.SourceConfig.TenantID
	}

	// Resolve domain configuration
	domains := resolveDomainConfig(oss.SourceConfig)

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: oss.SourceConfig.Hostname,
		Username:         oss.SourceConfig.Username,
		Password:         oss.SourceConfig.Password,
		DomainName:       domains.domainName,
		DomainID:         domains.domainID,
		AllowReauth:      true,
	}

	// Explicitly set scope if project details are available
	switch {
	case projectID != "":
		// ProjectID must be supplied alone in a Scope for Gophercloud v2
		opts.Scope = &gophercloud.AuthScope{
			ProjectID: projectID,
		}
	case projectName != "":
		// ProjectName requires domain info for disambiguation
		opts.Scope = &gophercloud.AuthScope{
			ProjectName: projectName,
			DomainName:  domains.projectDomainName,
			DomainID:    domains.projectDomainID,
		}
	case domains.domainID != "":
		opts.Scope = &gophercloud.AuthScope{
			DomainID: domains.domainID,
		}
	case domains.domainName != "":
		opts.Scope = &gophercloud.AuthScope{
			DomainName: domains.domainName,
		}
	}

	oss.Logger.Debugf(oss.Ctx, "OpenStack AuthOptions: Endpoint=%s, Username=%s, Project=%s, Domain=%s",
		opts.IdentityEndpoint, opts.Username, projectName, domains.domainName)

	// Setup custom HTTP client to respect validateCert
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !oss.SourceConfig.ValidateCert,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}

	provider, err := openstack.NewClient(opts.IdentityEndpoint)
	if err != nil {
		return fmt.Errorf("error creating OpenStack provider client: %s", err)
	}
	provider.HTTPClient = *httpClient

	err = openstack.Authenticate(oss.Ctx, provider, opts)
	if err != nil {
		return fmt.Errorf("error authenticating with OpenStack (check credentials/scope): %s", err)
	}

	region := oss.SourceConfig.Region
	if region == "" {
		region = "RegionOne"
	}

	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}
	oss.ComputeClient = computeClient

	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return fmt.Errorf("error creating network client: %s", err)
	}
	oss.NetworkClient = networkClient

	blockStorageClient, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return fmt.Errorf("error creating block storage client: %s", err)
	}
	oss.BlockStorageClient = blockStorageClient

	imageClient, err := openstack.NewImageV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return fmt.Errorf("error creating image client: %s", err)
	}
	oss.ImageClient = imageClient

	initFuncs := []func(context.Context) error{
		oss.initServers,
		oss.initFlavors,
		oss.initNetworks,
		oss.initVolumes,
		oss.initImages,
	}

	for _, initFunc := range initFuncs {
		startTime := time.Now()
		if err := initFunc(oss.Ctx); err != nil {
			return fmt.Errorf("openstack initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		oss.Logger.Infof(
			oss.Ctx,
			"Successfully initialized %s in %f seconds",
			utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"),
			duration.Seconds(),
		)
	}

	return nil
}

func (oss *Source) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		oss.syncServers,
	}
	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		funcName := utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync")
		err := syncFunc(nbi)
		if err != nil {
			if oss.SourceConfig.ContinueOnError {
				oss.Logger.Errorf(
					oss.Ctx,
					"Error syncing %s: %s (continuing due to continueOnError flag)",
					funcName,
					err,
				)
			} else {
				return err
			}
		} else {
			duration := time.Since(startTime)
			oss.Logger.Infof(
				oss.Ctx,
				"Successfully synced %s in %f seconds",
				funcName,
				duration.Seconds(),
			)
		}
	}
	return nil
}
