package app

import (
	"context"
	"go.uber.org/zap"
)

type AuditAction string

const (
	AuditActionCreate AuditAction = "create"
	AuditActionDelete AuditAction = "delete"
	AuditActionUpdate AuditAction = "update"
	AuditActionList   AuditAction = "list"
)

func (a *App) LogAuditEvent(ctx context.Context, action AuditAction, resourceName string, fields ...zap.Field) {
	base := []zap.Field{
		zap.String("log_type", "audit"),
		zap.String("action", string(action)),
		zap.String("resource", resourceName),
	}

	a.Logger.Info("audit", append(base, fields...)...)
}

//{"level":"info","ts":1740256467.1234,"msg":"audit","log_type":"audit","action":"create","resource":"user_profile","user_id":"12345"}
