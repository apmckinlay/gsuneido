// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// abstract base class
// derived classes must define AutoComplete
ScintillaAddon
	{
	idleTime: 400 // ms = .4 sec
	listSize: 10 // number of items in pop up list
	New(@args)
		{
		super(@args)
		.timer = IdleTimer(.idleTime, .charAdded)
		}
	Init()
		{
		.AutoCSetMaxHeight(.listSize)
		.AutoCSetTypeSeparator('`'.Asc()) // default of '?' conflicts
		}
	CharAdded(c)
		{
		if c =~ '[._[:alnum:]]'
			.timer.Reset()
		else
			.timer.Kill()
		}
	Backspace()
		{
		.timer.Reset()
		}
	charAdded()
		{
		if false isnt (word = .GetCurrentReference()) and
			not .GetWordChars().Has?(.GetAt(.GetSelect().cpMin))
			{
			.AutoComplete(word)
			}
		}

	AutoShow(word, matches)
		{
		if matches.Size() is 0 or
			(matches.Size() is 1 and matches[0] is word)
			{
			.AutocCancel()
			return
			}
		if matches.Size() > 0
			.AutocShow(word, matches)
		}
	AutocShow(word, matches) // overridden by Addon_auto_complete_code
		{
		.SCIAutocShow(word.Size(), matches)
		}

	MatchesExcludingSelf(list, word, matches = false)
		{
		return .Matches(list, word $ ' ', matches)
		}
	// TODO ignore internal capitals and underscores (like LibLocateList)
	Matches(list, word, matches = false)
		{
		if matches is false
			matches = Object()
		from = list.BinarySearch(word)
		to = Min(from + 100/*= limit*/, list.BinarySearch(word.RightTrim() $ '~'))
		for (i = from; i < to; ++i)
			matches.Add(list[i])
		return matches
		}

	Destroy()
		{
		.timer.Kill()
		}
	}
