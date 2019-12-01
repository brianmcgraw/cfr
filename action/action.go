package action

import (
	"math/rand"
	"time"
)

// Action = 0 means we play rock, 1 for paper, 2 for scissors
// A strategy is just a distribution of % Rock, % Paper, % Scissors
// GetAction takes in a strategy and returns a random action based on the provided probability distribution.
func GetAction(strategy [3]float64) int {
	rand.Seed(time.Now().UnixNano())
	action := 0
	r := rand.Float64()

	var cumulativeProbability float64 = 0

	for action < 2 {
		cumulativeProbability += strategy[action]

		if r < cumulativeProbability {
			break
		}
		action++
	}

	return action
}
