// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	indic_word: false
	styleLevel: 96
	Init()
		{
		if .indic_word isnt false
			.ClearIndicator(.indic_word)
		.indic_word = .IndicatorIdx(level: .styleLevel)
		.IndicSetAlpha(.indic_word, 10) /*= 0=transparent .. 255=opaque */
		.IndicSetOutlineAlpha(.indic_word, 150) /*= 0=transparent .. 255=opaque */
		}

	Styling()
		{
		return [
			[level: .styleLevel,
				indicator: [INDIC.STRAIGHTBOX, fore: .GetSchemeColor('occurrence')]]]
		}

	UpdateUI()
		{
		.ClearIndicator(.indic_word)
		selected = .GetSelText()
		textSelected? = selected isnt ''
		if not textSelected?
			selected = .GetCurrentWord()
		n = .mark?(selected) ? .mark_all_occurrences(.Get(), selected, textSelected?) : 0
		if n < 2
			.ClearIndicator(.indic_word)
		}

	mark?(find)
		{ return find isnt '' and not find.Has?('\n') and not find.White?() }

	// DOES NOT USE REGEX ON PURPOSE:
	// When evaluated, regex will throw errors if formatted wrong.
	// As this code runs while highlighting text, it is not safe to use with Regex.
	// ie: 'Find' =~ '())'
	mark_all_occurrences(text, find, textSelected?)
		{
		.SetIndicatorCurrent(.indic_word)
		searched = n = 0
		strLen = find.Size()
		while text.Size() isnt idx = text.Find(find)
			{
			currentIdx = searched + idx

			if textSelected? or .hasWordBoundaries?(text, idx, strLen)
				{
				.IndicatorFillRange(currentIdx, strLen)
				n++
				}
			searched += idx + strLen
			text = text[idx + strLen ..]
			}
		return n
		}

	hasWordBoundaries?(text, idx, strLen)
		{
		boundaryCharBefore = idx is 0 or text[idx - 1] =~ '\W'
		boundaryCharAfter = text[idx + strLen] =~ '\W'
		return boundaryCharBefore and boundaryCharAfter
		}
	}
