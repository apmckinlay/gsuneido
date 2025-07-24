// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// can optionally pass an ignore list or function
// e.g. Addon_speller: #(ignore: (apm))
/* or
c = class
	{
	New()
		{
		.list = QueryList('tables', 'table')
		}
	Call(word)
		{
		return .list.Has?(word)
		}
	}
ScintillaAddonsControl(Addon_speller: [ignore: c])
*/
ScintillaAddonForChanges
	{
	Ignore: ()
	styleLevel: 3
	Init()
		{
		.indic_word = .IndicatorIdx(level: .styleLevel)
		}

	Styling()
		{
		return [[level: .styleLevel, indicator: [INDIC.SQUIGGLE, fore: 0x0000ff]]]
		}

	getter_ignore?()
		{
		if String?(.Ignore)
			.Ignore = Global(.Ignore)

		ignore? = false
		if Object?(.Ignore)
			ignore? = .ignore_list
		else if Function?(.Ignore)
			ignore? = .Ignore
		else if Class?(.Ignore)
			ignore? = new .Ignore
		if ignore? is false
			throw "invalid Addon_speller ignore"
		return .ignore? = ignore?
		}

	ProcessChunk(text, pos)
		{
		if .GetReadonly() is 1 or not Speller.Open()
			return
		.ClearIndicator(.indic_word, pos, text.Size())
		typos = .collectTypos(pos, text, .ignore?)
		// Defer in order to allow other Addon's to add their indicators first.
		.Defer({ .indicateTypos(typos) })
		}

	collectTypos(pos, text, ignore?)
		{
		typos = Object()
		while text.Size() isnt idx = text.FindRx('[[:alpha:]]')
			{
			w = .check_word(pos + idx, typos, ignore?)
			pos = w.end
			text = text.AfterFirst(w.word)
			}
		return typos
		}

	check_word(pos, typos, ignore?)
		{
		w = .getCurrentWord(pos, ignore?)
		if w.ignore
			return w
		suggestions = Speller(w.word)
		if suggestions.NotEmpty?() and suggestions[0] isnt w.word.Lower()
			typos.Add([pos: w.org, len: w.len])
		return w
		}

	indicateTypos(typos)
		{
		for typo in typos
			if not .IndicatorAtPos?(typo.pos)
				.SetIndicator(.indic_word, typo.pos, typo.len)
		}

	ignore_list(word)
		{
		return .Ignore.Has?(word)
		}

	Set()
		{
		.nsuggestions = 0
		}

	getter_origMenu()
		{
		return .origMenu = .Parent.Context_Menu.Copy()
		}

	UpdateUI()
		{
		.clear_menu()
		if .GetReadonly() is 1
			.ClearIndicator(.indic_word)
		else
			{
			w = .getCurrentWord(.GetCurrentPos(), .ignore?)
			if not w.ignore
				.load_menu(w.word)
			}
		}

	getCurrentWord(pos, ignore?)
		{
		w = .wordBoundaries(pos)
		w.word = .GetRange(w.org, w.end)
		w.len = w.word.Size()
		w.ignore = .ignoreWord(w, ignore?)
		return w
		}

	// need this to be non-greedy, but still match contractions correctly
	// first section [[:alpha:]][a-z]*? non-greedaly matches any normal word or the first
	// part of a contraction, the second part (\>'[a-z]+?)? matches the contraction if it
	// exists
	word_pattern: "\<[[:alpha:]][a-z]*?(\>'[a-z]+?)?\>"
	// if a word contains one of these extra characters, DO NOT pass to the Speller
	// it can cause the Speller to end up out of sync
	specialChars: '[._]'
	// Limit based on the largest english word in a trusted dictionary
	maxWordSize: 45
	ignoreWord(wordOb, ignore?)
		{
		word = wordOb.word
		if wordOb.len > .maxWordSize or word !~ .word_pattern or word =~ .specialChars
			return true
		if wordOb.org > wordOb.end
			return true
		return ignore?(word)
		}

	// includes "extra" characters so things like xxx_yyy and www.url.com are one word
	WordChars: "'._0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	wordBoundaries(pos)
		{
		org = end = pos
		while .WordChars.Has?(.GetAt(org - 1))
			--org
		while .WordChars.Has?(.GetAt(end))
			++end
		word = .GetRange(org, end)
		if word[0] is "'"
			++org
		if word[-1] is "'"
			--end
		word = word.Trim()
		if word[-1] is '.'
			--end
		return [:org, :end]
		}

	clear_menu()
		{
		.nsuggestions = 0
		.Parent.Context_Menu = .origMenu.Copy()
		}

	load_menu(word)
		{
		suggestions = Speller(word)
		if suggestions.Empty?()
			{
			.clear_menu()
			return
			}
		menu = .Parent.Context_Menu = .origMenu.Copy()
		menu.Add("", at: 0)
		for suggested_word in suggestions.Reverse!()
			menu.Add('Change to ' $ suggested_word, at: 0)
		.nsuggestions = suggestions.Size()
		}

	nsuggestions: 0
	ContextMenuChoice(i)
		{
		if i >= .nsuggestions
			return
		w = .getCurrentWord(.GetCurrentPos(), .ignore?)
		if w.ignore
			return
		changeTo = .Parent.Context_Menu[i].RemovePrefix('Change to ')
		.SetSelect(w.org, w.end - w.org)
		.ReplaceSel(changeTo)
		.SetSelect(w.end + changeTo.Size() - w.len)
		}
	}
