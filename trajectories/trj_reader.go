package trajectories

import (
	"bufio"
	"strconv"
	"strings"
)

type _TrjReader struct {
	reader *bufio.Scanner
}

type _TrjItem struct {
	_TrjReader
	AtomID   int
	AtomType int
	S        [3]float64
	K        [3]int
}

func trjItemFromString(s string) *_TrjItem {
	parts := strings.Split(s, " ")
	atomID, _ := strconv.Atoi(parts[0])
	atomType, _ := strconv.Atoi(parts[1])
	xs, _ := strconv.ParseFloat(parts[2], 64)
	ys, _ := strconv.ParseFloat(parts[3], 64)
	zs, _ := strconv.ParseFloat(parts[4], 64)
	kx, _ := strconv.Atoi(parts[5])
	ky, _ := strconv.Atoi(parts[6])
	kz, _ := strconv.Atoi(parts[7])
	return &_TrjItem{
		AtomID:   atomID,
		AtomType: atomType,
		S:        [3]float64{xs, ys, zs},
		K:        [3]int{kx, ky, kz},
	}
}

type _TrjTimestep struct {
	_TrjReader
	timestep      int
	numberOfAtoms int
	boxBounds     [3][3]float64
}

func makeTrjTimestep(reader *bufio.Scanner) *_TrjTimestep {
	return &_TrjTimestep{
		_TrjReader: _TrjReader{
			reader: reader,
		},
		timestep:      -1,
		numberOfAtoms: 0,
	}
}

func (timestep *_TrjTimestep) GetTimestep() int {
	if timestep.timestep == -1 {
		for timestep.reader.Scan() {
			line := timestep.reader.Text()
			if strings.Contains(line, "TIMESTEP") {
				timestep.reader.Scan()
				line = timestep.reader.Text()
				line = strings.Trim(line, " ")
				if t, err := strconv.Atoi(line); err == nil {

					timestep.timestep = t
				} else {
					timestep.timestep = 0
				}
			}
		}
	}
	return timestep.timestep
}

type _BoxBounds [3][2]float64
