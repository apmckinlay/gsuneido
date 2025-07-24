// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
// WARNING:
// 	This class is designed to reduce an IDE to its minimum system tables.
//	This process cannot easily be undone.
class
	{
	mandatoryTables: #(tables, ide_settings, suneidolog, svc_settings, columns, configlib,
		indexes, persistent, views, suneidoc, imagebook, stdlib)
	New(.excludes = false, dropLibraries = false, dropBooks = false, .quiet = false)
		{
		if .excludes is false
			.excludes = Object()
		.excludes.Add(@.mandatoryTables)
		libraries = LibraryTables()
		if not dropLibraries
			.excludes.Add(@libraries)
		else
			for lib in libraries.Difference(.excludes)
				ServerEval(#Unuse, lib)
		if not dropBooks
			.excludes.Add(@BookTables())
		.excludes.Sort!().Unique!()
		.dropTables(0)
		}

	dropTables(runs)
		{
		rerun? = false
		.print('-----------------------------------------------------------------')
		.print('Run: ' $ runs)
		problemTables = Object()
		QueryApplyMulti(#tables, update:)
			{
			if .excludes.Has?(it.table)
				continue
			try
				Database('drop ' $ it.table)
			catch (e)
				{
				rerun? = true
				problemTables.Add(it.table)
				.print(e)
				}
			}
		if not .stop(runs, problemTables) and rerun?
			.dropTables(++runs)
		}

	print(msg)
		{
		if not .quiet
			Print(msg)
		}

	stop(runs, problemTables)
		{
		if stop = runs > 10 /*= max runs*/
			SuneidoLog('WARNING: Max runs reached, stopping',
				params: [:problemTables, excludes: ])
		return stop
		}
	}
