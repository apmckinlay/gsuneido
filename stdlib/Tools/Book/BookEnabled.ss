// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.

// -- always hides
// book options with more slashes "win"
// 	e.g. -/A/B ... +/A will leave B disabled
// otherwise LAST +/- one wins

class
	{
	CallClass(book, option)
		{
		enabled = .Enabled(book, option)
		if enabled is "not found"
			return false

		// must handle 'hidden' options even for default user since
		// option may be hidden because library that defines the option is
		// not being used (ex. CAD vc US Payroll options).
		// Must check for default after hidden options are checked.
		return Suneido.User isnt 'default' or enabled is "hidden" ? enabled : true
		}

	Enabled(book, option)
		{
		if option is "/Cover" or option is "/Contents"
			return true

		try
			options = .options(book)
		catch (e)
			{
			if not e.Has?("can't find") // book has no BookOptions record defined
				SuneidoLog('ERROR: (CAUGHT) ' $ e, calls:, caughtMsg: 'page not enabled')
			return "not found"
			}

		return .enabled(option, options)
		}

	enabled(option, options)
		{
		enabled = false
		specific = function (option)
			{ return option.Tr("^/").Size() }
		most_specific = 0
		for book_option in options.Copy()
			{
			if book_option.Prefix?('--') and
				(option.Suffix?(book_option[2..]) or
				option.Prefix?(book_option[2..] $ '/'))
				return 'hidden'
			if false is option.Has?(book_option[1..]) or
				specific(book_option) < most_specific
				continue
			most_specific = specific(book_option)
			if book_option.Prefix?('+')
				enabled = true
			else // book_option.Prefix?('-')
				enabled = false
			}
		return enabled
		}
	options(book) // overridden by test
		{
		fn = Global(book.Capitalize() $ "_BookOptions")
		return fn()
		}
	}