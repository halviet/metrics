package logger

import "go.uber.org/zap"

type Opts struct {
	Lvl string
}

func New(opts ...Opts) (*zap.Logger, error) {
	if len(opts) == 0 {
		return zap.NewNop(), nil
	}

	op := opts[0]

	lvl, err := zap.ParseAtomicLevel(op.Lvl)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	return cfg.Build()
}
