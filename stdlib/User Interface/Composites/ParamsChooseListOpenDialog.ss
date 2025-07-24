// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Load Values"
	CallClass(hwnd, fieldSaveName)
		{
		return OkCancel(Object(this, fieldSaveName), .Title, hwnd)
		}
	New(field)
		{
		.list = .Vert.ListBox
		.prefix = 'InlistValues - '
		QueryApply("params where report =~ 'InlistValues - '
			and report_options is " $ Display(field) $ " sort report")
			{ |x|
			.list.AddItem(x.report[.prefix.Size() ..])
			}
		// need delay or it scrolls for some unknown reason
		.Defer(.selectFirst)
		}
	selectFirst()
		{
		.list.SetCurSel(0)
		.list.SetFocus()
		}
	Controls:
		(Vert
			(ListBox xmin: 200, ymin: 200)
			(Skip 5)
			(Static 'Right click on list to delete.'))
	ListBoxDoubleClick(sel/*unused*/)
		{
		.Send('On_OK')
		}
	OK()
		{
		sel = .list.GetCurSel()
		if sel is -1
			{
			Alert("Please select a list.")
			return false
			}
		return .prefix $ .list.GetText(sel)
		}

	ListBox_ContextMenu(x, y)
		{
		if .list.GetCurSel() is -1
			return // click outside items
		ContextMenu(#("Delete")).ShowCall(this, x, y)
		}

	On_Context_Delete()
		{
		name = .list.GetText(.list.GetCurSel())
		if OkCancel('Delete the ' $ Display(name) $ ' list?',
			title: "Delete List",
			hwnd: .Window.Hwnd, flags: MB.ICONQUESTION)
			{
			.list.DeleteItem(.list.GetCurSel())
			.deleteList(name)
			}
		}
	deleteList(name)
		{
		QueryDo('delete ' $ .reportQuery(name))
		}

	reportQuery(name)
		{
		return 'params where report = ' $ Display(.prefix $ name)
		}
	}
