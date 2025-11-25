package types

type UserRole string

const (
	Admin      UserRole = "admin"
	TeamMember UserRole = "team_member"
)

type NewUser struct {
	Email    string `json:"email" doc:"user email" format:"email"`
	Username string `json:"username" doc:"username" minLength:"1" example:"username"`
	Password string `json:"password" doc:"password" minLength:"8" example:"password"`
	Role     string `json:"role" doc:"users role" example:"admin"`
}

type NewUserInput struct {
	Body NewUser
}

type NewUserResponse struct {
	Body struct {
		ID string
	}
}

type User struct {
	ID       string
	Email    string
	Username string
	Password string
	Role     string
}

type GetUserResponseBody struct {
	ID       string
	Email    string
	Username string
	Role     string
}

type GetUserResponse struct {
	Body GetUserResponseBody
}

type UserClaims struct {
	ID    string
	Email string
	Role  string
}
