package state

import (
	"sync"
)

type State struct {
	Step       int
	TotalSteps int
	Module     string
	Data       map[string]interface{}
}

type StateManager struct {
	states  map[int64]State
	history map[int64][]State
	mu      sync.RWMutex
}

func NewStateManager() *StateManager {
	return &StateManager{
		states:  make(map[int64]State),
		history: make(map[int64][]State),
	}
}

func (sm *StateManager) Set(chatID int64, state State) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.history[chatID] = append(sm.history[chatID], sm.states[chatID])
	sm.states[chatID] = state
}

func (sm *StateManager) Get(chatID int64) State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.states[chatID]
}

func (sm *StateManager) Back(chatID int64) State {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if len(sm.history[chatID]) > 0 {
		last := sm.history[chatID][len(sm.history[chatID])-1]
		sm.history[chatID] = sm.history[chatID][:len(sm.history[chatID])-1]
		sm.states[chatID] = last
		return last
	}
	return sm.states[chatID]
}

func (sm *StateManager) Clear(chatID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.states, chatID)
	delete(sm.history, chatID)
}
