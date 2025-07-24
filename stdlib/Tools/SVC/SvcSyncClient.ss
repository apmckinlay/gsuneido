// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
/*
SvcSyncClient is used to compare the local library with the master library.

Checksum goes: name, path, text (timmed with Trim()), lib committed date (string)
*/
class
	{
	New(.table, .master) // table is SvcLibrary or SvcBook
		{
		.calc_checksums()
		}

	calc_checksums()
		{
		.names = Object()
		.checksums = Object()
		if .table.Type is 'lib'
			queryType = ' sort name'
		else
			queryType = ' where not path.Has?("commitRec")
				extend srt = path $ "/" $ name sort srt'

		QueryApply(.table.Query() $ queryType)
			{|x|
			x.path = .table.GetPath(x)
			x.name = .table.Type is 'lib' ? x.name : x.path $ '/' $ x.name
			.names.Add(x.name)
			.checksums.Add(SvcCksum(x))
			}
		}

	Check()
		{
		.check('', '~', results = Object(), 0)
		return results
		}

	check(from, to, results, depth)
		{
		Assert(depth < 20) /*= max depth */
		for r in .master.GetChecksums(.table.Table(), from, to)
			{
			if not .check_range(r)
				if r.n is 1 // single item
					.results(r, results)
				else // range
					.check(r.from, r.to, results, depth + 1) // recurse
			}
		}

	check_range(range)
		{
		range.here = Object()
		cksum = Adler32()
		for (i = .names.BinarySearch(range.from);
			i < .names.Size() and .names[i] < range.to; ++i)
			{
			range.here.Add(i)
			cksum.Update(.checksums[i].Hex())
			}
		return range.cksum is (cksum.Value() & 0xffffffff) /*=
			need mask to be compatible with BuiltDate > 2025-05-09 */
		}

	results(r, results)
		{
		match = false
		for i in r.here
			results.Add(r.from is .names[i]
				? '# ' $ (match = r.from)
				: 'L ' $ .names[i])
		if match is false
			results.Add('M ' $ r.from)
		}
	}
