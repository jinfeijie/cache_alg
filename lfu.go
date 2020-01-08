package cache

import "time"

type LFUNode struct {
	Key        string
	Value      interface{}
	PreNode    *LFUNode
	NextNode   *LFUNode
	LatestTime int64
	Cnt        uint64
}

func NewLFUNode(key string, value interface{}, nextNode *LFUNode, preNode *LFUNode, latestTime int64, cnt uint64) *LFUNode {
	return &LFUNode{Key: key, Value: value, NextNode: nextNode, PreNode: preNode, LatestTime: latestTime, Cnt: cnt}
}

type LFUCache struct {
	Size        uint64
	IndexNodes  *LFUNode
	CurrentSize uint64
}

func NewLFUCache(size uint64, indexNodes *LFUNode, currentSize uint64) *LFUCache {
	return &LFUCache{Size: size, IndexNodes: indexNodes, CurrentSize: currentSize}
}

type LFU struct {
	lfu *LFUCache
}

func NewLFU(lfu *LFUCache) *LFU {
	return &LFU{lfu: lfu}
}

// get
func (c *LFU) Get(Key string) interface{} {

	var lfuVal interface{}

	lfuCurrent := c.lfu.IndexNodes
	for {
		if lfuCurrent.Key == Key {
			lfuVal = lfuCurrent.Value
			lfuCurrent.Cnt += 1
			lfuCurrent.LatestTime = time.Now().UnixNano()

			for {
				if lfuCurrent.PreNode == nil {
					goto ENDLFU
				}

				if lfuCurrent.PreNode.Cnt > lfuCurrent.Cnt {
					goto ENDLFU
				}

				if lfuCurrent.Cnt >= lfuCurrent.PreNode.Cnt {
					if lfuCurrent.NextNode != nil {
						lfuCurrent.PreNode.NextNode = lfuCurrent.NextNode
						lfuCurrent.NextNode.PreNode = lfuCurrent.PreNode
					} else {
						lfuCurrent.PreNode.NextNode = nil
					}

					lfuCurrent.NextNode = lfuCurrent.PreNode
					lfuCurrent.PreNode = lfuCurrent.PreNode.PreNode
				}
			}
		}
		if lfuCurrent.NextNode != nil {
			lfuCurrent = lfuCurrent.NextNode
		} else {
			goto ENDLFU
		}
	}
ENDLFU:
	if *lfuCurrent != *c.lfu.IndexNodes {
		if lfuCurrent.PreNode == nil {
			c.lfu.IndexNodes.PreNode = lfuCurrent
			c.lfu.IndexNodes = lfuCurrent
		}
	}
	return lfuVal
}

func (c *LFU) Set(key string, Value interface{}) uint {
	// lfu
	lfuNode := NewLFUNode(key, Value, c.lfu.IndexNodes, nil, time.Now().UnixNano(), 1)
	c.lfu.IndexNodes = lfuNode
	if c.lfu.IndexNodes.NextNode != nil {
		c.lfu.IndexNodes.NextNode.PreNode = c.lfu.IndexNodes
	}
	c.lfu.CurrentSize += 1
	if c.lfu.CurrentSize > c.lfu.Size {
		idx := 1
		lfuCurrentNode := c.lfu.IndexNodes
		for {
			if uint64(idx) == c.lfu.Size {
				lfuCurrentNode.NextNode = nil
			}

			if lfuCurrentNode.NextNode == nil {
				break
			}
			lfuCurrentNode = lfuCurrentNode.NextNode
			idx += 1
		}
	}
	return 1
}

func (c *LFU) Del(key string) uint {
	return 1
}

func (c *LFU) GetAll() interface{} {
	var lfu []LFUNode
	lfuCurrentNode := c.lfu.IndexNodes
	for {
		lfu = append(lfu, *lfuCurrentNode)
		if lfuCurrentNode.NextNode == nil {
			break
		}
		lfuCurrentNode = lfuCurrentNode.NextNode
	}

	return lfu
}
