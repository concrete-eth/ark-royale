package rts

type BuildingState uint8

const (
	BuildingState_Nil BuildingState = iota
	BuildingState_Unpaid
	BuildingState_Building
	BuildingState_Built
	BuildingState_Destroyed
	BuildingState_Count
)

func (s BuildingState) Uint8() uint8 {
	return uint8(s)
}

func (s BuildingState) IsNil() bool {
	return s == BuildingState_Nil
}

func (s BuildingState) IsActive() bool {
	return s == BuildingState_Built
}

func (s BuildingState) HasBeenBuilt() bool {
	return s == BuildingState_Built || s == BuildingState_Destroyed
}
