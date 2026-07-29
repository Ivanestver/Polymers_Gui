package atomistic

import "polymers/base"

type _Atom struct {
	Number  int
	Label   string
	Mass    int
	Charge  int
	Valence int
	Coords  base.Vector3DF
}

type _AtomNumber = int
type _BondValence = int

type _Monomer struct {
	Name          string
	Atoms         []_Atom
	Bonds         map[_AtomNumber]map[_AtomNumber]_BondValence
	Mass          int
	Head          *_Atom
	Tail          *_Atom
	RotationPivot *_Atom
}

type _Polymer struct {
	Monomers []*_Monomer
	Bonds    map[_AtomNumber]map[_AtomNumber]_BondValence
}

func NewPolymer() *_Polymer {
	return &_Polymer{
		Monomers: make([]*_Monomer, 0),
		Bonds:    make(map[_AtomNumber]map[_AtomNumber]_BondValence),
	}
}

func (molecule *_Monomer) GetMassCenter() base.Vector3DF {
	res := molecule.Atoms[0].Coords
	for i := 1; i < len(molecule.Atoms); i++ {
		res.AddF(molecule.Atoms[i].Coords)
	}
	res.MultiplyByConstantF(1.0 / float64(len(molecule.Atoms)))
	return res
}

func (molecule *_Monomer) MoveTo(point *base.Vector3DF) {
	massCenter := molecule.GetMassCenter()
	direction := base.SubtractVecF(*point, massCenter)
	for i := 0; i < len(molecule.Atoms); i++ {
		currentAtom := &molecule.Atoms[i]
		currentAtom.Coords.AddF(direction)
	}
}

func (molecule *_Monomer) GetBondsCount() int {
	count := 0
	for _, m := range molecule.Bonds {
		count += len(m)
	}
	return count
}

func (molecule *_Monomer) Copy() *_Monomer {
	newMolecule := &_Monomer{
		Name:  molecule.Name,
		Atoms: make([]_Atom, len(molecule.Atoms)),
		Bonds: make(map[_AtomNumber]map[_AtomNumber]_BondValence, len(molecule.Bonds)),
		Mass:  molecule.Mass,
	}
	copy(newMolecule.Atoms, molecule.Atoms)
	for key, m := range molecule.Bonds {
		newMolecule.Bonds[key] = make(map[_AtomNumber]_BondValence, len(m))
		newM := newMolecule.Bonds[key]
		for i, j := range m {
			newM[i] = j
		}
	}
	for i, atom := range molecule.Atoms {
		if molecule.Head != nil && molecule.Head.Number == atom.Number {
			newMolecule.Head = &newMolecule.Atoms[i]
		} else if molecule.Tail != nil && molecule.Tail.Number == atom.Number {
			newMolecule.Tail = &newMolecule.Atoms[i]
		} else if molecule.RotationPivot != nil && molecule.RotationPivot.Number == atom.Number {
			newMolecule.RotationPivot = &newMolecule.Atoms[i]
		} else {
			continue
		}
	}

	return newMolecule
}

type _Subtitution []*_Monomer

func (substitution *_Subtitution) GetMonomer(current int) *_Monomer {
	return (*substitution)[current%len(*substitution)]
}
