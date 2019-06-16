package main

import (
	"Go_Exo/mines"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type stat struct {
	nbMinor int
	cost int
}

func main()  {
	if len(os.Args) != 5 {
		fmt.Print("Usage: nb_turn nb_AluminiumMinor nb_IronMinor nb_TitaniumMinor\n")
	} else {
		nbTurn, err := strconv.Atoi(os.Args[1])
		nbAluminiumMinors, errAlu := strconv.Atoi(os.Args[2])
		nbIronMinors, errIron := strconv.Atoi(os.Args[3])
		nbTitaniumMinors, errTitanium := strconv.Atoi(os.Args[4])
		if err != nil || nbTurn <= 0 || errAlu != nil || errIron != nil || errTitanium != nil{
			fmt.Printf("Error: Wrong given arguments.\n")
		} else {
			simulate(nbTurn, nbAluminiumMinors, nbIronMinors, nbTitaniumMinors)
		}
	}
}

func simulate(nbTurn int, nbAluminiumMinors int, nbIronMinors int, nbTitaniumMinors int) {
	var mine = mines.Mines()
	var minors = mines.DefineMinors(nbAluminiumMinors, nbIronMinors, nbTitaniumMinors)
	var stat = stat{nbMinor: len(minors), cost: 0}

	go mines.CoordinateMinors(mine, minors)

	for i := 0; i < nbTurn; i++ {
		turn(mine, minors)
		stat.cost += 10 * stat.nbMinor
	}
	close(mines.Requests)
	printStats(stat, mine, minors)
}

func turn(mine map[mines.MineType]int, minors []mines.Minor) {
	var wg sync.WaitGroup

	for i := 0; i < len(minors); i++ {
		wg.Add(1)
		go mines.Product(minors[i], &wg)
	}
	wg.Wait()
}

func printStats(stats stat, mine map[mines.MineType]int, minors []mines.Minor) {
	fmt.Printf("Simulation finished.\nMinors: %d\nCost: %d\n", stats.nbMinor, stats.cost)
	fmt.Printf("Titanium in mine: %d\nAluminium in mine: %d\nIron in mine: %d\n", mine[mines.Titanium], mine[mines.Aluminium], mine[mines.Iron])

	for i := 0; i < len(minors); i++ {
	 	var productType = "Titanium"
	 	if minors[i].ProductType == mines.Aluminium {
			productType = "Aluminium"
		} else if minors[i].ProductType == mines.Iron {
			productType = "Iron"
		}
		fmt.Printf("Minor %d who produced %s has in stack %d\n", minors[i].Id, productType, minors[i].Stack)
	}
}