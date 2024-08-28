package tgbot

var SessionMgr *SessionManager

func init() {
	SessionMgr = &SessionManager{
		sessions: make(map[int64]*Session),
	}
}

type SessionManager struct {
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
