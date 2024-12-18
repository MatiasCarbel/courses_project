package models

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password,omitempty"`
    Admin    bool   `json:"admin"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type Response struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
    Token   string `json:"token,omitempty"`
} 
