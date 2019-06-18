package carriers

import (
	"Go_Exo/factories"
	"Go_Exo/mines"
	"fmt"
	"sync"
)

const cargoSize = 11

var Requests = make(chan Carrier)
var responses = make(chan bool)

type Carrier struct {
	id int
	loaded bool
	cargo []mines.MineType
	totalOreTaken map[mines.MineType]int
}

func DefineCarriers(nbCarriers int) []Carrier{
	var carriers = make([]Carrier, nbCarriers)

	for i := 0; i < nbCarriers; i++ {
		carriers[i].id = i + 1
		carriers[i].loaded = false
		carriers[i].cargo = make([]mines.MineType, cargoSize)
		for j := 0; j < cargoSize; j++ {
			carriers[i].cargo[j] = mines.Default
		}
		carriers[i].totalOreTaken = make(map[mines.MineType]int)
		carriers[i].totalOreTaken[mines.Titanium] = 0
		carriers[i].totalOreTaken[mines.Aluminium] = 0
		carriers[i].totalOreTaken[mines.Iron] = 0
	}

	return carriers
}

func checkRemainingPlaces(carrier Carrier) int {
	var remainingPlaces = 0

	for i := 0; i < cargoSize; i++ {
		if carrier.cargo[i] == mines.Default {
			remainingPlaces += 1
		}
	}
	return remainingPlaces
}

func updateCarrierCargo(carrier Carrier, oresToTake map[mines.MineType]int) Carrier {
	var i = 0
	for ; carrier.cargo[i] != mines.Default; i++ {

	}

	for j := 0; j < oresToTake[mines.Iron]; j++ {
		carrier.cargo[i] = mines.Iron
		i++
	}

	for j := 0; j < oresToTake[mines.Titanium]; j++ {
		carrier.cargo[i] = mines.Titanium
		i++
	}

	for j := 0; j < oresToTake[mines.Aluminium]; j++ {
		carrier.cargo[i] = mines.Aluminium
		i++
	}

	return carrier
}

func updateTotalOresTaken(carrier Carrier, oresToTake map[mines.MineType]int) Carrier {
	carrier.totalOreTaken[mines.Iron] += oresToTake[mines.Iron]
	carrier.totalOreTaken[mines.Aluminium] += oresToTake[mines.Aluminium]
	carrier.totalOreTaken[mines.Titanium] += oresToTake[mines.Titanium]
	return carrier
}

//Emulates work of a carrier
func Carry(carrier Carrier, wg *sync.WaitGroup) {
	var remainingPlaces = checkRemainingPlaces(carrier)

	if remainingPlaces == 0 {
		carrier.loaded = true
	} else {
		carrier.loaded = false
	}

	if carrier.loaded {
		factories.FactoryRequests <-factories.FactoryRequest{Deposit: true, ToProduce: factories.Default, ToDeposit: carrier.cargo}
		resp := <-factories.FactoryResponsesToCarriers
		carrier.cargo = resp
	} else {
		mines.MineRequests <-mines.MineRequest{Deposit: false, Ores: nil, OresToTake: make([]mines.MineType, remainingPlaces)}

		resp := <-mines.MineResponses
		carrier = updateCarrierCargo(carrier, resp)
		carrier = updateTotalOresTaken(carrier, resp)
	}
	Requests <-carrier
	respCarrier := <-responses
	if respCarrier {
		wg.Done()
	}
}

func CoordinateCarriers(carriers []Carrier) {
	for {
		req := <-Requests
		if req.id == 0 {
			continue
		}
		carriers[req.id - 1] = req
		responses <-true
	}
}

func DescribeCarrier(carrier Carrier) {
	var nbIron = 0
	var nbAluminium = 0
	var nbTitanium = 0

	for i := 0; i < cargoSize; i++ {
		if carrier.cargo[i] == mines.Titanium {
			nbTitanium += 1
		} else if carrier.cargo[i] == mines.Iron {
			nbIron += 1
		} else if carrier.cargo[i] == mines.Aluminium {
			nbAluminium += 1
		}
	}
	fmt.Printf("[Carrier %d]\n\tAluminium: %d\n\tIron: %d\n\tTitanium: %d\n\tTotal taken Aluminium: %d\n\tTotal taken Iron: %d\n\tTotal taken Titanium: %d\n", carrier.id, nbAluminium, nbIron, nbTitanium, carrier.totalOreTaken[mines.Aluminium], carrier.totalOreTaken[mines.Iron], carrier.totalOreTaken[mines.Titanium])
}