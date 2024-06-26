// Code generated by go generate; DO NOT EDIT.
package devices

var DeviceTypesMapSupermicro = map[string]*DeviceData{
    "AS-1114S-WN10RT": {
        Manufacturer: "Supermicro",
        Model: "AS-1114S-WN10RT",
        Slug: "supermicro-as-1114s-wn10rt",
        UHeight: 1,
        PartNumber: "",
        IsFullDepth: true,
        Airflow: "",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 0,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "Serial", Type: "de-9", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "Power 1", Label: "", Type: "iec-60320-c14", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Power 2", Label: "", Type: "iec-60320-c14", MaximumDraw: 0, AllocatedDraw: 0 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "Gig-E 1", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "Gig-E 2", Label: "", Type: "1000base-t", MgmtOnly: true },
        },
    },
    "AS-1123US-TR4": {
        Manufacturer: "Supermicro",
        Model: "AS-1123US-TR4",
        Slug: "supermicro-as-1123us-tr4",
        UHeight: 1,
        PartNumber: "",
        IsFullDepth: true,
        Airflow: "",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 0,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "Serial", Type: "de-9", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "Power 1", Label: "", Type: "iec-60320-c14", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Power 2", Label: "", Type: "iec-60320-c14", MaximumDraw: 0, AllocatedDraw: 0 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "Gig-E 1", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "Gig-E 2", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "Gig-E 3", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "Gig-E 4", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "BMC", Label: "", Type: "1000base-t", MgmtOnly: true },
        },
    },
    "IoT SuperServer SYS-510D-8C-FN6P": {
        Manufacturer: "Supermicro",
        Model: "IoT SuperServer SYS-510D-8C-FN6P",
        Slug: "supermicro-iot-superserver-sys-510d-8c-fn6p",
        UHeight: 1,
        PartNumber: "SYS-510D-8C-FN6P",
        IsFullDepth: false,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 4.54,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "PSW", Label: "", Type: "iec-60320-c14", MaximumDraw: 200, AllocatedDraw: 0 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "BMC", Label: "", Type: "1000base-t", MgmtOnly: true },
            { Name: "LAN1", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN2", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN3", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN4", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "LAN5", Label: "", Type: "25gbase-x-sfp28", MgmtOnly: false },
            { Name: "LAN6", Label: "", Type: "25gbase-x-sfp28", MgmtOnly: false },
        },
    },
    "SYS-1019P-WTR": {
        Manufacturer: "Supermicro",
        Model: "SYS-1019P-WTR",
        Slug: "supermicro-sys-1019p-wtr",
        UHeight: 1,
        PartNumber: "",
        IsFullDepth: true,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 0,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "COM1", Type: "de-9", Label: "Rear", Poe: false },
            { Name: "COM2", Type: "de-9", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "PSU1", Label: "", Type: "iec-60320-c14", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "PSU2", Label: "", Type: "iec-60320-c14", MaximumDraw: 0, AllocatedDraw: 0 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "Gig-E 1", Label: "", Type: "10gbase-t", MgmtOnly: false },
            { Name: "Gig-E 2", Label: "", Type: "10gbase-t", MgmtOnly: false },
            { Name: "Gig-E 3", Label: "", Type: "1000base-t", MgmtOnly: true },
        },
    },
    "SYS-2028U-E1CNR4T&#43;": {
        Manufacturer: "Supermicro",
        Model: "SYS-2028U-E1CNR4T&#43;",
        Slug: "supermicro-sys-2028u-e1cnr4t-plus",
        UHeight: 2,
        PartNumber: "",
        IsFullDepth: true,
        Airflow: "",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 0,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "Serial", Type: "de-9", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "Power 1", Label: "", Type: "iec-60320-c14", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Power 2", Label: "", Type: "iec-60320-c14", MaximumDraw: 0, AllocatedDraw: 0 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
            { Name: "PCI-E 1", Label: "", Position: "1" },
            { Name: "PCI-E 2", Label: "", Position: "2" },
            { Name: "PCI-E 3", Label: "", Position: "3" },
            { Name: "PCI-E 4", Label: "", Position: "4" },
            { Name: "PCI-E LP 1", Label: "", Position: "5" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "Te1", Label: "", Type: "10gbase-t", MgmtOnly: false },
            { Name: "Te2", Label: "", Type: "10gbase-t", MgmtOnly: false },
            { Name: "Te3", Label: "", Type: "10gbase-t", MgmtOnly: false },
            { Name: "Te4", Label: "", Type: "10gbase-t", MgmtOnly: false },
            { Name: "BMC", Label: "", Type: "1000base-t", MgmtOnly: true },
        },
    },
    "SuperServer 1018R-WC0R": {
        Manufacturer: "Supermicro",
        Model: "SuperServer 1018R-WC0R",
        Slug: "supermicro-superserver-1018r-wc0r",
        UHeight: 1,
        PartNumber: "SYS-1018R-WC0R",
        IsFullDepth: true,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 11.3,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "Serial", Type: "de-9", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
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
            { Name: "PSU 1", Label: "", Position: "PSU-1" },
            { Name: "PSU 2", Label: "", Position: "PSU-2" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "Gig-E 1", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "Gig-E 2", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "BMC", Label: "", Type: "1000base-t", MgmtOnly: true },
        },
    },
    "SuperServer 5018D-MTF": {
        Manufacturer: "Supermicro",
        Model: "SuperServer 5018D-MTF",
        Slug: "supermicro-superserver-5018d-mtf",
        UHeight: 1,
        PartNumber: "SYS-5018D-MTF",
        IsFullDepth: true,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 13.8,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "Serial", Type: "de-9", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "Power 1", Label: "", Type: "iec-60320-c14", MaximumDraw: 350, AllocatedDraw: 0 },
        },
        PowerOutlets: []PowerOutlet{
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
            { Name: "Gig-E 1", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "Gig-E 2", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "BMC", Label: "", Type: "1000base-t", MgmtOnly: true },
        },
    },
}
