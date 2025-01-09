package logging

import (
	"errors"
	"log/slog"
	"net/url"
	"os"
	"strings"
	// slogsyslog "github.com/samber/slog-syslog/v2"
)

func GetLogger(sink string, options *slog.HandlerOptions) (*slog.Logger, error) {
	handler, err := GetHandler(sink, options)
	if err != nil {
		return nil, err
	}

	return slog.New(handler), nil
}

func GetHandler(sink string, options *slog.HandlerOptions) (slog.Handler, error) {
	index := strings.IndexRune(sink, ':')

	var params string
	if index > -1 {
		params = sink[index+1:]
		sink = sink[:index]
	}

	switch sink {
	case "console":
		u, _ := url.Parse(params)
		prefix := u.Query().Get("prefix")
		if prefix != "" {
			return slog.NewTextHandler(LinePrefixWriter{Writer: os.Stdout, Prefix: prefix}, options), nil
		}

		return slog.NewTextHandler(os.Stdout, options), nil
	case "file":
		file, err := os.OpenFile(params, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}

		return slog.NewTextHandler(file, options), nil
	// case "syslog":
	// 	writer, err := net.Dial("udp", params)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return slogsyslog.Option{
	// 		Writer: writer,

	// 		Level:       options.Level,
	// 		AddSource:   options.AddSource,
	// 		ReplaceAttr: options.ReplaceAttr,
	// 	}.NewSyslogHandler(), nil
	default:
		return nil, errors.New("unknown sink")
	}
}
