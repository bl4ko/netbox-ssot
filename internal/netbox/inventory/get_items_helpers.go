package inventory

import (
	"context"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

func (nbi *NetboxInventory) GetContainerDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	if role, ok := nbi.GetDeviceRole(constants.DeviceRoleContainer); ok {
		return role, nil
	}
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

func (nbi *NetboxInventory) GetFirewallDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	if role, ok := nbi.GetDeviceRole(constants.DeviceRoleFirewall); ok {
		return role, nil
	}
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

func (nbi *NetboxInventory) GetSwitchDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	if role, ok := nbi.GetDeviceRole(constants.DeviceRoleSwitch); ok {
		return role, nil
	}
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

func (nbi *NetboxInventory) GetServerDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	if role, ok := nbi.GetDeviceRole(constants.DeviceRoleServer); ok {
		return role, nil
	}
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

func (nbi *NetboxInventory) GetVMDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	if role, ok := nbi.GetDeviceRole(constants.DeviceRoleVM); ok {
		return role, nil
	}
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

func (nbi *NetboxInventory) GetVMTemplateDeviceRole(ctx context.Context) (*objects.DeviceRole, error) {
	if role, ok := nbi.GetDeviceRole(constants.DeviceRoleVMTemplate); ok {
		return role, nil
	}
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
