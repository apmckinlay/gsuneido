// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function (view, args)
	{
	m = args[0]
	cols = #(
		GetHeaderSelectPrompt:,
		GetExcludeSelectFields:,
		SetHeaderChanged:,
		GetColumns:,
		SetColumns:,
		InsertColumn:,
		RemoveColumn:,
		SetColWidth:,
		GetColWidth:)

	addons = #(
		DeleteRow:,
		AfterDelete:,
		SaveFirst:,
		RemoveRowByKeyPair:,// NOTE: Do not use this when in an editable list state.
							// This method is designed to remove lines from the list,
							// NOT delete the record from the database.
							// IE: Remove row associated with a record deleted from
							// MultiView > AccessControl
		On_Edit:,
		On_New:, 			// Ctrl + N
		ForceEditMode:,
		SaveOutstandingChanges:,

		GetExpandCtrlAndRecord:,
		CollapseAll:,
		ExpandByField:,
		GetExpandedControl:,
		UpdateTotalSelected:,
		On_Count:
		)

	addonsPrefixes = #('On_VirtualListThumb', 'VirtualListThumb_', 'VirtualListGrid_',
		'On_Context_Global')
	cols.Member?(m)
		? view.GetModel().ColModel[m](@+1args)
		: addons.Member?(m) or addonsPrefixes.Any?({ m.Prefix?(it) })
			? view.Addons.SendToOneAddon(@args)
			: view.GetContextMenu().RedirectContextMenu(view, args)
	}