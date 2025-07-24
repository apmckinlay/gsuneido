// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
ListControl
	{
	New(@args)
		{
		super(@.enableResetColumns(args))
		// stretch the last column if stretchColumn is not specified
		.stretchColumn = args.GetDefault('stretchColumn', false)
		}

	enableResetColumns(args)
		{
		args.resetColumns = true
		return args
		}

	Header_AllowTrack(col)
		{
		if col is .GetColumns().Size() - 1
			return false
		return super.Header_AllowTrack(col)
		}
	HeaderTrack(col, width)
		{
		super.HeaderTrack(col, width)
		.Adjust(.GetColNum(col) >= .getStretchCol() ? .GetNumCols() - 1 : false)
		}
	Adjust(stretchCol = false)
		{
		.recalcCols(stretchCol)
		.Repaint()
		}
	recalcCols(stretchCol = false)
		{
		GetClientRect(.Hwnd, rect = Object())
		w = rect.right - rect.left
		otherwid = 0
		if stretchCol is false
			stretchCol = .getStretchCol()
		for (i = 0; i < .GetNumCols(); ++i)
			if i isnt stretchCol
				otherwid += .GetColWidth(i)
		// don't let stretched colums get negative width
		stretch = w - otherwid - 1
		.SetColWidth(stretchCol, stretch < 0 ? 0 : stretch)
		}
	getStretchCol()
		{
		if .stretchColumn is false or false is col = .GetColumns().Find(.stretchColumn)
			return .GetNumCols() - 1
		return col
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.Adjust()
		}
	Repaint()
		{
		if .Member?('ListStretchControl_stretchColumn') // constructed
			.recalcCols() // to handle vert scroll bars added/removed
		super.Repaint()
		}
	On_Context_Reset_Columns()
		{
		super.On_Context_Reset_Columns()
		.Adjust()
		}
	UpdateHorzScroll(sif)
		{
		sif.nPage = sif.nMax + 1
		super.UpdateHorzScroll(sif)
		}
	}
