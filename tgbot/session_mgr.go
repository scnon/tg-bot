package tgbot

type SessionEvent func(s *Session)

var SessionMgr *SessionManager

func init() {
	SessionMgr = &SessionManager{
		sessions: map[int64]*Session{},
	}
}

type SessionManager struct {
	events   []SessionEvent
	sessions map[int64]*Session
}

func (m *SessionManager) GetSession(uid int64) *Session {
	session := m.sessions[uid]
	if session == nil {
		session = &Session{}
		m.sessions[uid] = session
	}

	return session
}

func (m *SessionManager) LoadSession(sessions map[int64]*Session) {
	m.sessions = sessions
}

func (m *SessionManager) RegSessionEvent(event SessionEvent) {
	m.events = append(m.events, event)
}

func (m *SessionManager) OnChange(s *Session) {
	for _, event := range m.events {
		event(s)
	}
}
