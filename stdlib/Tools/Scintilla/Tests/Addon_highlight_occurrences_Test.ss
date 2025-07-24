// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_mark()
		{
		fn = Addon_highlight_occurrences.Addon_highlight_occurrences_mark?

		Assert(fn('') is: false)
		Assert(fn(' ') is: false)
		Assert(fn('\n') is: false)
		Assert(fn('\t') is: false)
		Assert(fn('word\n') is: false)
		Assert(fn('word'))
		Assert(fn('words with spaces'))
		Assert(fn('a'))
		}

	Test_mark_all_occurrences()
		{
		mock = Mock(Addon_highlight_occurrences)
		mock.When.SetIndicatorCurrent([anyArgs:]).Do({ })
		mock.When.IndicatorFillRange([anyArgs:]).Do({ })
		mock.When.mark_all_occurrences([anyArgs:]).CallThrough()
		text = 'This is test text, there will be a lot of random words ' $
			'some words will repeat. Like text, there, random. This test ' $
			'will also only find words using regex'

		.verifyMatches(mock, text, #words, #(49, 60, 135), true)
		.verifyMatches(mock, text, #words, #(49, 60, 135), false)
		.verifyMatches(mock, text, #is, #(2, 5, 107), true)
		.verifyMatches(mock, text, #is, #(5), false)
		.verifyMatches(mock, text, #This, #(0, 105), true)
		.verifyMatches(mock, text, #This, #(0, 105), false)
		.verifyMatches(mock, text, 'will repeat', #(66), true)
		.verifyMatches(mock, text, 'will repeat', #(66), false)
		.verifyMatches(mock, text, #this, #(), true)
		.verifyMatches(mock, text, #this, #(), false)
		.verifyMatches(mock, text, #a, #(33, 43, 75, 98, 120), true)
		.verifyMatches(mock, text, #a, #(33), false)
		}

	verifyMatches(mock, text, find, indices, selectedText?)
		{
		Assert(mock.mark_all_occurrences(text, find, selectedText?) is: indices.Size())
		mock.Verify.SetIndicatorCurrent(false)
		for index in indices
			{
			Assert(text[index .. index + find.Size()] is: find)
			mock.Verify.IndicatorFillRange(index, find.Size())
			}
		mock.Mock_calls.Delete(all:)
		}
	}