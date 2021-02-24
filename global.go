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
	"sort"
)

const (
	globalMinSize   = 3
	globalMaxSize   = 40
	globalMinNumVal = 2
	globalMaxNumVal = 5
	numSoundSet     = 2
)

var stateColors []color.Color = []color.Color{
	color.RGBA{192, 192, 192, 255},
	color.RGBA{255, 153, 51, 255},
	color.RGBA{153, 51, 255, 255},
	color.RGBA{153, 255, 51, 255},
	color.RGBA{51, 153, 255, 255},
}

var sounds [numSoundSet][globalMaxNumVal - 1][]byte

var tempos []int = genTempos()

func genTempos() []int {
	// 3600 frames per minute, these are the prime dividers of 3600
	primeDiv := []int{2, 2, 2, 2, 3, 3, 5, 5}
	res := genDiv(primeDiv)
	sort.Ints(res)
	return res
}

func genDiv(primeDiv []int) []int {
	if len(primeDiv) == 0 {
		return []int{1}
	}
	i := 1
	for cur := primeDiv[0]; i < len(primeDiv) && primeDiv[i] == cur; i++ {
	}
	divs := genDiv(primeDiv[i:])
	numDivs := len(divs)
	for j := 0; j < i; j++ {
		for k := 0; k < numDivs; k++ {
			divs = append(divs, divs[k+j*numDivs]*primeDiv[0])
		}
	}
	return divs
}
