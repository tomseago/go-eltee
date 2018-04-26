package eltee

import (
	"fmt"
)

type StateMapper struct {
	insts []ProfileControlInstance

	mapFn func(ws *WorldState, list []ProfileControlInstance) error
}

func NewStateMapper() *StateMapper {
	sm := &StateMapper{
		insts: make([]ProfileControlInstance, 0),
	}

	return sm
}

func (sm *StateMapper) AppendInstance(inst ProfileControlInstance) {
	sm.insts = append(sm.insts, inst)
}

func (sm *StateMapper) MapState(ws *WorldState) error {
	if sm.mapFn == nil {
		return fmt.Errorf("No mapFn defined")
	}

	sm.mapFn(ws, sm.insts)
}

type SMColor struct {
	colorName string
}

func NewColorMapper(colorName string) StateMapper {

	sm := &SMColor{
		colorName: colorName,
	}

}

func (sm *SMColor) MapStateTo(ws *WorldState, list []ProfileControlInstance) error {
	if ws == nil || list == nil {
		return nil
	}

	color := ws.GetColor(sm.colorName)

	for i := 0; i < len(list); i++ {
		settable, ok := list[i].(WorldColorSettable)
		if !ok {
			log.Warningf("Can not map color %v", sm.colorName)
			continue
		}

		settable.SetFromWorldColor(color)
	}

	return nil
}
