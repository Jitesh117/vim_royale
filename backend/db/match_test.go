package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"vim_royale/backend/models"
	//"math"
)

// using inmemory sqlite to test the db update

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Match{}, &models.MatchKeystroke{}); err != nil {
		t.Fatalf("failed to migrate test database %v", err)
	}

	return db
}

func TestCreateMatchAndSendRatingDelta(t *testing.T) {
	db := setupTestDB(t)

	winner := models.User{
		Provider:   "github",
		ProviderID: "1",
		Username:   "alice",
		Rating:     1500,
	}

	loser := models.User{
		Provider:   "github",
		ProviderID: "2",
		Username:   "bob",
		Rating:     1600,
	}

	if err := db.Create(&winner).Error; err != nil {
		t.Fatal(err)
	}

	if err := db.Create(&loser).Error; err != nil {
		t.Fatal(err)
	}

	matchId, winnerDelta, loserDelta, winnnerNewRating, loserNewRating, err := CreateMatchAndSendRatingDelta(db, "github:1", "github:2", "target code", "polluted code")

	if err != nil {
		t.Fatal(err)
	}

	if matchId == "" {
		t.Error("match id should not be empty")
	}

	if winnerDelta <= 0 {
		t.Error("winnder delta should not decrease")
	}

	if loserDelta >= 0 {
		t.Error("loser delta should not increase")
	}

	if winnnerNewRating <= 1500 {
		t.Error("Winner new rating should be greater than old rating")
	}

	if loserNewRating >= 1600 {
		t.Error("Loser new rating should be lesser than old rating")
	}

	var updatedWinner models.User
	if err := db.First(&updatedWinner, winner.ID).Error; err != nil {
		t.Fatal(err)
	}

	var updatedLoser models.User
	if err := db.First(&updatedLoser, loser.ID).Error; err != nil {
		t.Fatal(err)
	}

	if updatedWinner.Won != 1 {
		t.Errorf("got %d wins, want 1", updatedWinner.Won)
	}

	if updatedWinner.Matches != 1 {
		t.Errorf("got %d matches, want 1", updatedWinner.Matches)
	}

	if updatedLoser.Lost != 1 {
		t.Errorf("got %d losses, want 1", updatedLoser.Lost)
	}

	if updatedLoser.Matches != 1 {
		t.Errorf("got %d matches, want 1", updatedLoser.Matches)
	}
}

func TestCreateDrawMatch(t *testing.T) {
	db := setupTestDB(t)

	playerA := models.User{
		Provider:   "github",
		ProviderID: "1",
		Username:   "alice",
		Rating:     1500,
	}

	playerB := models.User{
		Provider:   "github",
		ProviderID: "2",
		Username:   "bob",
		Rating:     1600,
	}

	if err := db.Create(&playerA).Error; err != nil {
		t.Fatal(err)
	}

	if err := db.Create(&playerB).Error; err != nil {
		t.Fatal(err)
	}

	matchID, err := CreateDrawMatch(db, "github:1", "github:2", "target code", "polluted code")
	if err != nil {
		t.Fatal(err)
	}

	if matchID == "" {
		t.Fatal("match id should not be empty")
	}

	match, err := GetMatchByMatchID(db, matchID)
	if err != nil {
		t.Fatal(err)
	}

	if match.Outcome != "draw" {
		t.Errorf("got outcome %q, want draw", match.Outcome)
	}

	if match.WinnerID != 0 {
		t.Errorf("got winnerID %d, want 0", match.WinnerID)
	}

	if match.PlayerAID == nil || *match.PlayerAID != playerA.ID {
		t.Errorf("got playerAID %v, want %d", match.PlayerAID, playerA.ID)
	}

	if match.PlayerBID == nil || *match.PlayerBID != playerB.ID {
		t.Errorf("got playerBID %v, want %d", match.PlayerBID, playerB.ID)
	}

	var updatedPlayerA models.User
	if err := db.First(&updatedPlayerA, playerA.ID).Error; err != nil {
		t.Fatal(err)
	}

	var updatedPlayerB models.User
	if err := db.First(&updatedPlayerB, playerB.ID).Error; err != nil {
		t.Fatal(err)
	}

	if updatedPlayerA.Rating != 1500 {
		t.Errorf("got playerA rating %.2f, want 1500", updatedPlayerA.Rating)
	}

	if updatedPlayerB.Rating != 1600 {
		t.Errorf("got playerB rating %.2f, want 1600", updatedPlayerB.Rating)
	}

	if updatedPlayerA.Won != 0 {
		t.Errorf("got playerA wins %d, want 0", updatedPlayerA.Won)
	}

	if updatedPlayerA.Lost != 0 {
		t.Errorf("got playerA losses %d, want 0", updatedPlayerA.Lost)
	}

	if updatedPlayerB.Won != 0 {
		t.Errorf("got playerB wins %d, want 0", updatedPlayerB.Won)
	}

	if updatedPlayerB.Lost != 0 {
		t.Errorf("got playerB losses %d, want 0", updatedPlayerB.Lost)
	}

	if updatedPlayerA.Matches != 1 {
		t.Errorf("got playerA matches %d, want 1", updatedPlayerA.Matches)
	}

	if updatedPlayerB.Matches != 1 {
		t.Errorf("got playerB matches %d, want 1", updatedPlayerB.Matches)
	}
}

func TestApplyMixedMatchRatingDelta(t *testing.T) {
	tests := []struct {
		name           string
		opponentRating float64
		didWin         bool
		expectPositive bool
		expectedWins   int
		expectedLosses int
	}{
		{
			name:           "player wins",
			opponentRating: 1600,
			didWin:         true,
			expectPositive: true,
			expectedWins:   1,
			expectedLosses: 0,
		},
		{
			name:           "player loses",
			opponentRating: 1600,
			didWin:         false,
			expectPositive: false,
			expectedWins:   0,
			expectedLosses: 1,
		},
		{
			name:           "default opponent rating",
			opponentRating: 0,
			didWin:         true,
			expectPositive: true,
			expectedWins:   1,
			expectedLosses: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)

			player := models.User{
				Provider:   "github",
				ProviderID: "1",
				Username:   "alice",
				Rating:     1500,
			}

			if err := db.Create(&player).Error; err != nil {
				t.Fatal(err)
			}

			playerDelta, playerNewRating, err :=
				ApplyMixedMatchRatingDelta(
					db,
					"github:1",
					tc.opponentRating,
					tc.didWin,
				)

			if err != nil {
				t.Fatal(err)
			}

			if tc.expectPositive {
				if playerDelta <= 0 {
					t.Errorf("expected positive delta, got %.2f", playerDelta)
				}

				if playerNewRating <= 1500 {
					t.Errorf(
						"expected rating increase, got %.2f",
						playerNewRating,
					)
				}
			} else {
				if playerDelta >= 0 {
					t.Errorf("expected negative delta, got %.2f", playerDelta)
				}

				if playerNewRating >= 1500 {
					t.Errorf(
						"expected rating decrease, got %.2f",
						playerNewRating,
					)
				}
			}

			var updatedPlayer models.User
			if err := db.First(&updatedPlayer, player.ID).Error; err != nil {
				t.Fatal(err)
			}

			if updatedPlayer.Won != tc.expectedWins {
				t.Errorf(
					"got wins %d, want %d",
					updatedPlayer.Won,
					tc.expectedWins,
				)
			}

			if updatedPlayer.Lost != tc.expectedLosses {
				t.Errorf(
					"got losses %d, want %d",
					updatedPlayer.Lost,
					tc.expectedLosses,
				)
			}

			if updatedPlayer.Matches != 1 {
				t.Errorf(
					"got matches %d, want 1",
					updatedPlayer.Matches,
				)
			}
		})
	}
}
