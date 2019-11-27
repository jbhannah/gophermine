package console

import (
	tm "github.com/buger/goterm"
	log "github.com/sirupsen/logrus"
)

type TermFormatter struct {
	*log.TextFormatter
}

func (tf *TermFormatter) Format(entry *log.Entry) ([]byte, error) {
	defer print("> ")

	bytes, err := tf.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	reset := []byte(tm.RESET_LINE)
	bytes = append(reset, bytes...)
	bytes = append(bytes, []byte("> ")...)

	return bytes, nil
}
