// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Used by Inspect and by Debugger
Controller
	{
	Name: "Inspect"
	Xmin: 300
	Ymin: 200
	New(x, title = "", .hwnd = 0, themed? = false)
		{
		super(.layout(x, themed?))
		if false is .list = .FindControl('ListBox')
			{
			.editor_title = title
			return
			}
		.list.SetFont('@mono')
		.Reset(x, title)
		}
	layout(x, themed?)
		{
		if .displayAsObject?(x)
			return Object(
				"Vert",
				Object("ListBox", x, :themed?)
				#(Horz
					Fill
					(Button 'Copy Name', size: '-1')
					Skip
					(Button 'Copy Value', size: '-1')
					Fill))
		return Object('ScintillaAddons', set: String(x), readonly:, Addon_wrap:,
			xmin: 300, height: 15)
		}
	displayAsObject?(x)
		{
		return Object?(x) or Instance?(x) or Class?(x)
		}

	editor_title: 'Inspect'
	Startup()
		{
		.set_title(.list isnt false ? .stack.Top().title : .editor_title)
		}
	Reset(x, title = "")
		{
		.stack = Stack()
		.show(x, title)
		}
	show(x, title = "")
		{
		.list.DeleteAll()
		n = 6
		listSize = x.Size(list:)
		for m in x.Members()
			{
			m = .fmt(m, listSize)
			if m.Size() > n
				n = m.Size()
			}
		n += 2
		if .stack.Count() > 0
			.list.AddItem(" <<<".RightFill(n) $ .display(.stack.Top().x), -1)
		if false isnt base = BaseClass(x)
			.list.AddItem(" base".RightFill(n) $ .display(base), -2)
		for m in .sortedMembers(x)
			.list.AddItem(.fmt(m, listSize).RightFill(n) $ .display(x[m]))

		if title is ""
			title = Name(x)
		.set_title(title)

		.stack.Push(:x, :title)
		}
	sortedMembers(x)
		{
		return x.Members().SortWith!({ it is #this ? '' : it })
		}
	fmt(m, listSize)
		{
		if Number?(m) and 0 <= m and m <= listSize
			m = '(' $ m $ ')'
		else if not String?(m)
			m = .display(m)
		return m
		}
	GetMember(i)
		{
		currentOb = .stack.Top()
		members = .sortedMembers(currentOb.x)
		return members[i]
		}
	GetValue(i)
		{
		currentOb = .stack.Top()
		members = .sortedMembers(currentOb.x)
		return currentOb.x[members[i]]
		}
	ListBoxDoubleClick(sel)
		{
		if sel is -1
			return 0 // no currently selected item
		i = .list.GetSelected()
		if i is -1
			{
			.stack.Pop()
			.show(@.stack.Pop()) // pop because show will re-push
			return 0
			}
		x = .stack.Top().x
		if i is -2
			{
			x = x.Base()
			title = ""
			}
		else
			{
			m = .sortedMembers(x)[i]
			x = x[m]
			title = .stack.Top().title $ "." $ (String?(m) ? m : .display(m))
			}
		if .displayAsObject?(x)
			{
			if KeyPressed?(VK.SHIFT)
				Inspect.Window(x, title, .hwnd) // open in new window
			else
				.show(x, title) // open in same window
			}
		else if String?(x)
			Inspect(x, title, .hwnd)
		return 0
		}
	set_title(title)
		{
		.Send('SetTitle', title is "" ? "Inspect" : title)
		}
	display(a)
		{
		try
			return Display(a)
		catch (e)
			return "\xab" $ (e.Has?("buffer overflow")
				? "too large to display"
				: "Display() failed: " $ e) $ "\xbb"
		}
	On_Copy_Name()
		{
		.copySelected(member?:)
		}
	On_Copy_Value()
		{
		.copySelected()
		}
	copySelected(member? = false)
		{
		curSel = .list.GetCurSel()
		if curSel is -1
			return // no currently selected item

		i = .list.GetSelected()
		if i is -1 or i is -2
			return
		value = member?
			? String(.GetMember(i))
			: Display(.GetValue(i))
		ClipboardWriteString(value)
		}
	}
