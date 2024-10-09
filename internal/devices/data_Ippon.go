// Code generated by go generate; DO NOT EDIT.
package devices

var DeviceTypesMapIppon = map[string]*DeviceData{
    "Innova RT 3000": {
        Manufacturer: "Ippon",
        Model: "Innova RT 3000",
        Slug: "ippon-innova-rt-3000",
        UHeight: 2,
        PartNumber: "621781",
        IsFullDepth: true,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 28.8,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "Serial", Type: "de-9", Label: "", Poe: false },
            { Name: "USB", Type: "usb-b", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "Inlet", Label: "", Type: "iec-60320-c14", MaximumDraw: 2700, AllocatedDraw: 0 },
        },
        PowerOutlets: []PowerOutlet{
            { Name: "Outlet 1", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 2", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 3", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 4", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 5", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 6", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 7", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 8", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "16A Outlet 1", Type: "iec-60320-c19", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
            { Name: "IntelligentSlot", Label: "", Position: "Rear" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
        },
    },
    "Innova RT II 6000": {
        Manufacturer: "Ippon",
        Model: "Innova RT II 6000",
        Slug: "ippon-innova-rt-ii-6000",
        UHeight: 5,
        PartNumber: "1005639",
        IsFullDepth: true,
        Airflow: "front-to-rear",
        FrontImage: false,
        RearImage: false,
        SubdeviceRole: "",
        Weight: 59.1,
        WeightUnit: "",
        IsPowered: false,
        ConsolePorts: []ConsolePort{
            { Name: "Serial", Type: "de-9", Label: "", Poe: false },
            { Name: "USB", Type: "usb-b", Label: "", Poe: false },
        },
        ConsoleServerPorts: []ConsoleServerPort{
        },
        PowerPorts: []PowerPort{
            { Name: "Inlet", Label: "", Type: "iec-60309-3p-n-e-4h", MaximumDraw: 6000, AllocatedDraw: 0 },
        },
        PowerOutlets: []PowerOutlet{
            { Name: "Outlet 1", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 2", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 3", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 4", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 5", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "Outlet 6", Type: "iec-60320-c13", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "16A Outlet 1", Type: "iec-60320-c19", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
            { Name: "16A Outlet 2", Type: "iec-60320-c19", Label: "", PowerPort: "Inlet", FeedLeg: "", MaximumDraw: 0, AllocatedDraw: 0 },
        },
        FrontPorts: []FrontPort{
        },
        RearPorts: []RearPort{
        },
        ModuleBays: []ModuleBay{
            { Name: "IntelligentSlot", Label: "IntelligentSlot", Position: "Rear" },
        },
			  DeviceBays: []DeviceBay{
        },
        InventoryItems: []InventoryItem{
        },
        Interfaces: []Interface{
        },
    },
}