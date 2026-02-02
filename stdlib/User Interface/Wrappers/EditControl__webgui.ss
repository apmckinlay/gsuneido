// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	text: ''
	Name: 'Edit'
	Status: ''
	ComponentName: "Edit"
	New(.mandatory = false, .readonly = false, style/*unused*/ = 0,
		bgndcolor = "", textcolor = "", hidden = false,
		tabover = false, font = "", size = "", weight = "", underline = false,
		width = false, height = false, cue = false, readOnlyBgndColor = false,
		status = '')
		{
		.SuSetHidden(hidden)
		.Send("Data")
		.ContextExtra = Object()

		if cue isnt false
			.SetCue(cue)

		if status > ''
			.SetStatus(status)

		.ComponentArgs = Object(readonly, bgndcolor, textcolor, hidden, tabover,
			font, size, weight, underline,
			width, height, readOnlyBgndColor)
		}

	isReadonly?: false
	SetReadOnly(readOnly)
		{
		if .readonly
			return
		.Act('SetReadOnly', .isReadonly? = readOnly)
		}

	SetReadOnlyBrush(color = false)
		{
		if color is false
			return
		if Object?(color)
			color = RGB(color[0], color[1], color[2])
		.Act(#SetReadOnlyColor, color)
		}

	GetReadOnly()
		{
		return .readonly or .isReadonly?
		}

	Get()
		{
		return .text
		}

	Set(value)
		{
		if (not String?(value))
			value = Display(value)
		.text = value
		.Dirty?(false)
		.sel = [value.Size(), value.Size()]
		.Act('Set', value)
		}

	// used by SetWindowText to set value directly
	SetText(text)
		{
		.text = text
		.dirty? = false
		.Act('SetText', text)
		}

	EN_CHANGE(text)
		{
		if .GetReadOnly() is true
			ProgrammerError('Detect edit change in readonly',
				params: [:text, name: .Name], caughtMsg: 'unattended; msg not processed')

		.text = text
		.Send("Edit_Change")
		.dirty? = true
		}

	dirty?: false
	Dirty?(dirty = "")
		{
		Assert(dirty is true or dirty is false or dirty is "")
		if (dirty isnt "")
			.dirty? = dirty
		return .dirty?
		}

	EN_KILLFOCUS()
		{
		.KillFocus()
		if .Send("Dialog?") isnt true
			{
			if 0 is valid? = .Send('Edit_ParentValid?')
				valid? = .Valid?()
			.Act('SetValid', valid?)
			}
		if .Status > ""
			.Send("Status", "")
		}

	KILLFOCUS()
		{
		.EN_KILLFOCUS()
		}

	KillFocus()
		{
		}

	EN_SETFOCUS()
		{
		.Send('Field_SetFocus')
		if .Status > ""
			.Send("Status", .Status)
		.SetValid() // don't color invalid when focused
		return 0
		}

	SETFOCUS()
		{
		.EN_SETFOCUS()
		}

	SetValid(valid? = true, force = false)
		{
		.Act('SetValid', valid?, force)
		}

	Valid?() // derived classes should define this
		{
		return .validCheck?(.Get(), .mandatory)
		}

	validCheck?(data, mandatory)
		{
		return not (mandatory and data is "")
		}

	ValidData?(@args)
		{
		return .validCheck?(args[0], args.GetDefault('mandatory', false))
		}

	SetSel(start, end)
		{
		.sel = [start, end]
		.Act('SetSel', start, end)
		}

	sel: [0, 0]
	UpdateSel(.sel) {}

	GetSel()
		{
		return .sel
		}

	GetSelText()
		{
		range = .GetSel()
		return .text[range[0]..range[1]]
		}

	EnsureSelect()
		{
		.SetSel(.sel[0], .sel[1])
		}

	SelectAll()
		{
		.sel = [0, .text.Size()]
		.Act('SelectAll')
		}

	ReplaceSel(text)
		{
		.text = .text.ReplaceSubstr(.sel[0], .sel[1] - .sel[0], text)
		.Act('ReplaceSel', text)
		}

	SetCue(cue)
		{
		.Act('SetCue', cue)
		}

	// used by stdlib:GetWindowText__webgui
	GetWindowText()
		{
		return .text
		}

	On_Delete()
		{
		if .GetReadOnly() is true
			return
		.Act('On_Delete')
		}
	On_Cut()
		{
		if .GetReadOnly() is true
			return
		.Act('On_Cut')
		}
	On_Copy()
		{
		.Act('On_Copy')
		}
	On_Paste()
		{
		if .GetReadOnly() is true
			return
		.Act('On_Paste')
		}
	On_Undo()
		{
		if .GetReadOnly() is true
			return
		.Act('On_Undo')
		}
	On_Select_All()
		{
		.SelectAll()
		}

	AddContextMenuItem(name, runFunc, enabledFunc = function () { #(addToMenu:) })
		{
		if Object?(name)
			{
			// Handle Cascade menu options
			n = Object()
			for idx in name.Members()
				n.Add(Object(name: name[idx], runFunc: runFunc[idx]))
			.ContextExtra.Add(Object(name: n, :enabledFunc))
			}
		else
			.ContextExtra.Add(Object(:name, :runFunc, :enabledFunc))
		}

	FieldMenu: ('&Undo\tCtrl+Z', '',
		'Cu&t\tCtrl+X', '&Copy\tCtrl+C', '&Paste\tCtrl+V', '&Delete', '',
		'Select &All\tCtrl+A')
	ContextMenu(x, y)
		{
		.SetFocus()

		// when another ctrl loses focus, it could trigger list destroy (dynamic layouts)
		if .Destroyed?()
			return 0

		menu = .FieldMenu.Copy()
		for item in .ContextExtra
			{
			enabledOb = (item.enabledFunc)()
			if enabledOb.addToMenu is true
				{
				if Object?(item.name)
					menu.Add(item.name)
				else
					menu.Add(Object(name: item.name,
						state: enabledOb.GetDefault('state', MFS.ENABLED)
						runFunc: item.runFunc))
				}
			}

		if Suneido.User is 'default'
			menu.Add(@.DevMenu)
		i = ContextMenu(menu).Show(.Hwnd, x, y) - 1
		if i is -1
			{
			.EnsureSelect()
			return 0
			}
		.callContext(menu, i)
		return 0
		}

	callContext(menu, chosen, j = 0)
		{
		for (m = 0; j <= chosen and m < menu.Size(); ++j, ++m)
			{
			if .isCascadeMenu(menu[m])
				{
				if menu[m].Member?('name')
					j = .callContext(menu[m].name, chosen, j) - 1
				else
					j = .callContext(menu[m], chosen, j) - 1
				}
			else if (j is chosen)
				{
				// if we have runFunc call it, otherwise .send a msg
				if Object?(menu[m])
					{
					(menu[m].runFunc)(source: this)
					continue
					}
				.ContextMenuCall(menu[m])
				}

			}
		return j
		}
	isCascadeMenu(menu)
		{
		return (Object?(menu) and
				(not menu.Member?('name') or Object?(menu.name)))
		}

	SetStatus(status)
		{
		.Status = status
		.ToolTip(.Status)
		}

	// TODO: implement me
	ShowBalloonTip(msg/*unused*/, icon/*unused*/ = 'NONE')
		{
		}

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
