package references

type QuestItem int

func (q QuestItem) ID() int {
	return int(q)
}

func (q QuestItem) Type() ItemType {
	return ItemTypeQuestItem
}

const (
	Amulet  QuestItem = 0
	Crown   QuestItem = 1
	Sceptre QuestItem = 2
)
