// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		func = IdentifyRecord
		keys = #()
		rec = []
		Assert(func(keys,rec) is: '')
		keys = #('singleKey')
		Assert(func(keys,rec) is: '')

		rec = [singleKey: 'singleValue']
		Assert(func(keys,rec) is: 'singleKey singleValue')

		rec = [singleKey: 1000]
		Assert(func(keys,rec) is: 'singleKey 1000')

		keys = #('singleKey', 'comp,keys')
		rec = [comp: 'hello', keys: 'there']
		Assert(func(keys,rec) is: 'comp hello,keys there')

		rec.singleKey = 'priority'
		Assert(func(keys,rec) is: 'singleKey priority')

		time = Timestamp()
		rec.timestampKey = time
		keys = #('singleKey', 'comp,keys', 'timestampKey')
		Assert(func(keys,rec) is: 'singleKey priority')

		rec.Delete('singleKey')
		Assert(func(keys,rec) is: 'timestampKey ' $ Display(time))

		displayFunc = function (record, keyfield)
			{
			return 'here is the field: ' $ keyfield $ '\r\nhere is the value: ' $
				Display(record[keyfield])
			}
		Assert(func(keys, rec, displayFunc) is: 'here is the field: timestampKey\r\n' $
			'here is the value: ' $ Display(time))
		}
	}