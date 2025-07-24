// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// TODO: handle unquoted strings in constants
function (lib, table = "translatelanguage", quiet = false)
	{
	n = 0
	oldEntries = Object()
	QueryApply(table)
		{ |x|
		oldEntries[x.trlang_from] = true
		}
	QueryApply(lib $ ' where name !~ "Test$"')
		{ |x|
		scan = Scanner(x.text)
		for s in scan
			{
			if scan.Type() isnt #STRING or s.Size() <= 4 /*= min required string size */
				continue
			s = s[1 .. -1]
			if s =~ '[^.&a-zA-Z ]'
				continue
			if s =~ "(?i)^(www.|create|update|alter|delete|insert|destroy)"
				continue

			if s.Suffix?('...')
				s = s[.. -3] /*= to trim '...' */
			s = s.Tr('&')
			s = s.Trim()

			if not oldEntries.Member?(s)
				{
				if quiet isnt true
					Print(s)
				++n
				oldEntries[s] = true
				}
			}
		}
	return n
	}