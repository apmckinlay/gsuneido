// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	excludes: ('__protect', '__protect_Test', '__valid', '__valid_Test')
	RemoveTagFromName(name)
		{
		return .split(name).name
		}

	GetTagFromName(name)
		{
		return .split(name).tag
		}

	split(name)
		{
		if false isnt pos = name.FindLast('__')
			{
			if not .excludes.Has?(name[pos..])
				return Object(name: name[..pos], tag: name[pos..])
			}
		return Object(:name, tag: '')
		}

	// to remove trial tags, set trial to ''
	SetTrialTag(name, trial, trialTags)
		{
		split = .split(name)
		tags = split.tag.RemovePrefix('__').Split('_').Remove(@trialTags)
		if trial isnt ''
			tags.Add(trial)
		return split.name $ Opt('__', tags.Join('_'))
		}

	AddMode(mode, onlyClient? = false)
		{
		modes = Suneido.GetInit(#LibraryTags_Modes, Object)
		modes.AddUnique(mode)
		.Reset(:onlyClient?)
		}

	Reset(onlyClient? = false)
		{
		if Client?() and not onlyClient?
			ServerEval('LibraryTags.Reset')

		trials = LastContribution(#LibraryTags_Trials)
		modes = Suneido.GetDefault(#LibraryTags_Modes, #())
		tags = .BuildTags(trials, modes)
		if .ConvertTagInfo(.GetTagsInUse()) isnt tags
			Suneido.LibraryTags(@tags)
		}

	BuildTags(trials, modes)
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
		return Suneido.Info('library.tags').SafeEval()
		}

	GetRecord(name, lib, tags = false, exclude = false)
		{
		if tags is false
			{
			// #() means the client is using the server's tags
			if #() is tags = .GetTagsInUse()
				tags = ServerEval(LibraryTags.GetTagsInUse)
			}

		for (i = tags.Size() - 1; i >= 0; i--)
			{
			if exclude isnt false and tags[i].Has?(exclude)
				continue
			if false isnt rec = Query1(lib, group: -1, name: name $ tags[i])
				return rec
			}
		return false
		}
	}