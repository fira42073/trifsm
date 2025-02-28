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

	open --> closed: Door closed now
	closed --> open: Door opened now
*/
type Door struct {
	// FSM is the finite state machine for the door.
	FSM *fsm.FSM
	// To is the destination of the door.
	To string
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = NewDoorFSM(fsm.Callbacks{
		DoorFSMOpen__DoorFSMClosed: func(_ context.Context, e *fsm.Event) { d.enterState(e) },
		DoorFSMClosed__DoorFSMOpen: func(_ context.Context, e *fsm.Event) { d.enterState(e) },
	})

	return d
}

func (d *Door) enterState(e *fsm.Event) {
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
}

func TestDoor(t *testing.T) {
	t.Run("NewDoor", func(t *testing.T) {
		door := NewDoor("heaven")

		err := door.FSM.Event(context.Background(), DoorFSMOpen__DoorFSMClosed)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		err = door.FSM.Event(context.Background(), DoorFSMClosed__DoorFSMOpen)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		mmd, err := fsm.VisualizeForMermaidWithGraphType(door.FSM, fsm.StateDiagram)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		fmt.Println(mmd)
	})
}
