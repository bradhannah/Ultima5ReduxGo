package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

type NPCEnemy struct {
	EnemyReference references.EnemyReference
	mapUnitDetails MapUnitDetails
}

func NewEnemyNPC(enemyRef references.EnemyReference, npcNum int) NPCEnemy {
	enemy := NPCEnemy{}
	enemy.EnemyReference = enemyRef
	enemy.mapUnitDetails.NPCNum = npcNum

	enemy.mapUnitDetails.AStarMap = NewAStarMap()
	enemy.mapUnitDetails.CurrentPath = nil

	return enemy
}

func (enemy *NPCEnemy) GetMapUnitType() MapUnitType {
	return Enemy
}

func (enemy *NPCEnemy) Pos() references.Position {
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

func (enemy *NPCEnemy) Floor() references.FloorNumber {
	return enemy.mapUnitDetails.Floor
}

func (enemy *NPCEnemy) PosPtr() *references.Position {
	return &enemy.mapUnitDetails.Position
}

func (enemy *NPCEnemy) SetPos(position references.Position) {
	enemy.mapUnitDetails.Position = position
}

func (enemy *NPCEnemy) SetFloor(floor references.FloorNumber) {
	enemy.mapUnitDetails.Floor = floor
}
