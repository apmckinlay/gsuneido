// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// display line, column, and modified on status bar
CodeViewAddon
	{
	Name: 		Status
	Inject: 	exterior

	StatusBarControl: Statusbar
	Controls(addonControls)
		{ addonControls.Add(Object(.StatusBarControl), at: 10) }

	Addon_RedirMethods()
		{ return #(Status) }

	UpdateUI()
		{
		pos = .GetCurrentPos()
		line =  .LineFromPosition(pos)
		nlines = .GetLineCount()
		col = .GetColumn(pos)
		endcol = .GetColumn(.GetLineEndPosition(line))
		len = .GetLength()
		status = "\t\tLine " $ (line + 1) $ " of " $ nlines $
			" | Col " $ col $ " of " $ endcol $
			' | Pos ' $ pos $ " of " $ len $ "    "
		.Status(status)
		}

	Invalidate()
		{ .UpdateUI() }

	AfterSave()
		{ .UpdateUI() }

	Status(@status)
		{ .AddonControl.Set(@status) }
	}
