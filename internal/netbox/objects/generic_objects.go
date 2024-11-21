package objects

type IDItem interface {
	GetID() int
}

type OrphanItem interface {
	GetNetboxObject() *NetboxObject
	GetID() int
}
