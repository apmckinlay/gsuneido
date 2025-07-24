// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Make(tablename, keysize = 10, recsize = 100, nrecs = 1000)
		{
		Database("create " $ tablename $ " (key, fld) key(key)")
		.output(tablename, keysize, recsize, nrecs)
		}
	output(tablename, keysize, recsize, nrecs)
		{
		minLimit = 100
		for (i = 0; i < nrecs; )
			Transaction(update:)
				{|t|
				limit = Min(i + minLimit, nrecs)
				t.Query(tablename)
					{|q|
					for (; i < limit; ++i)
						.output1(q, .make_record(i, keysize, recsize))
					}
				}
		}
	output1(q, rec)
		{
		q.Output(rec)
		}
	make_record(i, keysize, recsize)
		{
		return [key: .make_key(i, keysize), fld: .make_field(recsize - keysize)]
		}
	make_key(i, keysize)
		{
		key = Display(i).LeftFill(keysize)
		Assert(key isSize: keysize)
		return key
		}
	make_field(size)
		{
		return " ".Repeat(size)
		}
	}