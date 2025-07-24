// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	OkForResetAll?: true
	Func(library)
		{
		suppressList = .suppressList(library)
		return suppressList isnt false
			? suppressList
			: #()
		}

	suppressList(library)
		{
		suppressRef = .suppressRef(library)
		if Function?(suppressRef) or Class?(suppressRef)
			try
				{
				// Must evaluate prior to the return in order for
				// the try block to catch/log: "no return value"
				result = suppressRef()
				return result
				}
			catch (error)
				{
				.log(library, error, 'suppression list')
				return false
				}
		return suppressRef
		}

	log(library, error, event)
		{
		SuneidoLog('ERROR: Unexpected issue evaluating ' $ event $ ': ' $ Display(error),
			params: [:library])
		}

	suppressRef(library)
		{
		try
			return Global(library.Capitalize() $ '_CheckLibrarySuppressions')
		catch (error)
			if not error.Has?(`can't find`)
				.log(library, error, 'reference')
		return false
		}

	ResetRequired?(library, name)
		{
		suppressRef = .suppressRef(library)
		if Class?(suppressRef) and suppressRef.Method?('Suppressed?')
			try
				{
				suppressedState = suppressRef.Suppressed?(name)
				cachedState = this(library).Has?(name)
				return suppressedState isnt cachedState
				}
			catch (error)
				.log(library, error, 'reset required')
		return false
		}

	Libraries()
		{
		return LibraryTables().Filter({ .suppressList(it) isnt false })
		}
	}
