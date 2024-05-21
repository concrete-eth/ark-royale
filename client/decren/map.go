package decren

// Implements a map with ordered keys.
// The keys will always be returned in the same order between writes.
// Adding or deleting a key might change the order of the keys.
type OrderedMap struct {
	size     int
	keys     []string
	keyIndex map[string]int
	data     map[string]interface{}
}

// Creates a new OrderedMap.
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		size:     0,
		keys:     make([]string, 0),
		keyIndex: make(map[string]int),
		data:     make(map[string]interface{}),
	}
}

func (om *OrderedMap) Set(key string, value interface{}) {
	if _, exists := om.data[key]; !exists {
		if len(om.keys) > om.size {
			om.keys[om.size] = key
		} else {
			om.keys = append(om.keys, key)
		}
		om.keyIndex[key] = om.size
		om.size++
	}
	om.data[key] = value
}

func (om *OrderedMap) Delete(key string) {
	if _, exists := om.data[key]; exists {
		// Delete the key from the data map
		delete(om.data, key)

		// Replace the deleted key in the keys slice
		om.keys[om.keyIndex[key]] = om.keys[om.size-1]
		om.keys[om.size-1] = ""

		// Update the index of the key that was moved
		om.keyIndex[om.keys[om.keyIndex[key]]] = om.keyIndex[key]

		// Delete the key from the keyIndex map
		delete(om.keyIndex, key)

		// Decrease the size
		om.size--
	}
}

func (om *OrderedMap) Get(key string) (interface{}, bool) {
	val, exists := om.data[key]
	return val, exists
}

func (om *OrderedMap) Size() int {
	return om.size
}

func (om *OrderedMap) Keys() []string {
	return om.keys[:om.size]
}
