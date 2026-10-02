package auth

type TokenManager interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(tokenStr string) (int64, error)
}
