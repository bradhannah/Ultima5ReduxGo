package references

import (
	"runtime"
	"sync"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/stretchr/testify/assert"
)

var (
	talkReferences *TalkReferences
	loadOnce       sync.Once
)

func getBaseTalkReferences() *TalkReferences {
	loadOnce.Do(func() {
		cfg := config.NewUltimaVConfiguration()

		dataOvl := NewDataOvl(cfg)
		var err error
		talkReferences = NewTalkReferences(cfg, dataOvl)

		if err != nil {
			return
		}
	})

	return talkReferences
}

func TestNewTalkReferences(t *testing.T) {
	t.Parallel()

	t.Run("ValidConfigAndDataOvl", func(t *testing.T) {
		talkReferences := getBaseTalkReferences()

		assert.NotNil(t, talkReferences)
		assert.NotNil(t, talkReferences.talkDataForSmallMapType)
		assert.NotNil(t, talkReferences.WordDict)
		assert.NotNil(t, talkReferences.talkScripts)

		assert.Len(t, talkReferences.talkScripts[Towne], 48)
	})

	t.Run("CorrectNumberOfScripts", func(t *testing.T) {
		talkReferences := getBaseTalkReferences()

		assert.Len(t, talkReferences.talkScripts[Towne], 48)
		assert.Len(t, talkReferences.talkScripts[Castle], 40)
		assert.Len(t, talkReferences.talkScripts[Keep], 32)
		// maybe is supposed to be 16, can't be sure yet
		if !assert.Len(t, talkReferences.talkScripts[Dwelling], 15) {
			runtime.Breakpoint()
		}
	})

	t.Run("NoFunnyCharactersInScript", func(t *testing.T) {
		talkReferences := getBaseTalkReferences()

		for _, scripts := range talkReferences.talkScripts {
			for _, script := range scripts {
				for nLines, lines := range script.Lines {
					for nLine, line := range lines {
						_ = nLine
						_ = nLines

						if !assert.Regexp(t, "([a-zA-Z0-9 !#'.-?,@]+|^$)", line.Str) {
							runtime.Breakpoint()
						}
					}
				}
			}
		}
	})
}
