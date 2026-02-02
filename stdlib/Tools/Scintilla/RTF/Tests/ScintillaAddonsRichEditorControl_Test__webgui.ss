// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_ScintillaAddonsRichEditorControl_Test
	{
	Test_Set()
		{
		sci = Mock(ScintillaAddonsRichEditorControl)
		sci.When.Act([anyArgs:])
		sci.Eval(ScintillaAddonsRichEditorControl.Set, 'Hello')
		sci.Verify.superSet('Hello')
		sci.Verify.Never().Act([anyArgs:])

		text = '<span style="font-weight:normal;font-style:normal">A </span>' $
			'<span style="font-weight:bold;font-style:normal">Test </span>' $
			'<span style="font-weight:normal;font-style:normal;' $
				'text-decoration: line-through">Class </span>' $
			'<span style="font-weight:normal;font-style:italic">String</span>'
		sci.Eval(ScintillaAddonsRichEditorControl.Set, text)
		sci.Verify.Act(#SetStyleObject, #(
			#(from: #(ch: 0, line: 0), to: #(ch: 2, line: 0), txt: 'A ',
				style: #(bold: false, italic: false, underline: false, strikeout: false)),
			#(from: #(ch: 2, line: 0), to: #(ch: 7, line: 0), txt: 'Test ',
				style: #(bold:, italic: false, underline: false, strikeout: false)),
			#(from: #(ch: 7, line: 0), to: #(ch: 13, line: 0), txt: 'Class ',
				style: #(bold: false, italic: false, underline: false, strikeout:)),
			#(from: #(ch: 13, line: 0), to: #(ch: 19, line: 0), txt: 'String',
				style: #(bold: false, italic:, underline: false, strikeout: false))))
		}

	Test_Get() { }
	Test_beforeDelete() { }
	Test_getStyleVal() { }
	Test_getText() { }
	Test_handleSCNotification() { }
	Test_nonDeleteAction() { }
	Test_parseStyledText() { }
	Test_setFontStyle() { }
	Test_systemPerformedUndoDelete() { }
	Test_userPerformedAction() { }
	Test_userPerformedUndoRedo() { }
	Test_setStyled() { }
	Test_buildStyled() { }
	}