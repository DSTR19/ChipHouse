package domain

type Role string

const (
	RoleUser     Role = "USER"
	RoleAdmin    Role = "ADMIN"
	RoleInvestor Role = "INVESTOR"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	Role         Role
	TenantID     *int
}
