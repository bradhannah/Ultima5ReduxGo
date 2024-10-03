package legacy

const (
	BRIT_CBT              = "BRIT.CBT"
	BRIT_DAT              = "BRIT.DAT"
	BRIT_OOL              = "BRIT.OOL"
	CASTLE_DAT            = "CASTLE.DAT"
	CASTLE_NPC            = "CASTLE.NPC"
	CASTLE_TLK            = "CASTLE.TLK"
	DATA_OVL              = "DATA.OVL"
	DUNGEON_CBT           = "DUNGEON.CBT"
	DUNGEON_DAT           = "DUNGEON.DAT"
	DWELLING_DAT          = "DWELLING.DAT"
	DWELLING_NPC          = "DWELLING.NPC"
	DWELLING_TLK          = "DWELLING.TLK"
	INIT_GAM              = "INIT.GAM"
	KEEP_DAT              = "KEEP.DAT"
	KEEP_NPC              = "KEEP.NPC"
	KEEP_TLK              = "KEEP.TLK"
	LOOK2_DAT             = "LOOK2.DAT"
	NEW_SAVE_FILE         = "save.json"
	NEW_SAVE_SUMMARY_FILE = "summary.json"
	SAVED_GAM             = "SAVED.GAM"
	SAVED_OOL             = "SAVED.OOL"
	SHOPPE_DAT            = "SHOPPE.DAT"
	SIGNS_DAT             = "SIGNS.DAT"
	TOWNE_DAT             = "TOWNE.DAT"
	TOWNE_NPC             = "TOWNE.NPC"
	TOWNE_TLK             = "TOWNE.TLK"
	UNDER_DAT             = "UNDER.DAT"
	UNDER_OOL             = "UNDER.OOL"
	MISCMAPS_DAT          = "MISCMAPS.DAT"
)

// Slices for grouped files
var (
	TalkFiles     = []string{CASTLE_TLK, TOWNE_TLK, DWELLING_TLK, KEEP_TLK}
	NpcFiles      = []string{CASTLE_NPC, TOWNE_NPC, DWELLING_NPC, KEEP_NPC}
	SmallMapFiles = []string{CASTLE_DAT, TOWNE_DAT, DWELLING_DAT, KEEP_DAT}
)
