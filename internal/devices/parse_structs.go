package devices

type ConsolePort struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Label string `yaml:"label"`
	Poe   bool   `yaml:"poe"`
}

type ConsoleServerPort struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Label string `yaml:"label"`
}

type PowerPort struct {
	Name          string  `yaml:"name"`
	Label         string  `yaml:"label"`
	Type          string  `yaml:"type"`
	MaximumDraw   float64 `yaml:"maximum_draw"`
	AllocatedDraw float64 `yaml:"allocated_draw"`
}

type PowerOutlet struct {
	Name          string  `yaml:"name"`
	Type          string  `yaml:"type"`
	Label         string  `yaml:"label"`
	PowerPort     string  `yaml:"power_port"`
	FeedLeg       string  `yaml:"feed_leg"`
	MaximumDraw   float64 `yaml:"maximum_draw"`
	AllocatedDraw float64 `yaml:"allocated_draw"`
}

type FrontPort struct {
	Name             string `yaml:"name"`
	Label            string `yaml:"label"`
	Type             string `yaml:"type"`
	RearPort         string `yaml:"rear_port"`
	RearPortPosition int    `yaml:"rear_port_position"`
}

type RearPort struct {
	Name      string `yaml:"name"`
	Label     string `yaml:"label"`
	Type      string `yaml:"type"`
	Positions int    `yaml:"positions"`
	Poe       bool   `yaml:"poe"`
}

type ModuleBay struct {
	Name     string `yaml:"name"`
	Label    string `yaml:"label"`
	Position string `yaml:"position"`
}

type DeviceBay struct {
	Name  string `yaml:"name"`
	Label string `yaml:"label"`
}

type InventoryItem struct {
	Name         string `yaml:"name"`
	Label        string `yaml:"label"`
	Manufacturer string `yaml:"manufacturer"`
	PartID       string `yaml:"part_id"`
}

type Interface struct {
	Name     string `yaml:"name"`
	Label    string `yaml:"label"`
	Type     string `yaml:"type"`
	MgmtOnly bool   `yaml:"mgmt_only"`
}

type DeviceData struct {
	Manufacturer  string `yaml:"manufacturer"`
	Model         string `yaml:"model"`
	Slug          string `yaml:"slug"`
	UHeight       int    `yaml:"u_height"`
	PartNumber    string `yaml:"part_number"`
	IsFullDepth   bool   `yaml:"is_full_depth"`
	Airflow       string `yaml:"airflow"`
	FrontImage    bool   `yaml:"front_image"`
	RearImage     bool   `yaml:"rear_image"`
	SubdeviceRole string `yaml:"subdevice_role"`
	// Comments           string              `yaml:"comments"` // This breaks formatting also not needed at all
	Weight             float64             `yaml:"weight"`
	WeightUnit         string              `yaml:"weight-unit"`
	IsPowered          bool                `yaml:"is-powered"`
	ConsolePorts       []ConsolePort       `yaml:"console-ports"`
	ConsoleServerPorts []ConsoleServerPort `yaml:"console-server-ports"`
	PowerPorts         []PowerPort         `yaml:"power-ports"`
	PowerOutlets       []PowerOutlet       `yaml:"power-outlets"`
	FrontPorts         []FrontPort         `yaml:"front-ports"`
	RearPorts          []RearPort          `yaml:"rear-ports"`
	ModuleBays         []ModuleBay         `yaml:"module-bays"`
	DeviceBays         []DeviceBay         `yaml:"device-bays"`
	InventoryItems     []InventoryItem     `yaml:"inventory-items"`
	Interfaces         []Interface         `yaml:"interfaces"`
}
