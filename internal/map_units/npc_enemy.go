package map_units

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type NPCEnemy struct {
	EnemyReference references2.EnemyReference
	mapUnitDetails MapUnitDetails
}

func NewEnemyNPC(enemyRef references2.EnemyReference, npcNum int) NPCEnemy {
	enemy := NPCEnemy{}
	enemy.EnemyReference = enemyRef
	enemy.mapUnitDetails.NPCNum = npcNum

	// enemy.mapUnitDetails.AStarMap = map_state.NewAStarMap()
	enemy.mapUnitDetails.CurrentPath = nil

	return enemy
}

func (enemy *NPCEnemy) GetMapUnitType() MapUnitType {
	return Enemy
}

func (enemy *NPCEnemy) Pos() references2.Position {
	return enemy.mapUnitDetails.Position
}

func (enemy *NPCEnemy) MapUnitDetails() *MapUnitDetails {
	return &enemy.mapUnitDetails
}

func (enemy *NPCEnemy) SetVisible(visible bool) {
	enemy.mapUnitDetails.Visible = visible
}
func (enemy *NPCEnemy) IsVisible() bool {
	return enemy.mapUnitDetails.Visible
}
func (enemy *NPCEnemy) IsEmptyMapUnit() bool {
	return enemy.EnemyReference.KeyFrameTile == nil
}

func (enemy *NPCEnemy) Floor() references2.FloorNumber {
	return enemy.mapUnitDetails.Floor
}

func (enemy *NPCEnemy) PosPtr() *references2.Position {
	return &enemy.mapUnitDetails.Position
}

func (enemy *NPCEnemy) SetPos(position references2.Position) {
	enemy.mapUnitDetails.Position = position
}

func (enemy *NPCEnemy) SetFloor(floor references2.FloorNumber) {
	enemy.mapUnitDetails.Floor = floor
}
