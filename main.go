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
			gD.automaton.genGrid(gD.fresh)
			gD.state++
		}
	case stateChooseNumVal:
		if gD.chooseNumValUpdate() {
			gD.automaton.genBasicRules(gD.fresh)
			currentRule = 0
			gD.state++
		}
	case stateChooseRules:
		if gD.chooseRulesUpdate() {
			currentCell = 0
			gD.state++
		}
	case stateChooseInitial:
		if gD.chooseInitialGridUpdate() {
			gD.automaton.init()
			gD.playSounds()
			gD.frame = 0
			gD.state++
		}
	case stateRunAutomaton:
		gD.frame++
		if 3600/tempos[gD.tempoPos] <= gD.frame {
			gD.automaton.update()
			gD.playSounds()
			gD.frame = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			gD.fresh = false
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
			ebitenutil.DebugPrintAt(screen, "Réglage du tempo, utiliser les flèches pour le faire varier et entrée pour valider.", 10, 575)
		}
	}

	if gD.state >= stateChooseSize || !gD.fresh {
		ebitenutil.DebugPrintAt(screen, fmt.Sprint("Nombre de cellules : ", gD.automaton.size), 10, 30)
		if gD.state == stateChooseSize {
			ebitenutil.DebugPrintAt(screen, "Réglage du nombre de cellules, utiliser les flèches pour le faire varier et entrée pour valider.", 10, 575)
		}
	}

	if gD.state >= stateChooseNumVal || !gD.fresh {
		ebitenutil.DebugPrintAt(screen, fmt.Sprint("Nombre d'états par cellule : ", gD.automaton.numVal), 10, 50)
		if gD.state == stateChooseNumVal {
			ebitenutil.DebugPrintAt(screen, "Réglage du nombre d'états possibles pour chaque cellule, utiliser les flèches pour le faire varier et entrée pour valider.", 10, 575)
		}
	}

	if gD.state >= stateChooseRules || !gD.fresh {
		ebitenutil.DebugPrintAt(screen, "Règles : ", 10, 70)
		if gD.state == stateChooseRules {
			gD.automaton.drawRules(screen, 20, 100, true)
			ebitenutil.DebugPrintAt(screen, "Choix des règles, utiliser les flèches pour sélectionner une règle, espace pour changer son résultat.", 10, 565)
			ebitenutil.DebugPrintAt(screen, "Quand toutes les règles sont fixées, appuyer sur entrée pour valider.", 10, 575)
		} else {
			gD.automaton.drawRules(screen, 20, 100, false)
		}
	}

	if gD.state == stateChooseInitial {
		gD.automaton.draw(screen, 700, 300, true)
		if gD.state == stateChooseInitial {
			ebitenutil.DebugPrintAt(screen, "Choix de l'état initial des cellules, utiliser les flèches (gauche, droite) pour sélectionner une cellule, espace pour changer son état.", 10, 565)
			ebitenutil.DebugPrintAt(screen, "Quand l'état initial de toutes les cellules est fixé, appuyer sur entrée pour valider.", 10, 575)
		}
	}

	if gD.state == stateRunAutomaton {
		gD.automaton.draw(screen, 700, 300, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprint("Génération actuelle : ", gD.automaton.generation), 10, (len(gD.automaton.rules)+7)/8*26+110)
		ebitenutil.DebugPrintAt(screen, "La simulation tourne. Utiliser entrée pour recommencer avec de nouveaux paramètres.", 10, 565)
		if ebiten.IsFullscreen() {
			ebitenutil.DebugPrintAt(screen, "Utiliser espace pour changer le jeu de sons. Utiliser échape pour quitter le mode plein écran.", 10, 575)
		} else {
			ebitenutil.DebugPrintAt(screen, "Utiliser espace pour changer le jeu de sons. Utiliser échape pour quitter le logiciel.", 10, 575)
		}
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
