// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(CountryFromStateProv('') is: '')
		Assert(CountryFromStateProv('~') is: '')
		for p in StateCodes
			Assert(CountryFromStateProv(p) is: 'US')
		for p in ProvinceCodes
			Assert(CountryFromStateProv(p) is: 'CA')
		Assert(CountryFromStateProv('PQ') is: 'CA')
		Assert(CountryFromStateProv('NU') is: 'CA')
		Assert(CountryFromStateProv('NF') is: 'CA')
		}
	}