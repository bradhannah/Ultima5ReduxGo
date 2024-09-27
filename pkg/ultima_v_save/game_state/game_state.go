package game_state

import (
	"fmt"
	util "github.com/bradhannah/Ultima5ReduxGo/pkg/ultima_v_save/util"
	"io"
	"os"
	"reflect"
	"unsafe"
)

func (g *GameState) getSavedGamRaw(savedGamFilePath string) ([]byte, error) {
	// Open the file in read-only mode and as binary
	file, err := os.OpenFile(savedGamFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	buffer := make([]byte, savedGamFileSize)
	n, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	if n != savedGamFileSize {
		return nil, fmt.Errorf("expected file of size 4192 but was %d", n)
	}

	return buffer, nil
}

func (g *GameState) LoadSaveGame(savedGamFilePath string) error {
	// Open the file in read-only mode and as binary
	rawSaveGameBytesFromDisk, err := g.getSavedGamRaw(savedGamFilePath)
	if err != nil {
		return err
	}

	//var saveGame = GameState{}
	g.RawSave = [savedGamFileSize]byte(rawSaveGameBytesFromDisk)

	// Overlay player characters over memory rawSaveGameBytesFromDisk to easily consume data
	characterPtr := (*[NPlayers]PlayerCharacter)(unsafe.Pointer(&g.RawSave[startPositionOfCharacters]))
	g.Characters = *characterPtr

	return nil
}

/*

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(playerCharacter)
	if err != nil {
		return err
	}
	newBytes := buffer.Bytes()
	readBuf := bytes.NewBuffer(newBytes)
	dec := gob.NewDecoder(readBuf)
	var newPlayerCharacter PlayerCharacter
	if err := dec.Decode(&newPlayerCharacter); err != nil {
		fmt.Println("Decoding failed:", err)
		return err
	}

*/

func overlayCharactersOnSave(playerCharacters *[MAX_CHARACTERS_IN_PARTY]PlayerCharacter, saveGame *GameState) {
	playerCharactersT := reflect.TypeOf(playerCharacters)
	byteSizePlayerCharacter := uint(playerCharactersT.Size())

	//savedGameT := reflect.TypeOf(*saveGame)
	//byteSizeSavedGame := savedGameT.Size()

	savedGamePtr := unsafe.Pointer(saveGame)
	savedGameBytes := util.MakeByteSliceFromUnsafePointer(savedGamePtr, savedGamFileSize)
	characterPtr := unsafe.Pointer(&playerCharacters)

	copy(savedGameBytes[startPositionOfCharacters:], util.MakeByteSliceFromUnsafePointer(characterPtr, int(byteSizePlayerCharacter)))
}

func SaveFileOnTopOfSave(savedGamFilePath string, saveGame *GameState) error {
	saveFile, err := os.OpenFile(savedGamFilePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	saveFileBytes, err := io.ReadAll(saveFile)
	if err != nil {
		return nil
	}
	_ = saveFileBytes
	return nil
}

func (g *GameState) SaveCharactersOnSave(savedGamFilePath string, nCharPos uint, playerCharacter PlayerCharacter) error {
	saveFile, err := os.OpenFile(savedGamFilePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	saveFileBytes, err := io.ReadAll(saveFile)
	if err != nil {
		return nil
	}

	var newName [9]byte
	newName[0] = 'B'
	newName[1] = 'R'
	newName[2] = 0

	playerCharacter.Name = newName //[9]byte("OOF")

	t := reflect.TypeOf(playerCharacter)
	byteSizePlayerCharacter := uint(t.Size())
	startPos := startPositionOfCharacters + (byteSizePlayerCharacter * nCharPos)

	//characterPtr := (*[NPlayers]PlayerCharacter)(unsafe.Pointer(&buffer[startPositionOfCharacters]))
	characterPtr := unsafe.Pointer(&playerCharacter)
	copy(saveFileBytes[startPos:], util.MakeByteSliceFromUnsafePointer(characterPtr, int(byteSizePlayerCharacter)))

	//var characterBuffer bytes.Buffer
	//enc := gob.NewEncoder(&characterBuffer)
	//err = enc.Encode(playerCharacter)
	//if err != nil {
	//	return err
	//}
	//
	//t := reflect.TypeOf(playerCharacter)
	//startPos := startPositionOfCharacters + (uint(t.Size()) * nCharPos)
	//newCharacterBytes := characterBuffer.Bytes()
	//copy(saveFileBytes[startPos:], newCharacterBytes)

	return nil
}

//func getClassStr(characterClass ultima_v_save.CharacterClass) string {
//	value, bExists := ultima_v_save.CharacterClassMap[characterClass]
//	if bExists {
//		return value
//	}
//	return ""
//}
//
//func getClassByStr(classStr string) (ultima_v_save.CharacterClass, bool) {
//	characterClass, bExists := ultima_v_save.FindKeyByValueT[ultima_v_save.CharacterClass, string](ultima_v_save.CharacterClassMap, classStr)
//	if bExists {
//		return characterClass, true
//	}
//	return ultima_v_save.Avatar, false
//}
