package algo

// Node Contain the song information and pointers to the previous and next nodes.
type Node struct {
	key   string // song identifier
	value string // song data
	prev  *Node
	next  *Node
}

// LRU cache contain the capacity,
// a hash map for quick access,
// and pointers to the head and tail of the doubly linked list.
type LRUCache struct {
	capacity int
	lookup   map[string]*Node
	head     *Node
	tail     *Node
}

// NewLRUCache initialize the cache with a given capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		lookup:   make(map[string]*Node),
	}
}

// Get retrieve a song from the cache and mark it as recently used.
func (cache *LRUCache) Get(key string) (string, bool) {
	if node, exists := cache.lookup[key]; exists {
		cache.moveToHead(node)
		return node.value, true
	}
	return "", false
}

// Put add a song to the cache. If the cache exceeds its capacity,
// remove the least recently used song.
func (cache *LRUCache) Put(key, value string) {
	if node, exists := cache.lookup[key]; exists {
		node.value = value
		cache.moveToHead(node)
		return
	}
	newNode := &Node{key: key, value: value}
	cache.lookup[key] = newNode
	cache.addToHead(newNode)
	if len(cache.lookup) > cache.capacity {
		removed := cache.removeTail()
		delete(cache.lookup, removed.key)
	}
}

// moveToHead moves an existing node to the head of the doubly linked list, marking it as the most recently used.
func (cache *LRUCache) moveToHead(node *Node) {
	if node == cache.head {
		return
	}
	// Remove the node from its current position
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == cache.tail {
		cache.tail = node.prev
	}
	// Add the node to the head
	node.prev = nil
	node.next = cache.head
	if cache.head != nil {
		cache.head.prev = node
	}
	cache.head = node
	if cache.tail == nil {
		cache.tail = node
	}
}

// addToHead adds a new node to the head of the list
func (cache *LRUCache) addToHead(node *Node) {
	node.prev = nil
	node.next = cache.head
	if cache.head != nil {
		cache.head.prev = node
	}
	cache.head = node
	if cache.tail == nil {
		cache.tail = node
	}
}

// removeTail function removes the node from the tail of the list, which is the least recently used item.
func (cache *LRUCache) removeTail() *Node {
	if cache.tail == nil {
		return nil
	}
	tail := cache.tail
	if tail.prev != nil {
		tail.prev.next = nil
	} else {
		// This was the only node
		cache.head = nil
	}
	cache.tail = tail.prev
	tail.prev = nil
	tail.next = nil
	return tail
}
