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
	Deposit bool // true if deposit from minor, false if request to take from carriers
	Ores map[MineType]int // Used for deposit
	OresToTake []MineType // Used for taking ores
}

type Minor struct {
	id int
	productType MineType
	totalProduced int
	stack int
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
		minors[i] = Minor{id: i + 1, productType: Aluminium, stack: 0, totalProduced: 0}
	}

	var maxMinorId = i + nbIronMinors
	for ; i < maxMinorId; i++ {
		minors[i] = Minor{id: i + 1, productType: Iron, stack: 0, totalProduced: 0}
	}

	maxMinorId = i + nbTitaniumMinors
	for ; i < maxMinorId; i++ {
		minors[i] = Minor{id: i + 1, productType: Titanium, stack: 0, totalProduced: 0}
	}

	return minors
}

// Emulates the work of a Minor
func Produce(minor Minor, wg *sync.WaitGroup) {
	minor.stack += 1 // Create ore each turn
	minor.totalProduced += 1

	MineRequests <-MineRequest{Deposit: true, Ores: map[MineType]int{minor.productType : minor.stack}, OresToTake:nil} //Deliver to mine
	resp := <-MineResponses

	minor.stack -= resp[minor.productType]

	Requests <-minor
	respMinor := <-responses
	if respMinor == true {
		wg.Done()
	}
}

func CoordinateMinors(minors []Minor) {
	for {
		req := <-Requests
		if req.id == 0 {
			continue
		}
		minors[req.id - 1].stack = req.stack
		minors[req.id - 1].totalProduced = req.totalProduced
		responses <-true
	}
}

// Checks if a silo has enough places to add one ore
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

// Checks if the ore exists and can be given
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
			var delivered = make(map[MineType]int)
			delivered[Aluminium] = checkDeposit(mines, Aluminium, req.Ores[Aluminium])
			delivered[Iron] = checkDeposit(mines, Iron, req.Ores[Iron])
			delivered[Titanium] = checkDeposit(mines, Titanium, req.Ores[Titanium])

			MineResponses <-delivered
		} else {
			var toGive = make(map[MineType]int)
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
			MineResponses <-toGive
		}
	}
}

// Init mines
func Mines() map[MineType]int {
	mines := make(map[MineType]int)
	mines[Iron] = 0
	mines[Titanium] = 0
	mines[Aluminium] = 0

	return mines
}

func DescribeMines(mines map[MineType]int) {
	fmt.Printf("[Mine]\n\tType: Aluminium\n\tStack: %d\n", mines[Aluminium])
	fmt.Printf("[Mine]\n\tType: Titanium\n\tStack: %d\n", mines[Titanium])
	fmt.Printf("[Mine]\n\tType: Iron\n\tStack: %d\n", mines[Iron])
}

func DescribeMinor(minor Minor) {
	var productType = "Titanium"
	if minor.productType == Aluminium {
		productType = "Aluminium"
	} else if minor.productType == Iron {
		productType = "Iron"
	}
	fmt.Printf("[Minor %d]\n\tType: %s\n\tStack: %d\n\tTotal: %d\n", minor.id, productType, minor.stack, minor.totalProduced)
}