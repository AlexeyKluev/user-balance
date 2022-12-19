package metrics

import (
	"context"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type RedisContextKey string

const (
	RedisContextKeyStart RedisContextKey = "start"
)

type RedisHook struct {
	Metrics *prometheus.HistogramVec
	Logger  *zap.Logger
}

var (
	_ redis.Hook = (*RedisHook)(nil)
)

func (r RedisHook) BeforeProcess(ctx context.Context, _ redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, RedisContextKeyStart, time.Now()), nil
}

func (r RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if since := getRedisSince(ctx); since != 0 {
		r.Metrics.WithLabelValues(cmd.Name()).Observe(since)

		if r.Logger != nil {
			r.Logger.Debug(
				"[REDIS]",
				zap.String("cmd", cmd.String()),
				zap.Float64("duration, ms", since),
			)
		}
	}

	return nil
}

func (r RedisHook) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, RedisContextKeyStart, time.Now()), nil
}

func (r RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if since := getRedisSince(ctx); since != 0 {
		var (
			name = make([]string, len(cmds))
			strs = make([]string, len(cmds))
		)

		for i, cmd := range cmds {
			name[i] = cmd.Name()
			strs[i] = cmd.String()
		}

		r.Metrics.WithLabelValues(strings.Join(name, "_")).Observe(since)

		if r.Logger != nil {
			r.Logger.Debug(
				"[REDIS]",
				zap.Strings("cmd", strs),
				zap.Float64("duration, ms", since),
			)
		}
	}

	return nil
}

func getRedisSince(ctx context.Context) float64 {
	var (
		since float64
	)

	if v := ctx.Value(RedisContextKeyStart); v != nil {
		if start, ok := v.(time.Time); ok {
			since = Since(start)
		}
	}

	return since
}
