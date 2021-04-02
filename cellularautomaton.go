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
	size          int
	numVal        int
	initialGrid   []int
	lastGrid      []int
	grid          []int
	nextGrid      []int
	score         [][]int
	previousRules []int
	keepOldRules  bool
	rules         []int
	generation    int
}

func initCellularAutomaton() CelAut {
	cA := CelAut{
		size:        globalDefaultSize,
		numVal:      globalDefaultNumVal,
		initialGrid: make([]int, globalDefaultSize, globalMaxSize),
		lastGrid:    make([]int, globalDefaultSize, globalMaxSize),
		grid:        make([]int, globalDefaultSize, globalMaxSize),
		nextGrid:    make([]int, globalDefaultSize, globalMaxSize),
		score:       make([][]int, globalDisplayLine-1),
		previousRules: make([]int,
			globalDefaultNumVal*globalDefaultNumVal*globalDefaultNumVal, globalMaxNumVal*globalMaxNumVal*globalMaxNumVal),
		rules: make([]int,
			globalDefaultNumVal*globalDefaultNumVal*globalDefaultNumVal,
			globalMaxNumVal*globalMaxNumVal*globalMaxNumVal),
	}
	for i := range cA.score {
		cA.score[i] = make([]int, globalDefaultSize, globalMaxSize)
	}
	return cA
}

func (cA *CelAut) genGrid(fresh bool) {
	if fresh || cA.size != len(cA.initialGrid) {
		cA.initialGrid = cA.initialGrid[:cA.size]
		cA.grid = cA.grid[:cA.size]
		for i := range cA.score {
			cA.score[i] = cA.score[i][:cA.size]
		}
	}
	copy(cA.grid, cA.initialGrid)
	cA.lastGrid = cA.lastGrid[:cA.size]
	cA.nextGrid = cA.nextGrid[:cA.size]
	for i := range cA.lastGrid {
		cA.lastGrid[i] = 0
		cA.nextGrid[i] = 0
	}
	for i := range cA.score {
		for j := range cA.score[i] {
			cA.score[i][j] = 0
		}
	}
}

func (cA *CelAut) genBasicRules(fresh bool) {
	numRules := cA.numVal * cA.numVal * cA.numVal
	if fresh {
		cA.rules = cA.rules[:numRules]
	} else if numRules != len(cA.rules) {
		if !cA.keepOldRules {
			cA.previousRules = cA.previousRules[:len(cA.rules)]
			copy(cA.previousRules, cA.rules)
			cA.keepOldRules = true
		}
		oldNumVal := 1
		for oldNumVal*oldNumVal*oldNumVal != len(cA.previousRules) {
			oldNumVal++
		}
		cA.rules = cA.rules[:numRules]
		for i := 0; i < len(cA.rules); i++ {
			left := i / (cA.numVal * cA.numVal)
			mid := (i / cA.numVal) % cA.numVal
			right := i % cA.numVal
			if left < oldNumVal && mid < oldNumVal && right < oldNumVal {
				oldPos := left*oldNumVal*oldNumVal + mid*oldNumVal + right
				cA.rules[i] = cA.previousRules[oldPos]
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
	cA.getScore()
}

func (cA *CelAut) update() {
	cA.generation++
	copy(cA.lastGrid, cA.grid)
	copy(cA.grid, cA.nextGrid)
	cA.getNextGrid()
	cA.updateScore()
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

func (cA *CelAut) getScore() {
	copy(cA.score[0], cA.initialGrid)
	for i := 1; i < len(cA.score); i++ {
		for j := 0; j < len(cA.score[i]); j++ {
			left := (j - 1 + len(cA.score[i])) % len(cA.score[i-1])
			mid := j
			right := (j + 1) % len(cA.score[i-1])
			ruleNum := cA.score[i-1][left]*cA.numVal*cA.numVal + cA.score[i-1][mid]*cA.numVal + cA.score[i-1][right]
			cA.score[i][j] = cA.rules[ruleNum]
		}
	}
}

func (cA *CelAut) updateScore() {
	for i := 0; i < len(cA.score)-1; i++ {
		for j := 0; j < len(cA.score[0]); j++ {
			cA.score[i][j] = cA.score[i+1][j]
		}
	}
	i := len(cA.score) - 1
	for j := 0; j < len(cA.score[i]); j++ {
		left := (j - 1 + len(cA.score[i])) % len(cA.score[i-1])
		mid := j
		right := (j + 1) % len(cA.score[i-1])
		ruleNum := cA.score[i-1][left]*cA.numVal*cA.numVal + cA.score[i-1][mid]*cA.numVal + cA.score[i-1][right]
		cA.score[i][j] = cA.rules[ruleNum]
	}
}

func (cA *CelAut) draw(screen *ebiten.Image, x, y int, drawCursor bool, drawFuturAndPast bool) {
	radius := float64(7 * len(cA.grid))
	xCenter := float64(x)
	yCenter := float64(y)
	for i := 0; i < len(cA.grid); i++ {
		cellX := xCenter + radius*math.Cos(2*math.Pi*float64(i)/float64(len(cA.grid)))
		cellY := yCenter + radius*math.Sin(2*math.Pi*float64(i)/float64(len(cA.grid)))
		cA.drawCell(i, screen, cellX, cellY, drawFuturAndPast, (i == currentCell) && drawCursor)
	}
}

func (cA *CelAut) drawPart(screen *ebiten.Image, x, y int, drawCursor bool, drawFuturAndPast bool) {

	lineSize := 16
	numLines := globalDisplayLine

	if drawFuturAndPast {
		if len(cA.lastGrid) > 0 {
			cA.drawLine(screen, x, y, false, cA.lastGrid, false)
		}
	}

	cA.drawLine(screen, x, y+lineSize, drawCursor, cA.grid, true)

	if drawFuturAndPast {
		for i := 1; i < numLines-1; i++ {
			cA.drawLine(screen, x, y+(i+1)*lineSize, false, cA.score[i], false)
		}
	}

}

func (cA *CelAut) drawLine(screen *ebiten.Image, x, y int, drawCursor bool, line []int, current bool) {

	bigSize := 12
	smallSize := 8
	colSize := 16
	cursorSize := 14

	for i := range line {
		if drawCursor && i == currentCell {
			ebitenutil.DrawRect(screen, float64(x-cursorSize/2+i*colSize), float64(y-cursorSize/2), float64(cursorSize), float64(cursorSize), color.White)
		}

		colorPos := line[i]
		if colorPos >= cA.numVal {
			colorPos = 0
		}
		cellColor := stateColors[colorPos]
		cellSize := smallSize
		if current {
			cellSize = bigSize
		}
		ebitenutil.DrawRect(screen, float64(x-cellSize/2+i*colSize), float64(y-cellSize/2), float64(cellSize), float64(cellSize), cellColor)
	}

}

func (cA *CelAut) drawCell(pos int, screen *ebiten.Image, x, y float64, drawOther, drawCursor bool) {
	smallSize := 5.0
	bigSize := 20.0
	cursorSize := 22.0
	if drawCursor {
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
