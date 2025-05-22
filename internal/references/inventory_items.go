package references

import "github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"

type InventoryItem struct {
	ItemName                   string              `json:"ItemName"`
	ItemIndex                  int                 `json:"ItemIndex"`
	ItemNameHighlight          string              `json:"ItemNameHighlight"`
	ItemSprite                 indexes.SpriteIndex `json:"ItemSprite"`
	ItemSpriteExposed          indexes.SpriteIndex `json:"ItemSpriteExposed"`
	ItemDescription            string              `json:"ItemDescription"`
	ItemDescriptionAttribution string              `json:"ItemDescriptionAttribution"`
}
