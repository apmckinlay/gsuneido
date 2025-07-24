// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Print Preview"

	Commands:
		(
		("First",		"Alt+F")
		("Prev",		"Alt+P")
		("Next",		"Alt+N")
		("Last",		"Alt+L")
		("Zoom In",		"Ctrl+Add",			tip: "Ctrl+Plus")
		("Zoom Out",	"Ctrl+Subtract", 	tip: "Ctrl+Minus")
		)

	New(.report, .params = false, extraButtons = #(Print, PDF))
		{
		super(.controls(extraButtons))
		.vbox = .Vert.Scroll.Center.PreviewPage
		.first_button = .FindControl(#First)
		.last_button = .FindControl(#Last)
		.prev_button = .FindControl(#Prev)
		.next_button = .FindControl(#Next)
		.showpage('first')
		.Defer(.startup)
		}

	Activate()
		{
		PreviewLimiter().Add(.Window.Hwnd)
		}

	startup()
		{
		if .vbox.Empty?()
			{
			.Window.SetVisible(false)
			status = .vbox.GetStatus()
			if not Report.IsGeneratingReportError?(status)
				{
				// cannot alert with proper hwnd since params window could be also destroyed
				if Instance?(.params) and .params.Method?('SetNoPage')
					.params.SetNoPage()
				else
					status = ReportStatus.NODATA
				}
			if Report.IsReportError?(status)
				Report.DisplayAlert(status)
			if not .Destroyed?()
				.Window.Destroy()
			}
		}
	controls(extraButtons)
		{
		buttons = Object('First', 'Prev', 'Next', 'Last', 'Zoom In', 'Zoom Out',
			'Reset Zoom')
		if Object?(extraButtons)
			{
			extraButtons = extraButtons.Copy().Remove('Preview')
			buttons.MergeUnion(extraButtons)
			}
		buttonLayout = Object('HorzEven')
		for (i = 0; i < buttons.Size(); i++)
			{
			buttonLayout.Add('Skip')
			cmdOb = .Commands.FindOne({ it[0] is buttons[i] })
			tip = ''
			if cmdOb isnt false
				tip = cmdOb.GetDefault('tip', cmdOb[1])
			buttonOb = Object(buttons[i].Has?('PDF') ? 'PDFButton' : 'Button',
				buttons[i], command: 'DynamicButton', xstretch: 1, :tip)
			buttonLayout.Add(buttonOb)
			}
		buttonLayout.Add('Skip')
		return Object("Vert",
			Object('Scroll',
				Object('Center',
					Object('PreviewPage', .report, .getScale()),
					border: 20)),
			#(Skip medium:),
			buttonLayout,
			#(Skip medium:))
		}
	getScale()
		{
		return .report.Member?('name') and
			false isnt (params = Query1("params",
				user: Suneido.User, report: .report.name)) and
			Object?(params.report_options)
			? params.report_options.GetDefault(#previewScale, 1)
			: 1
		}

	On_DynamicButton(subMenu = false, source = false)
		{
		if source is false
			return

		buttonMethod = 'On_' $ source.Name
		if subMenu isnt false
			buttonMethod $= '_' $ ToIdentifier(subMenu)

		if .Member?(buttonMethod)
			(this[buttonMethod])()
		else if .params isnt false and .params.Member?(buttonMethod)
			.ClickParamsButton(buttonMethod)
		}

	ClickParamsButton(buttonMethod)
		{
		// have to identify that we chose print from the preview so that
		// params won't override the paramsdata with empty record (param's
		// controls are most likely destroyed by this point
		.report.from_preview = true
		.report.previewWindow = .Window.Hwnd

		Finally({
			if .params[buttonMethod].Params() is "()"
				(.params[buttonMethod])()
			else
				(.params[buttonMethod])(@.report)},
			{
			// need to check for Window,
			// because sometimes it gets destroyed in the print process
			if not .Destroyed?()
				.Window.Destroy()
			})
		}

	zoomFactor: 1.41421356 // zoom by square root of 2 so two steps to half or double
	On_Zoom_In()
		{
		.zoom('In')
		}
	On_Zoom_Out()
		{
		.zoom('Out')
		}
	Zoom(dir)
		{
		.zoom(dir)
		return true
		}
	zoom(dir)
		{
		scale = dir is 'In' ? .zoomFactor : (1/.zoomFactor)
		.vbox.Scale(scale)
		.Resize(.x, .y, .w, .h)
		.Resize(.x, .y, .w, .h) // twice to propogate size change
		}

	On_Reset_Zoom()
		{
		.vbox.ResetScale()
		.Resize(.x, .y, .w, .h)
		.Resize(.x, .y, .w, .h) // twice to propogate size change
		}

	On_First()
		{
		.vbox.FirstPage()
		.showpage('first')
		}
	On_Next()
		{
		if false is .vbox.NextPage()
			.On_Last()
		else
			.showpage()
	}
	On_Prev()
		{
		if false is .vbox.PrevPage()
			.On_First()
		else
			.showpage()
		}
	On_Last()
		{
		// suneido.js returns false if user aborts the load
		finished? = .vbox.LastPage()
		.showpage(finished? is true ? 'last' : '')
		}
	showpage(firstlast = '')
		{
		.buttonSetEnabled(firstlast)
		.Send('SetTitle',
			TranslateLanguage('Print Preview - Page %1', Display(.vbox.GetPageNum() + 1)))
		}
	GetNumPages()
		{
		return .vbox.GetNumPages()
		}
	buttonSetEnabled(position)
		{
		.last_button.SetEnabled(position isnt 'last')
		.next_button.SetEnabled(position isnt 'last')
		.first_button.SetEnabled(position isnt 'first')
		.prev_button.SetEnabled(position isnt 'first')
		}
	Scroll(x, y)
		{
		.Vert.Scroll.Scroll(x, y)
		}
	Resize(.x, .y, .w, .h)
		{
		super.Resize(x, y, w, h)
		}
	updatePreviewScale()
		{
		if not .report.Member?('name')
			return
		QueryApply1("params", user: Suneido.User, report: .report.name)
			{|x|
			if not Object?(x.report_options)
				x.report_options = Object()
			x.report_options.previewScale = .vbox.GetScale()
			x.Update()
			}
		}
	Destroy()
		{
		.report.Delete("from_preview")
		.report.Delete("previewWindow")
		PreviewLimiter().Remove(.Window.Hwnd)
		.updatePreviewScale()
		super.Destroy()
		}
	}
