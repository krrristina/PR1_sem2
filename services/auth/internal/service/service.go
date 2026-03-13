package service

// AuthService содержит бизнес-логику авторизации
type AuthService struct {
	validToken string
}

func New() *AuthService {
	return &AuthService{validToken: "demo-token"}
}

func (s *AuthService) Login(username, password string) (string, bool) {
	if username == "student" && password == "student" {
		return s.validToken, true
	}
	return "", false
}

func (s *AuthService) Verify(token string) (string, bool) {
	if token == s.validToken {
		return "student", true
	}
	return "", false
}
