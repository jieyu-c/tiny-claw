package engine

import (
	"sync"
	"time"
	"tiny-claw/internal/schema"
)

type Session struct {
	Id        string
	WorkDir   string
	CreatedAt time.Time
	UpdatedAt time.Time
	history   []*schema.Message
	mu        sync.RWMutex
}

func NewSession(id string, workDir string) *Session {
	return &Session{
		Id:        id,
		WorkDir:   workDir,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		history:   make([]*schema.Message, 0),
	}
}

func (s *Session) Append(messages ...*schema.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append(s.history, messages...)
	s.UpdatedAt = time.Now()
}

func (s *Session) GetWorkingMemory(limit int) []*schema.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total := len(s.history)
	if total <= limit || limit <= 0 {
		res := make([]*schema.Message, total)
		copy(res, s.history)
		return res
	}
	res := make([]*schema.Message, total)
	copy(res, s.history[total-limit:])

	for len(res) > 0 {
		if res[0].Role == schema.RoleUser && res[0].ToolCallID != "" {
			res = res[1:]
		} else {
			break
		}
	}
	return res
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

var once sync.Once

var globalSessionMgr *SessionManager

func NewSessionManager() *SessionManager {
	once.Do(func() {
		globalSessionMgr = &SessionManager{
			sessions: make(map[string]*Session),
		}
	})
	return globalSessionMgr
}

func (s *SessionManager) GetOrCreate(id string, workDir string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sess, exists := s.sessions[id]; exists {
		return sess
	}

	sess := NewSession(id, workDir)
	s.sessions[id] = sess
	return sess

}
