package auth

type CreateUserBody struct {
	FirstName   string   `json:"first_name" validate:"required"`
	LastName    string   `json:"last_name" validate:"required"`
	Password    string   `json:"password" validate:"required"`
	Username    string   `json:"username" validate:"required"`
	Email       string   `json:"email" validate:"required,email"`
	SearchTerms []string `json:"search_terms" validate:"required"`
}

type LoginUserBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}
