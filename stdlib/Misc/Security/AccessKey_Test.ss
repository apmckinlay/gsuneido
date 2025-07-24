// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		for .. 100
			{
			nusers = Random(120) + 1
			expiry = Date().Plus(days: Random(100)).NoTime()
			key = AccessKey.MakeKey(nusers, expiry, mac: '123456')
			ob = AccessKey.SplitKey(key, mac: '123456')
			Assert(ob[0] is: nusers)
			Assert(ob[1] is: expiry)
			}
		bad = { |s| Assert(AccessKey.SplitKey(s, mac: '123456', quiet?:)[0] is: 0) }
		bad('')
		bad('abc')
		bad('123')
		bad('123 456')
		bad('<123')
		bad('77777-88888')
		}

	Test_getMacAddresses()
		{
		Assert(AccessKey.AccessKey_getMacAddresses("127001") is: Object("127001"))
		Assert(AccessKey.AccessKey_getMacAddresses(Object('111111', '222222', '123456'))
			is: Object('111111', '222222', '123456'))
		Assert(AccessKey.AccessKey_getMacAddresses(123456) is: Object(123456))
		Assert(AccessKey.AccessKey_getMacAddresses(false) isnt: Object(false) )
		}

	Test_multipleMacs()
		{
		macs = Object('112233445566'.FromHex(), 'aabbccddeeff'.FromHex(),
			'11aa22bb33cc'.FromHex())
		// valid key
		ob = AccessKey.SplitKey('20130930135jpmLm4vLH82ulPhnViTTBjQANwI',
			mac: macs)
		Assert(ob[0] is: 135)
		Assert(ob[1] is: #20130930)
		Assert(ob.mac.ToHex() is: 'aabbccddeeff')

		// invalid key
		ob = AccessKey.SplitKey('2013093035sz1+iosc19h0GqXxUN0JY+m6cY8',
			mac: macs)
		Assert(ob[0] is: 0)
		Assert(ob[1] is: #17000101)
		Assert(ob.mac is: false)
		}
	}
