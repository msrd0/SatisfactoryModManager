package logging

import (
	"context"
	"log/slog"

	"github.com/samber/lo"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/ficsitcli"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

func redactGamePathCredentialsMiddleware(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
	attrs := make([]slog.Attr, 0, record.NumAttrs())

	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, redactPaths(attr))
		return true
	})

	// new record with redacted paths
	record = slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.AddAttrs(attrs...)

	return next(ctx, record)
}

func redactPaths(attr slog.Attr) slog.Attr {
	k := attr.Key
	v := attr.Value
	kind := attr.Value.Kind()

	switch kind {
	case slog.KindGroup:
		attrs := v.Group()
		for i := range attrs {
			attrs[i] = redactPaths(attrs[i])
		}
		return slog.Group(k, lo.ToAnySlice(attrs)...)
	case slog.KindString:
		if isGamePath(v.String()) {
			return slog.String(k, utils.RedactPath(v.String()))
		}
	default:
		break
	}
	return attr
}

func isGamePath(str string) bool {
	return ficsitcli.FicsitCLI.GetInstallation(str) != nil
}
