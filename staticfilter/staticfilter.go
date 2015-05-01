/*

 This file implements the foundation for our CS51 final project.

 Joseph Kahn    josephkahn@college.harvard.edu
 Grace Lin      glin@college.harvard.edu
 Aron Szanto    aszanto@college.harvard.edu

*/

package StaticFilter

// Including library packages referenced in this file
import (
	"fmt"
	"github.com/willf/bitset"
	"hash"
	"hash/fnv"
	"math"
)

// type definition for standard bloom filter
type FilterBase struct {
	m uint        // size of bitset
	k uint        // number of hash functions
	h hash.Hash64 // hashing generator
	e float64
}

type Filter struct {
	params  *FilterBase    // needed for generation
	b       *[]bitset.BitSet // pointer to bitset slice
	Counter uint           // counts elements
}

/*
 calculates the length of the bitset and the number of
 required hash functions given the size of set being
 stored and the acceptable error bound for the task at hand
*/
// clean this up, make it one statement
func NewFilterBase(num uint, eps float64) *FilterBase {
	fb := new(FilterBase)
	fb.e = eps
	// calculating length
	fb.m = uint(math.Ceil(-1 * (float64(num) * math.Log(eps)) / (math.Log(2) * math.Log(2))))
	// calculate num hash functions
	fb.k = uint(math.Ceil(math.Log(eps)/math.Log(2)))
	// Pretty sure you can just do this
	// fb:= FilterBase{*insert equation for m here, inesrt equation for k here, h, eps}
	return fb
}

func (fb *FilterBase) CalcBits(d []byte) []uint32 {
	fb.h = fnv.New64a()
	fb.h.Reset()
	fb.h.Write(d)
	hash := fb.h.Sum64()
	h1 := uint32(hash & ((1 << 32) - 1))
	h2 := uint32(hash >> 32)
	//o := fmt.Sprint("h1 = ", h1, " and h2 = ", h2, "\n")
	//fmt.Printf(o)

	indices := make([]uint32, fb.k)
	for i := uint32(0); i <= uint32(fb.k); i++ {
		indices[i] = (h1 + h2*i) % uint32(fb.m) //changed this line to i-1 and the above to <= and 1
	}
	//op := fmt.Sprint(indices, " : bits set to be flipped\n")
	//fmt.Printf(op)
	return indices
}

func NewFilter(num uint, eps float64) *Filter {
	filter := new(Filter)
	filter.params = NewFilterBase(num, eps)
	fmt.Printf(fmt.Sprint("m = ", filter.params.m, "\n"))
	filter.b = bitset.New(filter.params.m)
	return filter
}

func (f *Filter) M() uint {
	return f.params.m
}

func (f *Filter) K() uint {
	return f.params.k
}

func (f *Filter) E() float64 {
	return f.params.e
}

// Takes in a slice of indexes
func (filter *Filter) Insert(data []byte) {
	//p := fmt.Sprint("\nOperating onfilter with k = ", filter.params.k, " and m = ", filter.params.m, "\n\n\n")
	//fmt.Printf(p)
	indices := filter.params.CalcBits(data)
	for i := uint(0); i < uint(len(indices)); i++ {
		filter.b = filter.b.Set(uint(indices[i]))
	}
	filter.Counter++
}

func (filter *Filter) Lookup(data []byte) bool {
	indices := filter.params.CalcBits(data)
	// might be there unless definitely not in set
	for i := 0; i < len(indices); i++ {
		if !filter.b.Test(uint(indices[i])) {
			// definitely not in set
			//op := fmt.Sprint("Bit #", i, " = ", indices[i], " would have been \n")
			//fmt.Printf(op)
			return false
		}
	}
	return true
}
