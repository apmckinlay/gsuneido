// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	left: false
	right: false
	ContextMenu: #('Snap Left/Right', ('Left', 'Right'))

	ContextCall(args, window)
		{
		if args[0] isnt 'On_Context_Snap_LeftRight'
			return
		pos = args.item[0]
		window.SetState(WindowPlacement.snapped)
		viewport = SuRender().GetClientRect()
		if pos is 'Left'
			.snapLeft(viewport, window)
		else
			.snapRight(viewport, window)
		}

	snapLeft(viewport, window)
		{
		percent = .right isnt false ? 100/*=100%*/ - .right.percent : 50
		min = window.Ctrl.Xmin + 2
		percent = Max(percent, (min / viewport.width).DecimalToPercent())
		window.SetStyles(
			Object(left: '0px', top: '0px', right: '',
				width: .getWidth(percent, min),
				height: 'calc(100% - (' $
					SuRender().Taskbar.GetTaskbarHeight() $ '))'),
			window.GetContainerEl())
		window.SetStyles(
			Object(width: '100%',
				height: 'calc(100% - ' $ window.NonClientHeight $ 'em)'))
		window.GetResizes().Each()
			{
			it.SetStyle('display',
				it.className isnt 'su-window-right-resize' ? 'none' : 'initial')
			}
		.left = Object(:window, :percent)
		}

	snapRight(viewport, window)
		{
		percent = .left isnt false ? 100/*=100%*/ - .left.percent : 50
		min = window.Ctrl.Xmin + 2
		percent = Max(percent, (min / viewport.width).DecimalToPercent())
		window.SetStyles(
			Object(left: '', top: '0px', right: '0px',
				width: .getWidth(percent, min),
				height: 'calc(100% - (' $
					SuRender().Taskbar.GetTaskbarHeight() $ '))'),
			window.GetContainerEl())
		window.SetStyles(
			Object(width: '100%',
				height: 'calc(100% - ' $ window.NonClientHeight $ 'em)'))
		window.GetResizes().Each()
			{
			it.SetStyle('display',
				it.className isnt 'su-window-left-resize' ? 'none' : 'initial')
			}
		.right = Object(:window, :percent)
		}

	getWidth(percent, min)
		{
		return 'max(calc(' $ percent $ '% - 2px), ' $ min $ 'px)'
		}

	Remove(window)
		{
		if .left isnt false and Same?(.left.window, window)
			.left = false
		else if .right isnt false and Same?(.right.window, window)
			.right = false
		}

	HorzResize(newWidth, window)
		{
		if .left isnt false and Same?(.left.window, window)
			{
			.resizeLeft(newWidth)
			return true
			}
		else if .right isnt false and Same?(.right.window, window)
			{
			.resizeRight(newWidth)
			return true
			}
		return false
		}

	resizeLeft(newWidth)
		{
		viewport = SuRender().GetClientRect()
		percent = ((newWidth + 2) / viewport.width).DecimalToPercent()
		if .canResize?(viewport, 100/*=100%*/ - percent, .left, .right)
			{
			leftMin = .left.window.Ctrl.Xmin + 2
			.left.window.GetContainerEl().SetStyle('width', .getWidth(percent, leftMin))
			if .right isnt false and .left.percent + .right.percent is 100/*=100%*/
				{
				rightMin = .right.window.Ctrl.Xmin + 2
				.right.window.GetContainerEl().SetStyle(
					'width', .getWidth(100/*=100%*/ - percent, rightMin))
				.right.percent = 100 - percent
				}
			.left.percent = percent
			}
		}

	resizeRight(newWidth)
		{
		viewport = SuRender().GetClientRect()
		percent = ((newWidth + 2) / viewport.width).DecimalToPercent()
		if .canResize?(viewport, 100/*=100%*/ - percent, .right, .left)
			{
			rightMin = .right.window.Ctrl.Xmin + 2
			.right.window.GetContainerEl().SetStyle('width', .getWidth(percent, rightMin))
			if .left isnt false and .left.percent + .right.percent is 100/*=100%*/
				{
				leftMin = .left.window.Ctrl.Xmin + 2
				.left.window.GetContainerEl().SetStyle(
					'width', .getWidth(100/*=100%*/ - percent, leftMin))
				.left.percent = 100 - percent
				}
			.right.percent = percent
			}
		}

	canResize?(viewport, percent, side, otherSide)
		{
		if otherSide is false
			return true
		if side.percent + otherSide.percent isnt 100/*=100%*/
			return true
		width = viewport.width * percent.PercentToDecimal()
		return width >= side.window.Ctrl.Xmin + 2
		}
	}
