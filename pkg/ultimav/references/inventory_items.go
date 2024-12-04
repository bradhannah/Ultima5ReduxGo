package references

type InventoryItem struct {
	ItemName                   string `json:"ItemName"`
	ItemIndex                  int    `json:"ItemIndex"`
	ItemNameHighlight          string `json:"ItemNameHighlight"`
	ItemSprite                 int    `json:"ItemSprite"`
	ItemSpriteExposed          int    `json:"ItemSpriteExposed"`
	ItemDescription            string `json:"ItemDescription"`
	ItemDescriptionAttribution string `json:"ItemDescriptionAttribution"`
}
