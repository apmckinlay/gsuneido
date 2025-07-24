// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	CallClass(layout)
		{
		ToolDialog(0, Object(this, layout), 'Page Setup')
		}

	New(layout)
		{
		pageSize = ReportPagePaperSpecs.GetPageSelection(layout.width, layout.height,
			defaultVal: 'Custom')
		size = pageSize is 'Custom'
			? [w: layout.width, h: layout.height]
			: pageSize
		rec = [
			supage_page_size: pageSize,
			supage_width: ReportPagePaperSpecs.GetWidth(size, 'Portrait'),
			supage_height: ReportPagePaperSpecs.GetHeight(size, 'Portrait'),
			supage_orientation: pageSize is 'Custom' or layout.width < layout.height
				? 'Portrait' : 'Landscape',
			supage_left: layout.left,
			supage_right: layout.right,
			supage_top: layout.top,
			supage_bottom: layout.bottom]
		.Data.Set(rec)
		}

	Controls()
		{
		return Object('Record', Object('Vert',
			#(GroupBox, 'Paper', ('Form'
				#('supage_page_size', group: 0), nl,
				#('supage_width', group: 0), #('Static', 'inches'), nl,
				#('supage_height', group: 0), #('Static', 'inches'))),
			#Skip,
			Object('Horz',
				#(GroupBox, 'Orientation', #('NoPrompt', 'supage_orientation')),
				#Skip,
				#(GroupBox, 'Margins (inches)', #('Form',
					#('supage_left', group: 0), #('supage_right', group: 1), nl,
					#('supage_top', group: 0), #('supage_bottom', group: 1)))),
			#OkCancel))
		}

	Record_NewValue(field, value)
		{
		if field isnt 'supage_page_size'
			return

		if ReportPagePaperSpecs.Options().Has?(value)
			{
			.Data.SetField(#supage_width,
				ReportPagePaperSpecs.GetWidth(value, 'Portrait'))
			.Data.SetField(#supage_height,
				ReportPagePaperSpecs.GetHeight(value, 'Portrait'))
			}
		}

	On_OK()
		{
		if not .Data.Dirty?()
			{
			.On_Cancel()
			return
			}

		if .Data.Valid() isnt true
			return

		rec = .Data.Get()
		size = rec.supage_page_size is 'Custom'
			? [w: rec.supage_width, h: rec.supage_height]
			: rec.supage_page_size
		.Window.Result([
			width: ReportPagePaperSpecs.GetWidth(size, rec.supage_orientation),
			height: ReportPagePaperSpecs.GetHeight(size, rec.supage_orientation),
			left: rec.supage_left,
			right: rec.supage_right,
			top: rec.supage_top,
			bottom: rec.supage_bottom])
		}

	On_Cancel()
		{
		.Window.Result(false)
		}
	}