package models

//CreateRoomRequest

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetRoomsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetClientsResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token    string `json:"token"`
	ID       string `json:"id"`
	Username string `json:"username"`
}

// {"token":"w234cdsdjhfskjdf.sfj398jojldsjfvn4lk.sdfj3oi4jkvfjdlsg","id":"234rlj34-4983vhfkdh44-fvdkfjvh","username":"petya"}
