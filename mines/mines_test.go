package mines

import "testing"

func TestDefineMinors(t *testing.T) {
	var nbAluminium = 0
	var nbIron = 0
	var nbTitanium = 0

	minors := DefineMinors(1, 2, 3)
	if len(minors) != 6 {
		t.Fail()
	}

	for i := 0; i < len(minors); i++ {
		if minors[i].productType == Iron {
			nbIron += 1
		} else if minors[i].productType == Aluminium {
			nbAluminium += 1
		} else if minors[i].productType == Titanium {
			nbTitanium += 1
		}

		if minors[i].stack != 0 {
			t.Fail()
		}
		if minors[i].totalProduced != 0 {
			t.Fail()
		}
	}

	if nbAluminium != 1 || nbIron != 2 || nbTitanium != 3 {
		t.Fail()
	}
}

func TestCheckDeposit(t *testing.T) {
	var mines = make(map[MineType] int)

	mines[Aluminium] = siloSize
	mines[Titanium] = siloSize - 2
	mines[Iron] = siloSize - 6

	if checkDeposit(mines, Aluminium, 3) != 0 {
		t.Fail()
	}
	if checkDeposit(mines, Titanium, 3) != 2 || mines[Titanium] != siloSize {
		t.Fail()
	}
	if checkDeposit(mines, Iron, 4) != 4 || mines[Iron] != siloSize - 2{
		t.Fail()
	}
}

func TestCheckToGive(t *testing.T) {
	var mines = make(map[MineType] int)

	mines[Aluminium] = siloSize
	mines[Titanium] = 0
	mines[Iron] = siloSize - 6

	if checkToGive(mines, Aluminium) == false || mines[Aluminium] != siloSize - 1 {
		t.Fail()
	}
	if checkToGive(mines, Titanium) == true || mines[Titanium] != 0 {
		t.Fail()
	}
	if checkToGive(mines, Iron) == false || mines[Iron] != siloSize - 7 {
		t.Fail()
	}
}

func TestMines(t *testing.T) {
	var mines = Mines()

	if mines[Aluminium] != 0 || mines[Titanium] != 0 || mines[Iron] != 0 {
		t.Fail()
	}
}

func TestCoordinateMinors(t *testing.T) {
	var minors = DefineMinors(1, 1, 1)
	minors[0].stack = 1
	go CoordinateMinors(minors)
	minors[0].stack += 1
	Requests <-minors[0]
	resp := <-responses
	if resp != true {
		t.Fail()
	}
	if minors[0].stack != 2 {
		t.Fail()
	}
}