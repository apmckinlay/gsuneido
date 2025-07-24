// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_states()
		{
		Assert(States isSize: StateCodes.Size())
		for i in .. States.Size()
			Assert(States[i].BeforeFirst(' - ') is: StateCodes[i])
		Assert(States isSize: StateNames.Size())
		for i in .. States.Size()
			Assert(States[i].AfterFirst(' - ') is: StateNames[i])
		}
	Test_provinces()
		{
		Assert(Provinces isSize: ProvinceCodes.Size())
		for i in .. Provinces.Size()
			Assert(Provinces[i].BeforeFirst(' - ') is: ProvinceCodes[i])
		Assert(Provinces isSize: ProvinceNames.Size())
		for i in .. Provinces.Size()
			Assert(Provinces[i].AfterFirst(' - ') is: ProvinceNames[i])
		}
	Test_combined()
		{
		Assert(StatesProvs is: [States, Provinces].Flatten().Sort!())
		Assert(StatesProvsMex
			is: [States, Provinces, MexicanStates].Flatten().Sort!())
		Assert(StateProvCodes is: [StateCodes, ProvinceCodes].Flatten().Sort!())
		Assert(StateProvMexCodes
			is: [StateCodes, ProvinceCodes, MexicanStateCodes].Flatten().Sort!())
		}
	}