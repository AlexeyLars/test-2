package main

import "fmt"

type SortingStrategy interface {
	Sort(array []int)
}

type BubbleSortStrategy struct{}

func (b *BubbleSortStrategy) Sort(array []int) {
	fmt.Println("Sorting using Bubble Sort")
	n := len(array)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if array[j] > array[j+1] {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
}

type QuickSortStrategy struct{}

func (q *QuickSortStrategy) Sort(array []int) {
	fmt.Println("Sorting using Quick Sort")
	q.quickSort(array, 0, len(array)-1)
}

func (q *QuickSortStrategy) quickSort(array []int, low, high int) {
	if low < high {
		pi := q.partition(array, low, high)
		q.quickSort(array, low, pi-1)
		q.quickSort(array, pi+1, high)
	}
}

func (q *QuickSortStrategy) partition(array []int, low, high int) int {
	pivot := array[high]
	i := low - 1
	for j := low; j < high; j++ {
		if array[j] < pivot {
			i++
			array[i], array[j] = array[j], array[i]
		}
	}
	array[i+1], array[high] = array[high], array[i+1]
	return i + 1
}

// Sorter - контекст
type Sorter struct {
	strategy SortingStrategy
}

func (s *Sorter) SetStrategy(strategy SortingStrategy) {
	s.strategy = strategy
}

func (s *Sorter) SortArray(array []int) {
	s.strategy.Sort(array)
}

func main() {
	fmt.Println("=== Strategy Pattern ===")
	sorter := &Sorter{}

	// Использование пузырьковой сортировки
	sorter.SetStrategy(&BubbleSortStrategy{})
	array1 := []int{5, 3, 8, 4, 2}
	sorter.SortArray(array1)
	fmt.Println("Result:", array1)

	// Использование быстрой сортировки
	sorter.SetStrategy(&QuickSortStrategy{})
	array2 := []int{5, 3, 8, 4, 2}
	sorter.SortArray(array2)
	fmt.Println("Result:", array2)
	fmt.Println()
}
