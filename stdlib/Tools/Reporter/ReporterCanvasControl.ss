// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	Title: "Reporter Canvas Control"
	Name: 'ReporterCanvas'

	ready?: false
	New()
		{
		super(.layout())
		.drawControl = .FindControl('Draw')
		.data = .FindControl('Data')
		.data.Set([report_page_size: ReportPagePaperSpecs.Default.page,
			report_page_orientation: ReportPagePaperSpecs.Default.orientation])
		.drawControl.On_Select()
		.Send(#Data)
		.ready? = true
		}

	layout()
		{
		extraItems = GetContributions('ReporterExtraCanvasItems')
		return Object('Vert',
			Object('Draw',
				width: .width = ReportPagePaperSpecs.Default.w,
				height: .height = ReportPagePaperSpecs.Default.h,
				:extraItems, noJustify?:)
			Object('Record'
				#(GroupBox, 'Canvas',
					#('Horz', 'report_page_size', Skip, 'report_page_orientation'))))
		}

	NewValue()
		{
		.onChanged()
		}

	onChanged()
		{
		.dirty? = true
		.Send('NewValue', .Get())
		}

	Record_NewValue(@unused)
		{
		if not .ready? or .data.Valid() isnt true
			return

		data = .data.Get()
		.width = ReportPagePaperSpecs.GetWidth(data.report_page_size,
			data.report_page_orientation)
		.height = ReportPagePaperSpecs.GetHeight(data.report_page_size,
			data.report_page_orientation)
		.drawControl.SetSize(.width, .height)

		.onChanged()
		}

	Get()
		{
		ob = .drawControl.Get()
		ob.canvasWidth = .width
		ob.canvasHeight = .height
		ob.orientation = .data.Get().report_page_orientation
		return ob
		}

	Set(value)
		{
		if value is ''
			return // empty canvas
		.dirty? = false
		width = value.GetDefault(#canvasWidth, ReportPagePaperSpecs.Default.w)
		height = value.GetDefault(#canvasHeight, ReportPagePaperSpecs.Default.h)
		.data.Set([
			report_page_size: ReportPagePaperSpecs.GetPageSelection(width, height),
			report_page_orientation:
				value.GetDefault(#orientation, ReportPagePaperSpecs.Default.orientation)
			])
		if width isnt .width or height isnt .height
			.drawControl.SetSize(.width = width, .height = height)
		.drawControl.Set(value)
		}

	GetDAFs(items = false)
		{
		if items is false
			items = .drawControl.GetAllItems()
		return .getDAFs(items).Sort!().Unique!()
		}

	getDAFs(items)
		{
		ob = Object()
		ForEachDAF(items)
			{
			ob.Add(it.GetField())
			}
		return ob
		}

	dirty?: false
	Dirty?(state = "")
		{
		if state isnt ''
			.dirty? = state
		return .dirty?
		}

	Valid?()
		{
		return true
		}

	CheckContent()
		{
		ForEachDAF(.drawControl.GetAllItems())
			{
			if not it.Valid?()
				return 'Invalid Data Field'
			}
		return ''
		}

	GetDAFAvailableCols()
		{
		return .Send('GetRptDesignCols')
		}

	PromptToField(prompt)
		{
		return .Send('PromptToField', prompt)
		}

	FieldToPrompt(field)
		{
		return .Send('FieldToPrompt', field)
		}

	DesignChanged()
		{
		ForEachDAF(.drawControl.GetAllItems(),
			{ it.DesignChanged(.drawControl.GetCanvas()) })
		}

	Destroy()
		{
		.Send(#NoData)
		super.Destroy()
		}
	}
