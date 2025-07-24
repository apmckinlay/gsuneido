// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// Contributions from Claudio Mascioni
ChooseControl
	{
	Name: 'ChooseList'

	New(list = false, width = 10, .mandatory = false, allowOther = false,
		set = false, selectFirst = false, .listField = false, .splitValue = ',',
		.listSeparator = ' - ', status = '', .otherListOptions = #(), font = "",
		trim = true, size = "", bgndcolor = "", textcolor = "", tabover = false,
		hidden = false, weight = '', field = #(), cue = false, .readonly = false)
		{
		super(Object('Field', name: 'Value', style: WS.CLIPSIBLINGS,
			:width, :mandatory, :status, :font, :size, :weight, :trim,
			:bgndcolor, :textcolor, :tabover, :hidden, :cue, :readonly).Merge(field))
		.listarg = list
		.hfont = Suneido.hfont
		.allowOther? = allowOther is true

		.Field.WithSelectObject(.hfont)
			{|hdc|
			GetTextMetrics(hdc, tm = Object())
			.lineheight = tm.Height + tm.ExternalLeading
			}

		if selectFirst isnt false and .list.Member?(0)
			{
			.Set(.list[0])
			.Send('NewValue', .Get())
			}

		if set isnt false
			{
			.Set(set)
			.Send('NewValue', .Get())
			}

		}
	GetList()
		{
		return .list
		}
	SetList(list)
		{
		.listarg = list
		}
	getter_list()
		{
		return .ListGet(this, .listField, .listarg, .splitValue)
		}

	// also used by ChooseManyControl, ChooseManyFieldControl and ChooseListTupleControl
	ListGet(ctrl, listfield, listarg, splitValue = ",")
		{
		if listfield is false
			return listarg
		list = ctrl.Send('GetField', listfield)
		if list is 0 // nothing handled GetField
			return #()
		return ChooseList_ValidDataListFromRules.SplitListFromString(list, splitValue)
		}

	listbox: Controller
		{
		Xstretch: 0
		Ystretch: 0
		New(list, sel, xmin, lineheight, .chooselist, listSeparator)
			{
			super(Object('Scroll',
				Object(ChooseList, list, sel, chooselist, xmin: xmin - 2,
					:listSeparator)
				xmin: xmin - 2,
				ymin: Min(ChooseList.NumLines, list.Size()) * lineheight,
				dyscroll: lineheight,
				noEdge:))
			.Xmin = .Scroll.Xmin
			.Ymin = .Scroll.Ymin
			}
		Activate()
			{
			if .fromListEditWindow?()
				.chooselist.Window.ClosingListEdit = false
			}
		Inactivate()
			{
			.chooselist.Result(false)
			if .Window.Hwnd isnt 0
				{
				chooselist = .chooselist
				if .fromListEditWindow?()
					{
					.ClearFocus() // so it does not activate previous window, e.g. Chrome
					chooselist.Window.ClosingListEdit = true
					return
					}
				hwnd = .Window.Hwnd
				.Window.Hwnd = 0
				if chooselist.Member?(#Field) and chooselist.Member?(#Window)
					PostMessage(chooselist.Window.Hwnd, WM.APP_SETFOCUS, 0,
						chooselist.Field.Hwnd)
				DestroyWindow(hwnd)
				}
			}
		fromListEditWindow?()
			{
			return .chooselist.Member?(#Window) and
				.chooselist.Window.Base?(ListEditWindow)
			}
		MouseWheel(wParam)
			{
			.Scroll.MOUSEWHEEL(wParam)
			}
		}
	Dialog?()
		{
		return .dialog
		}
	dialog: false
	On_DropDown()
		{
		if .readonly
			return
		if ((false is r = .InitDropDown()) or .list.Size() is 0)
			return
		x = r.left; y = r.bottom
		wr = GetWorkArea(r)
		height = Min(ChooseList.NumLines, .list.Size()) * .lineheight + 2
		if y + height > wr.bottom
			y -= height + (r.bottom - r.top)

		.Field.WithSelectObject(.hfont, {|hdc| width = .getWidth(hdc) })
		if .list.Size() > ChooseList.NumLines
			width += GetSystemMetrics(SM.CXVSCROLL)
		if .WndProc.Base?(ListEditWindow)
			{
			x--
			width = Max(width, r.right - r.left) + 2
			}
		else
			width = Max(width, .Xmin)
		if .Xstretch > 0
			width = Max(width, r.right - r.left)
		if wr.right < (x + width)
			x -= ((x + width) - wr.right) + 2

		sel = .matchPrefix(.allowOther?, .list, .listSeparator)
		.dialog = true

		// save and restore dirty flag since field losing focus (when list displayed)
		//  can cause the field to become not dirty
		dirty? = .Dirty?()
		Window(Object(.listbox, .list, sel, ScaleWithDpiFactor.Reverse(width),
			ScaleWithDpiFactor.Reverse(.lineheight), this, .listSeparator),
			parentHwnd: .Window.Hwnd,
			style: WS.POPUP | WS.BORDER, exStyle: WS_EX.TOOLWINDOW, :x, :y)
		.Dirty?(dirty?)
		}

	getWidth(hdc)
		{
		width = 0
		for item in .list
			{
			text = String(item)
			// allowing a width > screen width causes problems with hdc and painting
			// chose 175 as Max because it almost fills the screen width at 800 x 600
			GetTextExtentPoint32(hdc, text, Min(175/*= max size length*/, text.Size()),
				ex = Object())
			width = Max(width, ex.x + 8) /*= offset */
			}
		return width
		}

	Result(value)
		{
		.dialog = false
		if value isnt false
			{
			.Set(value)
			.NewValue(.Get())
			.Field.Valid?()
			}
		}
	matchPrefix(allowOther?, list, listSeparator)
		{
		value = .Field.Get()
		if value is "" or allowOther?
			return false
		if false isnt i = list.FindIf(
			{|val|
			val = String(val)
			if listSeparator isnt ''
				val = val.BeforeFirst(listSeparator)
			val is value
			})
			return i
		if false isnt i = list.FindIf(
			{|val|
			String(val).Prefix?(value)
			})
			return i
		if false isnt i = list.FindIf(
			{|val|
			String(val).Lower().Prefix?(value.Lower())
			})
			return i
		return false
		}
	itemInList?(value)
		{
		// does not loop through all the items when list contains numbers and
		// strings, thus .list.Copy()
		for item in .list.Copy()
			{
			item = String(item)
			if .listSeparator isnt ''
				item = item.BeforeFirst(.listSeparator)
			if item is value
				return true
			}
		if .otherListOptions.Has?(value)
			return true
		return false
		}

	FieldSetFocus()
		{
		super.FieldSetFocus()
		// refresh list on SETFOCUS without affecting dirty flag
		if .listField isnt false
			{
			.Send('DoWithoutDirty')
				{
				.Send("InvalidateFields", Object(.listField))
				}
			}
		}
	FieldKillFocus()
		{
		if not .otherListOptions.Empty?()
			{
			val = .Field.Get()
			if .Field.Dirty?() is true and .otherListOptions.Has?(val)
				.Field.Set("")
			}
		if not .dialog
			.killfocus()
		}
	FieldReturn()
		{
		.killfocus()
		}
	killfocus()
		{
		if .Dirty?() and
			(false isnt i = .matchPrefix(.allowOther?, .list, .listSeparator))
			{
			item = String(.list[i])
			if .listSeparator isnt ''
				item = item.BeforeFirst(.listSeparator)
			.Field.Set(item)
			.Send('NewValue', .Get()) // use Get in case sub-classes override
			.Send('NotifyKilledFocus')
			}
		}

	Valid?()
		{
		value = .Get()
		if .allowOther? and value isnt ''
			return true
		if not .mandatory and value is ''
			return true
		return .itemInList?(value)
		}

	ValidData?(@args)
		{
		if '' is value = args[0]
			return not args.Member?('mandatory') or args.mandatory is false

		if args.Member?('allowOther') and args.allowOther is true
			return true

		listOb = .GetValidDataList(args)
		otherOptions = args.GetDefault('otherListOptions', #())
		return listOb.Has?(value) or otherOptions.Has?(value)
		}

	GetValidDataList(args)
		{
		listOb = Object()
		if args.Member?('list')
			listOb = args.list
		else if args.Member?('listField') and
			Object?(ob = ChooseList_ValidDataListFromRules(args))
			listOb = ob
		else if args.Member?(1) and Object?(args[1])
			listOb = args[1]
		else
			listOb = .ValidationList()

		if '' isnt sep = args.GetDefault('listSeparator', ' - ')
			listOb = listOb.Map({ it.BeforeFirst(sep) })
		return listOb
		}

	ValidationList()
		{
		return Object()
		}

	SelectItem(i)
		{
		if i >= .list.Size()
			i = 0
		.Set(.list[i])
		.Send('NewValue', .Get())
		}

	SelectedItem()
		{
		return .list.Find(.Get())
		}
	}
