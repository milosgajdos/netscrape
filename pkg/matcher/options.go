package matcher

// Options are filter options.
type Options struct {
	Types      MatchFunc
	Names      MatchFunc
	Groups     MatchFunc
	Versions   MatchFunc
	Kinds      MatchFunc
	Namespaces MatchFunc
	Weights    MatchFunc
	Attrs      MatchFunc
}

// Option configures Options.
type Option func(*Options)

func WithTypes(t MatchFunc) Option {
	return func(o *Options) {
		o.Types = t
	}
}

func WithNames(n MatchFunc) Option {
	return func(o *Options) {
		o.Names = n
	}
}

func WithGroups(g MatchFunc) Option {
	return func(o *Options) {
		o.Groups = g
	}
}

func WithVersions(v MatchFunc) Option {
	return func(o *Options) {
		o.Versions = v
	}
}

func WithKinds(k MatchFunc) Option {
	return func(o *Options) {
		o.Kinds = k
	}
}

func WithNamespaces(n MatchFunc) Option {
	return func(o *Options) {
		o.Namespaces = n
	}
}

func WithWeights(w MatchFunc) Option {
	return func(o *Options) {
		o.Weights = w
	}
}

func WithAttrs(a MatchFunc) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}
