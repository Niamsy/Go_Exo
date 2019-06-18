package main

import (
	"Go_Exo/carriers"
	"Go_Exo/factories"
	"Go_Exo/mines"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type stat struct {
	nbMinor int
	nbCarrier int
	nbWorker int
	cost int
}

func main()  {
	if len(os.Args) != 9 {
		fmt.Print("Usage: nb_turn nb_AluminiumMinors nb_IronMinors nb_TitaniumMinors nb_Carriers nb_SteelWorkers nb_SignWorkers nb_HeatShieldWorkers\n")
	} else {
		nbTurn, err := strconv.Atoi(os.Args[1])
		nbAluminiumMinors, errAlu := strconv.Atoi(os.Args[2])
		nbIronMinors, errIron := strconv.Atoi(os.Args[3])
		nbTitaniumMinors, errTitanium := strconv.Atoi(os.Args[4])
		nbCarriers, errCarrier := strconv.Atoi(os.Args[5])
		nbSteelWorkers, errSteel := strconv.Atoi(os.Args[6])
		nbSignWorkers, errSign := strconv.Atoi(os.Args[6])
		nbHeatShieldWorkers, errShield := strconv.Atoi(os.Args[6])
		if err != nil || nbTurn <= 0 || errAlu != nil || errIron != nil || errTitanium != nil || errCarrier != nil || errSteel != nil || errSign != nil || errShield != nil{
			fmt.Printf("Error: Wrong given arguments.\n")
		} else {
			simulate(nbTurn, nbAluminiumMinors, nbIronMinors, nbTitaniumMinors, nbCarriers, nbSteelWorkers, nbSignWorkers, nbHeatShieldWorkers)
		}
	}
}

func simulate(nbTurn int, nbAluminiumMinors int, nbIronMinors int, nbTitaniumMinors int, nbCarriers int, nbSteelWorkers int, nbSignWorkers int, nbHeatShieldWorkers int) {
	var mine = mines.Mines()
	var minors = mines.DefineMinors(nbAluminiumMinors, nbIronMinors, nbTitaniumMinors)
	var carrierActors = carriers.DefineCarriers(nbCarriers)
	var workers = factories.DefineWorkers(nbSteelWorkers, nbSignWorkers, nbHeatShieldWorkers)
	var factory = factories.Factories()
	var stat = stat{nbMinor: len(minors), nbCarrier: len(carrierActors), nbWorker: len(workers),cost: 0}

	go mines.CoordinateMinors(minors)
	go factories.CoordinateWorkers(workers)
	go carriers.CoordinateCarriers(carrierActors)
	go mines.CoordinateMines(mine)
	go factories.CoordinateFactories(factory)

	for i := 0; i < nbTurn; i++ {
		turn(minors, carrierActors, workers)
		stat.cost += 10 * (stat.nbMinor + stat.nbCarrier + stat.nbWorker)
	}

	close(mines.Requests)
	close(mines.MineRequests)
	close(factories.FactoryRequests)
	close(factories.Requests)

	printStats(stat, mine, minors, carrierActors, workers, factory)
}

func turn(minors []mines.Minor, carrierActors []carriers.Carrier, workers []factories.Worker) {
	var wg sync.WaitGroup

	for i := 0; i < len(minors); i++ {
		wg.Add(1)
		go mines.Produce(minors[i], &wg)
	}
	for i := 0; i < len(workers); i++ {
		wg.Add(1)
		go factories.Produce(workers[i], &wg)
	}
	for i := 0; i < len(carrierActors); i++ {
		wg.Add(1)
		go carriers.Carry(carrierActors[i], &wg)
	}
	wg.Wait()
}

func printStats(stats stat, mine map[mines.MineType]int, minors []mines.Minor, carriersActor []carriers.Carrier, workers []factories.Worker, factory []factories.Factory) {
	fmt.Printf("Simulation finished.\nMinors: %d\nCarriers %d\nWorkers %d\nCost: %d\n", stats.nbMinor, stats.nbCarrier, stats.nbWorker, stats.cost)
	fmt.Printf("Titanium in mine: %d\nAluminium in mine: %d\nIron in mine: %d\n", mine[mines.Titanium], mine[mines.Aluminium], mine[mines.Iron])
	factories.DescribeFactories(factory)

	for i := 0; i < len(minors); i++ {
	 	mines.DescribeMinor(minors[i])
	}

	for i := 0; i < len(carriersActor); i++ {
		carriers.DescribeCarrier(carriersActor[i])
	}

	for i := 0; i < len(workers); i++ {
		factories.DescribeWorker(workers[i])
	}
}