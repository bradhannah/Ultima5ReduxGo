package party_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state/util"
)

/*
enum <ubyte> PartyStatus { InTheParty = 0x00, HasntJoinedYet = 0xFF, AtTheInn = 0x01 };
enum <ubyte> Location { };
enum <ubyte> Boolean { False = 0, True = 1 };
enum <ubyte> BritOrUnderworld { Britannia = 0, Underworld = 0xFF };
enum <ubyte> MoonstoneStatus { Buried = 0, InInventory = 0xFF };
enum <ubyte> DungeonOrientation { North = 0, East = 1, South = 2, West = 2 };
enum <ubyte> DungeonType { Cave = 1, Mine = 2, Dungeon = 3 };
*/

/*
CharacterClass

enum <ubyte> Class { Avatar = 'A', Fighter = 'F', Bard = 'B', Wizard = 'M' };
*/

const MaxCharactersInParty = 6

type CharacterClass byte

const (
	Avatar  CharacterClass = 'A'
	Fighter CharacterClass = 'F'
	Bard    CharacterClass = 'B'
	Wizard  CharacterClass = 'M'
)

var CharacterClasses = util.OrderedMapping[CharacterClass]{
	{Id: Avatar, FriendlyName: "Avatar"},
	{Id: Fighter, FriendlyName: "Fighter"},
	{Id: Bard, FriendlyName: "Bard"},
	{Id: Wizard, FriendlyName: "Wizard"},
}

type CharacterStatus byte

const (
	Good     CharacterStatus = 'G'
	Poisoned CharacterStatus = 'P'
	Charmed  CharacterStatus = 'C'
	Sleep    CharacterStatus = 'S'
	Dead     CharacterStatus = 'D'
)

var CharacterStatuses = util.OrderedMapping[CharacterStatus]{
	{Id: Good, FriendlyName: "Good"},
	{Id: Poisoned, FriendlyName: "Poisoned"},
	{Id: Charmed, FriendlyName: "Charmed"},
	{Id: Sleep, FriendlyName: "Sleep"},
	{Id: Dead, FriendlyName: "Dead"},
}

// Gender
// enum <ubyte> Gender { Male = 0x0B, Female = 0x0C };

type CharacterGender byte

const (
	Male   CharacterGender = 0x0B
	Female CharacterGender = 0x0C
)

var CharacterGenders = util.OrderedMapping[CharacterGender]{
	{Id: Male, FriendlyName: "Male"},
	{Id: Female, FriendlyName: "Female"},
}

/*
Other smaller definitions
*/
// type BritOrUnderworld byte
//
// const (
//	Britannia  BritOrUnderworld = 0xff
//	Underworld BritOrUnderworld = 0x00
// )
