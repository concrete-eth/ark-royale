package decren

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	client_utils "github.com/concrete-eth/ark-royale/client/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// Returns a string ID from a slice of interface{} values.
// IDs are not guaranteed to be unique.
func idToStr(id []interface{}) string {
	return fmt.Sprintf("%v", id)
}

type ShadowType uint

const (
	ShadowType_Box ShadowType = iota
	ShadowType_Cast
)

type Shadow struct {
	Enabled bool
	Type    ShadowType
	Color   color.RGBA
	Offset  image.Point
}

func newShadowColorMatrix(shadow Shadow) colorm.ColorM {
	colorM := colorm.ColorM{}
	switch shadow.Type {
	case ShadowType_Box:
	case ShadowType_Cast:
		colorM.Scale(0, 0, 0, 1)
		colorM.Translate(float64(shadow.Color.R)/255, float64(shadow.Color.G)/255, float64(shadow.Color.B)/255, 0)
		colorM.Scale(1, 1, 1, float64(shadow.Color.A)/255)
	}
	return colorM
}

// Holds the state and drawing information for a sprite.
type Sprite struct {
	id string // Unique ID within the layer.

	layer  *Layer          // The layer this sprite belongs to.
	rect   image.Rectangle // The rectangle the sprite occupies within the layer.
	image  *ebiten.Image   // The image to draw.
	colorM colorm.ColorM   // The color matrix to apply to the image when drawing.

	shadow  Shadow // The shadow to apply to the sprite's image.
	visible bool   // Whether the sprite is visible.

	// Relative overrides used for animation.
	positionDelta   image.Point   // The position delta to apply to the sprite's position.
	sizeDelta       image.Point   // The size delta to apply to the sprite's size.
	imageOverride   *ebiten.Image // The image to draw instead of the sprite's image.
	colorMultiplier colorm.ColorM // The color multiplier to apply to the sprite's image.
}

func (s *Sprite) setLayerDirty() {
	s.layer.dirty = true
}

func (s *Sprite) ID() string {
	return s.id
}

func (s *Sprite) Layer() *Layer {
	return s.layer
}

func (s *Sprite) Delete() {
	if _, ok := s.layer.sprites.Get(s.id); !ok {
		return
	}
	s.layer.sprites.Delete(s.id)
	s.setLayerDirty()
}

func (s *Sprite) SetPosition(position image.Point) *Sprite {
	if position.Eq(s.rect.Min) {
		return s
	}
	s.rect = s.rect.Sub(s.rect.Min).Add(position)
	s.setLayerDirty()
	return s
}

func (s *Sprite) Position() image.Point {
	return s.rect.Min
}

func (s *Sprite) SetSize(size image.Point) *Sprite {
	if size.Eq(s.rect.Size()) {
		return s
	}
	s.rect.Max = s.rect.Min.Add(size)
	s.setLayerDirty()
	return s
}

func (s *Sprite) Size() image.Point {
	return s.rect.Size()
}

func (s *Sprite) SetRect(rect image.Rectangle) *Sprite {
	if rect.Eq(s.rect) {
		return s
	}
	s.rect = rect
	s.setLayerDirty()
	return s
}

func (s *Sprite) Rect() image.Rectangle {
	return s.rect
}

func (s *Sprite) FitToImage() *Sprite {
	s.SetSize(s.image.Bounds().Size())
	return s
}

func (s *Sprite) SetColorMatrix(colorM colorm.ColorM) *Sprite {
	s.colorM = colorM
	s.setLayerDirty()
	return s
}

func (s *Sprite) ColorMatrix() colorm.ColorM {
	return s.colorM
}

func (s *Sprite) SetImage(image *ebiten.Image) *Sprite {
	if image == s.image {
		return s
	}
	s.image = image
	s.setLayerDirty()
	return s
}

func (s *Sprite) Image() *ebiten.Image {
	return s.image
}

func (s *Sprite) SetPositionDelta(positionDelta image.Point) *Sprite {
	if positionDelta.Eq(s.positionDelta) {
		return s
	}
	s.positionDelta = positionDelta
	s.setLayerDirty()
	return s
}

func (s *Sprite) PositionDelta() image.Point {
	return s.positionDelta
}

func (s *Sprite) SetSizeDelta(sizeDelta image.Point) *Sprite {
	if sizeDelta.Eq(s.sizeDelta) {
		return s
	}
	s.sizeDelta = sizeDelta
	s.setLayerDirty()
	return s
}

func (s *Sprite) SizeDelta() image.Point {
	return s.sizeDelta
}

func (s *Sprite) SetImageOverride(imageOverride *ebiten.Image) *Sprite {
	if imageOverride == s.imageOverride {
		return s
	}
	s.imageOverride = imageOverride
	s.setLayerDirty()
	return s
}

func (s *Sprite) ImageOverride() *ebiten.Image {
	return s.imageOverride
}

func (s *Sprite) SetColorMultiplier(colorMultiplier colorm.ColorM) *Sprite {
	s.colorMultiplier = colorMultiplier
	s.setLayerDirty()
	return s
}

func (s *Sprite) ColorMultiplier() colorm.ColorM {
	return s.colorMultiplier
}

func (s *Sprite) SetShadow(shadow Shadow) *Sprite {
	if shadow == s.shadow {
		return s
	}
	s.shadow = shadow
	s.setLayerDirty()
	return s
}

func (s *Sprite) Shadow() Shadow {
	return s.shadow
}

func (s *Sprite) ToggleShadow(enabled bool) *Sprite {
	if enabled == s.shadow.Enabled {
		return s
	}
	s.shadow.Enabled = enabled
	s.setLayerDirty()
	return s
}

func (s *Sprite) SetVisible(visible bool) *Sprite {
	if visible == s.visible {
		return s
	}
	s.visible = visible
	s.setLayerDirty()
	return s
}

func (s *Sprite) Visible() bool {
	return s.visible
}

// Wrapper around OrderedMap with Sprite types.
type SpriteOrderedMap struct {
	OrderedMap
}

// Creates a new SpriteOrderedMap.
func NewSpriteOrderedMap() *SpriteOrderedMap {
	return &SpriteOrderedMap{
		OrderedMap: *NewOrderedMap(),
	}
}

func (om *SpriteOrderedMap) Set(key string, sprite *Sprite) {
	om.OrderedMap.Set(key, sprite)
}

func (om *SpriteOrderedMap) Get(key string) (*Sprite, bool) {
	val, exists := om.OrderedMap.Get(key)
	if !exists {
		return nil, false
	}
	return val.(*Sprite), true
}

// Holds the state and drawing information for a layer.
type Layer struct {
	layerSet        *LayerSet         // The layer set this layer belongs to.
	cache           bool              // Whether to cache the layer's drawn image.
	depth           int               // The depth of the layer.
	backgroundColor color.Color       // The background color of the layer.
	visible         bool              // Whether the layer is visible.
	srcRect         image.Rectangle   // The source rectangle of the layer.
	destRect        image.Rectangle   // The destination rectangle of the layer.
	colorM          colorm.ColorM     // The color matrix to apply to the layer's image when drawing.
	sprites         *SpriteOrderedMap // The sprites in the layer.
	image           *ebiten.Image     // The cached image of the layer.
	dirty           bool              // Whether the layer needs to be redrawn.
}

func (l *Layer) SetDepth(depth int) *Layer {
	if l.depth == depth {
		return l
	}
	l.depth = depth
	l.layerSet.sortIds()
	return l
}

func (l *Layer) Depth() int {
	return l.depth
}

func (l *Layer) SetBackgroundColor(backgroundColor color.Color) *Layer {
	if l.backgroundColor == backgroundColor {
		return l
	}
	l.backgroundColor = backgroundColor
	l.dirty = true
	return l
}

func (l *Layer) BackgroundColor() color.Color {
	return l.backgroundColor
}

func (l *Layer) SetVisible(visible bool) *Layer {
	if l.visible == visible {
		return l
	}
	l.visible = visible
	l.dirty = true
	return l
}

func (l *Layer) Visible() bool {
	return l.visible
}

func (l *Layer) SetSourceRect(rect image.Rectangle) *Layer {
	if l.srcRect.Eq(rect) {
		return l
	}
	l.srcRect = rect
	l.dirty = true
	return l
}

func (l *Layer) SourceRect() image.Rectangle {
	return l.srcRect
}

func (l *Layer) SetDestinationRect(rect image.Rectangle) *Layer {
	if l.destRect.Eq(rect) {
		return l
	}
	l.destRect = rect
	l.dirty = true
	return l
}

func (l *Layer) DestinationRect() image.Rectangle {
	return l.destRect
}

func (l *Layer) SetColorMatrix(colorM colorm.ColorM) *Layer {
	l.colorM = colorM
	l.dirty = true
	return l
}

func (l *Layer) ColorMatrix() colorm.ColorM {
	return l.colorM
}

func (l *Layer) Sprite(id ...interface{}) *Sprite {
	idStr := idToStr(id)
	if sprite, ok := l.sprites.Get(idStr); ok {
		return sprite
	}
	sprite := &Sprite{id: idStr, layer: l, visible: true}
	l.sprites.Set(idStr, sprite)
	l.dirty = true
	return sprite
}

func (l *Layer) Sprites() []*Sprite {
	sprites := make([]*Sprite, l.sprites.Size())
	i := 0
	for _, key := range l.sprites.Keys() {
		sprites[i], _ = l.sprites.Get(key)
		i++
	}
	return sprites
}

func (l *Layer) RemoveSprite(id ...interface{}) {
	idStr := idToStr(id)
	if _, ok := l.sprites.Get(idStr); !ok {
		return
	}
	l.sprites.Delete(idStr)
	l.dirty = true
}

func (l *Layer) Cache(cache bool) *Layer {
	if l.cache == cache {
		return l
	}
	if !cache {
		l.image = nil
	}
	return l
}

func (l *Layer) Clear() {
	if l.sprites.Size() == 0 {
		return
	}
	l.sprites = NewSpriteOrderedMap()
	l.dirty = true
}

func (l *Layer) IsDirty() bool {
	return l.dirty
}

// Fills the given image with the layer's sprite that overlap the layer's source rectangle.
func (l *Layer) fillWithSource(dst *ebiten.Image) {
	if l.backgroundColor != nil {
		dst.Fill(l.backgroundColor)
	}
	sx, sy := client_utils.ScaleRatio(dst.Bounds(), l.srcRect)
	for _, key := range l.sprites.Keys() {
		sprite, ok := l.sprites.Get(key)
		if !ok {
			panic(fmt.Sprintf("sprite %s not found", key))
		}
		if !sprite.visible {
			continue
		}
		if !sprite.rect.Overlaps(l.srcRect) {
			continue
		}

		src := sprite.image
		rect := sprite.rect
		colorM := colorm.ColorM{}
		colorM.Concat(sprite.colorM)

		if sprite.imageOverride != nil {
			src = sprite.imageOverride
		}
		if src == nil {
			continue
		}

		colorM.Concat(sprite.colorMultiplier)
		colorM.Concat(l.colorM)

		rect = rect.Add(sprite.positionDelta)
		rect.Max = rect.Max.Add(sprite.sizeDelta)
		dstRect := client_utils.ScaleRectangle(rect.Sub(l.srcRect.Min), sx, sy).Add(dst.Bounds().Min)
		op := client_utils.NewDrawOptions(dstRect, src.Bounds())

		// rect = rect.Add(sprite.positionDelta)
		// rect.Max = rect.Max.Add(sprite.sizeDelta)
		// relRect := rect.Sub(l.srcRect.Min)
		// scaledRect := image.Rectangle{
		// 	Min: image.Point{
		// 		X: relRect.Min.X * dst.Bounds().Dx() / l.srcRect.Dx(),
		// 		Y: relRect.Min.Y * dst.Bounds().Dy() / l.srcRect.Dy(),
		// 	},
		// 	Max: image.Point{
		// 		X: relRect.Max.X * dst.Bounds().Dx() / l.srcRect.Dx(),
		// 		Y: relRect.Max.Y * dst.Bounds().Dy() / l.srcRect.Dy(),
		// 	},
		// }
		// dstRect := scaledRect.Add(dst.Bounds().Min)
		// op := client_utils.NewDrawOptions(dstRect, src.Bounds())

		if sprite.shadow.Enabled {
			shadowSrc := src
			if sprite.shadow.Type == ShadowType_Box {
				shadowSrc = ebiten.NewImage(32, 32)
				shadowSrc.Fill(sprite.shadow.Color)
			}

			shadowRect := dstRect.Add(sprite.shadow.Offset)
			op := client_utils.NewDrawOptions(shadowRect, shadowSrc.Bounds())
			// op.Filter = ebiten.FilterLinear

			colorM := newShadowColorMatrix(sprite.shadow)
			colorM.Concat(l.colorM)
			colorm.DrawImage(dst, shadowSrc, colorM, op)
		}

		colorm.DrawImage(dst, src, colorM, op)
	}
}

// Returns an image with the layer's sprites that overlap the layer's source rectangle drawn on it.
func (l *Layer) Image() *ebiten.Image {
	if l.cache && !l.dirty {
		// If the layer is cached and not dirty, return the cached image
		return l.image
	}
	// Otherwise, create a new image and fill it with the layer's sprites
	image := ebiten.NewImage(l.destRect.Dx(), l.destRect.Dy())
	l.fillWithSource(image)
	l.dirty = false
	if l.cache {
		// If the layer is cached, set the cached image to the new image
		l.image = image
	}
	return image
}

// Draws the layer's sprites that overlap the layer's source rectangle on the given image.
func (l *Layer) Draw(image *ebiten.Image) {
	if l.cache {
		// If the layer is cached, draw the cached image on the given image
		img := l.Image()
		op := client_utils.NewDrawOptions(l.destRect, img.Bounds())
		colorm.DrawImage(image, img, colorm.ColorM{}, op)
	} else {
		// Otherwise, fill the given image with the layer's sprites
		l.fillWithSource(image.SubImage(l.destRect).(*ebiten.Image))
	}
}

// Holds a set of layers.
type LayerSet struct {
	layers    map[string]*Layer
	sortedIds []string
}

// Creates a new LayerSet.
func NewLayerSet() *LayerSet {
	return &LayerSet{
		layers: make(map[string]*Layer),
	}
}

func (ls *LayerSet) sortIds() {
	sort.Slice(ls.sortedIds, func(i, j int) bool {
		iId := ls.sortedIds[i]
		jId := ls.sortedIds[j]
		return ls.layers[iId].depth < ls.layers[jId].depth
	})
}

func (ls *LayerSet) Layer(id ...interface{}) *Layer {
	idStr := idToStr(id)
	if _, ok := ls.layers[idStr]; ok {
		return ls.layers[idStr]
	}
	layer := &Layer{
		layerSet: ls,
		cache:    false,
		depth:    0,
		visible:  true,
		sprites:  NewSpriteOrderedMap(),
		dirty:    true,
	}
	ls.layers[idStr] = layer
	ls.sortedIds = append(ls.sortedIds, idStr)
	ls.sortIds()
	return layer
}

func (ls *LayerSet) Layers() []*Layer {
	layers := make([]*Layer, len(ls.sortedIds))
	for i, id := range ls.sortedIds {
		layers[i] = ls.layers[id]
	}
	return layers
}

func (ls *LayerSet) RemoveLayer(id ...interface{}) {
	idStr := idToStr(id)
	delete(ls.layers, idStr)
	for i, sortedId := range ls.sortedIds {
		if sortedId == idStr {
			ls.sortedIds = append(ls.sortedIds[:i], ls.sortedIds[i+1:]...)
			break
		}
	}
}

func (ls *LayerSet) Clear() {
	ls.layers = make(map[string]*Layer)
	ls.sortedIds = []string{}
}

// Draws the layer set's layers on the given image by increasing depth.
func (ls *LayerSet) Draw(image *ebiten.Image) {
	for _, id := range ls.sortedIds {
		layer := ls.layers[id]
		if !layer.visible {
			continue
		}
		layer.Draw(image)
	}
}

func (ls *LayerSet) IsDirty() bool {
	for _, layer := range ls.layers {
		if layer.dirty {
			return true
		}
	}
	return false
}
