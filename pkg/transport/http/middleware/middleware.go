package middleware

type MiddlewareList struct {
	Auth *AuthMiddleware
}

type Middleware struct {
	Auth *AuthMiddleware
}

func NewMiddleware(list *MiddlewareList) *Middleware {
	return &Middleware{
		Auth: list.Auth,
	}
}
