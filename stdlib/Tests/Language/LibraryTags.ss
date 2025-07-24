// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	excludes: ('__protect', '__valid')

	// To remove
	RemoveTagsFromName(name)
		{
		return .RemoveTagFromName(name)
		}

	RemoveTagFromName(name)
		{
		return name.RemoveSuffix(.GetTagFromName(name))
		}

	GetTagFromName(name)
		{
		if false isnt pos = name.FindLast('__')
			{
			if not .excludes.Has?(name[pos..])
				return name[pos..]
			}
		return ''
		}

	AddMode(mode, onlyClient? = false)
		{
		if BuiltDate() < #20250516
			return

		modes = Suneido.GetInit(#LibraryTags_Modes, Object)
		modes.AddUnique(mode)
		.Reset(:onlyClient?)
		}

	Reset(onlyClient? = false)
		{
		if BuiltDate() < #20250516
			return

		if Client?() and not onlyClient?
			ServerEval('LibraryTags.Reset')

		trials = OptContribution(#LibraryTags_Trials, #())
		modes = Suneido.GetDefault(#LibraryTags_Modes, #())
		tags = .buildTags(trials, modes)
		if .ConvertTagInfo(.GetTagsInUse()) isnt tags
			Suneido.LibraryTags(@tags)
		}

	buildTags(trials, modes)
		{
		tags = Object()
		modes = Object('').Append(modes)
		for mode in modes
			{
			tags.Add(mode)
			for trail in trials
				tags.Add(Opt(mode, '_') $ trail)
			}
		return tags.Remove('')
		}

	ConvertTagInfo(tags)
		{
		try
			{
			if tags.Empty?()
				return tags

			Assert(tags[0] is: '')
			return tags[1..].Map({ it.RemovePrefix('__') })
			}
		catch (e)
			{
			Print('ERROR: LibraryTags.ConvertTagInfo - ' $ e)
			return #()
			}
		}

	GetTagsInUse()
		{
		if BuiltDate() < #20250516
			return #('')

		return Suneido.Info('library.tags').SafeEval()
		}

	GetRecord(name, lib)
		{
		// #() means the client is using the server's tags
		if #() is curTags = .GetTagsInUse()
			curTags = ServerEval(LibraryTags.GetTagsInUse)

		for (i = curTags.Size() - 1; i >= 0; i--)
			{
			if false isnt rec = Query1(lib, group: -1, name: name $ curTags[i])
				return rec
			}
		return false
		}
	}