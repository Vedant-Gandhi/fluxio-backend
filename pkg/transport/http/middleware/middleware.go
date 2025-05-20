package middleware

type MiddlewareList struct {
	Auth *AuthMiddleware
}

// For now we keep it the same as middleware. In future if we need any private variable then we can detach/append it.
type Middleware struct {
	*MiddlewareList
}

func NewMiddleware(list *MiddlewareList) *Middleware {
	return &Middleware{
		list,
	}
}
