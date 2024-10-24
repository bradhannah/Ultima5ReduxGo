package indexes

type SpriteIndex int

const (
	Nothing                  SpriteIndex = 0
	Water                                = 1
	Swamp                                = 4
	Grass                                = 5
	Desert                               = 7
	SmallMountains                       = 12
	LeftDesert2                          = 30
	RightDesert2                         = 31
	PlowedField                          = 44
	WheatInField                         = 45
	Cactus                               = 47
	BrickFloor                           = 68
	LargeRockWall                        = 77
	StoneBrickWallSecret                 = 78
	StoneBrickWall                       = 79
	Telescope                            = 89
	TrollBridgeHoriz                     = 106
	TrollBridgeVert                      = 107
	Grate                                = 134
	Manacles                             = 133
	Stocks                               = 132
	SimpleCross                          = 137
	StoneHeadstone                       = 138
	BrickFloorHole                       = 140
	Lava                                 = 143
	ChairFacingUp                        = 144
	ChairFacingRight                     = 145
	ChairFacingDown                      = 146
	ChairFacingLeft                      = 147
	TableMiddle                          = 149
	MagicLockDoor                        = 151
	MagicLockDoorWithView                = 152
	Portcullis                           = 153
	TableFoodTop                         = 154
	TableFoodBottom                      = 155
	TableFoodBoth                        = 156
	Mirror                               = 157
	MirrorAvatar                         = 158
	MirrorBroken                         = 159
	Well                                 = 161
	HitchingPost                         = 162
	LeftBed                              = 171
	RightSconce                          = 176
	LeftScone                            = 177
	Brazier                              = 178
	CampFire                             = 179
	RegularDoor                          = 184
	LockedDoor                           = 185
	RegularDoorView                      = 186
	LockedDoorView                       = 187
	Fireplace                            = 188
	LampPost                             = 189
	CandleOnTable                        = 190
	CookStove                            = 191
	LadderUp                             = 200
	LadderDown                           = 201
	FenceHoriz                           = 202
	FenceVert                            = 203
	Clock1                               = 250
	Clock2                               = 251
	StarPattern                          = 256
	Chest                                = 257
	ItemMoney                            = 258
	ItemPotion                           = 259
	ItemScroll                           = 260
	ItemWeapon                           = 261
	ItemShield                           = 262
	ItemKey                              = 263
	ItemGem                              = 264
	ItemHelm                             = 265
	ItemRing                             = 266
	ItemArmour                           = 267
	ItemAnkh                             = 268
	ItemTorch                            = 269
	WoodenBox                            = 270
	ItemFood                             = 271
	HorseRight                           = 272
	HorseLeft                            = 273
	HolyFloorSymbol                      = 278
	AvatarOnLadderUp                     = 279
	AvatarOnLadderDown                   = 280
	AvatarSleepingInBed                  = 282
	Avatar                               = 284
	DeadBody                             = 286
	BloodSpatter                         = 287
	PirateShip_Up                        = 300
	PirateShip_Right                     = 301
	PirateShip_Left                      = 302
	PirateShip_Down                      = 303
	AvatarSittingFacingUp                = 304
	AvatarSittingFacingRight             = 305
	AvatarSittingFacingDown              = 306
	AvatarSittingFacingLeft              = 307
	Carpet2_MagicCarpet                  = 283
	PoisonField                          = 488
	MagicField                           = 489
	FireField                            = 490
	ElectricField                        = 491
	Shard                                = 436
	Crown                                = 437
	Sceptre                              = 438
	Amulet                               = 439
	BrickWallArchway                     = 135
)
const (
	Wizard_KeyIndex                  SpriteIndex = 320
	Bard_KeyIndex                                = 324
	Fighter_KeyIndex                             = 328
	Avatar_KeyIndex                              = 332
	AvatarSittingAndEatingFacingDown             = 308
	AvatarSittingAndEatingFacingUp               = 312

	BardPlaying_KeyIndex   = 348
	TownsPerson_KeyIndex   = 336
	Ray_KeyIndex           = 400
	Daemon1_KeyIndex       = 472
	StoneGargoyle_KeyIndex = 440
	Bat_KeyIndex           = 404
	ShadowLord_KeyIndex    = 508
	Waterfall_KeyIndex     = 212
	Fountain_KeyIndex      = 216
	Beggar_KeyIndex        = 365
	Guard_KeyIndex         = 368
	Rat_KeyIndex           = 400
	Troll_KeyIndex         = 484
	Whirlpool_KeyIndex     = 492
)

func (s SpriteIndex) IsDoor() bool {
	return s == RegularDoor || s == LockedDoor || s == RegularDoorView || s == LockedDoorView || s == MagicLockDoorWithView || s == MagicLockDoor
}

func (s SpriteIndex) IsWindowedDoor() bool {
	return s == RegularDoorView || s == LockedDoorView || s == MagicLockDoorWithView
}
