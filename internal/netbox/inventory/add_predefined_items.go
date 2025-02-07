package inventory

import (
	"context"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

func (nbi *NetboxInventory) AddContainerDeviceRole(
	ctx context.Context,
) (*objects.DeviceRole, error) {
	newRole, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			Description: constants.DeviceRoleContainerDescription,
		},
		Name:   constants.DeviceRoleContainer,
		Slug:   utils.Slugify(constants.DeviceRoleContainer),
		Color:  constants.DeviceRoleContainerColor,
		VMRole: true,
	})

	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (nbi *NetboxInventory) AddFirewallDeviceRole(
	ctx context.Context,
) (*objects.DeviceRole, error) {
	newRole, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			Description: constants.DeviceRoleFirewallDescription,
		},
		Name:   constants.DeviceRoleFirewall,
		Slug:   utils.Slugify(constants.DeviceRoleFirewall),
		Color:  constants.DeviceRoleFirewallColor,
		VMRole: false,
	})

	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (nbi *NetboxInventory) AddSwitchDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	newRole, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			Description: constants.DeviceRoleSwitchDescription,
		},
		Name:   constants.DeviceRoleSwitch,
		Slug:   utils.Slugify(constants.DeviceRoleSwitch),
		Color:  constants.DeviceRoleSwitchColor,
		VMRole: false,
	})

	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (nbi *NetboxInventory) AddServerDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	newRole, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			Description: constants.DeviceRoleServerDescription,
		},
		Name:   constants.DeviceRoleServer,
		Slug:   utils.Slugify(constants.DeviceRoleServer),
		Color:  constants.DeviceRoleServerColor,
		VMRole: false,
	})

	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (nbi *NetboxInventory) AddVMDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	newRole, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			Description: constants.DeviceRoleVMDescription,
		},
		Name:   constants.DeviceRoleVM,
		Slug:   utils.Slugify(constants.DeviceRoleVM),
		Color:  constants.DeviceRoleVMColor,
		VMRole: false,
	})

	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (nbi *NetboxInventory) AddVMTemplateDeviceRole(
	ctx context.Context,
) (*objects.DeviceRole, error) {
	newRole, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
		NetboxObject: objects.NetboxObject{
			Description: constants.DeviceRoleVMTemplateDescription,
		},
		Name:   constants.DeviceRoleVMTemplate,
		Slug:   utils.Slugify(constants.DeviceRoleVMTemplate),
		Color:  constants.DeviceRoleVMTemplateColor,
		VMRole: false,
	})

	if err != nil {
		return nil, err
	}
	return newRole, nil
}
