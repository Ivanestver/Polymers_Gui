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

type _Molecule struct {
	Name  string
	Atoms []_Atom
	Bonds map[_AtomNumber]map[_AtomNumber]_BondValence
	Mass  int
}

func (molecule *_Molecule) GetMassCenter() base.Vector3DF {
	res := molecule.Atoms[0].Coords
	for i := 0; i < len(molecule.Atoms); i++ {
		res.AddF(&molecule.Atoms[i].Coords)
	}
	res.MultiplyByConstantF(1.0 / float64(len(molecule.Atoms)))
	return res
}

func (molecule *_Molecule) MoveTo(point *base.Vector3DF) {
	massCenter := molecule.GetMassCenter()
	direction := base.SubtractVecF(point, &massCenter)
	for i := 0; i < len(molecule.Atoms); i++ {
		currentAtom := &molecule.Atoms[i]
		currentAtom.Coords.AddF(direction)
	}
}

func (molecule *_Molecule) GetBondsCount() int {
	count := 0
	for _, m := range molecule.Bonds {
		count += len(m)
	}
	return count
}

func (molecule *_Molecule) Copy() *_Molecule {
	newMolecule := &_Molecule{
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

	return newMolecule
}
