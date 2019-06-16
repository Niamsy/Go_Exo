package mines

import (
	"fmt"
	"sync"
)

type MineType int
const (
	Iron MineType = iota
	Titanium
	Aluminium
)

type Minor struct {
	Id int
	ProductType MineType
	Stack int
}

var Requests = make(chan Minor)
var responses = make(chan bool)

func DefineMinors(nbAluminumMinors int, nbIronMinors int, nbTitaniumMinors int) []Minor {
	var totalMinors = nbTitaniumMinors + nbIronMinors + nbAluminumMinors
	var minors = make([]Minor, totalMinors)

	var i = 0

	for ; i < nbAluminumMinors; i++ {
		minors[i] = Minor{Id : i, ProductType : Aluminium, Stack : 0}
		fmt.Printf("Just defined aluminium minor id %d\n", i)
	}

	var maxMinorId = i + nbIronMinors
	for ; i < maxMinorId; i++ {
		minors[i] = Minor{Id : i, ProductType : Iron, Stack : 0}
		fmt.Printf("Just defined iron minor id %d\n", i)
	}

	maxMinorId = i + nbTitaniumMinors
	for ; i < maxMinorId; i++ {
		minors[i] = Minor{Id : i, ProductType : Titanium, Stack : 0}
		fmt.Printf("Just defined titanium minor id %d\n", i)
	}
	return minors
}

func Product(minor Minor, wg *sync.WaitGroup){
	minor.Stack += 1 // Create ore each turn
	Requests <- minor // Deliver to mine
	resp := <- responses
	if resp == true {
		wg.Done()
	}
}

func CoordinateMinors(mines map[MineType] int, minors []Minor) {
	for {
		req := <-Requests

		if req.Stack <= 0 {
			continue
		}
		minors[req.Id].Stack = req.Stack

		var toDeliver int

		if mines[req.ProductType] + minors[req.Id].Stack > 30 {
			toDeliver = 30 - mines[req.ProductType]
		} else {
			toDeliver = minors[req.Id].Stack
		}

		mines[req.ProductType] += toDeliver
		minors[req.Id].Stack -= toDeliver
		responses <- true
	}
}

func Mines() map[MineType]int {
	mines := make(map[MineType]int)
	mines[Iron] = 0
	mines[Titanium] = 0
	mines[Aluminium] = 0

	return mines
}
