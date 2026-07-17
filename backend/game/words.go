package game

import (
	"bufio"
	"math/rand/v2"
	"os"
	"time"
)

func LoadWords(f string) ([]string, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func GetDailyWord(words []string) string {
	today := time.Now()
	seed := uint64(today.Year()*10000 + int(today.Month())*100 + today.Day())
	generator := rand.NewPCG(seed, 0)
	random := rand.New(generator)
	index := random.IntN(len(words))
	dailyWord := words[index]

	return dailyWord
}

func GetRandomWord(words []string) string {
	seed := uint64(time.Now().UnixNano())
	generator := rand.NewPCG(seed, 0)
	random := rand.New(generator)
	index := random.IntN(len(words))
	randomWord := words[index]

	return randomWord
}
