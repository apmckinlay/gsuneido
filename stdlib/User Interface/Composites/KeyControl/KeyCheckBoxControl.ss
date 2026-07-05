// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
KeyControl
	{
	New(@args)
		{
		super(@.processArgs(args))
		.Field.SetReadOnly(true)
		}

	processArgs(args)
		{
		args.width = 13
		return args
		}

	ProcessResults(checked)
		{
		.Set(checked)
		.NewValue(checked)
		}

	ReprocessValue()
		{
		}

	Set(.checked)
		{
		super.Set(.setStr(.checked))
		}

	setStr(checked)
		{
		count = Object?(checked)
			? checked.Size()
			: 0
		return count > 0
			? 'Selected: ' $ count
			: 'None'
		}

	Get()
		{
		return .checked
		}

	Valid?()
		{
		return .checked is '' or Object?(.checked)
		}

	SetReadOnly(readOnly)
		{
		super.SetReadOnly(readOnly)
		.Field.SetReadOnly(true)
		}

	BuildDialogControl(saveInfoName)
		{
		dialog = super.BuildDialogControl(saveInfoName)
		dialog[0] = .dialogControl
		dialog.checkBoxColumn = 'params_itemselected'
		dialog.closeButton? = false
		dialog.check = .checked
		return dialog
		}

	dialogControl: Controller
		{
		New(@args)
			{
			super(.layout(args))
			.checkList = .FindControl('checkList')
			.checkList.CheckRecordByKeys(.check)
			.checkedDisplay = .FindControl('checkedDisplay')
			.syncDisplay()
			}

		layout(args)
			{
			.field = args.prefixColumn = args.field
			.checkBoxColumn = args.checkBoxColumn
			.check = args.GetDefault('check', #()).Sort!()
			control = Object(KeyListCheckboxView, name: 'checkList').Append(args)
			control.value = .check.GetDefault(0, '')
			control.prefix = ''
			return Object('Vert',
				control,
				#(ScintillaAddonsEditor, readonly:, xstretch: 1, name: checkedDisplay),
				#(Horz, (Button, Clear), Fill, OkCancel))
			}

		syncDisplay()
			{
			.checkedDisplay.Set(.checked().Join('\r\n'))
			}

		checked()
			{
			checked = Object()
			for rec in .checkList.GetCheckedRecords().list
				checked.Add(rec[.field])
			return checked.Sort!()
			}

		VirtualList_LeftClick(rec, col)
			{
			if rec isnt false col is .checkBoxColumn
				.syncDisplay()
			return 0
			}

		VirtualList_Space()
			{
			.syncDisplay()
			return 0
			}

		Scintilla_DoubleClick()
			{
			line = .checkedDisplay.LineFromPosition()
			if '' isnt locate = .checkedDisplay.GetLine(line)
				{
				.checkedDisplay.SelectLine(line)
				.checkList.GetField().Set(locate.Trim())
				.checkList.FieldChange()
				}
			}

		On_Clear()
			{
			.checkList.UncheckAll()
			.syncDisplay()
			}

		On_OK()
			{
			.Window.Result(.checked())
			}
		}
	}