// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// Simulates the Scintilla behavior using a passed in string, rather than Editor.Get
	scintilla: class
		{
		Name: Editor
		New() { }

		Set(text)
			{ .text = text }

		Get()
			{ return .text }

		GetAt(pos)
			{
			// Mimicking the real scintilla control
			return pos < 0 or pos >= .GetTextLength()
				? '\x00'
				: .GetRange(pos, pos + 1)
			}

		GetRange(start, end)
			{
			start = .intoRange(start)
			end = .intoRange(end)
			s = ""
			for (i = start; i < end; i += .chunkSize)
				s $= .text[i .. Min(i + .chunkSize, end)]
			return s
			}

		intoRange(i)
			{ return Max(0, Min(i, .GetTextLength())) }

		chunkSize: 8192
		GetTextLength()
			{ return .text.Size() }
		}

	Test_getCurrentWord_basic()
		{
		sci = new .scintilla
		.addons = AddonManager(sci, [Addon_speller:])
		sci.Set('here is some text')

		.assertWord(pos: 0, word: 'here', org: 0, end: 4)
		.assertWord(pos: 9, word: 'some', org: 8, end: 12)
		.assertWord(pos: 17, word: 'text', org: 13, end: 17)
		}

	Test_getCurrentWord_complex()
		{
		sci = new .scintilla
		.addons = AddonManager(sci, [Addon_speller:])
		sci.Set(`here's a little more complex text
			including tabs, multi-line, www.urls.com, etc.
			Pneumonoultramicroscopicsilicovolcanoconiosis ` $ `a`.Repeat(46))

		.assertWord(pos: 0, word: `here's`, org: 0, end: 6)
		.assertWord(pos: 7, word: 'a', org: 7, end: 8)
		.assertWord(pos: 9, word: 'little', org: 9, end: 15)
		.assertWord(pos: 16, word: 'more', org: 16, end: 20)
		.assertWord(pos: 21, word: 'complex', org: 21, end: 28)
		.assertWord(pos: 29, word: 'text', org: 29, end: 33)

		// incremented by 5, 1 newline, 3 tabs, 1 to reach first char
		.assertWord(pos: 38, word: 'including', org: 38, end: 47)
		.assertWord(pos: 48, word: 'tabs', org: 48, end: 52)

		// hyphens do not get treated as one word
		.assertWord(pos: 54, word: 'multi', org: 54, end: 59)
		.assertWord(pos: 60, word: 'line', org: 60, end: 64)

		// urls are accurately calculate where they end, but are set to be ignored
		.assertWord(pos: 66, word: 'www.urls.com', org: 66, end: 78, ignore:)

		.assertWord(.addons, pos: 80, word: 'etc', org: 80, end: 83)

		// Word is at the limit of the largest word size
		// This is the largest English word in a trusted dictionary
		.assertWord(pos: 90, word: 'Pneumonoultramicroscopicsilicovolcanoconiosis',
			org: 89, end: 134)

		// Words larger than a reasonable amount of characters are treated as "invalid"
		// We still calculate where these "words" start / end
		.assertWord(pos: 136, word: `a`.Repeat(46), org: 135, end: 181, ignore:)
		}

	assertWord(pos, word, org, end, ignore = false)
		{
		wordOb = .addons.Collect(#Addon_speller_getCurrentWord, pos, {|unused| false})[0]
		Assert(wordOb.word is: word)
		Assert(wordOb.org is: org)
		Assert(wordOb.end is: end)
		Assert(wordOb.ignore is: ignore)
		}

	Test_collectTypos()
		{
		sci = new .scintilla
		addons = AddonManager(sci, [Addon_speller:])

		// NOTE: We override the Speller, so ANY errors we want to catch need to be
		// included in the below object
		typoOb = #(sentance, thise, eror, grammer, errirs, sintax)
		spellerSpy = .SpyOn(Speller)
		spellerSpy.Return(['typo'], when: { |word| typoOb.Has?(word) })
		spellerSpy.Return([], when: { |word| not typoOb.Has?(word) })
		spellerCalls = spellerSpy.CallLogs()

		ignore? = {|unused| false}
		sci.Set('here is  a   sentance with  one error.')
		typos = addons.Collect(#Addon_speller_collectTypos, 0, sci.Get(), ignore?)[0]
		Assert(spellerCalls isSize: 7)
		Assert(typos isSize: 1)
		Assert(typos[0].len is: 8)
		Assert(typos[0].pos is: 13)

		// Sentance starts and ends with a error
		spellerCalls.Delete(all:)
		sci.Set('thise sentance starts and ends with an eror')
		typos = addons.Collect(#Addon_speller_collectTypos, 0, sci.Get(), ignore?)[0]
		Assert(spellerCalls isSize: 8)
		Assert(typos isSize: 3)
		Assert(typos[0].len is: 5)
		Assert(typos[0].pos is: 0)
		Assert(typos[1].len is: 8)
		Assert(typos[1].pos is: 6)
		Assert(typos[2].len is: 4)
		Assert(typos[2].pos is: 39)

		// Similar test as before, but starts and ends with special characters
		// .<word> is not considered a word, for that reason, .thise is not considered a
		// typo OR passed to the speller
		spellerCalls.Delete(all:)
		sci.Set('.thise sentance starts and ends with an eror.')
		typos = addons.Collect(#Addon_speller_collectTypos, 0, sci.Get(), ignore?)[0]
		Assert(spellerCalls isSize: 7)
		Assert(typos isSize: 2)
		Assert(typos[0].len is: 8)
		Assert(typos[0].pos is: 7)
		Assert(typos[1].len is: 4)
		Assert(typos[1].pos is: 40)

		// Testing the ignore list properlly handles errors
		spellerCalls.Delete(all:)
		ignore? = {|word| word is 'sentance'}
		sci.Set(`this sentance isn't perfect (but we ignore the 1 typo)`)
		typos = addons.Collect(#Addon_speller_collectTypos, 0, sci.Get(), ignore?)[0]
		Assert(spellerCalls isSize: 8)
		Assert(typos isSize: 0)

		// Testing some of the more obscure errors / sitautions
		// n0t is not passed to Speller as it does not meet our word regex criteria
		// and is ignored as a result
		spellerCalls.Delete(all:)
		ignore? = {|unused| false}
		sci.Set(`www.testing.urls and "numbers" and n0t re@l 'words'.`)
		typos = addons.Collect(#Addon_speller_collectTypos, 0, sci.Get(), ignore?)[0]
		Assert(spellerCalls isSize: 6)
		Assert(typos isSize: 0)

		// Complex test with new lines, special characters and joins
		spellerCalls.Delete(all:)
		sci.Set(`here is a string,with weird grammer/errirs,
			new lines/a lot of spacing. Just general odd formatting
grammer\errirs,errirs

			sintax,oddities`)
		typos = addons.Collect(#Addon_speller_collectTypos, 0, sci.Get(), ignore?)[0]
		Assert(spellerCalls isSize: 23)
		Assert(typos isSize: 6)
		}
	}
