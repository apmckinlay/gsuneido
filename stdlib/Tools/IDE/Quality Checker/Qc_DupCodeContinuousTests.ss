// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	excludeLibs: #('demobookoptions', 'Kenlib')
	New()
		{
		.libs = GetContributions('ApplicationLibraries').Difference(.excludeLibs)
		.libs.RemoveIf({ not QueryColumns(it).Has?('group') })
		dupFinder = new DupCodeFinder
		.oldHash = dupFinder.FindDuplicates(.libs)
		}

	FindNewDuplicates()
		{
		newDups = Object()
		dupFinder = new DupCodeFinder
		newHash = dupFinder.FindDuplicates(.libs)
		.checkForDups(.oldHash, newHash, newDups)
		return newDups
		}

	checkForDups(oldHash, newHash, newDups)
		{
		oldRecs = .dupCount(oldHash)
		newRecs = .dupCount(newHash)
		for dup in newHash.Members()
			if oldHash.Member?(dup)
				{
				// need less than so we do not flag as 'new' if dup is removed
				if oldHash[dup].Split(',').Size() < newHash[dup].Split(',').Size()
					newDups.Add(newHash[dup])
				}
			else
				{
				// don't flag as a new dup if the number of dups per record did not change
				// i.e. line numbers and hash changed, but it just matched on a different
				// section of the same (original) duplication
				if .newDupCount(newHash, oldRecs, newRecs, dup)
					newDups.Add(newHash[dup])
				}
		}

	dupCount(hash)
		{
		recs = Object().Set_default(0)
		for hash in hash
			for rec in hash.Split(', ')
				recs[rec.BeforeLast(':')]++
		return recs
		}

	newDupCount(newHash, oldRecs, newRecs, dup)
		{
		for rec in newHash[dup].Split(', ')
			if oldRecs[rec.BeforeLast(':')] < newRecs[rec.BeforeLast(':')]
				return true
		return false
		}
	}
