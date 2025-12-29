package core

import (
	"container/heap"
	"context"
)

// PriorityItem represents an item in the priority queue
type PriorityItem struct {
	Password string
	Priority float64 // Higher priority = more likely password
	Index    int
}

// PriorityQueue implements a priority queue for passwords
type PriorityQueue []*PriorityItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority // Higher priority first
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PriorityItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// PriorityProcessor processes passwords by priority
type PriorityProcessor struct {
	queue *PriorityQueue
}

// NewPriorityProcessor creates a new priority processor
func NewPriorityProcessor() *PriorityProcessor {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &PriorityProcessor{queue: &pq}
}

// CalculatePriority calculates priority for a password
func CalculatePriority(password string) float64 {
	strength := CalculateStrength(password)
	analysis := AnalyzePassword(password)

	priority := float64(strength.Score)

	// Boost priority for realistic passwords
	if !analysis.CommonPattern && !analysis.KeyboardWalk && !analysis.Repeating && !analysis.Sequential {
		priority *= 1.2
	}

	// Boost for good entropy
	if analysis.Entropy > 40 {
		priority *= 1.1
	}

	// Penalty for very weak passwords
	if strength.Score < 20 {
		priority *= 0.5
	}

	return priority
}

// ProcessByPriority processes passwords and outputs by priority
func (p *PriorityProcessor) ProcessByPriority(ctx context.Context, passwords <-chan string, maxOutput int) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)

		// Collect all passwords with priorities
		for password := range passwords {
			select {
			case <-ctx.Done():
				return
			default:
				priority := CalculatePriority(password)
				item := &PriorityItem{
					Password: password,
					Priority: priority,
				}
				heap.Push(p.queue, item)
			}
		}

		// Output by priority
		count := 0
		for p.queue.Len() > 0 && (maxOutput <= 0 || count < maxOutput) {
			select {
			case <-ctx.Done():
				return
			default:
				item := heap.Pop(p.queue).(*PriorityItem)
				output <- item.Password
				count++
			}
		}
	}()

	return output
}

