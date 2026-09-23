// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Prefix: "Reporter - "

	HasPermission?(paramRec)
		{
		c = LastContribution(#Reporter)
		return c.HasPermission?(paramRec)
		}

	Rename(rptName, hwnd, reporterMode = #simple)
		{
		newname = OkCancel([.rename_dialog, rptName, .Prefix], "Rename Report", hwnd)
		if false is newname
			return

		if Reporter.CheckScheduled(rptName)
			return

		.renameReport(rptName, newname, reporterMode)
		}

	Delete(rptName, reporterMode = #simple, confirmationDialog = false, hwnd = 0)
		{
		if Reporter.CheckScheduled(rptName)
			return

		result = true
		if confirmationDialog
			result = OkCancel("Delete the " $ Display(rptName) $ " report?",
				title: "Delete Report",
				:hwnd, flags: MB.ICONQUESTION)

		if result isnt true
			return // cancelled

		.deleteReport(rptName, reporterMode)
		}

	deleteReport(name, reporterMode)
		{
		QueryDo("delete " $ .reportQuery(name))
		c = LastContribution(#Reporter)
		c.AfterDeleteReport(name, :reporterMode)
		}

	rename_dialog: Controller
		{
		New(name, .prefix)
			{
			super(.layout(name))
			.field = .Vert.Horz.Field
			.field.Set(name)
			.name = name
			}

		layout(name)
			{
			return [#Vert,
				[#Static, "Rename the " $ Display(name) $ " report"],
				#Skip,
				#(Horz, (Static, To), Skip, (Field, width: 40))]
			}

		OK()
			{
			newname = .field.Get()
			if newname is ""
				{
				.AlertInfo(#Reporter, "Please enter a name to rename the report to")
				return false
				}
			if false is Reporter.CheckName(newname)
				return false

			// duplicate key
			if false is QueryEmpty?('params
				where report = ' $ Display(.prefix $ newname))
				{
				.AlertInfo(#Reporter, "Report name already used")
				return false
				}
			return newname
			}
		}

	renameReport(name, newname, reporterMode)
		{
		Transaction(update:)
			{|t|
			t.QueryApply1(.reportQuery(name))
				{|x|
				x.report = .Prefix $ newname
				x.params.report_name = newname
				x.Update()
				}

			c = LastContribution(#Reporter)
			if c.ReporterBook is ""
				return

			t.QueryDo('delete ' $ c.ReporterBook $ '
				where path =~ ' $ Display(c.GetPath(reporterMode)) $
				' and name is ' $ Display(name))

			if false is rec = t.Query1(.reportQuery(newname))
				{
				t.Rollback()
				return
				}
			}

		CustomReportsMenu(newname, GetCustomReportsSource(rec.params.Source),
			ReporterModel(rec.report).BuildReportText(), :reporterMode)
		}

	reportQuery(name)
		{
		return "params where report = " $ Display(.Prefix $ name)
		}
	}
