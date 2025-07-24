// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
MemoizeSingle
	{
	table: ide_settings
	Columns: #(ide_logfont, ide_scifont, ide_qc_enabled, ide_color_scheme, ide_move_tab
		ide_book_auto_refresh, ide_show_line_numbers, ide_show_whitespace,
		ide_show_annotations, ide_show_fold_margin, ide_scroll_tabs,
		ide_selected_tab_bold, ide_selected_tab_color)
	Init() // Should only be called if in a development environment
		{
		if .ide?()
			.rebuildSettings()
		}

	ide?()
		{
		return Suneido.GetDefault(#Persistent, #(Set: '')).Set is #IDE
		}

	// Rebuild unpackable, unsaveable values. Also recalculates as needed
	settingsMap: #(
		ide_logfont: function(value)
			{
			if value.Member?(#fontPtSize)
				value.lfHeight = StdFonts.LfSize(value.fontPtSize)
			SetGuiFont(value.Copy())
			return value
			}
		ide_scifont: function (value)
			{
			if value.Member?(#fontPtSize)
				value.lfHeight = StdFonts.LfSize(value.fontPtSize)
			return value
			}
		)
	rebuildSettings()
		{
		for setting, fn in .settingsMap
			if '' isnt orig = .Get(setting)
				if orig isnt value = fn(Object?(orig) ? orig.Copy() : orig)
					.setValue(setting, value)
		.ResetCache()
		}

	setValue(setting, value)
		{
		QueryApply1(.table)
			{
			it[setting] = value
			it.Update()
			}
		}

	Func()
		{
		result = []
		try
			if false isnt settings = Query1(.table)
				result = settings
		catch (e)
			if not e.Has?('nonexistent table')
				SuneidoLog('ERROR: (CAUGHT) IDESettings - ' $ e)
		return result
		}

	Ensure()
		{
		Database('ensure ' $ .table $  ' (' $ .Columns.Join(', ') $ ') key()')
		if QueryEmpty?(.table)
			QueryOutput(.table, [])
		removeColumns = QueryColumns(.table).Difference(.Columns)
		if not removeColumns.Empty?()
			Database('alter ' $ .table $ ' drop (' $ removeColumns.Join(', ') $ ')')
		}

	Get(setting, defaultVal = '')
		{
		return .ide?() ? IDESettings().GetDefault(setting, defaultVal) : defaultVal
		}

	Set(setting, value, resetCache = false)
		{
		if not .ide?()
			{
			SuneidoLog('ERROR: IDESettings are not loaded, cannot set: ' $ setting $
				', to: ' $ Display(value))
			return
			}
		if .settingsMap.Member?(setting)
			(.settingsMap[setting])(value)
		.setValue(setting, value)
		if resetCache
			.ResetCache()
		}
	}
