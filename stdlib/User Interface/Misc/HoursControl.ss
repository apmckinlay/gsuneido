// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: 'Hours'
	// Note: if mask is less than 4 decimals, hours to decimal will not be accurate
	//       Also does not handle negative numbers
	New(width = 5, .readonly = false, .mandatory = false, justify = 'RIGHT')
		{
		super(:width, :readonly, :mandatory, :justify)
		.menu = Object('Decimal', 'Hours')
		.get_display_choice()
		.SubClass()
		.mask = '#####.####'
		}
	which: 0
	orig_which: 0
	ContextMenu(x, y)
		{
		dirty = .Dirty?()
		if x is 0 and y is 0 // keyboard
			{
			pt = Object(x: 10, y: 20)
			ClientToScreen(.Hwnd, pt)
			x = pt.x
			y = pt.y
			}
		.KillFocus()
		val = .Get()
		menu = .menu.Copy()
		menu[.which] = Object(name: menu[.which], state: MF.CHECKED, type: MFT.RADIOCHECK)
		i = ContextMenu(menu).Show(.Hwnd, x, y) - 1
		if menu.Member?(i)
			.which = i
		.Set(val)
		.Dirty?(dirty)
		return 0
		}
	get_display_choice()
		{
		if not TableExists?('params')
			return
		x = Query1(.display_param_query())
		if x isnt false and .menu.Has?(x.params)
			.which = .orig_which = .menu.Find(x.params)
		}
	display_param_query()
		{
		return 'params
			where user = ' $ Display(Suneido.User) $ '
			and report = "HoursControl Display"'
		}
	save_display_choice()
		{
		if .which is .orig_which or
			not TableExists?('params')
			return
		Transaction(update:)
			{ |t|
			t.QueryDo('delete ' $ .display_param_query())
			t.QueryOutput('params', Record(
				user: Suneido.User
				report: "HoursControl Display"
				params: .menu[.which]))
			}
		}
	Get()
		{
		text = GetWindowText(.Hwnd)
		if .menu[.which] is 'Hours'
			text = .convertToDecimal(text, .mask)
		if text is '#' or text is '-'
			return ''
		return Number(text)
		}
	Set(val)
		{
		text = ''
		if val is ''
			text = ''
		else if .menu[.which] is 'Hours'
			text = .convertToHours(val, .mask)
		else
			text = val.Format(.mask)
		SetWindowText(.Hwnd, text)
		}
	KillFocus()
		{
		text = GetWindowText(.Hwnd)
		if .menu[.which] is 'Hours'
			SetWindowText(.Hwnd, .convertToHours(text, .mask))
		else if .menu[.which] is 'Decimal'
			{
			text = .convertToDecimal(text, .mask)
			if text isnt ''
				text = text.Format(.mask)
			SetWindowText(.Hwnd, text)
			}
		}
	minutesPerHour: 60
	convertToHours(val, mask)
		{
		if String(val).Has?(':')
			return .checkHoursSplit(val, mask) ? val : ''
		if not Number?(val) and not val.Number?()
			return ''
		val = Number(val)
		minutes = (val.Frac() * .minutesPerHour).Round(0)
		return val.Int() $ ':' $ minutes.Pad(2, '0')
		}
	convertToDecimal(val, mask)
		{
		if Number?(val) or val.Number?()
			return Number(val)
		if not val.Has?(':')
			return ''
		if not .checkHoursSplit(val, mask)
			return ''
		ob = val.Split(':').Map!(Number)
		hour = ob[0]
		int = hour.Int()
		frac = hour.Frac() + ob[1] / .minutesPerHour
		if frac >= 1
			{
			int += frac.Int()
			if int.Format(mask) is '#'
				return ''
			frac = frac.Frac()
			}
		return Number(int $ (frac is 0 ? '' : frac))
		}
	checkHoursSplit(val, mask)
		{
		ob = val.Split(':')
		if ob.Size() isnt 2
			return false
		for item in ob
			if item isnt '' and not item.Number?()
				return false
		//check that hours not larger than mask (valid)
		if Number(ob[0]).Format(mask) is '#'
			return false
		// negative numbers not allowed
		if ob[0].Prefix?('-')
			return false

		return true
		}
	Valid?()
		{
		s = super.Get().Tr(',').Trim()
		if .readonly and s isnt "#"
			return true
		if s is ''
			return not .mandatory
		if not s.Number?()
			return .checkHoursSplit(s, .mask)
		return true
		}
	Destroy()
		{
		.save_display_choice()
		super.Destroy()
		}
	}