package factories

import (
	"Go_Exo/mines"
	"fmt"
	"sync"
)

type Type int
const (
	Steel Type = iota
	Sign
	HeatShield
	Default
)

type produceFunc func(map[mines.MineType]int) bool

type Factory struct {
	factoryType Type
	silos        map[mines.MineType]int
	production  int
	function produceFunc
}

type Worker struct {
	Id int
	ProductType Type
	Produced int
}

type FactoryRequest struct {
	Deposit bool // true if it's deposit in silos false if it's request to create a production
	ToDeposit []mines.MineType
	ToProduce Type
}

const siloSize = 15

var FactoryRequests = make(chan FactoryRequest)
var FactoryResponsesToCarriers = make(chan []mines.MineType)
var factoryResponsesToWorkers = make(chan bool)

var Requests = make(chan Worker)
var responses = make(chan bool)

func DefineWorkers(nbSteelWorkers int, nbSignWorkers int, nbHeatShieldWorkers int) []Worker{
	var totalWorkers = nbSteelWorkers + nbSignWorkers + nbHeatShieldWorkers
	var workers = make([]Worker, totalWorkers)

	var i = 0

	for ; i < nbSteelWorkers; i++ {
		workers[i] = Worker{Id : i + 1, ProductType : Steel, Produced : 0}
	}

	var maxMinorId = i + nbSignWorkers
	for ; i < maxMinorId; i++ {
		workers[i] = Worker{Id : i + 1, ProductType : Sign, Produced : 0}
	}

	maxMinorId = i + nbHeatShieldWorkers
	for ; i < maxMinorId; i++ {
		workers[i] = Worker{Id : i + 1, ProductType : HeatShield, Produced : 0}
	}

	return workers
}

// Emulates the work of a Worker
func Produce(worker Worker, wg *sync.WaitGroup) {
	FactoryRequests <-FactoryRequest{Deposit : false, ToDeposit : nil, ToProduce : worker.ProductType}
	resp := <-factoryResponsesToWorkers
	if resp == true {
		worker.Produced += 1
		Requests <-worker
		respMinor := <- responses
		if respMinor == true {
			wg.Done()
		}
	} else {
		wg.Done()
	}
}

func CoordinateWorkers(workers []Worker) {
	for {
		req := <-Requests
		if req.Id == 0 {
			continue
		}
		workers[req.Id - 1].Produced = req.Produced
		responses <- true
	}
}

func produceSteel(silos map[mines.MineType]int) bool {
	if silos[mines.Iron] >= 2 {
		silos[mines.Iron] -= 2
		return true
	}
	return false
}

func produceSign(silos map[mines.MineType]int) bool {
	if silos[mines.Iron] >= 1 && silos[mines.Titanium] >= 1 && silos[mines.Aluminium] >= 2 {
		silos[mines.Iron] -= 1
		silos[mines.Titanium] -= 1
		silos[mines.Aluminium] -= 2
		return true
	}
	return false
}

func produceHeatShield(silos map[mines.MineType]int) bool {
	if silos[mines.Iron] >= 3 && silos[mines.Titanium] >= 1 && silos[mines.Aluminium] >= 2 {
		silos[mines.Iron] -= 1
		silos[mines.Titanium] -= 1
		silos[mines.Aluminium] -= 2
		return true
	}
	return false
}

// Determines which factory has the more need for an ore
func chooseFactoryToDeliver(factories []Factory, ore mines.MineType) int{
	var factoryIndex = 1 // put in the sign factory by default

	for i := 0; i < len(factories); i++ {
		if _, exists := factories[i].silos[ore]; exists { // check if silo of certain ore exists
			if factories[factoryIndex].production > factories[i].production {
				factoryIndex = i
			}
		}
	}
	return factoryIndex
}

func deliverToFactories(factories []Factory, toDeliver []mines.MineType) {
	for i := 0; i < len(toDeliver); i++ {
		if toDeliver[i] != mines.Default {
			var idx = chooseFactoryToDeliver(factories, toDeliver[i])
			factories[idx].silos[toDeliver[i]] += 1
			toDeliver[i] = mines.Default
		}
	}
}

func CoordinateFactories(factories []Factory) {
	for {
		req := <-FactoryRequests

		if req.Deposit {
			deliverToFactories(factories, req.ToDeposit)
			FactoryResponsesToCarriers <-req.ToDeposit
		} else {
			var produced = factories[req.ToProduce].function(factories[req.ToProduce].silos)
			if produced == true {
				factories[req.ToProduce].production += 1
			}
			factoryResponsesToWorkers <-produced
		}
	}
}

func Factories() []Factory {
	var factories = make([]Factory, 3)

	factories[Steel] = Factory{factoryType : Steel, production : 0, silos : map[mines.MineType]int{mines.Iron : 0}, function : produceSteel}
	factories[Sign] = Factory{factoryType : Sign, production : 0, silos : map[mines.MineType]int{mines.Iron : 0, mines.Aluminium : 0, mines.Titanium : 0}, function : produceSign}
	factories[HeatShield] = Factory{factoryType : HeatShield, production : 0, silos : map[mines.MineType]int{mines.Iron : 0, mines.Aluminium : 0, mines.Titanium : 0}, function : produceHeatShield}

	return factories
}

func DescribeFactories(factories []Factory) {
	for _, value := range factories {
		var productType = "Steel"
		if value.factoryType == Sign {
			productType = "Sign"
		} else if value.factoryType == HeatShield {
			productType = "Heat Shield"
		}
		fmt.Printf("[Factory]\n\tType: %s\n\tTotal production: %d\n\tUnused Iron: %d\n\tUnused Titanium: %d\n\tUnused Aluminium: %d\n", productType, value.production, value.silos[mines.Iron], value.silos[mines.Titanium], value.silos[mines.Aluminium])
	}
}

func DescribeWorker(worker Worker) {
	var productType = "Steel"
	if worker.ProductType == Sign {
		productType = "Sign"
	} else if worker.ProductType == HeatShield {
		productType = "Heat Shield"
	}
	fmt.Printf("[Worker %d]\n\tType: %s\n\tTotal produced: %d\n", worker.Id, productType, worker.Produced)
}



