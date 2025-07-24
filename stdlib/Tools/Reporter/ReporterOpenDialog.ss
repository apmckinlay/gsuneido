// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Open Report"
	CallClass(hwnd, reporterMode = 'simple')
		{
		ToolDialog(hwnd, Object(this, reporterMode), border: 0)
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
		modes = .reporterMode is 'form' ? #(form) : #(simple, enhanced)
		QueryApply(.BuildListQuery(prefix))
			{ |x|
			if modes.Has?(x.params.GetDefault('reporterMode', 'simple')) and
				.HasPermission?(x)
				list.AddItem(x.report[prefix.Size() ..])
			}
		}

	selectFirst()
		{
		.list.SetCurSel(0)
		.list.SetFocus()
		}
	Controls:
		(Vert
			ListBox
			(Border 8 (Vert
				(Static 'Right click on reports to rename or delete.')
				(Skip 8)
				(HorzEqual Fill (Button Open) Skip (Button Cancel pad: 20) Fill)
			))
		)
	ListBoxDoubleClick(sel /*unused*/)
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
			menu.Add('Rename', 'Delete')
		ContextMenu(menu).ShowCall(this, x, y)
		}
	stdReport?(name)
		{
		stdPrefix = Opt('~', LastContribution('Reporter').StandardReportPrefix)
		return name.Prefix?(stdPrefix)
		}
	On_Context_Delete()
		{
		name = .list.GetText(.list.GetCurSel())
		if Reporter.CheckScheduled(name)
			return

		if OkCancel('Delete the ' $ Display(name) $ ' report?',
			title: "Delete Report",
			hwnd: .Window.Hwnd, flags: MB.ICONQUESTION)
			{
			.list.DeleteItem(.list.GetCurSel())
			.deleteReport(name)
			}
		}

	deleteReport(name)
		{
		QueryDo('delete ' $ .reportQuery(name))
		c = LastContribution('Reporter')
		c.AfterDeleteReport(name, reporterMode: .reporterMode)
		}

	reportQuery(name)
		{
		return 'params where report = ' $ Display(.prefix $ name)
		}

	On_Context_Rename()
		{
		i = .list.GetCurSel()
		name = .list.GetText(i)
		if false is newname = OkCancel(Object(.rename_dialog, name), "Rename Report",
			.Window.Hwnd)
			return

		if Reporter.CheckScheduled(name)
			return

		.renameReport(name, newname)
		.afterRenameReport(name, newname)
		.list.DeleteItem(i)
		.list.InsertItem(newname, i)
		.list.SetCurSel(i)
		}

	rename_dialog: Controller
		{
		New(name)
			{
			super(.layout(name))
			.field = .Vert.Horz.Field
			.field.Set(name)
			.name = name
			}
		layout(name)
			{
			return Object('Vert',
				Object('Static', 'Rename the ' $ Display(name) $ ' report'),
				'Skip',
				#(Horz, #(Static, To), Skip, (Field width: 40)))
			}
		OK()
			{
			newname = .field.Get()
			if newname is ''
				{
				.AlertInfo("Reporter", "Please enter a name to rename the report to")
				return false
				}
			if false is Reporter.CheckName(newname)
				return false
			return newname
			}
		}
	renameReport(name, newname)
		{
		try
			QueryApply1(.reportQuery(name))
				{ |x|
				x.report = .prefix $ newname
				x.params.report_name = newname
				x.Update()
				return ''
				}
		catch (unused, '*duplicate key')
			return 'already exists'
		}

	afterRenameReport(name, newname)
		{
		c = LastContribution('Reporter')
		if c.ReporterBook is ''
			return
		QueryDo('delete ' $ c.ReporterBook $ '
			where path =~ ' $ Display(c.GetPath(.reporterMode)) $
			' and name is ' $ Display(name))
		rec = Query1(.reportQuery(newname))
		CustomReportsMenu(newname, GetCustomReportsSource(rec.params.Source),
			ReporterModel(rec.report).BuildReportText(), reporterMode: .reporterMode)
		}

	HasPermission?(x)
		{
		c = LastContribution('Reporter')
		return c.HasPermission?(x)
		}
	}