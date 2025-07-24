// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_setFontStyle()
		{
		// uses bitFlags for style
		// bit 1 = Bold
		// bit 2 = Italic
		// bit 3 = Underline
		// bit 4 = Strikeout/through
		setFontStyle = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_setFontStyle
		Assert(setFontStyle(SCIRT.NORMAL)
			is: 'span style="font-weight:normal;font-style:normal"')
		Assert(setFontStyle(SCIRT.BOLD)
			is: 'span style="font-weight:bold;font-style:normal"')
		Assert(setFontStyle(SCIRT.ITALIC)
			is: 'span style="font-weight:normal;font-style:italic"')
		Assert(setFontStyle(SCIRT.BOLD_ITALIC)
			is: 'span style="font-weight:bold;font-style:italic"')
		Assert(setFontStyle(SCIRT.UNDERLINE)
			is: 'span style="font-weight:normal;font-style:normal;' $
				'text-decoration:underline"')
		}

	Test_parseStyledText()
		{
		parseStyledText = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_parseStyledText
		text = .buildText(Object(
			Object("A ", SCIRT.NORMAL),
			Object('Test ', SCIRT.BOLD),
			Object('String', SCIRT.ITALIC)))

		handleStyleBlock = MockObject(Object(
			Object(#Call, SCIRT.NORMAL, 0),
			Object(#Call, SCIRT.BOLD, 2),
			Object(#Call, SCIRT.ITALIC, 7)))
		handleCharBlock = MockObject(Object(
			Object(#Call, 'A',  0, SCIRT.NORMAL),
			Object(#Call, ' ',  1, SCIRT.NORMAL),
			Object(#Call, 'T',  2, SCIRT.BOLD),
			Object(#Call, 'e',  3, SCIRT.BOLD),
			Object(#Call, 's',  4, SCIRT.BOLD),
			Object(#Call, 't',  5, SCIRT.BOLD),
			Object(#Call, ' ',  6, SCIRT.BOLD),
			Object(#Call, 'S',  7, SCIRT.ITALIC),
			Object(#Call, 't',  8, SCIRT.ITALIC),
			Object(#Call, 'r',  9, SCIRT.ITALIC),
			Object(#Call, 'i', 10, SCIRT.ITALIC),
			Object(#Call, 'n', 11, SCIRT.ITALIC),
			Object(#Call, 'g', 12, SCIRT.ITALIC)))

		parseStyledText(text, handleStyleBlock, handleCharBlock)
		}

	buildText(text)
		{
		s = ''
		for txt in text
			s $= txt[0].Map({ it $ txt[1].Chr() })
		return s
		}

	Test_Get()
		{
		sci = Mock(ScintillaAddonsRichEditorControl)
		text = .buildText(Object(
			Object("'A' \r\n", SCIRT.NORMAL),
			Object('<Test> ', SCIRT.BOLD),
			Object('"String"', SCIRT.ITALIC)))
		sci.When.GetLength().Return(12)
		sci.When.getStyledText(0, 12).Return(text)
		sci.When.parseStyledText([anyArgs:]).CallThrough()
		sci.When.setFontStyle([anyArgs:]).CallThrough()
		sci.When.Get().CallThrough()

		Assert(sci.Get() like:
			'<span style="font-weight:normal;font-style:normal"></span>' $
			'<span style="font-weight:normal;font-style:normal">' $
				'&apos;A&apos; <br /></span>' $
			'<span style="font-weight:bold;font-style:normal">&lt;Test&gt; </span>' $
			'<span style="font-weight:normal;font-style:italic">' $
				'&quot;String&quot;</span>')
		}

	Test_getText()
		{
		sci = Mock(ScintillaAddonsRichEditorControl)
		text = .buildText(Object(
			Object("A ", SCIRT.NORMAL),
			Object('Test ', SCIRT.BOLD),
			Object('Class ', SCIRT.STRIKETHROUGH),
			Object('String', SCIRT.ITALIC)))
		sci.When.GetLength().Return(18)
		sci.When.GetRange(0, 18).Return('A Test Class String')
		sci.When.getStyledText(0, 18).Return(text)
		sci.When.parseStyledText([anyArgs:]).CallThrough()
		sci.When.setFontStyle([anyArgs:]).CallThrough()
		sci.When.getText([anyArgs:]).CallThrough()

		Assert(sci.getText(0, 18) is: 'A Test String')
		Assert(sci.getText(0, 18, includeStrikeThrough:) is: 'A Test Class String')
		}

	Test_getStyleVal()
		{
		getStyleVal = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_getStyleVal
		Assert(getStyleVal(bold:) is: SCIRT.BOLD)
		Assert(getStyleVal(bold:, italic:) is: SCIRT.BOLD_ITALIC)
		Assert(getStyleVal(bold:, underline:) is: SCIRT.BOLD_UNDERLINE)
		Assert(getStyleVal(strikethrough:) is: SCIRT.STRIKETHROUGH)
		Assert(getStyleVal() is: SCIRT.NORMAL)
		}

	Test_Set()
		{
		sci = Mock(ScintillaAddonsRichEditorControl)
		sci.Eval(ScintillaAddonsRichEditorControl.Set, 'Hello')
		sci.Verify.superSet('Hello')
		sci.Verify.Never().setStyled([anyArgs:])

		text = '<span style="font-weight:normal;font-style:normal">A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class </span>' $
			'<span style="font-weight:normal;font-style:italic">String</span>'
		sci.Eval(ScintillaAddonsRichEditorControl.Set, text)
		sci.Verify.setStyled([anyArgs:])
		}

	Test_buildStyled()
		{
		buildStyled = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_buildStyled
		child1 = FakeObject(Text: 'Hello ',
			Attributes: #(style: 'font-weight:normal;font-style:normal'))
		child2 = FakeObject(Text: 'World',
			Attributes: #(style:
				'font-weight:bold;font-style:italic;text-decoration:underline'))


		styleObject = Object()
		Assert(buildStyled(child1, styleObject, 0) is: 'Hello ')
		Assert(styleObject[0] is: Object(start: 0, length: 6, styleBit: SCIRT.NORMAL,
			indicator: false))

		Assert(buildStyled(child2, styleObject, 7) is: 'World')
		Assert(styleObject[1] is: Object(start: 7, length: 5,
			styleBit: SCIRT.BOLD_ITALIC_UNDERLINE, indicator: false))
		}

	Test_setStyled()
		{
		sci = Mock(ScintillaAddonsRichEditorControl)
		setStyled = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_setStyled
		sci.ScintillaAddonsRichEditorControl_stylingMask = 0x1f
		sci.STRIKETHROUGHINDIC = 1
		sci.When.buildStyled([anyArgs:]).CallThrough()

		text = '<span style="font-weight:normal;font-style:normal">A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class </span>' $
			'<span style="font-weight:normal;font-style:italic">String</span>'
		parsed = ParseHTMLRichText.GetParsedText(text)
		sci.Eval(setStyled, parsed)

		sci.Verify.superSet('A Test Class String')

		sci.Verify.Times(4).StartStyling([anyArgs:])
		sci.Verify.StartStyling(0, 0x1f)
		sci.Verify.StartStyling(2, 0x1f)
		sci.Verify.StartStyling(7, 0x1f)
		sci.Verify.StartStyling(13, 0x1f)

		sci.Verify.Times(4).SetStyling([anyArgs:])
		sci.Verify.SetStyling(2, SCIRT.NORMAL)
		sci.Verify.SetStyling(5, SCIRT.BOLD)
		sci.Verify.SetStyling(6, SCIRT.STRIKETHROUGH)
		sci.Verify.SetStyling(6, SCIRT.ITALIC)

		sci.Verify.Times(3).IndicatorClearRange([anyArgs:])
		sci.Verify.IndicatorFillRange(7, 6)
		}

	// undobuffer Tests go here

	Test_userPerformedAction()
		{
		userPerformedAction = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_userPerformedAction
		Assert(userPerformedAction(SC.MOD_BEFOREDELETE | SC.PERFORMED_UNDO) is: false)
		Assert(userPerformedAction(SC.MOD_BEFOREDELETE | SC.PERFORMED_USER))
		}

	Test_nonDeleteAction()
		{
		// delete action types are: MOD_DELETETEXT, MOD_BEFOREDELETE, MOD_CHANGEINDICATOR
		nonDeleteAction = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_nonDeleteAction
		Assert(nonDeleteAction(SC.MOD_DELETETEXT | SC.PERFORMED_USER) is: false)
		Assert(nonDeleteAction(SC.MOD_BEFOREDELETE | SC.PERFORMED_USER) is: false)
		Assert(nonDeleteAction(SC.MOD_CHANGEINDICATOR | SC.PERFORMED_USER) is false)
		Assert(nonDeleteAction(SC.MOD_INSERTTEXT | SC.PERFORMED_USER))
		}

	Test_userPerformedUndoRedo()
		{
		userPerformedUndoRedo = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_userPerformedUndoRedo
		Assert(userPerformedUndoRedo(SC.MOD_CONTAINER | SC.PERFORMED_USER))
		Assert(userPerformedUndoRedo(SC.MOD_BEFOREDELETE | SC.PERFORMED_USER) is: false)
		}

	Test_systemPerformedUndoDelete()
		{
		systemPerformedUndoDelete = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_systemPerformedUndoDelete
		Assert(systemPerformedUndoDelete(SC.PERFORMED_UNDO | SC.MOD_BEFOREDELETE)
			is: false)
		Assert(systemPerformedUndoDelete(SC.MOD_INSERTTEXT | SC.PERFORMED_REDO) is: false)
		Assert(systemPerformedUndoDelete(SC.PERFORMED_UNDO | SC.MOD_INSERTTEXT))
		}

	Test_handleSCNotification()
		{
		handleSCNotification = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_handleSCNotification
		mock = Mock(ScintillaAddonsRichEditorControl)
		mock.ScintillaAddonsRichEditorControl_beginDelete = ''

		// test sends in order:
		// a user performed delete,
		// a normal non-delete
		// a non-delete durring a delete
		// a user performed undo/redo
		// a system performed undo/redo
		mock.When.userPerformedAction([anyArgs:]).Return(true, false, false, false, false)
		mock.When.nonDeleteAction([anyArgs:]).Return(false, true, true, false, false)
		mock.When.userPerformedUndoRedo([anyArgs:]).
			Return(false, false, false, true, false)
		mock.When.systemPerformedUndoDelete([anyArgs:]).Return(false, false, false, true)

		// user performed delete
		scnote = Object(modificationType: 'testModType', length: 1, token: 2)
		mock.Eval(handleSCNotification, scnote)
		mock.Verify.handleUserDelete(scnote, 'testModType')

		// non-delete action when not in the middle of a delete action
		mock.ScintillaAddonsRichEditorControl_beginDelete = ""
		mock.Eval(handleSCNotification, scnote)
		Assert(mock.ScintillaAddonsRichEditorControl_beginDelete isnt: false)

		// non-delete action when in the middle of a delete action
		mock.ScintillaAddonsRichEditorControl_beginDelete = true
		mock.Eval(handleSCNotification, scnote)
		Assert(mock.ScintillaAddonsRichEditorControl_beginDelete is: false)

		// user performed undo/redo
		mock.Eval(handleSCNotification, scnote)
		mock.Verify.handleUserUndoRedo(scnote, 'testModType')

		// system performed undo/redo
		mock.ScintillaAddonsRichEditorControl_styles = #(style: 'style',
			styleVal: 'styleVal')
		mock.Eval(handleSCNotification, scnote)
		mock.Verify.undoFormatAction('styleVal')

		// ensure each method was only called once
		mock.Verify.handleUserDelete([anyArgs:])
		mock.Verify.handleUserUndoRedo([anyArgs:])
		mock.Verify.undoFormatAction([anyArgs:])
		}

	Test_beforeDelete()
		{
		beforeDelete = ScintillaAddonsRichEditorControl.
			ScintillaAddonsRichEditorControl_beforeDelete
		mock = Mock(ScintillaAddonsRichEditorControl)
		mock.When.getText(1, 2, includeStrikeThrough:).Return('e')
		mock.When.getText(0, 1, includeStrikeThrough:).Return('H')
		mock.When.getText(0, 5, includeStrikeThrough:).Return('Hello')

		// Single Char Deletes

		// Begin Delete
		scnote = Object(modificationType: 'testModType', position: 1 length: 1, token: 2)
		mock.Eval(beforeDelete, scnote)
		mock.Verify.getText([anyArgs:])
		mock.Verify.addDeletedFormatToBuffer(Object(start: 1, end: 2))

		// Mid Delete
		scnote = Object(mondificationType: 'testModType', position: 0, length: 1,
			token: 2)
		mock.Eval(beforeDelete, scnote)
		mock.Verify.Times(2).getText([anyArgs:])
		mock.Verify.Never().BeginUndoAction()
		mock.Verify.addDeletedFormatToBuffer(Object(start: 0, end: 1))
		mock.Verify.Never().EndUndoAction()

		// End Delete
		mock.ScintillaAddonsRichEditorControl_beginDelete = false

		// Block Delete
		scnote = Object(modificationType: 'testModType', position: 0, length: 5, token: 0)
		mock.Eval(beforeDelete, scnote)
		mock.Verify.Times(3).getText([anyArgs:])
		mock.Verify.BeginUndoAction()
		mock.Verify.addDeletedFormatToBuffer(Object(start: 0, end:5))
		mock.Verify.EndUndoAction()
		}
	}