package references

type Equipment int

func (e Equipment) ID() int {
	return int(e)
}

func (e Equipment) Type() ItemType {
	return ItemTypeEquipment
}

//goland:noinspection GoUnusedConst
const (
	BareHands        Equipment = -2
	LeatherHelm      Equipment = 0
	ChainCoif        Equipment = 1
	IronHelm         Equipment = 2
	SpikedHelm       Equipment = 3
	SmallShield      Equipment = 4
	LargeShield      Equipment = 5
	SpikedShield     Equipment = 6
	MagicShield      Equipment = 7
	JewelShield      Equipment = 8
	ClothArmour      Equipment = 9
	LeatherArmour    Equipment = 10
	RingMail         Equipment = 11
	ScaleMail        Equipment = 12
	ChainMail        Equipment = 13
	PlateMail        Equipment = 14
	MysticArmour     Equipment = 15
	Dagger           Equipment = 16
	Sling            Equipment = 17
	Club             Equipment = 18
	FlamingOil       Equipment = 19
	MainGauche       Equipment = 20
	Spear            Equipment = 21
	ThrowingAxe      Equipment = 22
	ShortSword       Equipment = 23
	Mace             Equipment = 24
	MorningStar      Equipment = 25
	Bow              Equipment = 26
	Arrows           Equipment = 27
	Crossbow         Equipment = 28
	Quarrels         Equipment = 29
	LongSword        Equipment = 30
	TwoHHammer       Equipment = 31
	TwoHAxe          Equipment = 32
	TwoHSword        Equipment = 33
	Halberd          Equipment = 34
	ChaosSword       Equipment = 35
	MagicBow         Equipment = 36
	SilverSword      Equipment = 37
	MagicAxe         Equipment = 38
	GlassSword       Equipment = 39
	JeweledSword     Equipment = 40
	MysticSword      Equipment = 41
	RingInvisibility Equipment = 42
	RingProtection   Equipment = 43
	RingRegeneration Equipment = 44
	AmuletOfTurning  Equipment = 45
	SpikedCollar     Equipment = 46
	Ankh             Equipment = 47
	NoEquipment      Equipment = 255 // 0xFF in decimal
)
