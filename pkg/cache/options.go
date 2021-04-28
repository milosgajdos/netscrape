package cache

// Options are cache options.
type Options struct {
	// Upsert option for upsert operations.
	Upsert bool
}

// Option configures Options.
type Option func(*Options)

// WithUpsert enables cache upsert.
func WithUpsert() Option {
	return func(o *Options) {
		o.Upsert = true
	}
}
