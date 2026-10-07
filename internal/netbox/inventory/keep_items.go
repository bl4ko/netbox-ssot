package inventory

import (
	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// KeepVMNetworkObjects keeps the network objects NetBox already has for vm (its
// interfaces, their IP and MAC addresses, and the prefixes of those IPs) out of
// the orphans. Sources call it when they could not read the VM's network data,
// so that unknown data is not mistaken for removed data.
func (nbi *NetboxInventory) KeepVMNetworkObjects(vm *objects.VM) {
	ifaceNames := nbi.keepVMInterfaces(vm)
	ipAddresses := nbi.keepVMIPAddresses(vm, ifaceNames)
	nbi.keepVMMACAddresses(vm, ifaceNames)
	nbi.keepPrefixesOfIPAddresses(ipAddresses)
}

func (nbi *NetboxInventory) keepVMInterfaces(vm *objects.VM) []string {
	nbi.vmInterfacesLock.Lock()
	defer nbi.vmInterfacesLock.Unlock()
	ifaceNames := make([]string, 0, len(nbi.vmInterfacesIndexByVMIdAndName[vm.ID]))
	for name, iface := range nbi.vmInterfacesIndexByVMIdAndName[vm.ID] {
		nbi.OrphanManager.RemoveItem(iface)
		ifaceNames = append(ifaceNames, name)
	}
	return ifaceNames
}

func (nbi *NetboxInventory) keepVMIPAddresses(vm *objects.VM, ifaceNames []string) []*objects.IPAddress {
	nbi.ipAddressesLock.Lock()
	defer nbi.ipAddressesLock.Unlock()
	ipAddresses := make([]*objects.IPAddress, 0)
	vmIndex := nbi.ipAddressesIndex[constants.ContentTypeVirtualizationVirtualMachine]
	for _, ifaceName := range ifaceNames {
		for _, ipAddress := range vmIndex[ifaceName][vm.Name] {
			nbi.OrphanManager.RemoveItem(ipAddress)
			ipAddresses = append(ipAddresses, ipAddress)
		}
	}
	return ipAddresses
}

func (nbi *NetboxInventory) keepVMMACAddresses(vm *objects.VM, ifaceNames []string) {
	nbi.macAddressesLock.Lock()
	defer nbi.macAddressesLock.Unlock()
	vmIndex := nbi.macAddressesIndex[constants.ContentTypeVirtualizationVirtualMachine]
	for _, ifaceName := range ifaceNames {
		for _, macAddress := range vmIndex[ifaceName][vm.Name] {
			nbi.OrphanManager.RemoveItem(macAddress)
		}
	}
}

func (nbi *NetboxInventory) keepPrefixesOfIPAddresses(ipAddresses []*objects.IPAddress) {
	nbi.prefixesLock.Lock()
	defer nbi.prefixesLock.Unlock()
	for _, ipAddress := range ipAddresses {
		prefix, _, err := utils.GetPrefixAndMaskFromIPAddress(ipAddress.Address)
		if err != nil {
			continue
		}
		vrfID := 0
		if ipAddress.VRF != nil {
			vrfID = ipAddress.VRF.ID
		}
		if nbPrefix, ok := nbi.prefixesIndexByPrefix[prefix][vrfID]; ok {
			nbi.OrphanManager.RemoveItem(nbPrefix)
		}
	}
}

// KeepPlatform keeps platform out of the orphans while a source keeps it on a VM
// without being able to read the VM's current platform.
func (nbi *NetboxInventory) KeepPlatform(platform *objects.Platform) {
	nbi.platformsLock.Lock()
	defer nbi.platformsLock.Unlock()
	nbi.OrphanManager.RemoveItem(platform)
}
