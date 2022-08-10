package present

import (
	"container/heap"
	"errors"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func (p PresentHeap) Len() int {
	return len(p)
}

func (p PresentHeap) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PresentHeap) Less(i int, j int) bool {
	if p[i].Value > p[j].Value {
		return true
	} else if p[i].Value == p[j].Value && p[i].Size < p[j].Size {
		return true
	}
	return false
}

func (p *PresentHeap) Push(x any) {
	*p = append(*p, x.(Present))
}

func (p *PresentHeap) Pop() any {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}

func GetNCoolestPresents(p []Present, n int) (best []Present, err error) {
	h := (*PresentHeap)(&p)
	heap.Init(h)
	if n > h.Len() || n < 0 {
		return nil, errors.New("number of presents is less then needed")
	}
	for ; n != 0; n-- {
		best = append(best, heap.Pop(h).(Present))
	}
	err = nil
	return
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

type Stat struct {
	SumVal  int
	Indices []int
}

func addStat(s1 Stat, s2 Stat) Stat {
	return Stat{s1.SumVal + s2.SumVal, append(s1.Indices, s2.Indices...)}
}

func generateResult(p []Present, indices []int) []Present {
	result := make([]Present, 0, len(indices))
	for _, v := range indices {
		result = append(result, p[v])
	}
	return result
}

func GrabPresents(p []Present, cap int) []Present {

	// Initializing array for DP solving (Tabulation method - bottom-up)
	K := make([][]Stat, len(p)+1)
	for i := range K {
		K[i] = make([]Stat, cap+1)
	}

	for i := 0; i <= len(p); i++ {
		for w := 0; w <= cap; w++ {
			if i == 0 || w == 0 {
				K[i][w] = Stat{0, []int{}}
			} else if p[i-1].Size <= w {
				prevStat := K[i-1][w-p[i-1].Size]
				if p[i-1].Value+prevStat.SumVal > K[i-1][w].SumVal {
					K[i][w] = addStat(Stat{p[i-1].Value, []int{i - 1}}, prevStat)
				} else {
					K[i][w] = K[i-1][w]
				}
			} else {
				K[i][w] = K[i-1][w]
			}
		}
	}
	return generateResult(p, K[len(p)][cap].Indices)
}
