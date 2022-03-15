package session

type Session struct {
	authenticated bool
	userID        uint
}

func Create(authenticated bool, userID uint) *Session {
	return &Session{
		authenticated: authenticated,
		userID:        userID,
	}
}

func (s Session) IsAuthenticated() bool {
	return s.userID > 0 && s.authenticated
}

func (s Session) GetUserID() uint {
	return s.userID
}
