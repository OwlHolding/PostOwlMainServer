package main

type LoadBalancer struct {
	data    []string
	weights []int16
}

func (queue *LoadBalancer) Get(weight int16) string {
	var m int16 = 32767
	var mi int
	for index, value := range queue.weights {
		if value < m {
			m = value
			mi = index
		} else {
			if queue.weights[index]-1 >= 0 {
				queue.weights[index] -= 1
			}
		}
	}

	queue.weights[mi] = weight
	return queue.data[mi]
}
