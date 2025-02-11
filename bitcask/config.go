package bitcask

type Options struct {
	dir string
}

type Config func(*Options) error

func DefaultOptions() *Options {
	return &Options{
		dir: ".",
	}
}

func WithDir(dir string) Config {
	return func(o *Options) error {
		o.dir = dir
		return nil
	}
}
