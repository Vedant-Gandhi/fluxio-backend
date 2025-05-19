package middleware

// For now we keep it the same as middleware. In future if we need any private variable then we can detach it.
type MiddlewareList struct {
	Auth *AuthMiddleware
}

type Middleware struct {
	*MiddlewareList
}

func NewMiddleware(list *MiddlewareList) *Middleware {
	return &Middleware{
		list,
	}
}
