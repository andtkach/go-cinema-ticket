package auth

import "context"

type contextKey string

const (
	userIDKey contextKey = "userID"
	usernameKey contextKey = "username"
	groupsKey contextKey = "groups"
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey).(string)
	return v
}

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

func UsernameFromContext(ctx context.Context) string {
	v, _ := ctx.Value(usernameKey).(string)
	return v
}

func WithGroups(ctx context.Context, groups []string) context.Context {
	return context.WithValue(ctx, groupsKey, groups)
}

func GroupsFromContext(ctx context.Context) []string {
	v, _ := ctx.Value(groupsKey).([]string)
	return v
}
