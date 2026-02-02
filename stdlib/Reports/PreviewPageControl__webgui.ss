// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: "PreviewPage"
	ComponentName: 'PreviewPage'
	New(report, scale = 1, dimens = false)
		{
		.report = Report(@report)
		.report.SetDriver(new HtmlPreviewDriver)

		.dimens = dimens is false ? .report.GetDimens() : dimens
		.pg = 0; // the page currently displayed
		.npages = 0; // the number of pages generated
		.NextPage()
		.Scale(scale)
		.Defer({ if .Member?('Hwnd') SetFocus(.Hwnd) })
		.ComponentArgs = Object(scale, dimens)
		}

	scale: 1
	minZoom: .25
	maxZoom: 4
	Scale(by) // 2 = double, .5 = half
		{
		x = .scale * by
		if x > .minZoom and x < .maxZoom
			.scale = x
		.Act('UpdateScale', .scale)
		}
	ResetScale()
		{
		.scale = 1
		.Act('UpdateScale', .scale)
		}
	GetScale()
		{
		return .scale
		}
	FirstPage()
		{
		.Act(#DisplayPage, .pg = 0)
		}
	LastPage()
		{
		finished? = DoTaskWithPause('Working...', .nextpage)
		if not .Empty?()
			.Act(#DisplayPage, .pg = Max(0, .npages - 1))
		return finished?
		}
	NextPage()
		{
		if (.pg + 1 < .npages)
			++.pg
		else if (not .nextpage())
			{
			Beep()
			return false
			}
		.Act(#DisplayPage, .pg)
//		.sortAccessPoints()
		return true
		}
	eof: false
	status: ''
	nextpage()
		{
		if (.eof)
			return false
		page = .report.AddPage(.dimens)
		if not .report.NextPageSuccess?(vbox = .report.NextPage())
			{
			.eof = true
			.report.EndPage()
			if Report.IsGeneratingReportError?(vbox)
				{
				if .pg isnt 0 // first page handled by PreviewControl
					.report.DisplayAlert(vbox)
				.status = vbox
				}
			return false
			}
		.report.Paint(vbox)
		.report.EndPage()

		.pg = .npages++

		page.accessPoints = .report.
			GetDefault(#AccessPoints, Object()).
			GetDefault(.pg, Object()).
			Map({ it.Project(#left, #top, #right, #bottom) })

		.Act(#SetPage, .pg, page)
		return true
		}

	AccessPointClicked(idx)
		{
		pagePoints = .report.GetDefault(#AccessPoints, Object()).GetDefault(.pg, Object())
		if not pagePoints.Member?(idx)
			return
		.access(pagePoints[idx])
		}

	access(target)
		{
		access = target.access
		if access.control is 'AccessGoTo'
			AccessGoTo(access.access, access.goto_field, access.goto_value, .Hwnd)
		else if access.control is 'ReportGoTo'
			ReportGoTo(access.report, access.params, .Hwnd)
		else if access.control is 'AttachmentGoTo'
			AttachmentGoTo(access.file, .Hwnd)
		else if access.control is 'DynamicGoTo'
			DynamicGoTo(access.data, access.func, .Hwnd)
		}

	GetStatus()
		{
		return .status
		}

	PrevPage()
		{
		if (.pg is 0)
			{
			Beep()
			return false
			}
		else
			{
			.Act(#DisplayPage, --.pg)
			return true
			}
		}

	Empty?()
		{
		return .npages is 0
		}

	GetPageNum()
		{
		return .pg
		}

	GetNumPages()
		{
		return .npages
		}

	Destroy()
		{
		.report.Close(false)
		super.Destroy()
		}
	}
