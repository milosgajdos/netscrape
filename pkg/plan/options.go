package plan

// Options are space options.
type Options struct {
	Origin Origin
}

// Option configures Options.
type Option func(*Options)

// WithOrigin configures Origin option.
func WithOrigin(org Origin) Option {
	return func(o *Options) {
		o.Origin = org
	}
}
