package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"

type Equipment int

func (e Equipment) ID() int {
	return int(e)
}

func (e Equipment) Type() ItemType {
	return ItemTypeEquipment
}

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
	// FlamPor          Equipment = 48
	// VasFlam          Equipment = 49
	// InCorp           Equipment = 50
	// UusNox           Equipment = 51
	// UusZu            Equipment = 52
	// UusFlam          Equipment = 53
	// UusSanct         Equipment = 54
	NoEquipment Equipment = 255 // 0xFF in decimal
)

func (e Equipment) GetSpriteIndex() indexes.SpriteIndex {
	if e >= LeatherHelm && e <= SpikedHelm {
		return indexes.ItemHelm
	}
	if e >= SmallShield && e <= JewelShield {
		return indexes.ItemShield
	}
	if e >= ClothArmour && e <= MysticArmour {
		return indexes.ItemArmour
	}
	if e >= Dagger && e <= MysticSword {
		return indexes.ItemWeapon
	}
	if e >= RingInvisibility && e <= RingRegeneration {
		return indexes.ItemRing
	}
	if e >= AmuletOfTurning && e <= Ankh {
		return indexes.ItemAnkh
	}
	// if e >= FlamPor && e <= UusFlam {
	// return indexes.ItemScroll
	// }
	return indexes.Avatar
}
