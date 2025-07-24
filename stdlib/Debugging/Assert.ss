// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(@args) /*usage: (condition, msg="") or (value matcher: expected, msg="")*/
		{
		if .oldStyle?(args)
			.assert(@args)
		else
			.assertThat(args)
		}
	oldStyle?(args)
		{
		members = args.Members()
		return members is #(0) or members is #(0,1) or members is #(0, msg)
		}
	assert(x, msg = "")
		{
		if x is false
			throw "Assert FAILED" $ Opt(': ', msg)
		if x isnt true
			throw "Assert FAILED: expected true, but it was " $ MatcherWas.DisplayValue(x) $
				Opt('\n(', msg, ')')
		}
	assertThat(args)
		{
		value = args.Extract(0)
		msg = args.Extract(#msg, "")
		if Type(value) is 'Block'
			value = Catch(value)
		matcher_name = args.Members()[0]
		matcher_args = args[matcher_name]
		matcher = Global("Matcher_" $ matcher_name)
		if not matcher.Match(value, matcher_args)
			throw "expected " $ matcher.Expected(matcher_args) $
				" but it " $ matcher.Actual(value, matcher_args) $
				Opt('\n(', msg, ')')
		}
	}