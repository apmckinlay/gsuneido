// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		QueryApply('indexes')
			{|x|
			if x.columns > ""
				.oneIndex(x.table, x.columns)
			}
		}
	oneIndex(table, fields)
		{
		flds = fields.Split(',')
		prev = ""
		n = 0
		totalKeyLen = 0
		totalPreKeyLen = 0
		QueryApply(table $ ' sort ' $ fields)
			{|x|
			n++
			key = flds.Map({ Pack(x[it]) }).Join('0x00')
			totalKeyLen += key.Size()
			totalPreKeyLen += key.Size() - .commonPrefixLen(prev, key)
			prev = key
			}
		if totalPreKeyLen / n > 16
			Print(table, fields, (totalKeyLen / n).Ceiling(), (totalPreKeyLen / n).Ceiling())
		}
	commonPrefixLen(s, t)
		{
		n = Min(s.Size(), t.Size())
		for i in ..n
			{
			if s[i] isnt t[i]
				return i
			}
		return n
		}
	}