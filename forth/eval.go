//go:build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	commands map[string][]func(res *[]int) error
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	var commands = map[string][]func(res *[]int) error{
		"+":    {plus},
		"-":    {minus},
		"*":    {multi},
		"/":    {div},
		"dup":  {dup},
		"over": {over},
		"drop": {drop},
		"swap": {swap},
	}
	return &Evaluator{commands: commands}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {

	var s stack = strings.Split(row, " ")
	var res []int
	for i := 0; i < len(s); i++ {
		if s[i] == ":" {
			nameCommand := strings.ToLower(s[1])
			combinatedCommands := s[2 : len(s)-1]

			if _, err := strconv.Atoi(s[1]); err == nil {
				if _, err := strconv.Atoi(s[2]); err == nil {
					return nil, fmt.Errorf("err")
				}
			}
			e.commands[nameCommand] = doCombinatedCommands(combinatedCommands, e.commands)
			break
		} else {
			command := strings.ToLower(s[i])
			err := doCommands(command, &res, e.commands)
			if err != nil {
				return nil, err
			}
		}
	}

	return res, nil
}

func doCombinatedCommands(combinatedCommands stack, newCommands map[string][]func(res *[]int) error) []func(res *[]int) error {
	resFuncs := []func(res *[]int) error{}
	for _, com := range combinatedCommands {
		com = strings.ToLower(com)
		if _, ok := newCommands[com]; !ok {
			resFuncs = append(resFuncs, addNumber(com))
		}
		resFuncs = append(resFuncs, newCommands[com]...)
	}
	return resFuncs
}

func addNumber(com string) func(*[]int) error {
	return func(res *[]int) error {
		iCom, err := strconv.Atoi(com)
		if err != nil {
			return fmt.Errorf("err")
		}
		(*res) = append((*res), iCom)
		return nil
	}
}

func doCommands(command string, res *[]int, commands map[string][]func(res *[]int) error) error {
	if _, ok := commands[command]; !ok {
		if v, err := strconv.Atoi(command); err == nil {
			(*res) = append((*res), v)
		} else {
			return fmt.Errorf("err")
		}
	}
	for _, c := range commands[command] {
		err := c(res)
		if err != nil {
			return err
		}
	}
	return nil
}

type stack []string

func dup(res *[]int) error {
	if len(*res) < 1 {
		return fmt.Errorf("err")
	}
	*res = append(*res, (*res)[len(*res)-1])
	return nil
}

func over(res *[]int) error {
	if len(*res) < 2 {
		return fmt.Errorf("err")
	}
	*res = append(*res, (*res)[len(*res)-2])
	return nil
}

func drop(res *[]int) error {
	if len(*res) < 1 {
		return fmt.Errorf("err")
	}
	(*res) = (*res)[:len(*res)-1]
	return nil
}

func swap(res *[]int) error {
	if len(*res) < 2 {
		return fmt.Errorf("err")
	}
	(*res)[len(*res)-1], (*res)[len(*res)-2] = (*res)[len(*res)-2], (*res)[len(*res)-1]
	return nil
}

func convertToIntSlice(slice []string) []int {
	iSlice := make([]int, 0, len(slice))
	for _, v := range slice {
		if v != "lol" {
			iV, _ := strconv.Atoi(v)
			iSlice = append(iSlice, iV)
		}
	}
	return iSlice
}

func plus(res *[]int) error {
	if len(*res) < 2 {
		return fmt.Errorf("small len")
	}
	a, b := (*res)[len(*res)-2], (*res)[len(*res)-1]
	(*res) = (*res)[:len(*res)-2]
	(*res) = append((*res), a+b)
	return nil
}

func minus(res *[]int) error {
	if len(*res) < 2 {
		return fmt.Errorf("small len")
	}
	a, b := (*res)[len(*res)-2], (*res)[len(*res)-1]
	(*res) = (*res)[:len(*res)-2]
	(*res) = append((*res), a-b)
	return nil
}

func multi(res *[]int) error {
	if len(*res) < 2 {
		return fmt.Errorf("small len")
	}
	a, b := (*res)[len(*res)-2], (*res)[len(*res)-1]
	(*res) = (*res)[:len(*res)-2]
	(*res) = append((*res), a*b)
	return nil
}

func div(res *[]int) error {
	if len(*res) < 2 {
		return fmt.Errorf("small len")
	}
	a, b := (*res)[len(*res)-2], (*res)[len(*res)-1]
	if b == 0 {
		return fmt.Errorf("err")
	}
	(*res) = (*res)[:len(*res)-2]
	(*res) = append((*res), a/b)
	return nil
}
