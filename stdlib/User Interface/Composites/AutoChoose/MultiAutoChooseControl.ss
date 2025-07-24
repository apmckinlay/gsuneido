// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// see AutoChooseControl for testing notes
AutoChooseControl
	{
	Name: MultiAutoChoose
	Unsortable:,
	New(.list = false, width = 50, status = "", readonly = false,
		.mandatory = false, .allowOther = false, height = 3, style = 0, .sort = false)
		{
		super(.list, :style, :status, :width, :height, :readonly, :mandatory, :allowOther)
		.lineheight = .Ymin / height
		}

	InsertChoice(s)
		{
		.SetSel(.start, .end)
		.ReplaceSel(s.Tr(',') $ ', ')
		textEnd = .set(trailingComma:).Size()
		.SetSel(textEnd, textEnd) // Ensure cursor is set to the END of the text
		}

	set(trailingComma = false)
		{
		before = super.Get()
		if before is after = .getUnique(before) // changed (like trimmed, duplicates)
			return
		.Set(set = after $ (trailingComma ? ', ' : ''))
		.Send("NewValue", after)
		return set
		}

	GetPrefix()
		{
		sel = .GetSel()
		text = GetWindowText(.Hwnd)
		if sel[0] < text.Trim().Size()
			return ""
		i = sel[0]
		while i > 0 and text[i] isnt ','
			--i
		if i > 0
			++i
		while text[i].White?()
			++i
		.start = i
		.end = sel[1]
		return text[i :: sel[1] - i]
		}

	GetListPos() // called by AutoChooseList
		{
		r = GetWindowRect(.Hwnd)
		p = .PosFromChar(.start)
		r.left += p.x - 1
		r.top += p.y
		bottomMargin = .9
		r.bottom = r.top + (.lineheight * bottomMargin).Int()
		return r
		}

	Get()
		{
		return .getUnique(super.Get())
		}

	getUnique(text)
		{
		unique = .getValues(text).UniqueValues()
		if .sort
			unique.Sort!()
		return unique.Join(',')
		}

	AlternateJoinChar: false
	getValues(text)
		{
		if .AlternateJoinChar isnt false
			text = text.Tr(.AlternateJoinChar, ',')
		return text.Replace('\r?\n', ',').Split(',').Map!(#Trim).Remove('')
		}

	KillFocus()
		{
		super.KillFocus()
		if not .GetReadOnly()
			.set()
		}

	LBUTTONDOWN()
		{
		if GetFocus() isnt .Hwnd
			{
			.SetFocus()
			return 0
			}

		return 'callsuper'
		}

	EN_SETFOCUS()
		{
		s = .Get()
		if s isnt '' and s[-1] isnt ','
			.Set(s $= ', ')
		.SetSel(s.Size(), -1)
		super.EN_SETFOCUS()
		}

	Valid?()
		{
		value = .Get()
		if .mandatory and value is ""
			return false
		return .validChoices?(.Get(), .allowOther, .list)
		}

	ValidData?(@args)
		{
		allowOther = args.GetDefault('allowOther', false)
		list = args.GetDefault('list', false)
		if not super.ValidData?(args)
			return false
		if String?(list) and args.Member?('record')
			list = args.record[list]
		return .validChoices?(args[0], allowOther, list)
		}

	validChoices?(choices, allowOther, list)
		{
		for s in choices.Split(',').Map!(#Trim).Filter({|x| x isnt ''})
			if not .ValidChoice?(s, allowOther, list)
				return false
		return true
		}
	}
