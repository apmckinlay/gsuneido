// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(prev, cur, query_columns, hwnd = 0, quiet? = false)
		{
		// check for modified timestamp field first (ends in "_TS")
		if false isnt tsField = query_columns.FindOne({ it.Suffix?("_TS") })
			return .TSCheck(prev, cur, tsField, hwnd, :quiet?)

		changes = .getChanges(cur, prev, query_columns)
		return changes isnt ''
			? .recordChanged(hwnd, [:changes], changes $ '\n', :quiet?)
			: false
		}

	getChanges(cur, prev, query_columns)
		{
		// otherwise we check all query columns to see if they changed

		// WARNING: extended columns will not have a value on new records
		// so if the record in memory is saved more than once, the record read
		// from the table the second time will have a value for the extended
		// column but the record in memory will not.
		changes = ""
		for field in query_columns
			if cur[field] isnt prev[field]
				changes $= SelectPrompt(field) $
					' changed from ' $ Display(prev[field]) $
					' to ' $ Display(cur[field]) $ '\n'
		return changes
		}

	TSCheck(prev, cur, ts_col, hwnd = 0, quiet? = false)
		{
		if prev[ts_col] isnt cur[ts_col]
			{
			params = Object(ts_field: ts_col,
				prev_val: prev[ts_col], cur_val: cur[ts_col])
			if cur.Member?('bizuser_user_modified')
				params.bizuser_user_modified = cur.bizuser_user_modified
			if prev.Member?('bizuser_user_modified')
				params.prev_bizuser_user_modified = prev.bizuser_user_modified
			.recordChanged(hwnd, params, :quiet?)
			return true
			}
		return false
		}

	recordChanged(hwnd, params, msg = '', quiet? = false)
		{
		if not quiet?
			AlertDelayed('Another user has modified this record.\n\n' $
				msg $ 'Please use Current > Restore and re-do your changes.',
				'Overwrite Warning', hwnd)
		SuneidoLog("INFO: Update Record Conflict - Restore Changes", :params)
		return true
		}
	}
