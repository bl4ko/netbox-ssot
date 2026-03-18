package f5

import (
	"fmt"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// syncVirtualServers syncs F5 BIG-IP virtual servers as VIP IP addresses in NetBox.
func (fs *F5Source) syncVirtualServers(nbi *inventory.NetboxInventory) error {
	// Resolve target object and interface if configured.
	// Looks up as VM first, then as Device — works for both virtual and physical F5.
	var assignedObjectType constants.ContentType
	var assignedObjectID int
	if err := fs.resolveTarget(nbi, &assignedObjectType, &assignedObjectID); err != nil {
		return err
	}

	for _, vs := range fs.VirtualServers {
		fs.Logger.Debugf(
			fs.Ctx,
			"virtual server %s: destination=%s, mask=%s",
			vs.Name, vs.Destination, vs.Mask,
		)
		ip, maskBits, err := parseDestination(vs.Destination, vs.Mask)
		if err != nil {
			fs.Logger.Warningf(fs.Ctx, "skipping virtual server %s: %s", vs.Name, err)
			continue
		}

		if !utils.IsPermittedIPAddress(ip, fs.SourceConfig.PermittedSubnets, fs.SourceConfig.IgnoredSubnets) {
			fs.Logger.Debugf(fs.Ctx, "virtual server %s IP %s is not permitted", vs.Name, ip)
			continue
		}

		vrf, err := common.MatchIPToVRF(fs.Ctx, nbi, ip, fs.SourceConfig.IPVrfRelations)
		if err != nil {
			return fmt.Errorf("match ip to vrf for virtual server %s: %s", vs.Name, err)
		}

		description := vs.Description
		if description == "" {
			description = fmt.Sprintf("F5 Virtual Server: %s", vs.Name)
		}

		address := fmt.Sprintf("%s/%d", ip, maskBits)
		ipAddr := &objects.IPAddress{
			NetboxObject: objects.NetboxObject{
				Tags:        fs.GetSourceTags(),
				Description: description,
				CustomFields: map[string]interface{}{
					constants.CustomFieldArpEntryName: false,
				},
			},
			Address: address,
			Role:    &objects.IPAddressRoleVIP,
			VRF:     vrf,
		}
		if assignedObjectID > 0 {
			ipAddr.AssignedObjectType = assignedObjectType
			ipAddr.AssignedObjectID = assignedObjectID
		}

		_, err = nbi.AddIPAddress(fs.Ctx, ipAddr)
		if err != nil {
			fs.Logger.Warningf(fs.Ctx, "add ip address for virtual server %s: %s", vs.Name, err)
			continue
		}
	}
	return nil
}

// resolveTarget looks up the source hostname IP in the NetBox inventory to find the
// associated VM or Device, then resolves the targetInterface on it.
func (fs *F5Source) resolveTarget(
	nbi *inventory.NetboxInventory,
	objType *constants.ContentType,
	objID *int,
) error {
	hostname := fs.SourceConfig.Hostname
	ifaceName := fs.SourceConfig.TargetInterface

	// Lookup the IP address of the hostname in the NetBox inventory
	// to find which VM or Device it belongs to.
	ipObj := nbi.GetIPAddressByAddress(hostname)
	if ipObj == nil {
		fs.Logger.Warningf(fs.Ctx, "hostname IP %s not found in NetBox inventory, VIPs will be unassigned", hostname)
		return nil
	}

	switch ipObj.AssignedObjectType {
	case constants.ContentTypeVirtualizationVMInterface:
		vmIface := nbi.GetVMInterfaceByID(ipObj.AssignedObjectID)
		if vmIface == nil {
			return fmt.Errorf("VM interface ID %d for hostname IP %s not found", ipObj.AssignedObjectID, hostname)
		}
		vm := nbi.GetVMByID(vmIface.VM.ID)
		if vm == nil {
			return fmt.Errorf("VM ID %d for hostname IP %s not found", vmIface.VM.ID, hostname)
		}
		fs.Logger.Infof(fs.Ctx, "resolved hostname %s to VM: %s (ID: %d)", hostname, vm.Name, vm.ID)
		if ifaceName != "" {
			targetIface := nbi.GetVMInterfaceByVMIDAndName(vm.ID, ifaceName)
			if targetIface == nil {
				return fmt.Errorf("target interface %q not found on VM %q", ifaceName, vm.Name)
			}
			*objType = constants.ContentTypeVirtualizationVMInterface
			*objID = targetIface.ID
			fs.Logger.Infof(fs.Ctx, "resolved target VM interface: %s (ID: %d)", targetIface.Name, targetIface.ID)
		}
	case constants.ContentTypeDcimInterface:
		iface := nbi.GetInterfaceByID(ipObj.AssignedObjectID)
		if iface == nil {
			return fmt.Errorf("Device interface ID %d for hostname IP %s not found", ipObj.AssignedObjectID, hostname)
		}
		device := nbi.GetDeviceByID(iface.Device.ID)
		if device == nil {
			return fmt.Errorf("Device ID %d for hostname IP %s not found", iface.Device.ID, hostname)
		}
		fs.Logger.Infof(fs.Ctx, "resolved hostname %s to Device: %s (ID: %d)", hostname, device.Name, device.ID)
		if ifaceName != "" {
			targetIface, ok := nbi.GetInterface(ifaceName, device.ID)
			if !ok {
				return fmt.Errorf("target interface %q not found on Device %q", ifaceName, device.Name)
			}
			*objType = constants.ContentTypeDcimInterface
			*objID = targetIface.ID
			fs.Logger.Infof(
				fs.Ctx, "resolved target Device interface: %s (ID: %d)", targetIface.Name, targetIface.ID,
			)
		}
	default:
		fs.Logger.Warningf(
			fs.Ctx, "hostname IP %s has unsupported assigned type %s, VIPs will be unassigned",
			hostname, ipObj.AssignedObjectType,
		)
	}
	return nil
}

// parseDestination extracts IP and mask bits from F5 destination format.
// F5 destination format: "/partition/address:port" or "/partition/address%route-domain:port".
func parseDestination(destination, mask string) (string, int, error) {
	// Remove partition prefix: "/Common/10.0.0.1:80" -> "10.0.0.1:80"
	parts := strings.Split(destination, "/")
	addrPort := parts[len(parts)-1]

	// Remove route domain if present: "10.0.0.1%1:80" -> "10.0.0.1:80"
	if idx := strings.Index(addrPort, "%"); idx != -1 {
		rest := ""
		if colonIdx := strings.LastIndex(addrPort, ":"); colonIdx > idx {
			rest = addrPort[colonIdx:]
		}
		addrPort = addrPort[:idx] + rest
	}

	// Split address and port: "10.0.0.1:80" -> "10.0.0.1"
	var ip string
	if strings.Contains(addrPort, ":") {
		lastColon := strings.LastIndex(addrPort, ":")
		ip = addrPort[:lastColon]
	} else {
		ip = addrPort
	}

	// Handle IPv6 brackets
	ip = strings.TrimPrefix(ip, "[")
	ip = strings.TrimSuffix(ip, "]")

	if ip == "" {
		return "", 0, fmt.Errorf("empty IP in destination: %s", destination)
	}

	maskBits, err := utils.MaskToBits(mask)
	if err != nil {
		// Default to /32 for individual VIPs
		maskBits = constants.MaxIPv4MaskBits
		if strings.Contains(ip, ":") {
			maskBits = constants.MaxIPv6MaskBits
		}
	}

	return ip, maskBits, nil
}
