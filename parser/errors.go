package parser

import "bytes"

type multierr []error

func (m multierr) err() error {
	if len(m) == 0 {
		return nil
	}

	return m
}

func (m multierr) Error() string {
	var buf bytes.Buffer

	for _, er := range m {
		buf.WriteString(er.Error() + "\n")
	}

	return buf.String()
}
