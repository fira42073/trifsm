//go:generate go run github.com/fira42073/trifsm@senpai -f=$GOFILE
package generator

import (
	"context"
	"fmt"
	"testing"

	"github.com/looplab/fsm"
)

/*
FSM
---
stateDiagram-v2

	green --> yellow: Proceed with caution
	yellow --> red: Stop
	red --> green: gogogo
*/
type TrafficLight struct {
	// FSM is the finite state machine for the door.
	FSM *fsm.FSM
}

func NewTrafficLight(to string) *TrafficLight {
	tl := &TrafficLight{}

	tl.FSM = NewTrafficLightFSM(nil)

	return tl
}

func TestTrafficLight(t *testing.T) {
	t.Run("NewTrafficLight", func(t *testing.T) {
		tl := NewTrafficLight("heaven")

		err := tl.FSM.Event(context.Background(), TrafficLightFSMGreen__TrafficLightFSMYellow)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		err = tl.FSM.Event(context.Background(), TrafficLightFSMYellow__TrafficLightFSMRed)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		err = tl.FSM.Event(context.Background(), TrafficLightFSMRed__TrafficLightFSMGreen)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		err = tl.FSM.Event(context.Background(), TrafficLightFSMYellow__TrafficLightFSMRed)
		if err == nil {
			t.Errorf("no error when expected one: %v", err)
			t.Fail()
		}

		mmd, err := fsm.VisualizeForMermaidWithGraphType(tl.FSM, fsm.StateDiagram)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		fmt.Println(mmd)
	})
}
