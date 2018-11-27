package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// через флаги реализуются опции: установка пользователем лимита времени на ответы(35 секунд по умолчанию),
	// а также подключение к выбранному файлу
	csvFilename := flag.String("csv", "q-a.csv", "csv-файл в формате вопрос-ответ")
	timeLimit := flag.Int("Лимит", 35, "Сколько секунд осталось на выполнение квиза")

	flag.Parse()
	// чтение и ошибки
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Ошибка при открытии csv-файла: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Ошибка при парсинге запрашиваемого файла")
	}
	questionanswers := parseLines(lines)
	// реализация таймера посредсвом стандартного пакета
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	// для расчета корректного времени завершения используются горутин(ы) и канал
	// т.е. программа не должна ждать ответа на вопрос, чтобы завершиться и вывести результат
	for i, qa := range questionanswers {
		fmt.Printf("Вопрос №%d: %s = \n", i+1, qa.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nВы набрали %d из %d.\n", correct, len(questionanswers))
			return
		case answer := <-answerCh:
			if answer == qa.a {
				correct++
			}
		}
	}

	fmt.Printf("Вы правы на %d из %d.\n", correct, len(questionanswers))
}

// вопросы перемешиваются при каждом запуске
func parseLines(lines [][]string) []questionanswer {
	ret := make([]questionanswer, len(lines))
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i, line := range lines {
		ret[i] = questionanswer{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
		newPosition := r.Intn(len(lines) - 1)
		lines[i], lines[newPosition] = lines[newPosition], lines[i]
	}
	for i, line := range lines {
		ret[i] = questionanswer{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type questionanswer struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
