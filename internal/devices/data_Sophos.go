// Code generated by go generate; DO NOT EDIT.
package devices

var DeviceTypesMapSophos = map[string]*DeviceData{
    "XG 650": {
        Manufacturer: "Sophos",
        Model: "XG 650",
        Slug: "sophos-xg-650",
        UHeight: 2,
        PartNumber: "",
        IsFullDepth: true,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "parent",
        Weight: 0,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
        },
        ConsoleServerPorts: []ConsoleServerPort{
            { Name: "COM", Type: "rj-45", Label: "" },
            { Name: "Front USB 1", Type: "usb-a", Label: "" },
            { Name: "Front USB 2", Type: "usb-a", Label: "" },
            { Name: "Rear USB 1", Type: "usb-a", Label: "" },
            { Name: "Rear USB 2", Type: "usb-a", Label: "" },
        },
        PowerPorts: []PowerPort{
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
            { Name: "A", Label: "", Position: "A" },
            { Name: "B", Label: "", Position: "B" },
            { Name: "C", Label: "", Position: "C" },
            { Name: "D", Label: "", Position: "D" },
            { Name: "Front Hard Disk 1", Label: "", Position: "FHDD 1" },
            { Name: "Front Hard Disk 2", Label: "", Position: "FHDD 2" },
            { Name: "Rear Fan 1", Label: "", Position: "RF 1" },
            { Name: "Rear Fan 2", Label: "", Position: "RF 2" },
            { Name: "Rear Fan 3", Label: "", Position: "RF 3" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "MGMT 1", Label: "", Type: "1000base-t", MgmtOnly: true },
            { Name: "MGMT 2", Label: "", Type: "1000base-t", MgmtOnly: true },
        },
    },
    "XGS 2100": {
        Manufacturer: "Sophos",
        Model: "XGS 2100",
        Slug: "sophos-xgs-2100",
        UHeight: 1,
        PartNumber: "",
        IsFullDepth: false,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 4.7,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "COM Serial", Type: "rj-45", Label: "", Poe: false },
            { Name: "COM USB", Type: "usb-micro-b", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "internal-PS", Label: "", Type: "iec-60320-c14", MaximumDraw: 201, AllocatedDraw: 50 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
            { Name: "A", Label: "", Position: "A" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "MGMT", Label: "", Type: "1000base-t", MgmtOnly: true },
            { Name: "F1", Label: "", Type: "1000base-x-sfp", MgmtOnly: false },
            { Name: "F2", Label: "", Type: "1000base-x-sfp", MgmtOnly: false },
            { Name: "LAN 1 (LAN)", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 2 (WAN)", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 3", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 4", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 5", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 6", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 7", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 8", Label: "", Type: "1000base-t", MgmtOnly: false },
        },
    },
    "XGS 2300": {
        Manufacturer: "Sophos",
        Model: "XGS 2300",
        Slug: "sophos-xgs-2300",
        UHeight: 1,
        PartNumber: "",
        IsFullDepth: false,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 4.7,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "COM Serial", Type: "rj-45", Label: "", Poe: false },
            { Name: "COM USB", Type: "usb-micro-b", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "internal-PS", Label: "", Type: "iec-60320-c14", MaximumDraw: 201, AllocatedDraw: 50 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
            { Name: "A", Label: "", Position: "A" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "MGMT", Label: "", Type: "1000base-t", MgmtOnly: true },
            { Name: "F1", Label: "", Type: "1000base-x-sfp", MgmtOnly: false },
            { Name: "F2", Label: "", Type: "1000base-x-sfp", MgmtOnly: false },
            { Name: "LAN 1 (LAN)", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 2 (WAN)", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 3", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 4", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 5", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 6", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 7", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 8", Label: "", Type: "1000base-t", MgmtOnly: false },
        },
    },
    "XGS 3100": {
        Manufacturer: "Sophos",
        Model: "XGS 3100",
        Slug: "sophos-xgs-3100",
        UHeight: 1,
        PartNumber: "",
        IsFullDepth: false,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 4.7,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "COM Serial", Type: "rj-45", Label: "", Poe: false },
            { Name: "COM USB", Type: "usb-micro-b", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "internal-PS", Label: "", Type: "iec-60320-c14", MaximumDraw: 201, AllocatedDraw: 50 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
            { Name: "A", Label: "", Position: "A" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "MGMT", Label: "", Type: "1000base-t", MgmtOnly: true },
            { Name: "F1", Label: "", Type: "10gbase-x-sfpp", MgmtOnly: false },
            { Name: "F2", Label: "", Type: "10gbase-x-sfpp", MgmtOnly: false },
            { Name: "F3", Label: "", Type: "1000base-x-sfp", MgmtOnly: false },
            { Name: "F4", Label: "", Type: "1000base-x-sfp", MgmtOnly: false },
            { Name: "LAN 1 (LAN)", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 2 (WAN)", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 3", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 4", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 5", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 6", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 7", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 8", Label: "", Type: "1000base-t", MgmtOnly: false },
        },
    },
    "XGS 3300": {
        Manufacturer: "Sophos",
        Model: "XGS 3300",
        Slug: "sophos-xgs-3300",
        UHeight: 1,
        PartNumber: "",
        IsFullDepth: false,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 4.7,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "COM Serial", Type: "rj-45", Label: "", Poe: false },
            { Name: "COM USB", Type: "usb-micro-b", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "internal-PS", Label: "", Type: "iec-60320-c14", MaximumDraw: 201, AllocatedDraw: 50 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
            { Name: "A", Label: "", Position: "A" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "MGMT", Label: "", Type: "1000base-t", MgmtOnly: true },
            { Name: "F1", Label: "", Type: "10gbase-x-sfpp", MgmtOnly: false },
            { Name: "F2", Label: "", Type: "10gbase-x-sfpp", MgmtOnly: false },
            { Name: "F3", Label: "", Type: "1000base-x-sfp", MgmtOnly: false },
            { Name: "F4", Label: "", Type: "1000base-x-sfp", MgmtOnly: false },
            { Name: "LAN 1 (LAN)", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 2 (WAN)", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 3", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 4", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 5", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 6", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 7", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN 8", Label: "", Type: "1000base-t", MgmtOnly: false },
        },
    },
}
