package actions

import (
	"context"
	"log/slog"
)

type Action struct {
	Name string
	Do   func(ctx context.Context) error
}

type State string

const (
	StatePending State = ""
	StateRunning State = "RUNNING"
	StateDone    State = "DONE"
	StateError   State = "ERROR"
)

type ActionNode struct {
	Action Action
	State  State
	Next   []*ActionNode
}

func (n *ActionNode) GetErrors() int {
	errors := 0
	if n.State == StateError {
		errors++
	}

	for _, child := range n.Next {
		errors += child.GetErrors()
	}

	return errors
}

type ActionMap struct {
	roots []*ActionNode
}

func NewActionMap() *ActionMap {
	return &ActionMap{
		roots: []*ActionNode{},
	}
}

func (am *ActionMap) Append(node ...*ActionNode) {
	am.roots = append(am.roots, node...)
}

func (am *ActionMap) Run(ctx context.Context, maxConcurrent int) int {
	torun := make(chan *ActionNode)
	finished := make(chan *ActionNode)

	defer close(torun)
	defer close(finished)

	// Worker pool to run actions
	go func() {
		for n := range torun {
			go func(n *ActionNode) {
				slog.Info("Starting action", "action", n.Action.Name)
				err := n.Action.Do(ctx)
				if err != nil {
					slog.Error("Action failed", "action", n.Action.Name, "error", err)
					n.State = StateError
				} else {
					slog.Info("Action completed", "action", n.Action.Name)
					n.State = StateDone
				}

				finished <- n
			}(n)
		}
	}()

	// Main loop to manage running actions
	running := 0
	for {
		eligible := getEligibleActions(am.roots)
		if len(eligible) == 0 && running == 0 {
			break
		}

		for running < maxConcurrent && len(eligible) > 0 {
			action := eligible[0]
			action.State = StateRunning
			torun <- action

			eligible = eligible[1:]
			running++
		}

		select {
		case <-ctx.Done():
			return 0
		case <-finished:
			running--
		}
	}

	// Check for errors
	errors := 0
	for _, root := range am.roots {
		errors += root.GetErrors()
	}

	return errors
}

func getEligibleActions(nodes []*ActionNode) []*ActionNode {
	var eligible []*ActionNode
	for _, node := range nodes {
		if node.State == StatePending {
			eligible = append(eligible, node)
		} else {
			eligible = append(eligible, getEligibleActions(node.Next)...)
		}
	}
	return eligible
}
