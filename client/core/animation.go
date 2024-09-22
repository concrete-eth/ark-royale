package core

import (
	"time"

	"github.com/concrete-eth/ark-royale/client/assets"
	"github.com/concrete-eth/ark-royale/client/decren"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Animation interface {
	Update(c *CoreRenderer, frame uint)
	NFrames() uint
}

type AnimationMode uint8

const (
	AnimationMode_Once AnimationMode = iota
	AnimationMode_Loop
	AnimationMode_PingPong
	AnimationMode_Reverse
)

type AnimationConfig struct {
	FPS  uint8
	Mode AnimationMode
}

// Holds a set of animations.
type AnimationSet struct {
	nonce      uint64
	animations map[uint64]Animation
	configs    map[uint64]AnimationConfig
	startTimes map[uint64]uint64
}

var _ UpdatableWithRenderer = (*AnimationSet)(nil)

// Creates a new animation set.
func NewAnimationSet() *AnimationSet {
	return &AnimationSet{
		animations: make(map[uint64]Animation),
		configs:    make(map[uint64]AnimationConfig),
		startTimes: make(map[uint64]uint64),
	}
}

// Returns the current unix time in milliseconds.
func (a *AnimationSet) timeNow() uint64 {
	return uint64(time.Now().UnixMilli())
}

// Removes an animation from the set.
func (a *AnimationSet) remove(nonce uint64) {
	delete(a.animations, nonce)
	delete(a.configs, nonce)
	delete(a.startTimes, nonce)
}

// Starts running an animation.
func (a *AnimationSet) RunAnimation(anim Animation, config AnimationConfig) {
	a.nonce++
	a.animations[a.nonce] = anim
	a.configs[a.nonce] = config
	a.startTimes[a.nonce] = a.timeNow()
}

// Calls Update on all animations in the set with the corresponding frame
// and removes them if they are finished.
func (a *AnimationSet) Update(c *CoreRenderer) {
	now := a.timeNow()
	for nonce, anim := range a.animations {
		var (
			config      = a.configs[nonce]
			startTime   = a.startTimes[nonce]
			totalFrames = uint64(anim.NFrames())
		)
		frame := (now - startTime) / (1000 / uint64(config.FPS))

		switch config.Mode {
		case AnimationMode_Once:
			if frame >= totalFrames {
				a.remove(nonce)
				continue
			}
		case AnimationMode_Loop:
			frame = frame % totalFrames
		case AnimationMode_PingPong:
			frame = frame % (totalFrames * 2)
			if frame >= totalFrames {
				frame = totalFrames - (frame - totalFrames) - 1
			}
		case AnimationMode_Reverse:
			if frame >= totalFrames {
				a.remove(nonce)
				continue
			}
			frame = totalFrames - frame - 1
		}

		anim.Update(c, uint(frame))
	}
}

// Doubles the brightness of a sprite every other frame.
type FlashAnimation struct {
	sprite *decren.Sprite
}

var _ Animation = (*FlashAnimation)(nil)

// TODO: fix failing redundant init tx

func NewFlashAnimation(sprite *decren.Sprite, imageOverride *ebiten.Image) *FlashAnimation {
	// if imageOverride != nil {
	// 	sprite.SetImageOverride(imageOverride)
	// }
	return &FlashAnimation{
		sprite: sprite,
	}
}

func (a *FlashAnimation) NFrames() uint {
	return 2
}

func (a *FlashAnimation) Update(c *CoreRenderer, frame uint) {
	flash := frame%2 == 0
	colorM := colorm.ColorM{}
	if flash {
		colorM.ChangeHSV(0, 0, 2)
		a.sprite.SetColorMultiplier(colorM)
	} else {
		// a.sprite.SetImageOverride(nil)
		a.sprite.SetColorMultiplier(colorM)
	}
}

// Plays a fire animation for a unit.
type FireAnimation struct {
	sprite       *decren.Sprite
	playerId     uint8
	unitType     uint8
	direction    assets.Direction
	spriteGetter assets.SpriteGetter
}

var _ Animation = (*FireAnimation)(nil)

func NewFireAnimation(sprite *decren.Sprite, playerId uint8, unitType uint8, direction assets.Direction, spriteGetter assets.SpriteGetter) *FireAnimation {
	return &FireAnimation{
		sprite:       sprite,
		playerId:     playerId,
		unitType:     unitType,
		direction:    direction,
		spriteGetter: spriteGetter,
	}
}

func (f *FireAnimation) NFrames() uint {
	return 3
}

func (f *FireAnimation) Update(c *CoreRenderer, frame uint) {
	switch frame {
	case 0:
		f.sprite.SetImageOverride(f.spriteGetter.GetUnitFireFrame(f.playerId, f.unitType, f.direction, 0))
	case 1:
		f.sprite.SetImageOverride(f.spriteGetter.GetUnitFireFrame(f.playerId, f.unitType, f.direction, 1))
	case 2:
		f.sprite.SetImageOverride(nil)
	}
}
