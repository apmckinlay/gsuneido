// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// used by SvcCore, runs on version control database
class
	{
	CallClass(master, from = '', to = '~')
		{
		cksums = SvcSyncCksums(master)
		return .getChecksums(cksums.names, cksums.checksums, from, to)
		}

	NSPLIT: 20
	getChecksums(names, checksums, from, to)
		{
		ranges = Object()
		org = names.BinarySearch(from)
		end = names.BinarySearch(to)
		i = org
		for n in Prorate_Equally(.NSPLIT, end - org, round: 0)
			{
			if n is 0
				continue
			start = i is 0 ? '' : names[i] // inclusive
			cksum = Adler32()
			for (j = 0; j < n; ++j, ++i)
				cksum.Update(checksums[i].Hex())
			to = names.GetDefault(i, '~') // exclusive
			ranges.Add(Object(from: start, :to, :n, cksum: cksum.Value() & 0xffffffff /*=
				need mask to be compatible with BuiltDate > 2025-05-09 */))
			}
		return ranges
		}
	ResetCache()
		{
		SvcSyncCksums.ResetCache()
		}
	}
