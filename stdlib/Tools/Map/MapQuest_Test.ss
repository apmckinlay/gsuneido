// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(MapQuest.BuildUrl('100 Broadway', '', 'San Diego', 'CA', '', '90210')
			is: "http://mapquest.com/maps/map.adp?address=100 Broadway &" $
				"city=San Diego&state=CA&country=US&zipcode=&cid=lfmaplink")
		Assert(MapQuest.BuildUrl('', '', '', '', '', '90210')
			is: "http://mapquest.com/maps/map.adp?address= &" $
				"city=&state=&country=&zipcode=90210&cid=lfmaplink")
		}
	}