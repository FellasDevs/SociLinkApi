package dto

type SignUpRequestDto struct {
	Email     string
	Name      string
	Nickname  string
	Birthdate string
	Password  string
}

type SignInRequestDto struct {
	Email    string
	Password string
}

type SignInResponseDto struct {
	User      UserResponseDto
	AuthToken string
}

type SignUpResponseDto struct {
	User      UserResponseDto
	AuthToken string
}

type AuthHeader struct {
	AuthToken string `header:"Authorization"`
}
