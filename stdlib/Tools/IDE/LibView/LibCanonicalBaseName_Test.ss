// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(LibCanonicalBaseName('Fred') is: 'fred')
		Assert(LibCanonicalBaseName('Field_fred') is: 'fred')
		Assert(LibCanonicalBaseName('Rule_fred') is: 'fred')
		Assert(LibCanonicalBaseName('Trigger_fred') is: 'fred')
		Assert(LibCanonicalBaseName('FredControl') is: 'fred')
		Assert(LibCanonicalBaseName('FredFormat') is: 'fred')
		}
	}