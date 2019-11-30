package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	ROCK        = 0
	PAPER       = 1
	SCISSORS    = 2
	NUM_ACTIONS = 3
)

var regretSum [3]float64
var strategy [3]float64
var strategySum [3]float64

var oppStrategy = [3]float64{0.8, 0.1, 0.1}

func getStrategy() [3]float64 {
	var normalizingSum float64 = 0
	for a := 0; a < NUM_ACTIONS; a++ {
		// fmt.Println(regretSum, " is the regret sum")
		if regretSum[a] > 0 {
			strategy[a] = regretSum[a]
		} else {
			strategy[a] = 0
		}
		normalizingSum += strategy[a]
	}

	for i := 0; i < NUM_ACTIONS; i++ {
		if normalizingSum > 0 {
			strategy[i] = strategy[i] / normalizingSum
		} else {
			strategy[i] = 1.0 / NUM_ACTIONS
		}
		strategySum[i] += strategy[i]
	}

	return strategy
}

func Train(iterations int) bool {

	var actionUtility [3]float64

	for i := 0; i < iterations; i++ {
		strategy := getStrategy()
		log.Println("this is the strategy", strategy, "for iteration", i)
		var myAction int = GetAction(strategy)
		var opponentAction int = GetAction(oppStrategy)
		// fmt.Println("my action is", myAction, "and the opponents action is", opponentAction)
		// the action utility is 0 if we play the same move as our opponent (tie)
		actionUtility[opponentAction] = 0

		if opponentAction == NUM_ACTIONS-1 { // if opponent plays 2 (scissors)
			actionUtility[0] = 1 // value of playing rock(0) equals 1, because rock beats scissors
		} else {
			actionUtility[opponentAction+1] = 1 // if opponent does not play 2, meaning they played 0 or 1
			// if opponents plays 0(rock), value of playing paper is 1 (paper beats rock)
			// if opponent plays 1(paper), value of players scissors is 1 (scissors beats paper)
		}

		if opponentAction == 0 { // if opponent plays 0, value of playing scissors is -1 (rock beats scissors)
			actionUtility[NUM_ACTIONS-1] = -1
		} else {
			actionUtility[opponentAction-1] = -1 // if opponent does not play 0 , meaning they played 2 or 1
			// if 1(paper), value of 0(rock) is -1, paper beats rock
			// if 2(scissors), value of 1 (paper) is -1, scissors beats paper
		}

		for a := 0; a < NUM_ACTIONS; a++ {
			regretSum[a] += actionUtility[a] - actionUtility[myAction]
		}

	}

	return true
}

func GetAverageStrategy() (avgStrategy [3]float64) {

	var normalizingSum float64 = 0
	for a := 0; a < NUM_ACTIONS; a++ {
		normalizingSum += strategySum[a]
	}
	for i := 0; i < NUM_ACTIONS; i++ {
		if normalizingSum > 0 {
			avgStrategy[i] = strategySum[i] / normalizingSum
		} else {
			avgStrategy[i] = 1.0 / NUM_ACTIONS
		}
	}
	return avgStrategy
}

func main() {
	var err = os.Remove("test_log.txt")
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.OpenFile("test_log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer file.Close()

	//set output of logs to f
	log.SetOutput(file)

	f := Train(1000)
	log.Println(f)
	avgStrat := GetAverageStrategy()

	log.Println("The average strategy is", avgStrat)
}

// A strategy is just a distribution of % Rock, % Paper, % Scissors
// GetAction takes in a strategy and returns a random action based on probability distribution.
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
