// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
// NOTE: This Control could be made more generic if needed to be used in other places
// DESC: This is technically an override of DrawControl()
Controller
	{
	Title: 'Multi Canvas Control'
	Name: 'MultiCanvas'
	ready?: false

	New()
		{
		.canvas  = false // to be set by parent
		.palette = .FindControl('commonPalette')
		.data = .FindControl('Data')

		.data.Set([report_page_size: ReportPagePaperSpecs.Default.page,
			report_page_orientation: ReportPagePaperSpecs.Default.orientation])
		.page_header = .FindControl('page_header')
		.section_header = .FindControl('section_header')
		.section_footer = .FindControl('section_footer')
		.page_footer = .FindControl('page_footer')

		.addDefaultHeader()

		.Redir('On_Cut', this)
		// On_Copy does not work here and should be redirected to CanvasControl
		.Redir('On_Copy', .canvas)
		.Redir('On_Paste', this)
		.Redir('On_Delete', this)
		.Redir('On_Select_All', this)
		.Redir('On_Undo', this)
		.Redir('On_Redo', this)

		.undoHistory = Object(page_header: Object(), section_header: Object(),
			section_footer: Object(), page_footer: Object())
		.redoHistory = .undoHistory.DeepCopy()

		.FindControl('accordion_columns').ExpandAll()
		.expandIfNotEmpty()

		.On_Select()
		.Send(#Data)
		.ready? = true
		}
	Controls()
		{
		skipNum = 20
		.width = ReportPagePaperSpecs.Default.w
		.height = ReportPagePaperSpecs.Default.h
		canvasWidth = .width.InchesInCanvasUnit()
		.canvasHeight = 1.InchesInCanvasUnit()
		.extraItems = GetContributions('ReporterExtraCanvasItems')
		return Object('Vert'
				Object('DrawPalette', horizontal?:, extraItems: .extraItems,
					name: 'commonPalette')
				Object('Scroll'
				Object('Vert'
				Object('Accordion'
					Object('Page Header'
						Object('DrawCanvas', name: 'page_header',
							xmin: canvasWidth, ymin: .canvasHeight,
							xstretch: false, ystretch: false)
					),
				name: 'accordion_page_header')
				Object('Accordion'
					Object('Section Header'
						Object('DrawCanvas', name: 'section_header',
							xmin: canvasWidth, ymin: .canvasHeight,
							xstretch: false, ystretch: false)
					),
				name: 'accordion_section_header')
				Object('Accordion'
					Object('Reporter Columns'
					Object('Vert',
						Object('WndPane', Object('Border'
							Object('ReporterColumns'))),
						Object('Skip', skipNum)
						Object('Horz'
							#('Button', 'Add/Remove Columns...' pad: 30)
							Object('Skip', skipNum)
							#('Static', 'Drag columns to resize or rearrange.')
							Object('Skip', skipNum)
							#('Static',
								'Click on a column to change the heading or to total.')
							)
						)
					)
				name: 'accordion_columns')
				Object('Accordion'
					Object('Section Footer'
						Object('DrawCanvas', name: 'section_footer',
							xmin: canvasWidth, ymin: .canvasHeight,
							xstretch: false, ystretch: false)
					),
				name: 'accordion_section_footer')
				Object('Accordion'
					Object('Page Footer'
						Object('DrawCanvas', name: 'page_footer',
							xmin: canvasWidth, ymin: .canvasHeight,
							xstretch: false, ystretch: false)
					),
				name: 'accordion_page_footer')
			) noEdge:)
			Object('Record'
				#(GroupBox, 'Canvas',
					#('Horz', 'report_page_size', Skip, 'report_page_orientation'))))
		}
	Commands:
		(
		(Print, 			"Ctrl+P", 	"Print the current image")
		(Cut,				"Ctrl+X",	"Cut selected items")
		(Copy,				"Ctrl+C",	"Copy selected items")
		(Paste,				"Ctrl+V", 	"Paste items from clipboard")
		(Delete,			"Del", 		"Paste items from clipboard")
		(Select_All,		"Ctrl+A", 	"Select all items")
		(Group,				"", 		"Group selected items")
		(Ungroup,			"",			"Ungroup selected group")
		(Move_To_Front,		"",			"Move selected items to front")
		(Move_To_Back,		"", 		"Move selected items to back")
		(Resize,			"",			"Resize the selected item")
		)
	Menu:
		(
		("&Edit",
			"Cu&t", "&Copy", "&Paste", "Resize", "&Delete", "Select &All", "",
			"Group", "Ungroup", "", "Move To Front", "Move To Back", "", "Print",
			"Lock/Unlock")
		)
	Data() {}
	NoData() {}
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
		.forEachCanvas()
			{ |canvas|
			ForEachDAF(canvas.GetAllItems())
				{
				if not it.Valid?()
					return 'Invalid Data Field'
				}
			}
		return ''
		}

	addDefaultHeader()
		{
		builtReportHeader = OptContribution('MultiCanvasDefaultHeader', { #() })()
		for i in builtReportHeader
			.page_header.AddItem(i)
		}
	set(items)
		{
		if .canvas is false or items is ''
			return

		_canvas = .canvas
		items = DrawControl.BuildItems(items)
		.canvas.DeleteAll()
		if not items.Empty?()
			for item in items
				.canvas.AddItem(item)
		}
	Set(canvasses, defaultHeader = false)
		{
		if canvasses is '' or canvasses is .cur
			return

		.cur = canvasses

		width = canvasses.GetDefault(#canvasWidth, ReportPagePaperSpecs.Default.w)
		height = canvasses.GetDefault(#canvasHeight, ReportPagePaperSpecs.Default.h)
		.data.Set([
			report_page_size: ReportPagePaperSpecs.GetPageSelection(width, height),
			report_page_orientation:
				canvasses.GetDefault(#orientation,
					ReportPagePaperSpecs.Default.orientation)
			])
		if .width isnt width
			{
			.width = width
			canvasWidth = .width.InchesInCanvasUnit()
			.forEachCanvas({ it.SetXminYmin(canvasWidth, .canvasHeight) })
			}

		for canvas in .allCanvasses
			{
			_canvas = canvas
			items = DrawControl.BuildItems(canvasses.GetDefault(canvas.Name, #()))
			canvas.DeleteAll()
			if not items.Empty?()
				for item in items
					canvas.AddItem(item)
			}
		if defaultHeader
			.addDefaultHeader()

		.expandIfNotEmpty()
		}
	get()
		{
		if .canvas is false
			return #()

		return Object(items: .canvas.Get())
		}
	Get()
		{
		allCanvassesGet = Object()
		for canvas in .allCanvasses
			{
			allCanvassesGet[canvas.Name] = Object(
				items: canvas.Get(),
				canvasWidth: .width,
				canvasHeight: .height)
			}
		allCanvassesGet.orientation = .data.Get().report_page_orientation
		allCanvassesGet.canvasWidth = .width
		allCanvassesGet.canvasHeight = .height
		return allCanvassesGet
		}
	cur: false
	CanvasChanged()
		{
		.dirty? = true
		.Send('NewValue', .cur = .Get())
		.addToUndoHistory()
		}

	NewValue()
		{
		.onChanged()
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
		canvasWidth = .width.InchesInCanvasUnit()
		.forEachCanvas({ it.SetXminYmin(canvasWidth, .canvasHeight) })

		.onChanged()
		}

	onChanged()
		{
		.dirty? = true
		.Send('NewValue', .Get())
		}

	expandIfNotEmpty()
		{
		.forEachCanvas()
			{|canvas|
			.FindControl('accordion_' $ canvas.Name).ContractAll()
			if canvas.GetAllItems() isnt #()
				{
				.FindControl('accordion_' $ canvas.Name).ExpandAll()
				}
			.canvas.ClearSelect()
			}
		}
	addToUndoHistory()
		{
		if .undoing? is true or .canvas is false
			return

		if .undoHistory[.canvas.Name].Empty?() or
			.undoHistory[.canvas.Name].Last() isnt .get()
			{
			.undoHistory[.canvas.Name].Add(.get())
			.redoHistory[.canvas.Name] = Object()
			if .undoHistory[.canvas.Name].Size() > 20 /*=length of undo history object*/
				.undoHistory[.canvas.Name].PopFirst()
			}
		}
	Repaint()
		{
		.allCanvasses.Each({|x| x.Repaint() })
		}
	WhichDrawCanvasClicked(name /*unused*/, source)
		{
		.canvas = source
		.forEachCanvas(){|canvas|
			if canvas.Name isnt .canvas.Name {canvas.ClearSelect()} }
		}
	getter_allCanvasses()
		{
		return Object(.page_header, .section_header, .section_footer, .page_footer)
		}
	forEachCanvas(block)
		{
		for canvas in .allCanvasses
			block(canvas)
		}
	Expand(source)
		{
		drawCanvasName = source.Label.Lower().Tr(' ', '_')
		.canvas = .FindControl(drawCanvasName)
		}
	On_Edit()
		{ .canvas.EditItem() }
	On_Cut()
		{ .canvas.CutItems() }
	On_Paste()
		{ .canvas.PasteItems() }
	On_Delete()
		{ .canvas.DeleteSelected() }
	On_Select_All()
		{ .canvas.SelectAll() }
	On_Select()
		{
		.palette.SetButtons('select')
		.forEachCanvas() {|canvas| canvas.SetTracker(DrawSelectTracker, false) }
		}
	On_Line()
		{
		.palette.SetButtons('line')
		.forEachCanvas() {|canvas| canvas.SetTracker(DrawLineTracker, CanvasLine) }
		}
	On_Ellipse()
		{
		.palette.SetButtons('ellipse')
		.forEachCanvas() {|canvas| canvas.SetTracker(DrawRectTracker, CanvasEllipse) }
		}
	On_Arc()
		{
		.palette.SetButtons('arc')
		.forEachCanvas() {|canvas| canvas.SetTracker(DrawRectTracker, DrawArcAdapter) }
		}
	On_Text()
		{
		.palette.SetButtons('text')
		.forEachCanvas() {|canvas| canvas.SetTracker(DrawClickTracker, DrawTextAdapter) }
		}
	On_Image()
		{
		.palette.SetButtons('image')
		.forEachCanvas() {|canvas| canvas.SetTracker(DrawClickTracker, DrawImageAdapter) }
		}
	On_Rectangle()
		{
		.palette.SetButtons('rectangle')
		.forEachCanvas() {|canvas| canvas.SetTracker(DrawRectTracker, CanvasRect) }
		}
	On_RoundRectangle()
		{
		.palette.SetButtons('roundrect')
		.forEachCanvas() {|canvas| canvas.SetTracker(DrawRectTracker, CanvasRoundRect) }
		}
	On_LockUnlock()
		{
		for i in .canvas.GetSelected()
			{
			i.ToggleLock()
			}
		}
	On_Group()
		{
		.palette.SetButtons(false)
		items = .canvas.GetSelected()
		if (items.Size() <= 1)
			{
			Alert('You must select at least two items to group.'title: 'Error',
				flags: MB.ICONERROR)
			return
			}
		group = Object()
		for (item in items.Copy())
			{
			group.Add(item)
			.canvas.RemoveItem(item)
			}
		newitem = CanvasGroup(group)
		.canvas.AddItemAndSelect(newitem)
		}
	On_Ungroup()
		{
		.palette.SetButtons(false)
		group = .canvas.GetSelected()
		if (group.Size() isnt 1)
			{
			Alert('You must select one item to ungroup.', title: 'Error',
				flags: MB.ICONERROR)
			return
			}
		if not group[0].Base?(CanvasGroup)
			return
		for (item in group[0].GetItems())
			.canvas.AddItemAndSelect(item)
		.canvas.RemoveItem(group[0])
		}
	On_Move_To_Back()
		{
		.palette.SetButtons(false)
		for item in .canvas.GetSelected().Copy()
			.canvas.MoveToBack(item)
		}
	On_Move_To_Front()
		{
		.palette.SetButtons(false)
		for item in .canvas.GetSelected().Copy()
			.canvas.MoveToFront(item)
		}
	On_Resize()
		{
		.palette.SetButtons(false)
		item = Object()
		item = .canvas.GetSelected()
		if (item.Size() isnt 1)
			{
			Alert('You must select one item to resize.', title: 'Error',
				flags: MB.ICONERROR)
			return
			}
		.canvas.ResetSize(item[0])
		SetFocus(.canvas.Hwnd)
		}
	On_Print()
		{
		.palette.SetButtons(false)
		items = .canvas.GetAllItems()
		fmt = .WrapItems(items)
		justifiedFmt = .ApplyJustify(fmt)
		Params.On_Preview(justifiedFmt, previewWindow: GetFocus()
			/*, default_orientation: 'Landscape'*/)
		}
	undoing?: false
	doWithoutUndo(block)
		{
		.undoing? = true
		block()
		.undoing? = false
		}
	On_Undo()
		{
		len = .undoHistory[.canvas.Name].Size()
		if len >= 2
			{
			.redoHistory[.canvas.Name].Add(.undoHistory[.canvas.Name][len - 1])
			.doWithoutUndo()
				{
				.set(.undoHistory[.canvas.Name][len - 2])
				}
			.undoHistory[.canvas.Name].PopLast()
			}
		}
	On_Redo()
		{
		if .redoHistory[.canvas.Name].Empty?()
			return
		redoHistory = .redoHistory[.canvas.Name][..-1]
		.set(.redoHistory[.canvas.Name].Last())
		.redoHistory[.canvas.Name] = redoHistory
		}

	Recv(@args)
		{
		if .extraItems.Empty?() or not args[0].Prefix?('On_')
			return 0

		name = args[0].RemovePrefix('On_')
		if false isnt item = .extraItems.FindOne({ it.name is name })
			{
			canvases = item.needData?
				? Object(.section_header, .section_footer)
				: .allCanvasses
			for canvas in canvases
				(item.handler)(canvas, .palette)
			}

		return 0
		}

	WndPane_ContextMenu(x, y)
		{
		return .Send('WndPane_ContextMenu', x, y)
		}

	On_AddRemove_Columns()
		{
		return .Send('On_AddRemove_Columns')
		}

	Canvas_LButtonUp()
		{
		.On_Select()
		}

	GetDAFAvailableCols()
		{
		return .Send('GetRptSortCols')
		}

	PromptToField(prompt)
		{
		return .Send('PromptToField', prompt)
		}

	FieldToPrompt(field, source)
		{
		if source not in (.section_header, .section_footer)
			return false
		return .Send('FieldToPrompt', field)
		}

	DesignChanged()
		{
		Object(.section_header, .section_footer).Each()
			{ |canvas|
			ForEachDAF(canvas.GetAllItems(), { it.DesignChanged(canvas) })
			}
		}

	WrapItems(items)
		{
		return Object('DrawItem', CanvasGroup(items))
		}
	ApplyJustify(fmt, drawOb = false)
		{
		if drawOb is false
			drawOb = .Get()

		switch drawOb.GetDefault('justify', 'Left')
			{
		case 'Left' :
			return Object('Horz' fmt 'Hfill' xstretch:1)
		case 'Right' :
			return Object('Horz' 'Hfill' fmt xstretch:1)
		case 'Center' :
			return Object('Horz' 'Hfill' fmt 'Hfill' xstretch:1)
			}
		}
	Destroy()
		{
		.Send(#NoData)
		super.Destroy()
		}
	}
