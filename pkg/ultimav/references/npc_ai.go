package references

type AiType int

const (
	// Fixed stay put
	Fixed AiType = 0

	// Wander wander a small radius
	Wander = 1

	// BigWander wander a wide radius
	BigWander = 2

	// ChildRunAway children are playful and run away from you
	ChildRunAway = 3

	// CustomAi special conversation for merchants
	CustomAi = 4

	// ExtortOrAttackOrFollow you better pay the guard or they will attack!
	ExtortOrAttackOrFollow = 6

	// DrudgeWorthThing he's a jerk - he will attack at first opportunity, and get as close to you as possible
	DrudgeWorthThing = 7

	// HorseWander horses wander if they aren't tied up
	HorseWander = 100

	// FollowAroundAndBeAnnoyingThenNeverSeeAgain people who are freed walk close to you and are annoying - but won't talk
	FollowAroundAndBeAnnoyingThenNeverSeeAgain = 101

	// SmallWanderWantsToChat they will wander in a small radius, but if they are next to you, then they will want to chat
	SmallWanderWantsToChat = 102

	// Begging beggars wander, and when next to Avatar they want to talk
	Begging = 103

	// GenericExtortingGuard They generally ask for 30-60gp tribute
	GenericExtortingGuard = 104

	// HalfYourGoldExtortingGuard These jerks want half your gold! Generally found in Minoc
	HalfYourGoldExtortingGuard = 105

	// MerchantBuyingSelling Merchant that doesn't move
	MerchantBuyingSelling = 106

	// MerchantBuyingSellingCustom I THINK this is a merchant that sells, but also wanders a bit
	MerchantBuyingSellingCustom = 107

	// MerchantBuyingSellingWander Merchants that wander their own store, but still sell
	MerchantBuyingSellingWander = 108

	// FixedExceptAttackWhenIsWantedByThePoPo They are Fixed unless wanted by the popo, at which time they seek and attack
	FixedExceptAttackWhenIsWantedByThePoPo = 109

	// BlackthornGuardWander Blackthorn guards will check in for password, with a small wander
	BlackthornGuardWander = 110

	// BlackthornGuardFixed Blackthorn guard stays put, but will ask password if prompted
	BlackthornGuardFixed = 111
	StoneGargoyleTrigger = 112
)
