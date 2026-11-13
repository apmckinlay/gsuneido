// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Open Report"
	CallClass(hwnd, reporterMode = #simple)
		{
		ToolDialog(hwnd, [this, reporterMode], border: 0)
		}

	New(.reporterMode)
		{
		.list = .Vert.ListBox
		.prefix = "Reporter - "
		.buildOpenList(.prefix, .list)
		// need delay or it scrolls for some unknown reason
		.Defer(.selectFirst)
		}

	BuildListQuery(prefix)
		{
		return "params where report =~ '^(?q)" $ prefix $ "'
			sort report"
		}

	buildOpenList(prefix, list)
		{
		modes = .reporterMode is #form ? #(form) : #(simple, enhanced)
		QueryApply(.BuildListQuery(prefix))
			{|x|
			if modes.Has?(x.params.GetDefault(#reporterMode, #simple)) and
				ReporterActions.HasPermission?(x)
				list.AddItem(x.report[prefix.Size() ..])
			}
		}

	selectFirst()
		{
		.list.SetCurSel(0)
		.list.SetFocus()
		}

	Controls: (Vert,
		ListBox,
		(Border, 8,
			(Vert,
				(Static, "Right click on reports to rename or delete."),
				(Skip, 8),
				(HorzEqual, Fill, (Button, Open), Skip, (Button, Cancel, pad: 20),
					Fill))))
	ListBoxDoubleClick(sel/*unused*/)
		{
		.On_Open()
		}

	On_Open()
		{
		sel = .list.GetCurSel()
		if sel is -1
			{
			Alert("Please select a report.")
			return
			}
		.Window.Result(.prefix $ .list.GetText(sel))
		}

	ListBox_ContextMenu(x, y)
		{
		if -1 is i = .list.GetCurSel()
			return // click outside items
		menu = Object()
		name = .list.GetText(i)
		if not .stdReport?(name)
			menu.Add(#Rename, #Delete)
		ContextMenu(menu).ShowCall(this, x, y)
		}

	stdReport?(name)
		{
		stdPrefix = Opt('~', LastContribution(#Reporter).StandardReportPrefix)
		return name.Prefix?(stdPrefix)
		}

	On_Context_Delete()
		{
		name = .list.GetText(.list.GetCurSel())
		result = ReporterActions.Delete(name, confirmationDialog:, hwnd: .Window.Hwnd)
		if result isnt true
			return
		.list.DeleteItem(.list.GetCurSel())
		}

	On_Context_Rename()
		{
		i = .list.GetCurSel()
		name = .list.GetText(i)

		if false is newname = ReporterActions.Rename(name, .Window.Hwnd)
			return

		.list.DeleteItem(i)
		.list.InsertItem(newname, i)
		.list.SetCurSel(i)
		}
	}
