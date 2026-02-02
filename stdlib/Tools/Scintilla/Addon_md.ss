// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Init()
		{
		.SetILexer(0, CreateLexer("markdown"))
		defaultBack = .GetSchemeColor("defaultBack")
		defaultFore = .GetSchemeColor("defaultFore")

		.DefineStyle(0, defaultFore, back: defaultBack)
		.DefineStyle(1/*=BLOCK_SYMBOL*/, .GetSchemeColor("keyword"), back: defaultBack)
		.DefineStyle(2/*=INLINE_SYMBOL*/, .GetSchemeColor("string"), back: defaultBack)
		.DefineStyle(3/*=HTML*/, .GetSchemeColor("number"), back: defaultBack)

		.DefineStyle(SC.STYLE_DEFAULT, defaultFore, back: defaultBack)

		.patterns = Object(
			// Headers
			Object(`\A\s*#+\s+`, .styles.BLOCK_SYMBOL, newline:),
			// Horizontal Rule
			Object(`\A[-_*]+$`, .styles.BLOCK_SYMBOL, newline:),
			// Blockquote
			Object(`\A\s*>\s+`, .styles.BLOCK_SYMBOL, newline:),
			// List Item
			Object(`\A\s*([-+*]|(\d+))\s+`, .styles.BLOCK_SYMBOL, newline:),
			// Fenced Code
			Object(`\A\s*` $  '```' $ `\s*\w*\s*$`, .styles.BLOCK_SYMBOL, newline:),
			// inline
			// Bold or Italic
			Object(`\A[*_]+`, .styles.INLINE_SYMBOL)
			// Inline Code
			Object('\\A`', .styles.INLINE_SYMBOL)
			// LInks
			Object(`\A!?\[[^\]]*?\]\([^)]*?\)`, .styles.INLINE_SYMBOL)

			// html
			Object(`\A<.*+?>`, .styles.HTML)
			)

		.SetWrapMode(SC.WRAP_WORD)
		}

	styles: (
		DEFAULT:		'\x00',
		BLOCK_SYMBOL:	'\x01',
		INLINE_SYMBOL:	'\x02',
		HTML:			'\x03',
		)

	Style(from/*unused*/, to/*unused*/)
		{
		.setStyles(.Hwnd, .style(.Get(), .patterns))
		}

	style(src, patterns)
		{
		styles = ''
		pos = 0
		newline = true
		while pos < src.Size()
			{
			found = false
			subStr = src[pos..]
			for pattern in patterns
				{
				if false isnt match = .matchPattern(pattern, subStr, newline)
					{
					styles $= pattern[1].Repeat(match[0][1])
					pos += match[0][1]
					found = true
					break
					}
				}
			if not found
				{
				res = src.Find1of('\n*_![`', pos: pos + 1)
				newline = src[res] is '\n'
				next = res isnt src.Size() and newline
					? res + 1
					: res
				styles $= .styles.DEFAULT.Repeat(next - pos)
				pos = next
				}
			}
		return styles
		}

	matchPattern(pattern, s, newline)
		{
		if newline is false and pattern.GetDefault(#newline, false) is true
			return false

		return s.Match(pattern[0])
		}

	chunk: 10000
	setStyles(hwnd, styles)
		{
		SendMessage(hwnd, SCI.STARTSTYLING, 0, 0x1f/*=unused*/)
		for (i = 0; i < styles.Size(); i += .chunk)
			SendMessageTextIn(hwnd, SCI.SETSTYLINGEX,
				Min(.chunk, styles.Size() - i), styles[i::.chunk])
		}
	}
