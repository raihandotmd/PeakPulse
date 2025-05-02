package models

// ExerciseListResponse is returned by the ListExercises endpoint.
// swagger:model ExerciseListResponse
type ExerciseListResponse struct {
	// The list of all exercises
	// required: true
	Data []Exercise `json:"data"`
	// The total number of exercises
	// required: true
	Count int64 `json:"count"`
}

// SignupResponse holds the result of a successful signup.
// swagger:model SignupResponse
type SignupResponse struct {
	// Success message
	Message string `json:"message"`
	// Created user object
	User User `json:"user"`
}

// LoginResponse holds the JWT token for a successful login.
// swagger:model LoginResponse
type LoginResponse struct {
	// JWT access token
	Token string `json:"token"`
}

// ErrorResponse holds a generic error message.
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error message
	Error string `json:"error"`
}
