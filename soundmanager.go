/*
GRAC, rhythm generation using cellular automata
Copyright (C) 2021 Loïg Jezequel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/loig/grac/assets"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type soundManager struct {
	context  *audio.Context
	soundset int
	use      bool
	players  []*audio.Player
	test     *audio.Player
}

func (gD *GameDisplay) playSounds() {
	if gD.audio.use {
		for _, player := range gD.audio.players {
			if player != nil {
				err := player.Close()
				if err != nil {
					log.Panic(err)
				}
			}
		}
		gD.audio.players = make([]*audio.Player, len(gD.automaton.grid))
		for i := 0; i < len(gD.automaton.grid); i++ {
			if gD.automaton.grid[i] != 0 {
				gD.playSound(i, gD.automaton.grid[i])
			}
		}
	}
}

func (gD *GameDisplay) playSound(playerpos, soundpos int) {
	soundBytes := sounds[gD.audio.soundset][soundpos-1]
	gD.audio.players[playerpos] = audio.NewPlayerFromBytes(gD.audio.context, soundBytes)
	gD.audio.players[playerpos].Play()
}

func (gD *GameDisplay) initSound() {
	soundBytes := sounds[0][2]
	player := audio.NewPlayerFromBytes(gD.audio.context, soundBytes)
	player.Play()
}

func initAudio() soundManager {
	context := audio.NewContext(44100)

	theSounds := [][][]byte{
		[][]byte{assets.Natural0, assets.Natural1, assets.Natural2, assets.Natural3},
		[][]byte{assets.Sound0, assets.Sound1, assets.Sound2, assets.Sound3},
	}

	for k := 0; k < numSoundSet; k++ {
		for i := 0; i < len(sounds[0]); i++ {
			sound, error := mp3.Decode(context, bytes.NewReader(theSounds[k][i]))
			if error != nil {
				log.Panic(error)
			}
			sounds[k][i], error = ioutil.ReadAll(sound)
		}
	}

	return soundManager{
		context:  context,
		soundset: 0,
		use:      false,
	}
}
