package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	STORAGE_FILEPATH      = "./bets.csv"
	LOTTERY_WINNER_NUMBER = 7574
)

type Bet struct {
	agency     int
	first_name string
	last_name  string
	document   string
	birthdate  time.Time
	number     int
}

func NewBet(agency string, firstName string, lastName string, document string, birthdate string, number string) (*Bet, error) {

	ag, err := strconv.Atoi(agency)
	if err != nil {
		return nil, fmt.Errorf("invalid agency: %v", err)
	}

	bd, err := time.Parse(time.DateOnly, birthdate)
	if err != nil {
		return nil, fmt.Errorf("invalid birthdate: %v", err)
	}

	num, err := strconv.Atoi(number)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %v", err)
	}

	return &Bet{
		agency:     ag,
		first_name: firstName,
		last_name:  lastName,
		document:   document,
		birthdate:  bd,
		number:     num,
	}, nil
}

func (b *Bet) HasWon() bool {
	return b.number == LOTTERY_WINNER_NUMBER
}

func StoreBets(bets []*Bet) error {
	file, err := os.OpenFile(STORAGE_FILEPATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, bet := range bets {
		record := []string{
			strconv.Itoa(bet.agency),
			bet.first_name,
			bet.last_name,
			bet.document,
			bet.birthdate.Format("2006-01-02"),
			strconv.Itoa(bet.number),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record to csv: %v", err)
		}
	}

	return nil
}

func LoadBets() ([]*Bet, error) {
	file, err := os.Open(STORAGE_FILEPATH)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var bets []*Bet
	for _, record := range records {
		if len(record) != 6 {
			return nil, fmt.Errorf("invalid record format")
		}
		bet, err := NewBet(record[0], record[1], record[2], record[3], record[4], record[5])
		if err != nil {
			return nil, fmt.Errorf("failed to create bet: %v", err)
		}
		bets = append(bets, bet)
	}

	return bets, nil
}
