// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package diagnostics

import "fmt"

type Level string

const (
	LevelOff   Level = "off"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func ParseLevel(s string) (Level, error) {
	switch l := Level(s); l {
	case LevelOff, LevelWarn, LevelError:
		return l, nil
	default:
		return "", fmt.Errorf("invalid level %q (want off|warn|error)", s)
	}
}

type Flag int

const (
	FlagNone Flag = iota
	FlagStrictStringConcat
	FlagStrictCrossTypeCompares
)

func (f Flag) String() string {
	switch f {
	case FlagStrictStringConcat:
		return "strictStringConcat"
	case FlagStrictCrossTypeCompares:
		return "strictCrossTypeCompares"
	}
	return ""
}

type Config struct {
	StrictStringConcat      Level
	StrictCrossTypeCompares Level
}

func DefaultConfig() Config {
	return Config{
		StrictStringConcat:      LevelError,
		StrictCrossTypeCompares: LevelWarn,
	}
}
