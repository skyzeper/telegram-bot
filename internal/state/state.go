package state

import (
	"sync"
)

// State represents the state of a user interaction
type State struct {
	Module     string                 `json:"module"`
	Step       int                    `json:"step"`
	TotalSteps int                    `json:"total_steps"`
	Data       map[string]interface{} `json:"data"`
}

// Manager manages user states
type Manager struct {
	states map[int64]State
	mutex  sync.Mutex
}

// NewManager creates a new state manager
func NewManager() *Manager {
	return &Manager{
		states: make(map[int64]State),
	}
}

// Get retrieves the state for a chat ID
func (m *Manager) Get(chatID int64) State {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.states[chatID]
}

// Set updates the state for a chat ID
func (m *Manager) Set(chatID int64, state State) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.states[chatID] = state
}

// Clear removes the state for a chat ID
func (m *Manager) Clear(chatID int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.states, chatID)
}