// Code generated by go generate; DO NOT EDIT.
package devices

var DeviceTypesMapTeltonika = map[string]*DeviceData{
    "RUT240": {
        Manufacturer: "Teltonika",
        Model: "RUT240",
        Slug: "teltonika-rut240",
        UHeight: 0,
        PartNumber: "RUT240",
        IsFullDepth: false,
        Airflow: "",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 125,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "PS0", Label: "", Type: "dc-terminal", MaximumDraw: 7, AllocatedDraw: 0 },
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
            { Name: "lan", Label: "", Type: "100base-tx", MgmtOnly: false },
            { Name: "mob1s1a1", Label: "", Type: "lte", MgmtOnly: false },
            { Name: "wan", Label: "", Type: "100base-tx", MgmtOnly: false },
        },
    },
    "TRB500": {
        Manufacturer: "Teltonika",
        Model: "TRB500",
        Slug: "teltonika-trb500",
        UHeight: 0,
        PartNumber: "TRB500",
        IsFullDepth: false,
        Airflow: "",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 241,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "PS0", Label: "", Type: "dc-terminal", MaximumDraw: 6, AllocatedDraw: 0 },
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
            { Name: "lan", Label: "", Type: "1000base-t", MgmtOnly: false },
            { Name: "mob1s1a1", Label: "", Type: "lte", MgmtOnly: false },
        },
    },
}
