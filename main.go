package mute

import (
	"encoding/json"
	"fmt"
)

type format string

const (
	// JSON dictates that an Event should be represented as JSON.
	JSON format = "json"
	// Text dictates that an Event should be represented as a formatted string.
	Text format = "text"
)

// Init will return a new logger with the given Route.
func Init(r ...Route) *logger {
	var rs []Route
	rs = append(rs, r...)

	l := logger{
		Routes: rs,
	}

	return &l
}

// Event describes an event and can be sent by a Logger.
type Event struct {
	// Message is the short text presented by an event.
	Message string
	// Data is a map of key-value pairs that can add context to an event.
	Data map[string]string
}

// Route describes a strategy that a logger will use to deliver an event.
type Route struct {
	Memory *[]string
	File   string
	Format format
}

// deliver will handle the conversion and storage or delivery of an event.
func (r Route) deliver(e Event) error {
	message, err := convert(e, r.Format)
	if err != nil {
		return err
	}

	if r.Memory != nil {
		*r.Memory = append(*r.Memory, message)
	}

	if r.File != "" {
		panic("file delivery is not implemented")
	}

	return nil
}

// Send will deliver the given events to all routes associated with this logger.
func (l *logger) Send(e ...Event) error {
	for _, r := range l.Routes {
		for _, e := range e {
			err := r.deliver(e)
			if err != nil {
				return err
				break
			}
		}
	}

	return nil
}

// logger is a type that contains a set of routes. A logger can receive events
// via the logger.Send(...Event) method, and will deliver the received events to
// all registered routes.
type logger struct {
	Routes []Route
}

// convert facilitates converting events to different formats.
func convert(e Event, format format) (string, error) {
	if format == JSON {
		result, err := toJSON(e)
		if err != nil {
			return "", err
		}

		return result, nil
	}

	if format == Text {
		return toText(e), nil
	}

	return "", fmt.Errorf(
		"Attempted to convert an invalid log type, only accept \"json\" and \"text\" values: %v", format,
	)
}

// toJson facilitates marshalling an event to JSON.
func toJSON(e Event) (string, error) {
	result, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf(
			"Failed to convert an event to JSON: \"%v\"", e.Message,
		)
	}

	return string(result), nil
}

// toText facilitates marshalling an event to a formatted string.
func toText(e Event) string {
	result := e.Message

	if len(e.Data) > 0 {
		for k, v := range e.Data {
			result += fmt.Sprintf(" [%v: %v]", k, v)
		}
	}

	return result
}
