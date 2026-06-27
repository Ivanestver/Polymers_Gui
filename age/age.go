package age

import (
	"errors"
	"fmt"
	"math/rand"
	"polymers/base"
	dt "polymers/datatypes"
	"polymers/outputformat"
	"polymers/views"
	"sort"
	"strconv"
)

const crosslinksCount = 0.5

type AgeAlg struct {
	globula *views.GlobulaView
}

func NewAgeAlg(globula *views.GlobulaView) *AgeAlg {
	return &AgeAlg{
		globula: globula,
	}
}

func (alg *AgeAlg) DoAging1(groupsCount int) {
	/*
		The aging process consists of 4 stages:
		1. Break connections and randomly assign the new ends as B and C so that we have 26% of B and 26% of C
		2. A part of Bs (3% of all groups to take) turn into C so that 29% of all groups are C
		3. Randomly turn other bins into C (6% in the first time and 6% in the second time) so that 41% of all groups are C
		4. Randomly create connections between bins untouched so that 59% of groups are B

		Example:
		Imagine we have 1000 bins at all and 100 aged groups required. The process is going to be the following:
		1. Break randomly (26/2)% of 100 = 26 connections and create 26 Bs and 26 Cs. 48 bins are left untouched
		2. Turn 3% of 100 = 3 Bs into C. 23 Bs left, 29 Cs left. 48 bins are left untouched
		3. Turn 6% of 100 = 6 bins randomly taken in globula into C twice. For now, we don't take Bs into account. 23 Bs left, 41 Cs left. 36 bins are left untouched
		4. Turn 59%-23%=36% of 100 = 36 bins untouched into Hs and create connections. 2 turns happen for 1 connection so that the number of connections is 36/2 = 18
	*/

	// =================DEBUG=================
	// // 4. Create connections
	// if turn == 3 {
	// 	globula.createBConnections(int(float64(groupsCount)*0.59), &Bs_)
	// 	turn++
	// }

	// // 3. Turn random bins into Cs (excluding Bs)
	// if turn == 2 {
	// 	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// 	// Do it twice
	// 	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// 	turn++
	// }

	// // 2. Turn some Bs into C
	// if turn == 1 {
	// 	turnBIntoC(&Bs_, int(float64(groupsCount)*0.03))
	// 	turn++
	// }

	// // 1. Break connections
	// if turn == 0 {
	// 	Bs_ = globula.breakConnections(int(float64(groupsCount) * 0.26))
	// 	turn++
	// }
	// =================DEBUG=================

	// 1. Break connections
	Cs_ := alg.breakConnections(int(float64(groupsCount)*0.26), base.O)

	// 2. Turn some Bs into C
	turnIntoAnotherGroup(&Cs_, int(float64(groupsCount)*0.03), base.N)

	// 3. Turn random bins into Cs (excluding Bs)
	alg.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// Do it twice
	alg.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))

	// 4. Create connections
	alg.createCrosslinks1(int(float64(groupsCount)*0.59), &Cs_)
	alg.globula.SetProperty(views.GlobulaAged)
}

func (alg *AgeAlg) DoAging2(groupsCount int, doCrosslinks bool) {
	// 1. Break connections
	Bs_ := alg.breakConnections(int(float64(groupsCount)*0.44), base.N)

	// 2. Turn some Bs into C
	turnIntoAnotherGroup(&Bs_, int(float64(groupsCount)*0.15), base.N)

	// 3. Turn random bins into Cs (excluding Bs)
	alg.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// Do it twice
	alg.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))

	// 4. Create connections
	if doCrosslinks {
		alg.createCrosslinks2(int(float64(groupsCount) * crosslinksCount))
	}
	alg.globula.SetProperty(views.GlobulaAged)
}

func (alg *AgeAlg) breakConnections(groupsCount int, monomerTypeToGather base.MendeleevTableElement) []*dt.Monomer {
	Bs := make([]*dt.Monomer, 0)
	triesNumber := 0
	for len(Bs) != groupsCount && triesNumber < groupsCount {
		chosenPoly := rand.Intn(alg.globula.Len()) // Take a random poly
		poly := alg.globula.GetPolymerByIdx(chosenPoly)
		if poly.Len() < 2 {
			continue
		}
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		nextMonomerNumber := chosenMonomerNumber + 1
		chosenMonomer := poly.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != base.C {
			triesNumber++
			continue
		}
		nextMonomer := poly.GetMonomerByIdx(nextMonomerNumber)
		if nextMonomer.MonomerType != base.C {
			triesNumber++
			continue
		}
		triesNumber = 0
		side := chosenMonomer.GetSideOfSibling(nextMonomer)
		dt.TierConnection(chosenMonomer, nextMonomer, side)
		if rand.Intn(2) == 0 {
			chosenMonomer.MonomerType = base.O
			nextMonomer.MonomerType = base.N
			if monomerTypeToGather == base.O {
				Bs = append(Bs, chosenMonomer)
			} else {
				Bs = append(Bs, nextMonomer)
			}
		} else {
			chosenMonomer.MonomerType = base.N
			nextMonomer.MonomerType = base.O
			if monomerTypeToGather == base.N {
				Bs = append(Bs, chosenMonomer)
			} else {
				Bs = append(Bs, nextMonomer)
			}
		}
	}

	return Bs
}

func turnIntoAnotherGroup(Bs *[]*dt.Monomer, groupsCount int, monomerTypeToTurn base.MendeleevTableElement) {
	mapUsedBs := make(map[int]bool)
	triesCount := 0
	for len(mapUsedBs) != groupsCount && triesCount < 3*groupsCount {
		mapUsedBs[rand.Intn(len(*Bs))] = true
		triesCount++
	}

	arr := make([]int, 0)
	for key := range mapUsedBs {
		arr = append(arr, key)
	}
	sort.Ints(arr)

	for i := len(arr) - 1; i >= 0; i-- {
		(*Bs)[arr[i]].MonomerType = monomerTypeToTurn
		*Bs = append((*Bs)[:arr[i]], (*Bs)[arr[i]+1:]...)
	}
}

func (alg *AgeAlg) turnRandomBinsIntoC(groupsCount int) {
	i := 0
	triesNumber := 0
	for i < groupsCount && triesNumber < groupsCount {
		chosenPoly := rand.Intn(alg.globula.Len()) // Take a random poly
		poly := alg.globula.GetPolymerByIdx(chosenPoly)
		if poly.Len() < 2 {
			triesNumber++
			continue
		}
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType == base.C {
			chosenMonomer.MonomerType = base.N
			i++
			triesNumber = 0
		} else {
			triesNumber++
		}
	}
}

func (alg *AgeAlg) createCrosslinks1(groupsCount int, Bs *[]*dt.Monomer) {
	for len(*Bs) != groupsCount {
		chosenPoly := rand.Intn(alg.globula.Len()) // Take a random poly
		poly := alg.globula.GetPolymerByIdx(chosenPoly)
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != base.C {
			continue
		}

		movementSides := dt.GetMovementSides()
		for i := 0; i < len(movementSides); i++ {
			chosenSide := movementSides[rand.Intn(len(movementSides))]
			nextMonomer, err := chosenMonomer.GetSibling(chosenSide)
			if err == nil && nextMonomer != nil && nextMonomer.MonomerType == base.C {
				dt.MakeConnection(chosenMonomer, nextMonomer, dt.ConnectionTypeCrosslinks)
				chosenMonomer.MonomerType = base.O
				nextMonomer.MonomerType = base.O
				*Bs = append(*Bs, chosenMonomer, nextMonomer)
				break
			}
		}
	}
}

func (alg *AgeAlg) createCrosslinks2(crosslinksCount int) {
	currentCount := 0
	timesRepeated := 0
	const maxTimesRepeated = 1000
	for currentCount != crosslinksCount && timesRepeated != maxTimesRepeated {
		chosenPoly := rand.Intn(alg.globula.Len()) // Take a random poly
		poly := alg.globula.GetPolymerByIdx(chosenPoly)
		if poly.Len() < 2 {
			continue
		}
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != base.C {
			timesRepeated++
			continue
		}

		movementSides := dt.GetMovementSides()
		for i := 0; i < len(movementSides); i++ {
			chosenSide := movementSides[rand.Intn(len(movementSides))]
			nextMonomer, err := chosenMonomer.GetSibling(chosenSide)
			if err == nil &&
				nextMonomer != nil &&
				nextMonomer.MonomerType == base.C &&
				(!dt.MonomersAreEqual(chosenMonomer.NextMonomer, nextMonomer) &&
					!dt.MonomersAreEqual(chosenMonomer.PrevMonomer, nextMonomer)) {
				dt.MakeConnection(chosenMonomer, nextMonomer, dt.ConnectionTypeCrosslinks)
				chosenMonomer.MonomerType = base.H
				nextMonomer.MonomerType = base.H
				currentCount += 1
				timesRepeated = 0
				break
			}
		}
	}
	if timesRepeated == maxTimesRepeated {
		print("Didn't make all crosslinks! The number of crosslinks done: " + strconv.FormatInt(int64(currentCount), 10))
	}
}

func (alg *AgeAlg) DoAging3(ncut, OcontainingCount, ncross int) error {
	if ncut < OcontainingCount {
		return errors.New("ncut is less that the O-containing monomers count")
	}
	printer := outputformat.GetPrint()
	// First, distribute O containing monomers
	// if warning, err := globula.aging3DistributeCutMonomers(OcontainingCount,
	// 	dt.MONOMER_TYPE_O_CONTAINING,
	// 	dt.MONOMER_TYPE_VYNIL); warning != nil {
	// 	printer.PrintlnWarning(warning.Error())
	// } else if err != nil {
	// 	return err
	// }
	alg.breakConnections(OcontainingCount, base.N)

	// Then, distribute what's left
	// if warning, err := globula.aging3DistributeCutMonomers(ncut-OcontainingCount,
	// 	dt.MONOMER_TYPE_VYNIL,
	// 	dt.MONOMER_TYPE_VYNIL); warning != nil {
	// 	printer.PrintlnWarning(warning.Error())
	// } else if err != nil {
	// 	return err
	// }
	alg.turnRandomBinsIntoC(ncut - OcontainingCount)

	// At the end, distribute crosslinks
	if warning, err := alg.aging3DistributeCrosslinks(ncross); warning != nil {
		printer.PrintlnWarning(warning.Error())
	} else if err != nil {
		return err
	}
	return nil
}

func (alg *AgeAlg) aging3DistributeCutMonomers(OcontainingCount int, typePrev, typeNext base.MendeleevTableElement) (warning, err error) {
	warning = nil
	err = nil
	trialsCount := 0
	for OcontainingCount > 0 && trialsCount < 10 {
		chosenPolymerNumber := rand.Intn(alg.globula.Len())
		chosenPolymer := alg.globula.GetPolymerByIdx(chosenPolymerNumber)
		chosenMonomerNumber := rand.Intn(chosenPolymer.Len()-2) + 1
		chosenMonomer := chosenPolymer.GetMonomerByIdx(chosenMonomerNumber)
		nextChosenMonomer := chosenPolymer.GetMonomerByIdx(chosenMonomerNumber + 1)
		if chosenMonomer.MonomerType != base.C ||
			nextChosenMonomer.MonomerType != base.C {
			trialsCount++
			continue
		}
		dt.BreakConnection1(chosenMonomer, nextChosenMonomer)
		chosenMonomer.MonomerType = typePrev
		nextChosenMonomer.MonomerType = typeNext
		OcontainingCount--
		trialsCount = 0
	}
	if trialsCount != 0 {
		return fmt.Errorf("%d O containing monomers weren't distributed", OcontainingCount), nil
	}
	return nil, nil
}

func (alg *AgeAlg) aging3DistributeCrosslinks(ncross int) (warning, err error) {
	warning = nil
	err = nil
	trialsCount := 0
	for ncross > 0 && trialsCount < 10 {
		chosenPolymerNumber := rand.Intn(alg.globula.Len())
		chosenPolymer := alg.globula.GetPolymerByIdx(chosenPolymerNumber)
		chosenMonomerNumber := rand.Intn(chosenPolymer.Len()-2) + 1
		chosenMonomer := chosenPolymer.GetMonomerByIdx(chosenMonomerNumber)
		nextChosenMonomer := chosenPolymer.GetMonomerByIdx(chosenMonomerNumber + 1)
		if chosenMonomer.MonomerType != base.C ||
			nextChosenMonomer.MonomerType != base.C {
			trialsCount++
			continue
		}
		dt.MakeConnection(chosenMonomer, nextChosenMonomer, dt.ConnectionTypeCrosslinks)
		chosenMonomer.MonomerType = base.H
		nextChosenMonomer.MonomerType = base.H
		ncross--
		trialsCount = 0
	}
	if trialsCount != 0 {
		return fmt.Errorf("%d crosslinks weren't distributed", ncross), nil
	}
	return nil, nil
}

func (alg *AgeAlg) aging4DistributeCrosslinks(ncross int) (warning, err error) {
	warning = nil
	err = nil
	trialsCount := 0
	for ncross > 0 && trialsCount < 10 {
		chosenPolymerNumber := rand.Intn(alg.globula.Len())
		chosenPolymer := alg.globula.GetPolymerByIdx(chosenPolymerNumber)
		chosenMonomerNumber := rand.Intn(chosenPolymer.Len())
		chosenMonomer := chosenPolymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != base.C {
			trialsCount++
			continue
		}
		movementSides := dt.GetMovementSides()
		for i := 0; i < len(movementSides); i++ {
			chosenSide := movementSides[rand.Intn(len(movementSides))]
			nextMonomer, _ := chosenMonomer.GetSibling(chosenSide)
			getBorderMonomer := func(side dt.Side, startMonomer *dt.Monomer) *dt.Monomer {
				curr := startMonomer
				prev, _ := curr.GetSibling(side)
				for prev != nil && prev.MonomerType != base.MendeleevTableElementUndefined {
					t := prev
					prev, _ = curr.GetSibling(side)
					curr = t
				}
				return curr
			}
			if (chosenMonomer.NextMonomer != nil && dt.MonomersAreEqual(chosenMonomer.NextMonomer, nextMonomer)) ||
				(chosenMonomer.PrevMonomer != nil && dt.MonomersAreEqual(chosenMonomer.PrevMonomer, nextMonomer)) {
				continue
			}
			if nextMonomer == nil ||
				nextMonomer.MonomerType == base.MendeleevTableElementUndefined {
				switch chosenSide {
				case dt.SideForward:
					nextMonomer = getBorderMonomer(dt.SideBackward, chosenMonomer)
				case dt.SideBackward:
					nextMonomer = getBorderMonomer(dt.SideForward, chosenMonomer)
				case dt.SideUp:
					nextMonomer = getBorderMonomer(dt.SideDown, chosenMonomer)
				case dt.SideDown:
					nextMonomer = getBorderMonomer(dt.SideUp, chosenMonomer)
				default:
					continue
				}
				dt.MakeConnectionUnsafe(chosenMonomer, nextMonomer, chosenSide, dt.ConnectionTypeCrosslinks)
			} else {
				dt.MakeConnection(chosenMonomer, nextMonomer, dt.ConnectionTypeCrosslinks)
			}
			chosenMonomer.MonomerType = base.H
			nextMonomer.MonomerType = base.H
			trialsCount = 0
			break
		}
	}
	if trialsCount != 0 {
		return fmt.Errorf("%d crosslinks weren't distributed", ncross), nil
	}
	return nil, nil
}

func (alg *AgeAlg) DoAgingSurface(ncut, nOContaining, ncross int) error {
	if ncut < nOContaining {
		return errors.New("ncut is less that the O-containing monomers count")
	}
	printer := outputformat.GetPrint()
	// First, distribute O containing monomers
	alg.breakConnections(nOContaining, base.N)

	// Then, distribute what's left
	alg.turnRandomBinsIntoC(ncut - nOContaining)

	// At the end, distribute crosslinks
	if warning, err := alg.aging4DistributeCrosslinks(ncross); warning != nil {
		printer.PrintlnWarning(warning.Error())
	} else if err != nil {
		return err
	}
	return nil
}
