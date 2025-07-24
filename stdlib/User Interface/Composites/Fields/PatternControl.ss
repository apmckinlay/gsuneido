// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// contributions from Jean Charles jessihash@wanadoo.fr
// patterns consist of:
//	a - letter converted to lower case
//	A - letter converted to upper case
//	# - digit
//  > - digit or letter converted to upper case
//  < - digit or letter converted to lower case
//	^c - a literal c
//	several choices can be supplied separated by |
//	other characters are taken literally
HandleEnterControl
	{
	New(pattern, width = 10, status = "", .mandatory = false, readonly = false,
		hidden = false, tabover = false)
		{
		super(:width, :status, :readonly, :hidden, :tabover)
		.pattern = pattern.Split('|')
		.regexp = .regexps(.pattern)
		.SubClass()
		}

	regexps(pattern)
		{
		regexp = Object()
		for p in pattern
			regexp.Add(.convertPattern(p))
		return regexp
		}

	convertPattern(pattern)
		{
		regexp = "^"
		for (i = 0; (c = pattern[i]) isnt ""; ++i)
			{
			if c is 'a' or c is 'A'
				elem = "[a-zA-Z]"
			else if c is '<' or c is '>'
				elem = "[a-zA-Z0-9]"
			else if c is '#'
				elem = "[0-9]"
			else if c is '^' // escape
				elem = '[' $ pattern[++i] $ ']?'
			else // turn off all other special chars with []
				elem = '[' $ c $ ']?'
			regexp $= elem
			}
		return regexp $ '$'
		}

	Valid?()
		{
		return .valid(.Get(), .mandatory, .regexp)
		}

	valid(s, mandatory, regexp)
		{
		valid = false
		if s is "" and not mandatory
			valid = true
		for i in regexp.Members()
			if s =~ regexp[i]
				valid = true
		return valid
		}

	ValidData?(@args)
		{
		value = args[0]
		// handle when this control is used directly
		pattern = args.GetDefault('pattern', args.GetDefault(1, ''))
		regexp = .regexps(.Pattern(:pattern).Split('|'))
		return .valid(value, args.GetDefault('mandatory', false), regexp)
		}

	Pattern(pattern = '')
		{
		return pattern
		}

	KillFocus()
		{
		s = .Get()
		for i in .pattern.Members()
			{
			if false isnt t = .match(s, .pattern[i], .regexp[i])
				{
				dirty? = .Dirty?()
				.Set(t)
				.Dirty?(dirty?)
				return
				}
			}
		}

	match(s, pat, regex)
		{
		if s !~ regex
			return false
		result = ""
		si = 0
		for (pi = 0; (pc = pat[pi]) isnt ""; ++pi)
			{
			sc = s[si++]
			if .lowerChar?(pc)
				result $= sc.Lower()
			else if .upperChar?(pc)
				result $= sc.Upper()
			else if .directTranslate?(pc, sc)
				result $= sc
			else if .literalMatch?(pc, sc, pat, pi)
				result $= pat[++pi]
			else // fill in literal char
				{
				if pc is '^'
					pc = pat[++pi]
				result $= pc
				--si // undo increment
				}
			}
		return result
		}

	lowerChar?(pc) { return pc is 'a' or pc is '<' }
	upperChar?(pc) { return pc is 'A' or pc is '>' }
	directTranslate?(pc, sc) { return pc is '#' or (pc is sc and pc isnt '^') }
	literalMatch?(pc, sc, pat, pi) { return pc is '^' and sc is pat[pi + 1] }
	}