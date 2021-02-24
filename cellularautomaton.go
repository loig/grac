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
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type CelAut struct {
	size        int
	numVal      int
	initialGrid []int
	lastGrid    []int
	grid        []int
	nextGrid    []int
	rules       []int
	generation  int
}

func initCellularAutomaton() CelAut {
	return CelAut{
		size:   25,
		numVal: 2,
	}
}

func (cA *CelAut) genGrid(fresh bool) {
	if fresh || cA.size != len(cA.initialGrid) {
		cA.initialGrid = make([]int, cA.size)
		cA.grid = make([]int, cA.size)
	} else {
		copy(cA.grid, cA.initialGrid)
	}
	cA.lastGrid = make([]int, cA.size)
	cA.nextGrid = make([]int, cA.size)
}

func (cA *CelAut) genBasicRules(fresh bool) {
	if fresh {
		cA.rules = make([]int, cA.numVal*cA.numVal*cA.numVal)
	} else if cA.numVal*cA.numVal*cA.numVal != len(cA.rules) {
		oldRules := cA.rules // copy?
		oldNumVal := 1
		for oldNumVal*oldNumVal*oldNumVal != len(oldRules) {
			oldNumVal++
		}
		cA.rules = make([]int, cA.numVal*cA.numVal*cA.numVal)
		for i := 0; i < len(cA.rules); i++ {
			left := i / (cA.numVal * cA.numVal)
			mid := (i / cA.numVal) % cA.numVal
			right := i % cA.numVal
			if left < oldNumVal && mid < oldNumVal && right < oldNumVal {
				oldPos := left*oldNumVal*oldNumVal + mid*oldNumVal + right
				cA.rules[i] = oldRules[oldPos]
				if cA.rules[i] >= cA.numVal {
					cA.rules[i] = 0
				}
			}
		}
	}
}

func (cA *CelAut) init() {
	cA.generation = 0
	for i := 0; i < len(cA.grid); i++ {
		if cA.initialGrid[i] >= cA.numVal {
			cA.initialGrid[i] = 0
		}
		cA.grid[i] = cA.initialGrid[i]
	}
	cA.getNextGrid()
}

func (cA *CelAut) update() {
	cA.generation++
	for i := 0; i < len(cA.lastGrid); i++ {
		cA.lastGrid[i] = cA.grid[i]
		cA.grid[i] = cA.nextGrid[i]
	}
	cA.getNextGrid()
}

func (cA *CelAut) getNextGrid() {
	for i := 0; i < len(cA.nextGrid); i++ {
		left := (i - 1 + len(cA.nextGrid)) % len(cA.grid)
		mid := i
		right := (i + 1) % len(cA.grid)
		ruleNum := cA.grid[left]*cA.numVal*cA.numVal + cA.grid[mid]*cA.numVal + cA.grid[right]
		cA.nextGrid[i] = cA.rules[ruleNum]
	}
}

func (cA *CelAut) draw(screen *ebiten.Image, x, y int, drawCursor bool) {
	radius := float64(7 * len(cA.grid))
	xCenter := float64(x)
	yCenter := float64(y)
	for i := 0; i < len(cA.grid); i++ {
		cellX := xCenter + radius*math.Cos(2*math.Pi*float64(i)/float64(len(cA.grid)))
		cellY := yCenter + radius*math.Sin(2*math.Pi*float64(i)/float64(len(cA.grid)))
		cA.drawCell(i, screen, cellX, cellY, !drawCursor, i == currentCell)
	}
}

func (cA *CelAut) drawCell(pos int, screen *ebiten.Image, x, y float64, drawOther, drawCursor bool) {
	smallSize := 5.0
	bigSize := 20.0
	cursorSize := 22.0
	if drawCursor && !drawOther {
		ebitenutil.DrawRect(screen, x-cursorSize/2, y-cursorSize/2, cursorSize, cursorSize, color.White)
	}
	colorPos := cA.grid[pos]
	if colorPos >= cA.numVal {
		colorPos = 0
	}
	cellColor := stateColors[colorPos]
	ebitenutil.DrawRect(screen, x-bigSize/2, y-bigSize/2, bigSize, bigSize, cellColor)
	if drawOther {
		lastColor := stateColors[cA.lastGrid[pos]]
		nextColor := stateColors[cA.nextGrid[pos]]
		ebitenutil.DrawRect(screen, x-bigSize/2-2, y-bigSize/2-smallSize-1, smallSize, smallSize, lastColor)
		ebitenutil.DrawRect(screen, x+bigSize/2-smallSize+2, y+bigSize/2+1, smallSize, smallSize, nextColor)
	}
}

func (cA *CelAut) drawRules(screen *ebiten.Image, x, y int, drawCursor bool) {
	xOffset := 37
	yOffset := 26
	for i := 0; i < len(cA.rules); i++ {
		cA.drawRule(i, screen, float64(x+(i%8)*xOffset), float64(y+(i/8)*yOffset), i == currentRule && drawCursor)
	}
}

func (cA *CelAut) drawRule(ruleNum int, screen *ebiten.Image, x, y float64, drawCursor bool) {
	size := 10.0
	if drawCursor {
		ebitenutil.DrawRect(screen, x-2, y-2, 3*size+6, 2*size+5, color.White)
		ebitenutil.DrawRect(screen, x-1, y-1, 3*size+4, 2*size+3, color.Black)
	}
	actualNumVal := 1
	for actualNumVal*actualNumVal*actualNumVal != len(cA.rules) {
		actualNumVal++
	}
	leftState := ruleNum / (actualNumVal * actualNumVal)
	midState := (ruleNum / (actualNumVal)) % actualNumVal
	rightState := ruleNum % actualNumVal
	ebitenutil.DrawRect(screen, x, y, size, size, stateColors[leftState])
	ebitenutil.DrawRect(screen, x+size+1, y, size, size, stateColors[midState])
	ebitenutil.DrawRect(screen, x+2*size+2, y, size, size, stateColors[rightState])
	state := cA.rules[ruleNum]
	ebitenutil.DrawRect(screen, x+size+1, y+size+1, size, size, stateColors[state])
}
