// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
WindowBaseComponent
	{
	CallClass(@args)
		{
		_ctrlspec = args
		new this(@args)
		}

	styles: `
		.su-lit-edit-window-container {
			position: absolute;
			top: 0;
			left: 0;
			overflow: visible;
			color: black;
			user-select: text;
			z-index: 1;
		}
		`
	New(control, readonly/*unused*/, col, row, parent)
		{
		if false is parentComponent = SuRender().GetRegisteredComponent(parent)
			throw "Cannot find parentComponent: " $ Display(parent)

		.xmin0 = .Xmin
		.ymin0 = .Ymin

		LoadCssStyles('su_list_edit_window.css', .styles)
		.ParentEl = parentComponent.GetCellEl(row, col)
		.Window = .Controller = this
		.CreateElement('div', className: 'su-lit-edit-window-container')
		.El.SetStyle('background-color', ToCssColor(CLR.ButtonFace))

		.HwndMap = SuUI.HtmlElMap(this) // Element map
		.Ctrl = .Construct(control)

		.Recalc()
		.ensureView(parentComponent)
		if .FocusFirst() is false
			.Ctrl.SetFocus()

		.El.AddEventListener('keydown', .hook, useCapture:)
		.RegisterActiveWindow()
		}

	New2()
		{
		.InitUniqueId()
		}

	Recalc()
		{
		if .Ctrl is false
			return

		if not (.Ctrl.Member?("HtmlDivComponent_ctrl") and
			.Ctrl.HtmlDivComponent_ctrl isnt false)
			{
			rect = SuRender().GetClientRect(.ParentEl)
			.Ctrl.Resize(rect.width, Max(.Ctrl.Ymin, rect.height))
			}
		.Xmin = Max(.xmin0, .Ctrl.Xmin)
		.Ymin = Max(.ymin0, .Ctrl.Ymin)
		.SetMinSize()
		}

	ensureView(parentComponent)
		{
		cellEl = .ParentEl
		scrollEl = parentComponent.GetScrollContainerEl()
		diff = (cellEl.offsetLeft + .Xmin) - (scrollEl.scrollLeft + scrollEl.clientWidth)
		if diff > 0
			.El.SetStyle('left', -diff $ 'px')
		}

	hook(event)
		{
		if event.key is 'Tab'
			{
			if event.shiftKey is true
				.sendToParent(event, -1)
			else if .Ctrl.HandleTab() is false	// if control doesn't handle tabs
				.sendToParent(event, 1)
			else
				{
				// .Ctrl should have handled the tab and set the focus
				// so prevent the default
				event.PreventDefault()
				event.StopPropagation()
				}
			}
		else if event.key is 'Enter' and event.ctrlKey is false
			.sendToParent(event, 0)
		}

	sendToParent(event, dir)
		{
		.Event('ListEditWindow_SendToParent', dir)
		event.PreventDefault()
		event.StopPropagation()
		}

	SetDefaultButton(uniqueId/*unused*/) {}
	CallDefaultButton() {}
	HighlightDefaultButton(highlight?/*unused*/) {}

	On_Cancel()
		{
		.Event('On_Cancel')
		}

	Destroy()
		{
		.Ctrl.Destroy()
		super.Destroy()
		}
	}
