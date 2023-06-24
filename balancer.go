package main

type LoadBalancer struct {
	data    []string
	weights []int16
}

func (queue *LoadBalancer) Get(weight int16) string {
	var m int16 = -1
	var mi int
	for index, value := range queue.weights {
		if value < m {
			m = value
			mi = index
		}
	}

	queue.weights[mi] = weight
	return queue.data[mi]
}
