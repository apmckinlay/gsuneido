// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_initLogFont()
		{
		suneidoLogFont = #(
			lfWidth: 0, lfEscapement: 0,
			fontPtSize: 9, lfCharSet: 0,
			lfUnderline: 0, lfQuality: 1,
			lfWeight: 400, lfOutPrecision: 3,
			lfHeight: -8, lfFaceName: 'Consolas',
			lfClipPrecision: 2, lfOrientation: 0,
			lfPitchAndFamily: 49, lfItalic: 0,
			lfStrikeOut: 0)
		mock = Mock(Hwnd)
		mock.When.suneidoLogFont().Return(suneidoLogFont)
		mock.When.initLogFont([anyArgs:]).CallThrough()

		// All default values, simply returns a copy of suneidoLogFont
		factor = 1
		size = weight = ''
		underline = italic = strikeout = false
		orientation = 0
		result = mock.
			initLogFont(factor, size, weight, underline, italic, strikeout, orientation)
		Assert(result equalsSet: suneidoLogFont)

		// "size" is not the default value "", calculates lfHeight
		expected = suneidoLogFont.Copy()
		expected.lfHeight = -12
		factor = 2
		size = 6
		result = mock.
			initLogFont(factor, size, weight, underline, italic, strikeout, orientation)
		Assert(result equalsSet: expected)

		// Adjust / test remaining arguments
		expected.lfUnderline = true
		expected.lfWeight = FW.BOLD
		expected.lfItalic = true
		expected.lfStrikeOut = true
		expected.lfWeight = FW.BOLD
		expected.lfEscapement = expected.lfOrientation = 900
		weight = FW.BOLD
		italic = strikeout = underline = true
		orientation = 900
		result = mock.
			initLogFont(factor, size, weight, underline, italic, strikeout, orientation)
		Assert(result equalsSet: expected)
		}
	}
