// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// TODO: Make DrawColorControl more intuitive and add to this control
	Title: "Draw"
	Name: "Draw"
	readOnly: false
	extraValues: #()
	New(valueOb = false, .readOnly = false, .extraItems = #(),
		.width = false, .height = false, .noJustify? = false)
		{
		.canvas = .FindControl('Canvas')
		.palette = .FindControl('Palette')
		.palette.SetButtons('select')
		.justify = .FindControl('RadioButtons')
		.Redir('On_Cut', this)
		// On_Copy does not work here and should be redirected to CanvasControl
		.Redir('On_Copy', .canvas)
		.Redir('On_Paste', this)
		.Redir('On_Delete', this)
		.Redir('On_Select_All', this)
		.Redir('On_Undo', this)
		.Redir('On_Redo', this)
		.undoHistory = Object()
		.redoHistory = Object()
		if valueOb isnt false
			.Set(valueOb)
		if 0 isnt editable = .Send('EditMode?')
			.SetReadOnly(editable is false)
		}
	ready?: false
	Startup()
		{
		.ready? = true
		}
	SetReadOnly(readOnly)
		{
		if not readOnly = readOnly or .readOnly
			.On_Select()
		else
			.palette.SetButtons(false)
		.palette.SetEnabled(not readOnly)
		super.SetReadOnly(readOnly)
		}
	SetJustification(justification)
		{
		.justification = justification
		}
	Controls()
		{
		canvas = Object('DrawCanvas')
		if .height isnt false and .width isnt false
			{
			canvas = Object('Scroll',
				Object('DrawCanvas',
					xmin: .width.InchesInCanvasUnit(),
					ymin: .height.InchesInCanvasUnit(),
					xstretch: false,
					ystretch: false)
				)
			}

		lay = Object('Vert'
				Object('EtchedLine', before: 0 after: 0)
				Object('Horz'
					Object('DrawPalette', .extraItems)
					canvas)
				Object('EtchedLine', before: 0 after: 0)
			)

		if .noJustify? is false
			lay.Add(Object('RadioButtons' 'Left' 'Center' 'Right' horz:))

		return lay
		}

	SetSize(.width, .height)
		{
		.canvas.SetXminYmin(
			.width.InchesInCanvasUnit(),
			.height.InchesInCanvasUnit())
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
		(Undo               "Ctrl+Z",   "Undo on the canvas")
		(Redo               "Ctrl+Y",   "Redo on the canvas")
		)
	Menu:
		(
		("&Edit",
			"Cu&t", "&Copy", "&Paste", "Resize", "&Delete", "Select &All", "",
			"Group", "Ungroup", "", "Move To Front", "Move To Back", "", "Print",
			"Undo", "Redo", "Lock/Unlock")
		)
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
	On_Line()
		{
		.palette.SetButtons('line')
		.canvas.SetTracker(DrawLineTracker, CanvasLine)
		}
	On_Rectangle()
		{
		.palette.SetButtons('rectangle')
		.canvas.SetTracker(DrawRectTracker, CanvasRect)
		}
	On_RoundRectangle()
		{
		.palette.SetButtons('roundrect')
		.canvas.SetTracker(DrawRectTracker, CanvasRoundRect)
		}
	On_LockUnlock()
		{
		for i in .canvas.GetSelected()
			{
			i.ToggleLock()
			}
		}
	Repaint()
		{
		.canvas.Repaint()
		}
	On_Ellipse()
		{
		.palette.SetButtons('ellipse')
		.canvas.SetTracker(DrawRectTracker, CanvasEllipse)
		}
	On_Arc()
		{
		.palette.SetButtons('arc')
		.canvas.SetTracker(DrawRectTracker, DrawArcAdapter)
		}
	On_Text()
		{
		.palette.SetButtons('text')
		.canvas.SetTracker(DrawClickTracker, DrawTextAdapter)
		}
	On_Image()
		{
		.palette.SetButtons('image')
		.canvas.SetTracker(DrawClickTracker, DrawImageAdapter)
		}
	On_Select()
		{
		.palette.SetButtons('select')
		.canvas.SetTracker(DrawSelectTracker, false)
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
	On_Undo()
		{
		len = .undoHistory.Size()
		if len >= 2
			{
			.redoHistory.Add(.undoHistory[len - 1])
			.doWithoutUndo()
				{
				.Set(.undoHistory[len - 2])
				}
			.undoHistory.PopLast()
			}
		}
	On_Redo()
		{
		if .redoHistory.Empty?()
			return
		redoHistory = .redoHistory[..-1]
		.Set(.redoHistory.Last())
		.redoHistory = redoHistory
		}
	undoing?: false
	doWithoutUndo(block)
		{
		.undoing? = true
		block()
		.undoing? = false
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

	Recv(@args)
		{
		if .extraItems.Empty?() or not args[0].Prefix?('On_')
			return 0

		name = args[0].RemovePrefix('On_')
		if false isnt item = .extraItems.FindOne({ it.name is name })
			{
			(item.handler)(.canvas, .palette)
			}

		return 0
		}

	CanvasChanged()
		{
		if .ready?
			.Send('NewValue')
		.addToUndoHistory()
		}

	GetDAFAvailableCols()
		{
		return .Send('GetDAFAvailableCols')
		}
	PromptToField(prompt)
		{
		return .Send('PromptToField', prompt)
		}
	FieldToPrompt(field)
		{
		return .Send('FieldToPrompt', field)
		}

	addToUndoHistory()
		{
		if .undoing? is true
			return

		if .undoHistory.Empty?() or .undoHistory.Last() isnt .Get()
			{
			.undoHistory.Add(.Get())
			.redoHistory = Object()
			if .undoHistory.Size() > 20 /* = length of undo history object */
				.undoHistory.PopFirst()
			}
		}
	NewValue(value /*unused*/, source)
		{
		if .ready? and source is .justify // radio button changed
			.Send('NewValue')
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
	GetAllItems()
		{
		return .canvas.GetAllItems()
		}
	Get()
		{
		return Object(items: .canvas.Get(),
			justify: .justify is false ? 'Left' : .justify.Get()).
			MergeNew(.extraValues)
		}
	BuildItems(valueOb, ignoreOld?/*unused*/ = false)
		{
		items = Object()
		if valueOb.Member?('items') and valueOb.items.Every?(Object?)
			{
			for item in valueOb.items
				items.Add(Construct(item).SetupScale())
			return items
			}

		if ((not valueOb.Member?('resources')) or (not valueOb.Member?('items')))
			return items

//		if not ignoreOld?
//			ProgrammerError('Found old canvas format', params: valueOb,
//				caughtMsg: 'handled')

		valueOb = valueOb.Copy()
		valueOb.resources = valueOb.resources.Copy()
		valueOb.items.Each()
			{
			// CanvasImage saves .references so we need to use Compile and object.Eval
			fn = ("function() { " $ it $ ".SetupScale() }").Compile()
			item = valueOb.resources.Eval(fn)
			items.Add(item)
			}
		return items
		}
	Set(valueOb, delAll = true)
		{
		if delAll
			.canvas.DeleteAll()

		if .justify isnt false
			.justify.Set(valueOb.GetDefault('justify', 'Left'))
		_canvas = .canvas
		items = .BuildItems(valueOb)
		if items.Empty?() is false
			{
			for item in items
				.canvas.AddItem(item)
			.canvas.Repaint()
			}
		.extraValues = valueOb.Copy().Remove(#items, #justify, #resources)
		}
	Canvas_LButtonUp()
		{
		.On_Select()
		}
	SetColor(color)
		{
		.canvas.SetColor(color)
		}
	SetLineColor(color)
		{
		.canvas.SetLineColor(color)
		}
	AddItem(item)
		{
		.canvas.AddItem(item).SetColor(.canvas.GetColor())
		}
	GetCanvas()
		{
		return .canvas
		}
	}
