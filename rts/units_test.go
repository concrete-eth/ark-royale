package rts

import (
	"image"
	"testing"
)

func testSetTargetPosition(t *testing.T, cmd interface {
	SetTargetPosition(image.Point)
	TargetPosition() image.Point
}) {
	targetPosition := image.Point{X: 123, Y: 456}
	cmd.SetTargetPosition(targetPosition)
	if cmd.TargetPosition() != targetPosition {
		t.Errorf("expected beta %v, got %v", targetPosition, cmd.TargetPosition())
	}
}

func testSetTargetBuilding(t *testing.T, cmd interface {
	SetTargetBuilding(uint8, uint8)
	TargetBuilding() (uint8, uint8)
}) {
	playerId, buildingId := uint8(12), uint8(34)
	cmd.SetTargetBuilding(playerId, buildingId)
	_playerId, _buildingId := cmd.TargetBuilding()
	if _playerId != playerId {
		t.Errorf("expected target player ID %v, got %v", playerId, _playerId)
	}
	if _buildingId != buildingId {
		t.Errorf("expected target building ID %v, got %v", buildingId, _buildingId)
	}
}

func TestWorkerCommandData(t *testing.T) {
	var cmdType WorkerCommandType
	var cmd WorkerCommandData

	// Test setting and getting the command type
	cmdType = WorkerCommandType_Build
	cmd = NewWorkerCommandData(cmdType)
	if cmd.Type() != cmdType {
		t.Errorf("expected command type %v, got %v", cmdType, cmd.Type())
	}

	// Test target building
	testSetTargetBuilding(t, &cmd)
}

func TestFighterCommandData(t *testing.T) {
	var cmdType FighterCommandType
	var cmd FighterCommandData

	// Test setting and getting the command type
	cmdType = FighterCommandType_HoldPosition
	cmd = NewFighterCommandData(cmdType)
	if cmd.Type() != cmdType {
		t.Errorf("expected command type %v, got %v", cmdType, cmd.Type())
	}

	// Test target position
	testSetTargetPosition(t, &cmd)

	cmdType = FighterCommandType_AttackBuilding
	cmd = NewFighterCommandData(cmdType)
	if cmd.Type() != cmdType {
		t.Errorf("expected command type %v, got %v", cmdType, cmd.Type())
	}

	// Test target building
	testSetTargetBuilding(t, &cmd)
}

func TestCommandPathMeta(t *testing.T) {
	meta := commandPathMeta(0)
	meta.SetPathLen(2)

	if meta.PathLen() != 2 {
		t.Errorf("expected path length 2, got %v", meta.PathLen())
	}
	if meta.Pointer() != 0 {
		t.Errorf("expected pointer 0, got %v", meta.Pointer())
	}
	meta.IncPointer()
	if meta.Pointer() != 1 {
		t.Errorf("expected pointer 1, got %v", meta.Pointer())
	}
	meta.IncPointer()
	if meta.Pointer() != 2 {
		t.Errorf("expected pointer 2, got %v", meta.Pointer())
	}
	meta.IncPointer()
	if meta.Pointer() != 2 {
		t.Errorf("expected pointer 2, got %v", meta.Pointer())
	}
}

func TestCommandPath(t *testing.T) {
	path := &CommandPath{}
	path.SetPath([]image.Point{{X: 1, Y: 2}, {X: 3, Y: 4}})
	if path.PathLen() != 2 {
		t.Errorf("expected path length 2, got %v", path.PathLen())
	}
	if len(path.Path()) != 2 {
		t.Errorf("expected path length 2, got %v", len(path.Path()))
	}
	if path.Path()[0] != (image.Point{X: 1, Y: 2}) {
		t.Errorf("expected path point 1,2, got %v", path.Path()[0])
	}
	if path.Path()[1] != (image.Point{X: 3, Y: 4}) {
		t.Errorf("expected path point 3,4, got %v", path.Path()[1])
	}

	if path.Pointer() != 0 {
		t.Errorf("expected pointer 0, got %v", path.Pointer())
	}
	if path.CurrentPoint() != (image.Point{X: 1, Y: 2}) {
		t.Errorf("expected current point 1,2, got %v", path.CurrentPoint())
	}
	path.IncPointer()
	if path.Pointer() != 1 {
		t.Errorf("expected pointer 1, got %v", path.Pointer())
	}
	if path.CurrentPoint() != (image.Point{X: 3, Y: 4}) {
		t.Errorf("expected current point 3,4, got %v", path.CurrentPoint())
	}
	path.IncPointer()
	if path.Pointer() != 2 {
		t.Errorf("expected pointer 2, got %v", path.Pointer())
	}
	if path.CurrentPoint() != (image.Point{}) {
		t.Errorf("expected current point 0,0, got %v", path.CurrentPoint())
	}
}
