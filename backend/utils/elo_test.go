package utils

import (
	"testing"
	"math"
)

func TestCalculateElo(t *testing.T) {
	tests := []struct {
		playerRating   float64
		opponentRating float64
		playerWins     bool
		expectedPlayer float64
		expectedOpp    float64
	}{
		{
			1500,
			1600,
			true,
			1520.4820799936924,
			1579.5179200063076,
		},
		{
			1500,
			1600,
			false,
			1488.4820799936924,
			1611.5179200063076,
		},
	}

	const epsilon = 1e-9

	for _, testCase := range tests {
		newPlayerRating, newOpponentRating := CalculateElo(
			testCase.playerRating,
			testCase.opponentRating,
			testCase.playerWins,
		)

		if math.Abs(newPlayerRating-testCase.expectedPlayer) > epsilon {
			t.Errorf(
				"player rating mismatch: got %.12f, want %.12f",
				newPlayerRating,
				testCase.expectedPlayer,
			)
		}

		if math.Abs(newOpponentRating-testCase.expectedOpp) > epsilon {
			t.Errorf(
				"opponent rating mismatch: got %.12f, want %.12f",
				newOpponentRating,
				testCase.expectedOpp,
			)
		}
	}
}

func TestCalculateElo_ConservesRating(t *testing.T) {
    newPlayerRating, newOpponentRating := CalculateElo(1500, 1600, true)

    const eps = 1e-9
    if math.Abs((newPlayerRating+newOpponentRating)-(1500+1600)) > eps {
        t.Fatal("rating should be conserved")
    }
}

func TestCalculateElo_WinnerGainsRating(t *testing.T) {
    newPlayerRating, newOpponentRating := CalculateElo(1500, 1600, true)

    if newPlayerRating <= 1500 {
        t.Fatal("winner should gain rating")
    }

    if newOpponentRating >= 1600 {
        t.Fatal("loser should lose rating")
    }
}

func TestCalculateElo_LoserLosesRating(t *testing.T) {
    newPlayerRating, newOpponentRating := CalculateElo(1500, 1600, false)

    if newPlayerRating >= 1500 {
        t.Fatal("loser should lose rating")
    }

    if newOpponentRating <= 1600 {
        t.Fatal("winner should gain rating")
    }
}



