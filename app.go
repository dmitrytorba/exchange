package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

type command struct {
	from   string
	to     string
	amount int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	exchange := createExchange()
	errorMsg := ""
	history := make([]string, 0, 10)

	for {
		clear()

		// display an error message if one is available
		if len(errorMsg) != 0 {
			fmt.Println("---", errorMsg, "---\n")
			errorMsg = ""
		}

		// display exchange state
		exchange.printout()

		// print out the history
		fmt.Println("\n  history")
		fmt.Println("----------------------------")
		for i := 0; i < len(history); i++ {
			fmt.Println(history[i])
		}

		// interpret a command
		fmt.Println("\nEnter Command (ex. BUY 100 USD WITH BTC):")
		text, _ := reader.ReadString('\n')
		text = strings.ToUpper(text)
		fields := strings.Fields(text)
		if len(fields) != 5 {
			errorMsg = "invalid command format"
			continue
		}

		command := &command{
			from: fields[4],
			to:   fields[2],
		}

		// check the second argument is an int
		if amount, err := strconv.Atoi(fields[1]); err == nil {
			command.amount = amount
		} else {
			errorMsg = "second argument is not an integer"
			continue
		}

		// check for the "BUY"
		if fields[0] != "BUY" {
			errorMsg = "first argument is not 'BUY'"
			continue
		}

		// check for the "WITH"
		if fields[3] != "WITH" {
			errorMsg = "fourth argument is not 'WITH'"
			continue
		}

		if _, ok := exchange.reserves[currency(command.from)]; !ok {
			errorMsg = "buying currency is not a valid currency"
			continue
		}

		if _, ok := exchange.reserves[currency(command.to)]; !ok {
			errorMsg = "the selling currency is not a valid currency"
			continue
		}

		exchange.execute(command.amount, currency(command.to), currency(command.from))
		history = append(history, fmt.Sprintf("bought %v %v with %v", command.amount, command.to, command.from))
	}
}
