package main

import (
	"cfr/action"
	"log"
)

// Interpretation of java code from http://modelai.gettysburg.edu/2013/cfr/RPSTrainer.java

// You can manipulate the percentages here to demonstrate the AI's average strategy does change
// e.g., with [R, P, S] = [.34, .33, .33], we should expect the average strategy in response to be the same.
// However, if the opponent plays a strategy of [.8, .1, .1], the average strategy in response
// should play paper a majority of the time.
var oppStrategy = [3]float64{0.8, 0.1, 0.1}

const (
	ROCK        = 0
	PAPER       = 1
	SCISSORS    = 2
	NUM_ACTIONS = 3
)

type RPS struct {
	RegretSum   [3]float64
	Strategy    [3]float64
	StrategySum [3]float64
}

func NewRPS() (rps RPS) {
	rps = RPS{
		RegretSum:   [3]float64{},
		Strategy:    [3]float64{},
		StrategySum: [3]float64{},
	}
	return
}

func (rps *RPS) GetStrategy() {
	var normalizingSum float64 = 0
	for a := 0; a < NUM_ACTIONS; a++ {

		// If we have positive regret from the previous round, that should be incorporated into current strategy
		if rps.RegretSum[a] > 0 {
			rps.Strategy[a] = rps.RegretSum[a]
		} else {
			rps.Strategy[a] = 0
		}
		normalizingSum += rps.Strategy[a]

	}

	for i := 0; i < NUM_ACTIONS; i++ {
		if normalizingSum > 0 {
			rps.Strategy[i] = rps.Strategy[i] / normalizingSum
		} else {
			rps.Strategy[i] = 1.0 / NUM_ACTIONS
		}
		rps.StrategySum[i] += rps.Strategy[i]
	}

}

func (rps *RPS) Train(iterations int) {

	var actionUtility [3]float64

	for i := 0; i < iterations; i++ {
		rps.GetStrategy()

		var myAction int = action.GetAction(rps.Strategy)
		var opponentAction int = action.GetAction(oppStrategy)

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
			rps.RegretSum[a] += actionUtility[a] - actionUtility[myAction]
		}

	}

}

func (rps *RPS) GetAverageStrategy() (avgStrategy [3]float64) {

	var normalizingSum float64 = 0
	for a := 0; a < NUM_ACTIONS; a++ {
		normalizingSum += rps.StrategySum[a]
	}
	for i := 0; i < NUM_ACTIONS; i++ {
		if normalizingSum > 0 {
			avgStrategy[i] = rps.StrategySum[i] / normalizingSum
		} else {
			avgStrategy[i] = 1.0 / NUM_ACTIONS
		}
	}
	return avgStrategy
}

func main() {
	strategy := NewRPS()

	strategy.Train(10000)

	avgStrat := strategy.GetAverageStrategy()

	log.Println("The average strategy is", avgStrat)
}
