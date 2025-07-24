// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Parse()
		{
		parse = ScintillaRichEditorHelper.Parse
		text = '<span style="font-weight:normal;font-style:normal">A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class </span>' $
			'<span style="font-weight:normal;font-style:italic">String</span>'
		result = parse(text)
		Assert(result isSize: 2)
		Assert(result.s is: 'A Test Class String')
		Assert(result.styles is: #(
			#(from: #(ch: 0, line: 0), to: #(ch: 2, line: 0), txt: 'A ',
				style: #(bold: false, italic: false, underline: false, strikeout: false)),
			#(from: #(ch: 2, line: 0), to: #(ch: 7, line: 0), txt: 'Test ',
				style: #(bold:, italic: false, underline: false, strikeout: false)),
			#(from: #(ch: 7, line: 0), to: #(ch: 13, line: 0), txt: 'Class ',
				style: #(bold: false, italic: false, underline: false, strikeout:)),
			#(from: #(ch: 13, line: 0), to: #(ch: 19, line: 0), txt: 'String',
				style: #(bold: false, italic:, underline: false, strikeout: false))))

		plainText = 'hello world\r\ntest string'
		result = parse(plainText)
		Assert(result.s is: 'hello world\r\ntest string')
		Assert(result.styles is: #())
		}

	Test_buildStyled()
		{
		buildStyled = ScintillaRichEditorHelper.ScintillaRichEditorHelper_buildStyled
		child1 = FakeObject(Text: 'Hello ',
			Attributes: #(style: 'font-weight:normal;font-style:normal'))
		child2 = FakeObject(Text: 'Worl\r\nd',
			Attributes: #(style:
				'font-weight:bold;font-style:italic;text-decoration:underline'))


		styleObject = Object()
		Assert(buildStyled(child1, styleObject, Object(line: 0, ch: 0)) is: 'Hello ')
		Assert(styleObject[0]
			is: #(from: #(line: 0, ch: 0), to: #(line: 0, ch: 6), txt: 'Hello '
				style: #(bold: false, italic: false, underline: false, strikeout: false)))

		Assert(buildStyled(child2, styleObject, Object(line: 0, ch: 7)) is: 'Worl\r\nd')
		Assert(styleObject[1]
			is: #(from: #(line: 0, ch: 7), to: #(line: 1, ch: 1), txt: 'Worl\r\nd'
				style: #(bold:, italic:, underline:, strikeout: false)))
		}

	Test_TrimStyledText()
		{
		trim = ScintillaRichEditorHelper.TrimStyledText
		plainText = '     this is plain text\r\n   '
		parsed = ScintillaRichEditorHelper.Parse(plainText)
		Assert({ trim(parsed.s, parsed.styles) }
			throws: 'styles expected but not found')

		stylesOb = Object(
			Object(from: #(ch: 0, line: 0), to: #(ch: 3, line: 1),
				style: #(bold: false, underline: false, strikeout: false, italic: false),
				txt: "     this is plain text\r\n   "))
		result = trim(parsed.s, stylesOb)
		Assert(result.s is: 'this is plain text')

		styledText = '<span style="font-weight:normal;font-style:normal">    A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class </span>' $
			'<span style="font-weight:normal;font-style:normal">' $
				'&lt;encoding&gt; </span>' $
			'<span style="font-weight:normal;font-style:italic">String<br /><br /></span>'
		parsed = ScintillaRichEditorHelper.Parse(styledText)
		result = trim(parsed.s, parsed.styles)
		Assert(result isSize: 2)
		Assert(result.s is: 'A Test Class <encoding> String')
		Assert(result.styles is: #(
			#(from: #(ch: 0, line: 0), to: #(ch: 2, line: 0), txt: "A ",
				style: #(bold: false, italic: false, underline: false, strikeout: false)),
			#(from: #(ch: 2, line: 0), to: #(ch: 7, line: 0), txt: 'Test ',
				style: #(bold:, italic: false, underline: false, strikeout: false)),
			#(from: #(ch: 7, line: 0), to: #(ch: 13, line: 0), txt: 'Class '
				style: #(bold: false, italic: false, underline: false, strikeout:)),
			#(from: #(ch: 13, line: 0), to: #(ch: 24, line: 0), txt: '<encoding> '
				style: #(bold: false, italic: false, underline: false, strikeout: false)),
			#(from: #(ch: 24, line: 0), to: #(ch: 30, line: 0), txt: 'String'
				style: #(bold: false, italic:, underline: false, strikeout: false))))
		}

	Test_Build()
		{
		build = ScintillaRichEditorHelper.Build
		parse = ScintillaRichEditorHelper.Parse

		text = '<span style="font-weight:normal;font-style:normal">A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class </span>' $
			'<span style="font-weight:normal;font-style:italic">String</span>'
		parsed = parse(text)
		rebuilt = build(parsed.s, parsed.styles)
		Assert(rebuilt is: text)

		textToTrim = '<span style="font-weight:normal;font-style:normal">     A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class </span>' $
			'<span style="font-weight:normal;font-style:italic">String<br /><br /></span>'
		parsed = parse(textToTrim)
		trimmed = ScintillaRichEditorHelper.TrimStyledText(parsed.s, parsed.styles)
		result = build(trimmed.s, trimmed.styles)
		Assert(result is: text)

		removeSpan = '<span style="font-weight:normal;font-style:normal">     A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class </span>' $
			'<span style="font-weight:normal;font-style:italic"><br /><br /></span>'
		parsed = parse(removeSpan)
		trimmed = ScintillaRichEditorHelper.TrimStyledText(parsed.s, parsed.styles)
		result = build(trimmed.s, trimmed.styles)
		Assert(result is: '<span style="font-weight:normal;font-style:normal">A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class</span>')
		}
	}
