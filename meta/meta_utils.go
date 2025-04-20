package meta

type MergeRowOptions struct {
	rows     []Row
	prefixes []string
}

func NewMergeRowOptions() *MergeRowOptions {
	return &MergeRowOptions{
		rows:     []Row{},
		prefixes: []string{},
	}
}

type MergeRowOpts func(*MergeRowOptions) error

func WithMergeRow(row Row, prefix string) MergeRowOpts {
	return func(o *MergeRowOptions) error {
		o.rows = append(o.rows, row)
		o.prefixes = append(o.prefixes, prefix)

		return nil
	}
}

func MergeRows2(opts ...MergeRowOpts) Row {
	mro := NewMergeRowOptions()
	final := make(map[string]interface{})

	for _, opt := range opts {
		opt(mro)
	}

	for idx, row := range mro.rows {
		if len(mro.prefixes[idx]) > 0 {
			for key, value := range row {
				withPrefix := mro.prefixes[idx] + "." + key
				final[withPrefix] = value
			}
		} else {
			for key, value := range row {
				final[key] = value
			}
		}
	}

	return final
}
