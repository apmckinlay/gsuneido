// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ToWordSpanish()
		{
		Assert(0.ToWordSpanish() is: 'cero')
		Assert(8.ToWordSpanish() is: 'ocho')
		Assert(18.ToWordSpanish() is: 'dieciocho')
		Assert(20.ToWordSpanish() is: 'veinte')
		Assert(100.ToWordSpanish() is: 'cien')
		Assert(650.ToWordSpanish() is: 'seiscientos cincuenta')
		Assert(750.ToWordSpanish() is: 'setecientos cincuenta')
		Assert(1000.ToWordSpanish() is: 'mil')
		Assert(1525.ToWordSpanish() is: 'mil quinientos veinticinco')
		Assert(9999.ToWordSpanish() is: 'nueve mil novecientos noventa y nueve')
		Assert(10000.ToWordSpanish() is: 'diez mil')
		Assert(80750.ToWordSpanish() is: 'ochenta mil setecientos cincuenta')
		Assert(100000.ToWordSpanish() is: 'cien mil')
		Assert(785694.ToWordSpanish()
			is: 'setecientos ochenta y cinco mil seiscientos noventa y cuatro')
		Assert(1000000.ToWordSpanish() is: 'un millon')
		Assert(1555555.ToWordSpanish()
			is: 'un millon quinientos cincuenta y cinco mil quinientos cincuenta y cinco')
		Assert(7856094.ToWordSpanish()
			is: 'siete millones ochocientos cincuenta y seis mil noventa y cuatro')
		Assert(35065781.ToWordSpanish()
			is: 'treinta y cinco millones sesenta y cinco mil setecientos ochenta y uno')
		Assert(41391129.ToWordSpanish()
			is: 'cuarenta y un millones trescientos noventa y un mil ciento veintinueve')
		Assert(171568094.ToWordSpanish()
			is: 'ciento setenta y un millones quinientos sesenta y ocho mil ' $
				'noventa y cuatro')
		Assert(2534515551.ToWordSpanish()
			is: 'dos mil quinientos treinta y cuatro millones quinientos ' $
				'quince mil quinientos cincuenta y uno')
		Assert(25345155551.ToWordSpanish()
			is: 'veinticinco mil trescientos cuarenta y cinco millones ciento ' $
				'cincuenta y cinco mil quinientos cincuenta y uno')
		Assert(225345155551.ToWordSpanish()
			is: 'doscientos veinticinco mil trescientos cuarenta y cinco millones ' $
				'ciento cincuenta y cinco mil quinientos cincuenta y uno')
		Assert(1000000000000.ToWordSpanish() is: 'Esa cifra es demasiado alta')
		}
	}