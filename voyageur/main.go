package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/MaxHalford/gago"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type City struct {
	X   int
	Y   int
	Nom byte
}

type Country struct {
	Size   int
	Cities []City
}

func (c *City) Distance(c2 *City) float64 {
	d := math.Sqrt(float64((c2.X-c.X)*(c2.X-c.X) + (c2.Y-c.Y)*(c2.Y-c.Y)))
	return d
}

func NewCountry(N int, size int) *Country {
	if N > 26 {
		return nil
	}

	c := make([]City, N)
	Nom := 'A'
	Nom--
	m := map[string]struct{}{}
	var x, y int
	var xy string

	gen := func() City {
		for g := true; g; {
			x = rand.Intn(size)
			y = rand.Intn(size)
			xy = fmt.Sprintf("%d|%d", x, y)
			_, g = m[xy]
		}
		Nom++
		return City{X: x, Y: y, Nom: byte(Nom)}
	}

	for i := 0; i < N; i++ {
		c[i] = gen()
	}

	return &Country{Size: size, Cities: c}
}

func (C *Country) Display() {
	m := make([][]byte, C.Size)
	for i := 0; i < C.Size; i++ {
		r := make([]byte, C.Size)
		for k := 0; k < C.Size; k++ {
			r[k] = '.'
		}
		m[i] = r
	}
	for _, c := range C.Cities {
		m[c.Y][c.X] = c.Nom
	}

	for _, l := range m {
		fmt.Printf("%s\n", string(l))
	}
}

var MyCountry = NewCountry(20, 20)

func main() {
	MyCountry.Display()
	var ga = gago.Generational(VectorFactory)
	ga.NPops = 10
	ga.PopSize = 1000
	ga.Initialize()
	fmt.Printf("Best fitness at generation 0: %f\n", ga.Best.Fitness)
	for i := 1; i < 50; i++ {
		ga.Enhance()
		fmt.Printf("Best fitness at generation %d: %f\n", i, ga.Best.Fitness)
	}

	ga.Best.Genome.(Vector).Print()
}

// A Vector contains float64s.
type Vector []int

// At method from Slice
func (X Vector) Print() {
	for i := 0; i < len(X); i++ {
		fmt.Print(string(MyCountry.Cities[X[i]].Nom))
	}
}

// The bigger the distance the worse the score
func (X Vector) Evaluate() float64 {

	distance := 0.0
	i := 0

	for ; i < len(MyCountry.Cities)-1; i++ {
		distance += MyCountry.Cities[X[i]].Distance(&MyCountry.Cities[X[i+1]])
	}
	distance += MyCountry.Cities[X[i]].Distance(&MyCountry.Cities[X[0]])
	return distance
}

// Mutate a Vector by permutation of cities
func (X Vector) Mutate(rng *rand.Rand) {
	gago.MutPermuteInt(X, 1, rng)
}

// Crossover a Vector with another Vector by applying uniform crossover.
func (X Vector) Crossover(Y gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossPMXInt(X, Y.(Vector), rng)
	return Vector(o1), Vector(o2)
}

// Clone a Vector to produce a new one that points to a different slice.
func (X Vector) Clone() gago.Genome {
	var Y = make(Vector, len(X))
	copy(Y, X)
	return Y
}

// VectorFactory returns a random vector by generating 2 values uniformally
// distributed between -10 and 10.
func VectorFactory(rng *rand.Rand) gago.Genome {
	v := Vector(rand.New(rand.NewSource(time.Now().Unix())).Perm(len(MyCountry.Cities)))
	return v
}
