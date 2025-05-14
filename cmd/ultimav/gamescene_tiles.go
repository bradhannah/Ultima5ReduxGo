package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const (
	xCenter = references.Coordinate(xTilesVisibleOnGameScreen / 2)
	yCenter = references.Coordinate(yTilesVisibleOnGameScreen / 2)
)

func (g *GameScene) getSmallCalculatedAvatarTileIndex(ogSpriteIndex indexes.SpriteIndex) indexes.SpriteIndex {
	if g.gameState.PartyVehicle != references.NoPartyVehicle {
		return g.gameState.PartyVehicle.GetSpriteByDirection(g.gameState.PreviousPartyVehicleDirection, g.gameState.PartyVehicleDirection)
	}
	return g.getSmallCalculatedNPCTileIndex(ogSpriteIndex, indexes.Avatar_KeyIndex, g.gameState.Position)
}

func (g *GameScene) getSmallCalculatedNPCTileIndex(ogSpriteIndex indexes.SpriteIndex, npcIndex indexes.SpriteIndex, spritePosition references.Position) indexes.SpriteIndex {
	switch ogSpriteIndex {
	case indexes.LeftBed:
		return indexes.AvatarSleepingInBed
	case indexes.ChairFacingRight, indexes.ChairFacingLeft, indexes.ChairFacingUp, indexes.ChairFacingDown:
		return g.getCorrectAvatarOnChairTile(ogSpriteIndex, &spritePosition)
	case indexes.LadderUp:
		return indexes.AvatarOnLadderUp
	case indexes.LadderDown:
		return indexes.AvatarOnLadderDown
	case indexes.Manacles:
		return indexes.Manacles_Prisoner
	}
	return npcIndex
}

func (g *GameScene) getSmallCalculatedTileIndex(ogSpriteIndex indexes.SpriteIndex, pos *references.Position) indexes.SpriteIndex {
	switch ogSpriteIndex {
	case indexes.Mirror:
		// is avatar in front of it
		pos := pos.GetPositionDown()
		if g.gameState.IsAvatarAtPosition(pos) {
			return indexes.MirrorAvatar
		}
	}
	return sprites.GetSpriteIndexWithAnimationBySpriteIndex(ogSpriteIndex, pos.GetHash())
}

func (g *GameScene) getCorrectAvatarOnChairTile(spriteIndex indexes.SpriteIndex, position *references.Position) indexes.SpriteIndex {
	switch spriteIndex {
	case indexes.ChairFacingRight:
		return indexes.AvatarSittingFacingRight
	case indexes.ChairFacingLeft:
		return indexes.AvatarSittingFacingLeft
	case indexes.ChairFacingUp, indexes.ChairFacingDown:
		return g.getCorrectAvatarEatingInChairTile(spriteIndex, position)
	}
	return spriteIndex
}

func (g *GameScene) setDrawBridge(theMap *game_state.LayeredMap, pos *references.Position, spriteIndex indexes.SpriteIndex) bool {
	const (
		leftXDrawBridge   = 14
		rightXDrawBridge  = 16
		topYDrawBridge    = 28
		bottomYDrawBridge = 29
	)

	if (pos.X >= leftXDrawBridge && pos.X <= rightXDrawBridge) && (pos.Y >= topYDrawBridge && pos.Y <= bottomYDrawBridge) {
		theMap.SetTileByLayer(game_state.MapOverrideLayer, pos, g.gameState.GetDrawBridgeWaterByTime(spriteIndex))
		return true
	}
	return false
}

func (g *GameScene) getCorrectAvatarEatingInChairTile(avatarChairTileIndex indexes.SpriteIndex, pos *references.Position) indexes.SpriteIndex {
	switch avatarChairTileIndex {
	case indexes.ChairFacingDown:
		downPos := pos.GetPositionDown()
		downPosTile := g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor).GetTopTile(downPos)
		if downPosTile.Index == indexes.TableFoodBoth || downPosTile.Index == indexes.TableFoodTop {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingDown, 0)
		}
		return indexes.AvatarSittingFacingDown
	case indexes.ChairFacingUp:
		upPos := pos.GetPositionUp()
		upPosTile := g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor).GetTopTile(upPos)
		if upPosTile.Index == indexes.TableFoodBoth || upPosTile.Index == indexes.TableFoodBottom {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingUp, 0)
		}
		return indexes.AvatarSittingFacingUp
	}

	return avatarChairTileIndex
}

// refreshSpecialTileOverrideExceptions
// Refreshes the special tiles that are not in the map like Portcullis, drawbridge and mirrors
func (g *GameScene) refreshSpecialTileOverrideExceptions(pos *references.Position, layer *game_state.LayeredMap) {
	tile := layer.GetTileTopMapOnlyTile(pos)
	if tile == nil {
		return
	}
	switch tile.Index {
	case indexes.Portcullis, indexes.BrickWallArchway:
		layer.SetTileByLayer(game_state.MapOverrideLayer, pos, g.gameState.GetArchwayPortcullisSpriteByTime())
	case indexes.WoodenPlankVert1Floor, indexes.WoodenPlankVert2Floor:
		g.setDrawBridge(layer, pos, tile.Index)
	case indexes.Mirror, indexes.MirrorAvatar:
		if g.gameState.IsAvatarAtPosition(pos.GetPositionDown()) {
			layer.SetTileByLayer(game_state.MapOverrideLayer, pos, indexes.MirrorAvatar)
		} else {
			layer.SetTileByLayer(game_state.MapOverrideLayer, pos, indexes.Mirror)
		}
	}
}

func (g *GameScene) refreshProvisionsAndEquipmentMapTiles(pos *references.Position, layer *game_state.LayeredMap) {
	if !layer.IsPositionVisible(pos, g.gameState.DateTime, &g.gameState.Lighting, g.gameState.Floor < 0) {
		layer.UnSetTileByLayer(game_state.EquipmentAndProvisionsLayer, pos)
		return
	}

	if !g.gameState.ItemStacksMap.HasItemStackAtPosition(pos) {
		layer.UnSetTileByLayer(game_state.EquipmentAndProvisionsLayer, pos)
		return
	}

	item := g.gameState.ItemStacksMap.Peek(pos)
	if item == nil {
		log.Fatal("Unexpected: item should exist since we checked ahead of it")
	}
	tileIndex := g.gameReferences.InventoryItemReferences.GetReferenceByItem(item.Item).ItemSprite

	layer.SetTileByLayer(game_state.EquipmentAndProvisionsLayer, pos, tileIndex)
}

func (g *GameScene) getTileVisibilityIndexByPosition(pos *references.Position) int {
	//if g.gameState.DateTime.

	// if g.gameState.HasTorchLit() {

	// }
	//g.gameState.XTilesInMap
	//g.gameState.
	return 1
}

func (g *GameScene) refreshMapUnitMapTiles(pos *references.Position, layer *game_state.LayeredMap, do *ebiten.DrawImageOptions) {
	mapUnitTile := layer.GetTileByLayer(game_state.MapUnitLayer, pos)
	underTile := layer.GetTileTopMapOnlyTile(pos)
	if mapUnitTile == nil || mapUnitTile.Index == 0 {
		mapUnitTile = layer.GetTileByLayer(game_state.EquipmentAndProvisionsLayer, pos)
		if mapUnitTile != nil && mapUnitTile.Index >= 512 {
			log.Fatalf("Unepexted tile index for map unit = %d", mapUnitTile.Index)
		}
		if mapUnitTile == nil || mapUnitTile.Index == indexes.NoSprites {
			return
		}
	}
	var tileIndex indexes.SpriteIndex
	if layer.IsPositionVisible(pos, g.gameState.DateTime, &g.gameState.Lighting, g.gameState.Floor < 0) && g.getTileVisibilityIndexByPosition(pos) > 0 {
		tileIndex = g.getSmallCalculatedNPCTileIndex(underTile.Index, mapUnitTile.Index, *pos)
		tileIndex = g.getSmallCalculatedTileIndex(tileIndex, pos)
		if mapUnitTile != nil && mapUnitTile.Index >= 512 {
			log.Fatalf("Unexpected map unit index = %d", mapUnitTile.Index)
		}
		if tileIndex != indexes.NoSprites {
			g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(tileIndex), do)
		}
	}

	// o := text.NewOutput(uf, 20, 1, 10)

	// o.DrawText(g.unscaledMapImage, fmt.Sprintf("x=%d y=%d", pos.X, pos.Y), &do)
}

func (g *GameScene) refreshStaticMapTiles(pos *references.Position, mapLayer *game_state.LayeredMap, do *ebiten.DrawImageOptions) {
	if !mapLayer.IsPositionVisible(pos, g.gameState.DateTime, &g.gameState.Lighting, g.gameState.Floor < 0) {
		return
	}

	tile := mapLayer.GetTileTopMapOnlyTile(pos)
	var spriteIndex indexes.SpriteIndex
	if tile == nil {
		if g.gameState.IsOutOfBounds(*pos) {
			spriteIndex = g.gameState.GetCurrentSmallLocationReference().GetOuterTile()
		} else {
			log.Fatal("bad index")
		}
	} else {
		spriteIndex = tile.Index
	}

	// get from the reference
	spriteIndex = g.getSmallCalculatedTileIndex(spriteIndex, pos)

	g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(spriteIndex), do)
}

// refreshAllMapLayerTiles
func (g *GameScene) refreshAllMapLayerTiles() {
	layer := g.gameState.GetLayeredMapByCurrentLocation()
	layer.RecalculateVisibleTiles(g.gameState.Position, &g.gameState.Lighting, g.gameState.DateTime, g.gameState.Floor < 0)

	if g.unscaledMapImage == nil {
		g.unscaledMapImage = ebiten.NewImage(sprites.TileSize*xTilesVisibleOnGameScreen, sprites.TileSize*yTilesVisibleOnGameScreen)
	}

	g.unscaledMapImage.Fill(image.Black)
	mapType := g.gameState.Location.GetMapType()

	do := ebiten.DrawImageOptions{}

	var avatarPos references.Position
	var avatarDo ebiten.DrawImageOptions

	var x, y references.Coordinate
	var pos = &references.Position{}

	for x = 0; x < xTilesVisibleOnGameScreen; x++ {
		for y = 0; y < yTilesVisibleOnGameScreen; y++ {
			pos.X = x + g.gameState.Position.X - xCenter
			pos.Y = y + g.gameState.Position.Y - yCenter
			if mapType == references.LargeMapType {
				pos = pos.GetWrapped(references.XLargeMapTiles, references.YLargeMapTiles)
			}

			do.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))
			g.refreshSpecialTileOverrideExceptions(pos, layer)
			g.refreshProvisionsAndEquipmentMapTiles(pos, layer)
			g.refreshStaticMapTiles(pos, layer, &do)
			if g.gameState.Position.Equals(pos) {
				avatarPos = *pos

				avatarDo = ebiten.DrawImageOptions{}
				avatarDo.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))
			}
			g.refreshMapUnitMapTiles(pos, layer, &do)
			do.GeoM.Reset()
		}
	}

	avatarSpriteIndex := layer.GetTileTopMapOnlyTile(&avatarPos).Index
	avatarSpriteIndex = g.getSmallCalculatedAvatarTileIndex(avatarSpriteIndex)

	g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(avatarSpriteIndex), &avatarDo)
}
