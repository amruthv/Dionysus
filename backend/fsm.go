package main

import "fmt"

type Input int

type State struct {
	Id         int
	Transition func(Input) int
}

type FSM struct {
	CurrentStateId int
	States         map[int]State
}

func NewFSM() *FSM {
	return &FSM{States: make(map[int]State)}
}

func (fsm *FSM) AddState(id int, transition func(Input) int) {
	fsm.States[id] = State{id, transition}
}

func (fsm *FSM) Start(startStateId int) {
	fsm.CurrentStateId = startStateId
}

func (fsm *FSM) Transition(input Input) {
	curr := fsm.States[fsm.CurrentStateId]
	next := curr.Transition(input)
	fmt.Printf("FSM transition (%d) -> (%d) -> (%d)\n", fsm.CurrentStateId, input, next)
	fsm.CurrentStateId = next
}
