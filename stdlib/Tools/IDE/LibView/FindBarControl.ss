// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// TODO: allow Escape from anywhere in FindBar (not just field)
//		(without stopping Escape from getting to Scintilla)
Controller
	{
	Name: FindBar
	New(data)
		{
		.occurrence = .FindControl('occurrence')
		.Data.Set(data)
		.kill_timer()
		.findtext = .FindControl('find')
		}
	GetText()
		{ return .findtext.Get() }
	SetText(text)
		{ .findtext.Set(text) }
	SetFocus()
		{
		.findtext.SetFocus()
		}
	Select()
		{
		.SetFocus()
		.findtext.Field.SelectAll()
		}
	Controls: (Record
		(Vert
			(HorzEqualHeight
				(EnhancedButton command: 'FindClose', image: 'cross.emf',
					imageColor: 0x737373, mouseOverImageColor: 0x0000ff,
					imagePadding: 0.15, tip: 'Close find bar (Esc)')
				Skip
				(FindStatic 'Find' weight: semibold color: 0x808080, justify: RIGHT)
				(Skip 4)
				FindBarEdit
				Skip
				(EnhancedButton command: 'Previous', image: 'back.emf',
					imageColor: 0x737373, mouseOverImageColor: 0x00cc00,
					imagePadding: 0.15, tip: 'Find previous occurrence (Shift+F3)')
				(EnhancedButton command: 'Next', image: 'forward.emf',
					imageColor: 0x737373, mouseOverImageColor: 0x00cc00,
					imagePadding: 0.15, tip: 'Find next occurrence (F3)')
				Skip
				(EnhancedButton command: 'Mark', image: 'plus.emf',
					imageColor: 0x737373, mouseOverImageColor: 0x00cc00,
					imagePadding: 0.15,
					tip: 'Mark all occurrences (Shift or Ctrl to add)')
				(EnhancedButton command: 'Clear', image: 'minus.emf',
					imageColor: 0x737373, mouseOverImageColor: 0x00cc00,
					imagePadding: 0.15, tip: 'Clear marks')
				(Skip 15)
				(CheckBox, 'Case', tip: "Match case", name: "case")
				Skip
				(CheckBox, 'Word', tip: "Whole words", name: "word")
				Skip
				(CheckBox, 'Regex', tip: "Regular expression", name: "regex")
				Skip
				(CheckBox, 'Expr', tip: "Expression", name: "expr")
				Skip)
			(Skip 2))
		)
	FieldReturn()
		{
		.Send('Find_Return')
		}
	FieldEscape()
		{
		.On_FindClose()
		}
	On_FindClose()
		{
		.Send('On_FindBar_Close')
		}
	On_Next()
		{
		.Send('On_FindBar_Next')
		}
	On_Previous()
		{
		.Send('On_FindBar_Previous')
		}
	On_Mark()
		{
		if not KeyPressed?(VK.SHIFT) and not KeyPressed?(VK.CONTROL)
			.Send('On_FindBar_Clear')
		.Send('On_FindBar_Mark')
		}
	On_Clear()
		{
		.Send('On_FindBar_Clear')
		}

	SetStatus(status)
		{
		color = false is status and .findtext.Get() isnt '' ? CLR.LIGHTRED : CLR.WHITE
		.findtext.Field.SetBgndColor(color)
		}

	// use timer to wait till typing stops
	timer: false
	Record_NewValue(name /*unused*/, value /*unused*/)
		{
		.Edit_Change()
		}
	Edit_Change()
		{
		.UpdateOccurrenceMsg()
		.kill_timer()
		.timer = Delay(200, .send_change) /*= delay before highlighting */
		}
	kill_timer()
		{
		if .timer is false
			return
		.timer.Kill()
		.timer = false
		}
	UpdateOccurrenceMsg()
		{
		.occurrence.Set(
			.occurNum is false or .occurCount is false or '' is .findtext.Get()
				? ''
				: .occurNum $ ' of ' $ .occurCount
		)
		}
	occurNum: false
	occurCount: false
	UpdateOccurrenceInfo(num, count)
		{
		if num isnt false
			.occurNum = num
		if count isnt false
			.occurCount = count
		}
	findtext: false
	send_change()
		{
		if .findtext is false
			return
		.Data.HandleFocus() // make sure record is updated
		.Send('Find_Change')
		}
	Destroy()
		{
		.kill_timer()
		super.Destroy()
		}
	}
