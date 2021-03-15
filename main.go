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
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameDisplay struct {
	state     int
	automaton CelAut
	tempoPos  int
	frame     int
	fresh     bool
	part      bool
	audio     soundManager
}

const (
	stateInit int = iota
	stateChooseTempo
	stateChooseSize
	stateChooseNumVal
	stateChooseRules
	stateChooseInitial
	stateRunAutomaton
)

func (gD *GameDisplay) initUpdate() bool {
	gD.frame++
	if gD.frame >= 6 {
		gD.initSound()
		gD.frame = 0
		return true
	}
	return false
}

func (gD *GameDisplay) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		gD.part = !gD.part
	}
	switch gD.state {
	case stateInit:
		if gD.initUpdate() {
			gD.state++
		}
	case stateChooseTempo:
		if gD.chooseTempoUpdate() {
			gD.state++
		}
	case stateChooseSize:
		if gD.chooseSizeUpdate() {
			gD.state++
		}
		gD.automaton.genGrid(gD.fresh)
	case stateChooseNumVal:
		if gD.chooseNumValUpdate() {
			currentRule = 0
			gD.state++
		}
		gD.automaton.genBasicRules(gD.fresh)
	case stateChooseRules:
		if gD.chooseRulesUpdate() {
			currentCell = 0
			gD.state++
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			gD.automaton.init()
			gD.playSounds()
			gD.frame = 0
			gD.state += 2
		}
		gD.automaton.init()
	case stateChooseInitial:
		if gD.chooseInitialGridUpdate() {
			gD.playSounds()
			gD.frame = 0
			gD.state++
		} else if inpututil.IsKeyJustPressed(ebiten.KeyShift) {
			gD.state--
		}
		gD.automaton.init()
	case stateRunAutomaton:
		gD.frame++
		if 3600/tempos[gD.tempoPos] <= gD.frame {
			gD.automaton.update()
			gD.playSounds()
			gD.frame = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			gD.fresh = false
			gD.automaton.genGrid(gD.fresh)
			gD.state = stateChooseTempo
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if !gD.audio.use {
				gD.audio.use = true
			} else {
				gD.audio.soundset = (gD.audio.soundset + 1) % numSoundSet
				if gD.audio.soundset == 0 {
					gD.audio.use = false
				}
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			return errors.New("Finished")
		}
	}
	return nil
}

func (gD *GameDisplay) Draw(screen *ebiten.Image) {
	if gD.state >= stateChooseTempo {
		ebitenutil.DebugPrintAt(screen, fmt.Sprint("Tempo : ", tempos[gD.tempoPos]), 10, 10)
		if gD.state == stateChooseTempo {
			ebitenutil.DebugPrintAt(screen, "Réglage du tempo", 10, 490)
			ebitenutil.DebugPrintAt(screen, "   Flèches : faire varier le tempo", 10, 505)
			ebitenutil.DebugPrintAt(screen, "   Entrée : valider le tempo", 10, 520)
		}
	}

	if gD.state >= stateChooseSize || !gD.fresh {
		ebitenutil.DebugPrintAt(screen, fmt.Sprint("Nombre de cellules : ", gD.automaton.size), 10, 25)
		if gD.state == stateChooseSize {
			ebitenutil.DebugPrintAt(screen, "Réglage du nombre de cellules", 10, 490)
			ebitenutil.DebugPrintAt(screen, "   Flèches : faire varier le nombre de cellules", 10, 505)
			ebitenutil.DebugPrintAt(screen, "   Entrée : valider le nombre de cellules", 10, 520)
		}
	}

	if gD.state >= stateChooseNumVal || !gD.fresh {
		ebitenutil.DebugPrintAt(screen, fmt.Sprint("Nombre d'états par cellule : ", gD.automaton.numVal), 10, 40)
		if gD.state == stateChooseNumVal {
			ebitenutil.DebugPrintAt(screen, "Réglage du nombre d'états possibles pour chaque cellule", 10, 490)
			ebitenutil.DebugPrintAt(screen, "   Flèches : faire varier le nombre d'états", 10, 505)
			ebitenutil.DebugPrintAt(screen, "   Entrée : valider le nombre d'états", 10, 520)
		}
	}

	if gD.state >= stateChooseNumVal || !gD.fresh {
		ebitenutil.DebugPrintAt(screen, "Règles : ", 10, 55)
		if gD.state == stateChooseRules {
			gD.automaton.drawRules(screen, 20, 75, true)
			ebitenutil.DebugPrintAt(screen, "Choix des règles", 10, 490)
			ebitenutil.DebugPrintAt(screen, "   Flèches : sélectionner une règle", 10, 505)
			ebitenutil.DebugPrintAt(screen, "   Espace : changer la règle sélectionnée", 10, 520)
			ebitenutil.DebugPrintAt(screen, "   Majuscule : passer au choix de l'état initial", 10, 535)
			ebitenutil.DebugPrintAt(screen, "   Entrée : lancer la simulation", 10, 550)
		} else {
			gD.automaton.drawRules(screen, 20, 75, false)
		}
	}

	if gD.state != stateRunAutomaton && (gD.state >= stateChooseSize || !gD.fresh) {
		if gD.part {
			gD.automaton.drawPart(screen, 350+(globalMaxSize-len(gD.automaton.grid))*8, 20, gD.state == stateChooseInitial, !gD.fresh || gD.state >= stateChooseRules)
		} else {
			gD.automaton.draw(screen, 700, 300, gD.state == stateChooseInitial, !gD.fresh || gD.state >= stateChooseRules)
		}
		if gD.state == stateChooseInitial {
			ebitenutil.DebugPrintAt(screen, "Choix de l'état initial des cellules", 10, 490)
			ebitenutil.DebugPrintAt(screen, "   Flèches (gauche, droite) : sélectionner une cellule", 10, 505)
			ebitenutil.DebugPrintAt(screen, "   Espace : changer l'état de la cellule sélectionnée", 10, 520)
			ebitenutil.DebugPrintAt(screen, "   Majuscule : passer au choix des règles", 10, 535)
			ebitenutil.DebugPrintAt(screen, "   Entrée : lancer la simulation", 10, 550)
		}
	}

	if gD.state == stateRunAutomaton {
		if gD.part {
			gD.automaton.drawPart(screen, 350+(globalMaxSize-len(gD.automaton.grid))*8, 20, false, true)
		} else {
			gD.automaton.draw(screen, 700, 300, false, true)
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprint("Simulation en cours (génération ", gD.automaton.generation, ")"), 10, 490)
		ebitenutil.DebugPrintAt(screen, "   Entrée : recommencer avec de nouveaux paramètres", 10, 505)
		ebitenutil.DebugPrintAt(screen, "   Espace : changer le jeu de sons", 10, 520)
	}

	if gD.part {
		ebitenutil.DebugPrintAt(screen, "   Tabulation : passer en mode visualisation", 10, 565)
	} else {
		ebitenutil.DebugPrintAt(screen, "   Tabulation : passer en mode partition", 10, 565)
	}
	if ebiten.IsFullscreen() {
		ebitenutil.DebugPrintAt(screen, "   Echape : quitter le mode plein écran", 10, 580)
	} else {
		ebitenutil.DebugPrintAt(screen, "   Echape : quitter le logiciel", 10, 580)
	}
}

func (gD *GameDisplay) Layout(width, height int) (int, int) {
	return 1000, 600
}

func main() {

	gD := GameDisplay{
		state:     0,
		automaton: initCellularAutomaton(),
		tempoPos:  25,
		fresh:     true,
		audio:     initAudio(),
	}

	ebiten.SetWindowSize(1000, 600)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("GRAC: Génération de Rythmes à l'aide d'Automates Cellulaires")

	ebiten.RunGame(&gD)

}
