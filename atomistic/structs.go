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
	Name  string
	Atoms []_Atom
	Bonds map[_AtomNumber]map[_AtomNumber]_BondValence
	Mass  int
	Head  *_Atom
	Tail  *_Atom
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
	for i := 0; i < len(molecule.Atoms); i++ {
		res.AddF(&molecule.Atoms[i].Coords)
	}
	res.MultiplyByConstantF(1.0 / float64(len(molecule.Atoms)))
	return res
}

func (molecule *_Monomer) MoveTo(point *base.Vector3DF) {
	massCenter := molecule.GetMassCenter()
	direction := base.SubtractVecF(point, &massCenter)
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
		switch atom.Number {
		case molecule.Head.Number:
			newMolecule.Head = &newMolecule.Atoms[i]
		case molecule.Tail.Number:
			newMolecule.Tail = &newMolecule.Atoms[i]
		}
	}

	return newMolecule
}

type _Subtitution struct {
	left  *_Monomer
	right *_Monomer
}

func _NewSubstitution(left, right *_Monomer) *_Subtitution {
	return &_Subtitution{
		left:  left,
		right: right,
	}
}

func (substitution *_Subtitution) GetLeft() *_Monomer {
	return substitution.left
}

func (substitution *_Subtitution) GetRight() *_Monomer {
	if substitution.right != nil {
		return substitution.right
	} else {
		return substitution.left
	}
}
