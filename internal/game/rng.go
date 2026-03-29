package game

type RNG struct {
	State uint64
}

func NewRNG(seed int64) *RNG {
	state := uint64(seed)
	if state == 0 {
		state = 0x9e3779b97f4a7c15
	}
	return &RNG{State: state}
}

func (r *RNG) next() uint64 {
	x := r.State
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	r.State = x
	return x * 2685821657736338717
}

func (r *RNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(r.next() % uint64(n))
}

func (r *RNG) Float64() float64 {
	return float64(r.next()>>11) / (1 << 53)
}

func (r *RNG) Perm(n int) []int {
	perm := make([]int, n)
	for index := range perm {
		perm[index] = index
	}
	for index := n - 1; index > 0; index-- {
		swap := r.Intn(index + 1)
		perm[index], perm[swap] = perm[swap], perm[index]
	}
	return perm
}
