package security

import "context"

// Puts token into context
func (token *JwtToken) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKey{}, token)
}

// Get token from context
func TokenFromContext(ctx context.Context) *JwtToken {
	value := ctx.Value(ContextKey{})
	if value != nil {
		return value.(*JwtToken)
	}

	return &JwtToken{}
}
