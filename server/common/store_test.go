package common

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBetMustKeepFields(t *testing.T) {
	bet, err := NewBet("1", "first", "last", "10000000", "2000-12-20", "7500")
	assert.Nil(t, err)

	assert.Equal(t, bet.agency, 1)
	assert.Equal(t, bet.first_name, "first")
	assert.Equal(t, bet.last_name, "last")
	assert.Equal(t, bet.document, "10000000")
	assert.Equal(t, bet.birthdate, time.Date(2000, 12, 20, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, bet.number, 7500)
}

func TestHasWonWithWinnerNumberMustBeTrue(t *testing.T) {
	bet, err := NewBet("1", "first", "last", "10000000", "2000-12-20", strconv.Itoa(LOTTERY_WINNER_NUMBER))
	assert.Nil(t, err)

	assert.True(t, bet.HasWon())
}

func TestHasWonWithLoserNumberMustBeFalse(t *testing.T) {
	bet, err := NewBet("1", "first", "last", "10000000", "2000-12-20", strconv.Itoa(LOTTERY_WINNER_NUMBER+1))
	assert.Nil(t, err)

	assert.False(t, bet.HasWon())
}

func TestStoreBetsAndLoadBetsKeepsFieldsData(t *testing.T) {
	toStore := []*Bet{
		{
			agency:     1,
			first_name: "first",
			last_name:  "last",
			document:   "10000000",
			birthdate:  time.Date(2000, 12, 20, 0, 0, 0, 0, time.UTC),
			number:     7500,
		},
	}

	StoreBets(toStore)
	storedBets, err := LoadBets()
	assert.Nil(t, err)

	assert.Equal(t, toStore, storedBets)
}

// test_store_bets_and_load_bets_keeps_registry_order
func TestStoreBetsAndLoadBetsKeepsRegistryOrder(t *testing.T) {
	toStore := []*Bet{
		{
			agency:     1,
			first_name: "first_0",
			last_name:  "last_0",
			document:   "10000000",
			birthdate:  time.Date(2000, 12, 20, 0, 0, 0, 0, time.UTC),
			number:     7500,
		},
		{
			agency:     2,
			first_name: "first_1",
			last_name:  "last_1",
			document:   "10000001",
			birthdate:  time.Date(2000, 12, 21, 0, 0, 0, 0, time.UTC),
			number:     7501,
		},
	}

	StoreBets(toStore)
	storedBets, err := LoadBets()
	assert.Nil(t, err)

	assert.Equal(t, toStore[0], storedBets[0])
	assert.Equal(t, toStore[1], storedBets[1])
}

func TestMain(m *testing.M) {
	_ = os.Remove(STORAGE_FILEPATH)
	code := m.Run()
	_ = os.Remove(STORAGE_FILEPATH)
	os.Exit(code)
}
