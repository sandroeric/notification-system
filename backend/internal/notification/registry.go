package notification

import (
	"sync"

	"github.com/sandro/notification-system/internal/domain"
)

type Registry struct {
	mu         sync.RWMutex
	strategies map[domain.Channel]domain.NotifierStrategy
}

func NewRegistry() *Registry {
	return &Registry{
		strategies: make(map[domain.Channel]domain.NotifierStrategy),
	}
}

func (r *Registry) Register(strategy domain.NotifierStrategy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.strategies[strategy.GetChannelType()] = strategy
}

func (r *Registry) Get(channel domain.Channel) (domain.NotifierStrategy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	strategy, exists := r.strategies[channel]
	return strategy, exists
}
