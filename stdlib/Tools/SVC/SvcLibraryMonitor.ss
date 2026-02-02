// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	New()
		{
		.className = Name(.Base())
		.subs = [
			PubSub.Subscribe('LibraryTreeChange', .libraryTreeChange)
			PubSub.Subscribe('LibraryRecordChange', .libraryRecordChange)
			]
		}

	libraryTreeChange(args)
		{
		.libraryChange(args, resetCache?:)
		}

	libraryChange(args, .resetCache? = false, .resetClass? = false)
		{
		try
			{
			.runningTests? = TestRunner.RunningTests?()
			args.Each(.processChange)
			if .resetClass?
				.Reset()
			}
		catch (e)
			SuneidoLog('ERROR: (CAUGHT) ' $ Name(this) $ ': ' $ e, caughtMsg: 'IDE error')
		}

	processChange(change)
		{
		if '' isnt name = change.GetDefault('name', '')
			{
			LibUnload(name)
			if name is .className
				.resetClass? = true
			}
		if .resetRequired?(change.GetDefault('table', ''), name)
			.resetCache? = true
		}

	resetRequired?(table, name)
		{
		return not .runningTests? and table isnt '' and name isnt ''
			? .resetCache?
			: false
		}

	libraryRecordChange(args)
		{
		.libraryChange(args)
		}

	Reset()
		{
		.subs.Each(#Unsubscribe)
		super.Reset()
		SvcLibraryMonitor()
		}
	}
