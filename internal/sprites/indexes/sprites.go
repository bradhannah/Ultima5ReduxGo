package indexes

const StandardNumberOfAnimationFrames = 4

type SpriteIndex int

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	NoSprites                SpriteIndex = -1
	Bang                                 = 0
	Water1                               = 1
	Water2                               = 2
	WaterShallow                         = 3
	Swamp                                = 4
	Grass                                = 5
	Desert                               = 7
	SmallMountains                       = 12
	LeftDesert2                          = 30
	RightDesert2                         = 31
	PathUpDown                           = 32
	PathLeftRight                        = 33
	PathUpRight                          = 34
	PathDownRight                        = 35
	PathLeftDown                         = 36
	PathUpLeft                           = 37
	PathAllWays                          = 38
	PlowedField                          = 44
	WheatInField                         = 45
	Cactus                               = 47
	WoodenPlankHorizFloor                = 64
	BrickFloor                           = 68
	HexMetalGridFloor                    = 69
	WoodenPlankVert1Floor                = 72
	WoodenPlankVert2Floor                = 73
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
	Barrel                               = 166
	LeftBed                              = 171
	RightBed                             = 172
	RightSconce                          = 176
	LeftScone                            = 177
	Brazier                              = 178
	CampFire                             = 179
	CannonFacingUp                       = 180
	CannonFacingRight                    = 181
	CannonFacingDown                     = 182
	CannonFacingLeft                     = 183
	RegularDoor                          = 184
	LockedDoor                           = 185
	RegularDoorView                      = 186
	LockedDoorView                       = 187
	Fireplace                            = 188
	LampPost                             = 189
	CandleOnTable                        = 190
	CookStove                            = 191
	Stairs1                              = 196
	Stairs2                              = 197
	Stair3                               = 198
	Stairs4                              = 199
	LadderUp                             = 200
	LadderDown                           = 201
	FenceHoriz                           = 202
	FenceVert                            = 203
	Clock1                               = 250
	Clock2                               = 251
	StarPattern                          = 256
	Chest                                = 257
	ItemGold                             = 258
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
	AvatarRidingHorseRight               = 274
	AvatarRidingHorseLeft                = 275
	AvatarRidingCarpetRight              = 276
	AvatarRidingCarpetLeft               = 277
	HolyFloorSymbol                      = 278
	AvatarOnLadderUp                     = 279
	AvatarOnLadderDown                   = 280
	AvatarSleepingInBed                  = 282
	Avatar                               = 284
	DeadBody                             = 286
	BloodSpatter                         = 287
	FrigateUpUnfurled                    = 288
	FrigateLeftUnfurled                  = 289
	FrigateRightUnfurled                 = 290
	FrigateDownUnfurled                  = 291
	FrigateUpFurled                      = 292
	FrigateRightFurled                   = 293
	FrigateDownFurled                    = 294
	FrigateLeftFurled                    = 295
	SkiffUp                              = 296
	SkiffRight                           = 297
	SkiffDown                            = 298
	SkiffLeft                            = 299
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

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
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
	Manacles_Prisoner      = 356
)

func (s SpriteIndex) IsDoor() bool {
	return s == RegularDoor || s == LockedDoor || s == RegularDoorView || s == LockedDoorView || s == MagicLockDoorWithView || s == MagicLockDoor
}

func (s SpriteIndex) IsUnlockedDoor() bool {
	return s == RegularDoor || s == RegularDoorView
}

func (s SpriteIndex) IsWindowedDoor() bool {
	return s == RegularDoorView || s == LockedDoorView || s == MagicLockDoorWithView
}

func (s SpriteIndex) IsPushableFloor() bool {
	return s == BrickFloor || s == HexMetalGridFloor
}

func (s SpriteIndex) IsBed() bool {
	return s == LeftBed || s == AvatarSleepingInBed || s == RightBed
}

func (s SpriteIndex) IsStairs() bool {
	return s == Stair3 || s == Stairs4 || s == Stairs1 || s == Stairs2
}

func (s SpriteIndex) IsFrigateUnfurled() bool {
	return s == FrigateLeftUnfurled || s == FrigateRightUnfurled || s == FrigateDownUnfurled || s == FrigateUpUnfurled
}

func (s SpriteIndex) IsFrigateFurled() bool {
	return s == FrigateLeftFurled || s == FrigateRightFurled || s == FrigateDownFurled || s == FrigateUpFurled
}

func (s SpriteIndex) IsHorseUnBoarded() bool {
	return s == HorseLeft || s == HorseRight
}

func (s SpriteIndex) IsSkiff() bool {
	return s == SkiffDown || s == SkiffLeft || s == SkiffRight || s == SkiffUp
}

func (s SpriteIndex) IsMagicCarpetUnboarded() bool {
	return s == Carpet2_MagicCarpet
}

func (s SpriteIndex) IsPartOfAnimation(keyFrameIndex SpriteIndex) bool {
	const lowestPossibleAnimationIndex = 304
	if s < lowestPossibleAnimationIndex {
		return s == keyFrameIndex
	}

	if s >= keyFrameIndex && s <= keyFrameIndex+StandardNumberOfAnimationFrames {
		return true
	}
	return false
}
