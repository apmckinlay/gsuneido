// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package annotations

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

// Registration is one signature entry: the exe's builtin declarations
// converted by builtinRegistrations, and methods defined in Suneido code
// like Objects and Strings supplied via TypeChecker.RegisterSignatures.
type Registration struct {
	Kind     string // "method" (default), "free", "static"
	Receiver string // methods: user-facing type name, see registrationReceiver
	Class    string // statics: global class name, e.g. "Database"
	Name     string
	Sig      string // same format as builtin declarations: "(params) :type"
}

// user-facing type names, unlike the interpreter's internal prefixes
// (ob, num, int, func). record folds to object like everywhere else.
var registrationReceiver = map[string]typealgebra.DynType{
	"string":   typealgebra.TString,
	"number":   typealgebra.TNumber,
	"object":   typealgebra.TObject,
	"record":   typealgebra.TObject,
	"date":     typealgebra.TDate,
	"class":    typealgebra.TClass,
	"function": typealgebra.TFunction,
	"sequence": typealgebra.TSequence,
}

func LoadRegistrations(regs []Registration) (Set, error) {
	return LoadLayered(nil, regs)
}

func LoadLayered(base, over []Registration) (Set, error) {
	set := Set{}
	type slot struct {
		name string
		idx  int
	}
	pos := map[string]slot{}
	layer := func(regs []Registration, label string) error {
		seen := map[string]bool{}
		for i, r := range regs {
			name, s, err := parseRegistration(r)
			if err != nil {
				return fmt.Errorf("%s[%d]: %v", label, i, err)
			}
			k := name + "\x00" + receiverKey(s.Receiver)
			if seen[k] {
				return fmt.Errorf("%s[%d]: duplicate registration %s",
					label, i, registrationLabel(r))
			}
			seen[k] = true
			if p, ok := pos[k]; ok { // an over entry replaces the base one
				set[p.name][p.idx] = s
				continue
			}
			pos[k] = slot{name, len(set[name])}
			set[name] = append(set[name], s)
		}
		return nil
	}
	if err := layer(base, "builtin"); err != nil {
		return nil, err
	}
	if err := layer(over, "signatures"); err != nil {
		return nil, err
	}
	return set, nil
}

func parseRegistration(r Registration) (string, Signature, error) {
	if r.Name == "" {
		return "", Signature{}, fmt.Errorf("missing name")
	}
	s, err := parseSig(r.Sig)
	if err != nil {
		return "", Signature{}, fmt.Errorf("%s: bad sig %q: %v",
			registrationLabel(r), r.Sig, err)
	}
	if s.Returns == nil {
		s.Returns = typealgebra.TUnknown
	}
	switch r.Kind {
	case "", "method":
		if r.Class != "" {
			return "", Signature{}, fmt.Errorf("%s: method must not have class",
				registrationLabel(r))
		}
		recv, ok := registrationReceiver[r.Receiver]
		if !ok {
			return "", Signature{}, fmt.Errorf("%s: unknown receiver %q",
				registrationLabel(r), r.Receiver)
		}
		s.Receiver = recv
		return r.Name, s, nil
	case "free":
		if r.Receiver != "" || r.Class != "" {
			return "", Signature{}, fmt.Errorf(
				"%s: free function must not have receiver or class",
				registrationLabel(r))
		}
		return r.Name, s, nil
	case "static":
		if r.Receiver != "" {
			return "", Signature{}, fmt.Errorf("%s: static must not have receiver",
				registrationLabel(r))
		}
		if r.Class == "" {
			return "", Signature{}, fmt.Errorf("%s: static missing class",
				registrationLabel(r))
		}
		return r.Class + "." + r.Name, s, nil
	}
	return "", Signature{}, fmt.Errorf("%s: unknown kind %q",
		registrationLabel(r), r.Kind)
}

func registrationLabel(r Registration) string {
	switch {
	case r.Class != "":
		return r.Class + "." + r.Name
	case r.Receiver != "":
		return r.Receiver + "." + r.Name
	}
	return r.Name
}
