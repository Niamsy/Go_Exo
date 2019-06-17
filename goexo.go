package main

import (
	"Go_Exo/carriers"
	"Go_Exo/mines"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type stat struct {
	nbMinor int
	nbCarrier int
	cost int
}

func main()  {
	if len(os.Args) != 6 {
		fmt.Print("Usage: nb_turn nb_AluminiumMinors nb_IronMinors nb_TitaniumMinors nb_Carriers\n")
	} else {
		nbTurn, err := strconv.Atoi(os.Args[1])
		nbAluminiumMinors, errAlu := strconv.Atoi(os.Args[2])
		nbIronMinors, errIron := strconv.Atoi(os.Args[3])
		nbTitaniumMinors, errTitanium := strconv.Atoi(os.Args[4])
		nbCarriers, errCarrier := strconv.Atoi(os.Args[5])
		if err != nil || nbTurn <= 0 || errAlu != nil || errIron != nil || errTitanium != nil || errCarrier != nil {
			fmt.Printf("Error: Wrong given arguments.\n")
		} else {
			simulate(nbTurn, nbAluminiumMinors, nbIronMinors, nbTitaniumMinors, nbCarriers)
		}
	}
}

func simulate(nbTurn int, nbAluminiumMinors int, nbIronMinors int, nbTitaniumMinors int, nbCarriers int) {
	var mine = mines.Mines()
	var minors = mines.DefineMinors(nbAluminiumMinors, nbIronMinors, nbTitaniumMinors)
	var carrierActors = carriers.DefineCarriers(nbCarriers)
	var stat = stat{nbMinor: len(minors), nbCarrier: len(carrierActors), cost: 0}

	go mines.CoordinateMinors(minors)
	go carriers.CoordinateCarriers(carrierActors)
	go mines.CoordinateMine(mine)

	for i := 0; i < nbTurn; i++ {
		turn(minors, carrierActors)
		stat.cost += 10 * (stat.nbMinor + stat.nbCarrier)
	}
	close(mines.Requests)
	close(mines.MineRequests)
	printStats(stat, mine, minors, carrierActors)
}

func turn(minors []mines.Minor, carrierActors []carriers.Carrier) {
	var wg sync.WaitGroup

	for i := 0; i < len(minors); i++ {
		wg.Add(1)
		go mines.Product(minors[i], &wg)
	}
	for i := 0; i < len(carrierActors); i++ {
		wg.Add(1)
		go carriers.Carry(carrierActors[i], &wg)
	}
	wg.Wait()
}

func printStats(stats stat, mine map[mines.MineType]int, minors []mines.Minor, carriersActor []carriers.Carrier) {
	fmt.Printf("Simulation finished.\nMinors: %d\nCarriers %d\nCost: %d\n", stats.nbMinor, stats.nbCarrier, stats.cost)
	fmt.Printf("Titanium in mine: %d\nAluminium in mine: %d\nIron in mine: %d\n", mine[mines.Titanium], mine[mines.Aluminium], mine[mines.Iron])

	for i := 0; i < len(minors); i++ {
	 	mines.DescribeMinor(minors[i])
	}

	for i := 0; i < len(carriersActor); i++ {
		carriers.DescribeCarrier(carriersActor[i])
	}
}