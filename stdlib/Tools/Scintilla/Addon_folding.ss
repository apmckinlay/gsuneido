// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonIDE
	{
	Setting: ide_show_fold_margin
	foldMargin: 2
	Init()
		{
		.SetMarginTypeN(.foldMargin, SC.MARGIN_SYMBOL)
		.setFolding(.Set)
		.SetMarginMaskN(.foldMargin, SC.MASK_FOLDERS)
		.SetMarginSensitiveN(.foldMargin, true)
		.SetFoldFlags(SC.FOLDFLAG_LINEAFTER_CONTRACTED)
		.SetFoldMarginColour(true, .GetSchemeColor('foldMargin'))
		.SetFoldMarginHiColour(true, .GetSchemeColor('foldMargin'))
		.defineMarker(SC.MARKNUM_FOLDEROPEN, SC.MARK_ARROWDOWN)
		.defineMarker(SC.MARKNUM_FOLDER, SC.MARK_ARROW)
		.defineMarker(SC.MARKNUM_FOLDERSUB, SC.MARK_EMPTY)
		.defineMarker(SC.MARKNUM_FOLDERTAIL, SC.MARK_EMPTY)
		.defineMarker(SC.MARKNUM_FOLDEREND, SC.MARK_EMPTY)
		.defineMarker(SC.MARKNUM_FOLDEROPENMID, SC.MARK_EMPTY)
		.defineMarker(SC.MARKNUM_FOLDERMIDTAIL, SC.MARK_EMPTY)

		.SetAutomaticFold(SC.AUTOMATICFOLD_CLICK)
		}

	defineMarker(num, type)
		{ .DefineMarker(num, type, CLR.NONE, CLR.DARKGRAY) }

	ContextMenu()
		{ #("Fold &All", "&Unfold All", "Show/Hide Fold &Margin") }

	On_ShowHide_Fold_Margin()
		{ .setFolding(not .Set) }

	setFolding(.Set)
		{ .SetMarginWidthN(.foldMargin, .Set ? ScaleWithDpiFactor(12) /*= width*/ : 0) }

	// Using toggleRange in order to ensure ALL children lines are collapsed in the record
	On_Fold_All()
		{
		for (line = 0; line < .GetLineCount(); ++line)
			if 0 isnt (.GetFoldLevel(line) & SC.FOLDLEVELHEADERFLAG)
				if SC.FOLDACTION_CONTRACT isnt .GetFoldExpanded(line)
					.ToggleFold(line)
		}

	On_Unfold_All()
		{ .FoldAll(SC.FOLDACTION_EXPAND) }
	}
