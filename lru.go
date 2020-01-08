package cache

type LRUNode struct {
	Key      string
	Value    interface{}
	PreNode  *LRUNode
	NextNode *LRUNode
}

func NewLRUNode(key string, value interface{}, nextNode *LRUNode, preNode *LRUNode) *LRUNode {
	return &LRUNode{Key: key, Value: value, NextNode: nextNode, PreNode: preNode}
}

type LRUCache struct {
	Size        uint64
	IndexNodes  *LRUNode
	CurrentSize uint64
}

func NewLRUCache(size uint64, indexNodes *LRUNode, currentSize uint64) *LRUCache {
	return &LRUCache{Size: size, IndexNodes: indexNodes, CurrentSize: currentSize}
}

type LRU struct {
	lru *LRUCache
}

func NewLRU(lru *LRUCache) *LRU {
	return &LRU{lru: lru}
}

// get
func (c *LRU) Get(Key string) interface{} {

	var lruVal interface{}

	lruCurrent := c.lru.IndexNodes
	for {
		if lruCurrent.Key == Key {
			lruVal = lruCurrent.Value
			if lruCurrent.PreNode != nil {
				if lruCurrent.NextNode != nil {
					lruCurrent.PreNode.NextNode = lruCurrent.NextNode
					lruCurrent.NextNode.PreNode = lruCurrent.PreNode
				} else {
					lruCurrent.PreNode.NextNode = nil
				}
				// 去接到头节点上
				c.lru.IndexNodes.PreNode = lruCurrent
				lruCurrent.NextNode = c.lru.IndexNodes
				lruCurrent.PreNode = nil
				c.lru.IndexNodes = lruCurrent
			}
			goto ENDLRU
		}
		if lruCurrent.NextNode != nil {
			lruCurrent = lruCurrent.NextNode
		} else {
			goto ENDLRU
		}
	}
ENDLRU:
	return lruVal
}

func (c *LRU) Set(key string, Value interface{}) uint {
	// lru
	// 先判断是否存在
	lruNode := NewLRUNode(key, Value, c.lru.IndexNodes, nil)
	c.lru.IndexNodes = lruNode
	if c.lru.IndexNodes.NextNode != nil {
		c.lru.IndexNodes.NextNode.PreNode = c.lru.IndexNodes
	}
	c.lru.CurrentSize += 1
	if c.lru.CurrentSize > c.lru.Size {
		idx := 1
		lruCurrentNode := c.lru.IndexNodes
		for {
			if uint64(idx) == c.lru.Size {
				lruCurrentNode.NextNode = nil
			}
			if lruCurrentNode.NextNode == nil {
				break
			}
			lruCurrentNode = lruCurrentNode.NextNode
			idx += 1
		}
	}
	return 1
}

func (c *LRU) Del(key string) uint {
	return 1
}

func (c *LRU) GetAll() interface{} {
	lruCurrentNode := c.lru.IndexNodes
	var lru []LRUNode
	for {
		lru = append(lru, *lruCurrentNode)
		if lruCurrentNode.NextNode == nil {
			break
		}
		lruCurrentNode = lruCurrentNode.NextNode
	}

	return lru
}
