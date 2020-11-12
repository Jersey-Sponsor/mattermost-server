// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	ActionKey = "-action"
)

// responsef creates an ephemeral command response using printf syntax.
func responsef(format string, args ...interface{}) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf(format, args...),
		Type:         model.POST_DEFAULT,
	}
}

// parseNamedArgs parses a command string into a map of arguments. It is assumed the
// command string is of the form `<action> --arg1 value1 ...` Supports empty values.
// Arg names are limited to [0-9a-zA-Z_].
func parseNamedArgs(cmd string) map[string]string {
	m := make(map[string]string)

	split := strings.Fields(cmd)

	// check for optional action
	if len(split) >= 2 && !strings.HasPrefix(split[1], "--") {
		m[ActionKey] = split[1] // prefix with hyphen to avoid collision with arg named "action"
	}

	for i := 0; i < len(split); i++ {
		if !strings.HasPrefix(split[i], "--") {
			continue
		}
		var val string
		arg := trimSpaceAndQuotes(strings.Trim(split[i], "-"))
		if i < len(split)-1 && !strings.HasPrefix(split[i+1], "--") {
			val = trimSpaceAndQuotes(split[i+1])
		}
		if arg != "" {
			m[arg] = val
		}
	}
	return m
}

func trimSpaceAndQuotes(s string) string {
	trimmed := strings.TrimSpace(s)
	trimmed = strings.TrimPrefix(trimmed, "\"")
	trimmed = strings.TrimPrefix(trimmed, "'")
	trimmed = strings.TrimSuffix(trimmed, "\"")
	trimmed = strings.TrimSuffix(trimmed, "'")
	return trimmed
}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "1", "t", "true", "yes", "y":
		return true, nil
	case "0", "f", "false", "no", "n":
		return false, nil
	}
	return false, fmt.Errorf("cannot parse '%s' as a boolean", s)
}