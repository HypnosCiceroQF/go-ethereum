package hypermodel

import (
	"sync"

	"github.com/ethereum/go-ethereum/p2p/costspace"
	"github.com/ethereum/go-ethereum/p2p/costspace/hyperbolic"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

type Model struct {
	mu     sync.RWMutex
	coords map[enode.ID]hyperbolic.Polar
	center hyperbolic.Polar
}

func New() *Model {
	return &Model{
		coords: make(map[enode.ID]hyperbolic.Polar),
		center: hyperbolic.Polar{R: 0, Theta: 0},
	}
}

// 外部喂坐标（今晚 PoC 用这个）
func (m *Model) Set(id enode.ID, p hyperbolic.Polar) {
	m.mu.Lock()
	m.coords[id] = p
	m.mu.Unlock()
}

func (m *Model) Cost(peer enode.ID) float64 {
	m.mu.RLock()
	p, ok := m.coords[peer]
	m.mu.RUnlock()
	if !ok {
		// 没坐标就给一个中性代价（别让它永远排最后导致永远不用）
		// return 1e9
		return float64(peer[0]) // 后续添加为TTP
	}
	// 这里先用“离中心的距离”作为 cost（越靠中心越优先）
	// 后面你要“目标导向”再扩展
	return hyperbolic.Distance(p, m.center)
}

func (m *Model) Coord(peer enode.ID) (float64, float64, bool) {
	m.mu.RLock()
	p, ok := m.coords[peer]
	m.mu.RUnlock()
	if !ok {
		return 0, 0, false
	}
	return p.R, p.Theta, true
}

var _ costspace.CostModel = (*Model)(nil)
