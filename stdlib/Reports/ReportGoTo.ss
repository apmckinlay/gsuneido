// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(rpt, params, hwnd)
		{
		if false is AccessPermissions(rpt)
			{
			Alert("You do not have permission to access this option.", "Report",
				hwnd, MB.ICONWARNING)
			return
			}

		result = .get_format(rpt)
		if result.fmt is false
			return

		preview = Object(Object(result.fmt), PreviewParams: params,
			previewDialog: , hwndOwner: hwnd, previewWindow: hwnd)
		if result.title isnt false
			preview.title = result.title
		if result.name isnt false
			preview.name = result.name
		Params.On_Preview(@preview)
		}
	get_format(rpt)
		{
		report = rpt.Eval() // needs Eval
		if not Object?(report)
			report = Global(rpt)()

		fmt = .getFormatClass(report)
		name = report.Member?('name') and String?(report.name)
			? report.name : false
		title = report.Member?('header') and report.header is false
			? false
			: report.Member?('title') ? report.title : false
		return Object(:fmt, :title, :name)
		}
	getFormatClass(report)
		{
		return report.Member?(1)
			? Class?(report[1])
				? report[1]
				: String?(report[1])
					? Global(report[1])
					: false
			: false
		}
	}