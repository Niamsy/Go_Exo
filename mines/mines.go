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
	Default
)

type MineRequest struct {
	Deposit bool
	Ores map[MineType]int // Used for deposit
	OresToTake []MineType // Used for taking ores
}

type Minor struct {
	Id int
	ProductType MineType
	Stack int
}

const siloSize = 30

var MineRequests = make(chan MineRequest)
var MineResponses = make(chan map[MineType]int)

var Requests = make(chan Minor)
var responses = make(chan bool)

func DefineMinors(nbAluminumMinors int, nbIronMinors int, nbTitaniumMinors int) []Minor {
	var totalMinors = nbTitaniumMinors + nbIronMinors + nbAluminumMinors
	var minors = make([]Minor, totalMinors)

	var i = 0

	for ; i < nbAluminumMinors; i++ {
		minors[i] = Minor{Id : i + 1, ProductType : Aluminium, Stack : 0}
	}

	var maxMinorId = i + nbIronMinors
	for ; i < maxMinorId; i++ {
		minors[i] = Minor{Id : i + 1, ProductType : Iron, Stack : 0}
	}

	maxMinorId = i + nbTitaniumMinors
	for ; i < maxMinorId; i++ {
		minors[i] = Minor{Id : i + 1, ProductType : Titanium, Stack : 0}
	}

	return minors
}

// Emulates the work of a Minor
func Produce(minor Minor, wg *sync.WaitGroup) {
	minor.Stack += 1 // Create ore each turn
	MineRequests <- MineRequest{Deposit: true, Ores: map[MineType]int{minor.ProductType : minor.Stack}, OresToTake:nil} //Deliver to mine
	resp := <- MineResponses
	minor.Stack -= resp[minor.ProductType]

	Requests <- minor
	respMinor := <- responses
	if respMinor == true {
		wg.Done()
	}
}

func CoordinateMinors(minors []Minor) {
	for {
		req := <-Requests
		if req.Id == 0 {
			continue
		}
		minors[req.Id - 1].ProductType = req.ProductType
		minors[req.Id - 1].Stack = req.Stack
		responses <- true
	}
}

func checkDeposit(mines map[MineType] int, productType MineType, productNb int) int {
	var toDeliver int

	if mines[productType] + productNb > siloSize {
		toDeliver = siloSize - mines[productType]
	} else {
		toDeliver = productNb
	}

	mines[productType] += toDeliver

	return toDeliver
}

func checkToGive(mines map[MineType] int, productType MineType) bool {
	if mines[productType] > 1 {
		mines[productType] -= 1
		return true
	} else {
		return false
	}
}

func CoordinateMines(mines map[MineType] int) {
	for {
		req := <-MineRequests

		if req.Deposit {
			delivered := make(map[MineType]int)
			delivered[Aluminium] = checkDeposit(mines, Aluminium, req.Ores[Aluminium])
			delivered[Iron] = checkDeposit(mines, Iron, req.Ores[Iron])
			delivered[Titanium] = checkDeposit(mines, Titanium, req.Ores[Titanium])

			MineResponses <- delivered
		} else {
			toGive := make(map[MineType]int)
			for i := 0; i < len(req.OresToTake); i++ {
				if checkToGive(mines, Aluminium) {
					toGive[Aluminium] += 1
					continue
				} else if checkToGive(mines, Iron) {
					toGive[Iron] += 1
					continue
				} else if checkToGive(mines, Titanium) {
					toGive[Titanium] += 1
					continue
				}
			}
			MineResponses <- toGive
		}
	}
}

func Mines() map[MineType]int {
	mines := make(map[MineType]int)
	mines[Iron] = 0
	mines[Titanium] = 0
	mines[Aluminium] = 0

	return mines
}

func DescribeMinor(minor Minor) {
	var productType = "Titanium"
	if minor.ProductType == Aluminium {
		productType = "Aluminium"
	} else if minor.ProductType == Iron {
		productType = "Iron"
	}
	fmt.Printf("Minor %d who produced %s has in stack %d\n", minor.Id, productType, minor.Stack)
}