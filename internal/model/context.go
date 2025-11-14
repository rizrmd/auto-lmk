package model

import (
	"context"
	"errors"
)

type contextKey string

const (
	TenantIDKey contextKey = "tenant_id"
	UserIDKey   contextKey = "user_id"
)

// WithTenantID adds tenant ID to context
func WithTenantID(ctx context.Context, tenantID int) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(ctx context.Context) (int, error) {
	tenantID, ok := ctx.Value(TenantIDKey).(int)
	if !ok {
		return 0, errors.New("tenant ID not found in context")
	}
	return tenantID, nil
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}
