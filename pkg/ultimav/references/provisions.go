package references

type Provision int

const NumberOfProvisions = 5

const (
	NoProvision Provision = -1
	Food                  = iota
	Gold
	Key
	Gem
	Torches
)
