// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	N: 8 // block size (number of lines)
	CallClass(libs)
		{
		if String?(libs)
			libs = [libs]
		return (new this).Detect(libs)
		}
	Detect(libs)
		{
		.init()
		for lib in libs
			.process(lib)
		.print()
		}

	FindDuplicates(libs)
		{
		if String?(libs)
			libs = [libs]

		.init()
		for lib in libs
			.process(lib)
		return .output()
		}

	init()
		{
		.hashes = Object().Set_default(Object())
		}

	process(lib)
		{
		QueryApply(lib $ " where name !~ 'Test$'", group: -1)
			{|x|
			.process1(lib, x)
			}
		}
	process1(lib, x)
		{
		last = -999
		lines = x.lib_current_text.Lines().Map!(#Trim)
		for (i = 0; i < lines.Size() - .N; ++i)
			{
			// don't start block with blank or '}' or ')' line
			if lines[i] is '' or lines[i] is '}' or lines[i] is ')'
				continue
			name = lib $ ':' $ x.name $ ':' $ (i + 1)
			hash = .hash(lines, i)
			v = .getValue(hash)
			if v is false
				.setValue(hash, name)
			else if i - last >= .N
				{
				last = i
				.setValue(hash, v $ ', ' $ name)
				}
			}
		}
	getValue(hash)
		{
		return .hashes[hash[0]].GetDefault(hash, false)
		}
	setValue(hash, value)
		{
		.hashes[hash[0]][hash] = value
		}
	// hash .N non-blank lines
	hash(lines, i)
		{
		hasher = .hasher()
		for (j = i, n = 0; n < .N and j < lines.Size(); ++j)
			if lines[j] isnt '' // skip blank lines
				{
				hasher.Update(lines[j])
				++n
				}
		return hasher.Value()
		}
	hasher()
		{
		return Sha1()
		}
	print()
		{
		// TODO ignore spurious duplicate hashes
		// TODO find longest duplicate range (may be longer than .N lines)
		// TODO merge overlapping ranges
		dups = .collectDups()
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())
		}

	collectDups()
		{
		dups = Object()
		for sub in .hashes.Members()
			dups.Append(.hashes[sub].Values().Filter({ it.Has?(',') }))
		return dups
		}

	output()
		{
		hashes = Object()
		for sub in .hashes.Members()
			for m in .hashes[sub].Members()
				{
				v = .hashes[sub][m]
				if v.Has?(',')
					hashes[m] = v
				}
		return hashes
		}
	}
