// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_EnFrancais()
		{
		Assert(0.EnFrancais() is: 'zero')
		Assert(8.EnFrancais() is: 'huit')
		Assert(18.EnFrancais() is: 'dix-huit')
		Assert(20.EnFrancais() is: 'vingt')
		Assert(80.EnFrancais() is: 'quatre-vingts')
		Assert(100.EnFrancais() is: 'cent')
		Assert(200.EnFrancais() is: 'deux cents')
		Assert(201.EnFrancais() is: 'deux cent un')
		Assert(750.EnFrancais() is: 'sept cent cinquante')
		Assert(1000.EnFrancais() is: 'mille')
		Assert(1525.EnFrancais() is: 'mille cinq cent vingt cinq')
		Assert(9999.EnFrancais() is: 'neuf mille neuf cent quatre-vingt dix-neuf')
		Assert(10000.EnFrancais() is: 'dix mille')
		Assert(80750.EnFrancais() is: 'quatre-vingt mille sept cent cinquante')
		Assert(100000.EnFrancais() is: 'cent mille')
		Assert(785694.EnFrancais()
			is: 'sept cent quatre-vingt cinq mille six cent quatre-vingt quatorze')
		Assert(1000000.EnFrancais() is: 'un million')
		Assert(1555555.EnFrancais()
			is: 'un million cinq cent cinquante cinq mille cinq cent cinquante cinq')
		Assert(2000000.EnFrancais() is: 'deux millions')
		Assert(1000000000.EnFrancais() is: 'un milliard')
		Assert(2147483647.EnFrancais()
			is: 'deux milliards cent quarante sept millions quatre cent quatre-vingt ' $
				'trois mille six cent quarante sept')
		}
	}