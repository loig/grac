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
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (gD *GameDisplay) chooseSizeUpdate() bool {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyRight):
		if gD.automaton.size < globalMaxSize {
			gD.automaton.size++
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		if gD.automaton.size > globalMinSize {
			gD.automaton.size--
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		return true
	}
	return false
}

func (gD *GameDisplay) chooseTempoUpdate() bool {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyRight):
		if gD.tempoPos < len(tempos)-1 {
			gD.tempoPos++
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		if gD.tempoPos > 0 {
			gD.tempoPos--
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		return true
	}
	return false
}

func (gD *GameDisplay) chooseNumValUpdate() bool {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyRight):
		if gD.automaton.numVal < globalMaxNumVal {
			gD.automaton.numVal++
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		if gD.automaton.numVal > globalMinNumVal {
			gD.automaton.numVal--
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		return true
	}
	return false
}

var currentRule int

func (gD *GameDisplay) chooseRulesUpdate() bool {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		if (currentRule+7)%8 < currentRule%8 {
			currentRule--
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		if (currentRule+1)%8 > currentRule%8 && currentRule+1 < len(gD.automaton.rules) {
			currentRule++
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		if currentRule-8 >= 0 {
			currentRule -= 8
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		if currentRule+8 < len(gD.automaton.rules) {
			currentRule += 8
		}
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		gD.automaton.rules[currentRule] = (gD.automaton.rules[currentRule] + 1) % gD.automaton.numVal
	case inpututil.IsKeyJustPressed(ebiten.KeyShift):
		return true
	}
	return false
}

var currentCell int

func (gD *GameDisplay) chooseInitialGridUpdate() bool {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		currentCell = (currentCell + len(gD.automaton.initialGrid) - 1) % len(gD.automaton.initialGrid)
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		currentCell = (currentCell + 1) % len(gD.automaton.initialGrid)
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		gD.automaton.initialGrid[currentCell] = (gD.automaton.initialGrid[currentCell] + 1) % gD.automaton.numVal
		gD.automaton.grid[currentCell] = gD.automaton.initialGrid[currentCell]
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		return true
	}
	return false
}
