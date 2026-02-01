package costspace

import "github.com/ethereum/go-ethereum/p2p/enode"

type CostModel interface {
	// 越小越优先
	Cost(peer enode.ID) float64
	// 给可视化用
	Coord(peer enode.ID) (r float64, theta float64, ok bool)
}
