// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// REFERENCE: http://www.marijn.org/everything-is-4/counting-0-to-100/italian/
	// REFERENCE: http://www.lexisrex.com/Italian-Numbers/1021
	tests: (
		0:	zero
		1:	uno
		2:	due
		3:	tre
		4:	quattro
		5:	cinque
		6:	sei
		7:	sette
		8:	otto
		9:	nove
		10:	dieci
		11:	undici
		12:	dodici
		13:	tredici
		14:	quattordici
		15:	quindici
		16:	sedici
		17:	diciassette
		18:	diciotto
		19:	diciannove
		20:	venti
		21:	ventuno
		22:	ventidue
		23:	ventitre
		24:	ventiquattro
		25:	venticinque
		26:	ventisei
		27:	ventisette
//		28:	ventotto
		29:	ventinove
		60:	sessanta
		61:	sessantuno
		62:	sessantadue
		63:	sessantatre
		64:	sessantaquattro
		65:	sessantacinque
		66:	sessantasei
		67:	sessantasette
		68:	sessantotto
		69:	sessantanove
		90:	novanta
		91:	novantuno
		92:	novantadue
		93:	novantatre
		94:	novantaquattro
		95:	novantacinque
		96:	novantasei
		97:	novantasette
		98:	novantotto
		99:	novantanove
		100: cento
		1021: milleventuno
		10245: diecimiladuecentoquarantacinque
		102455: centoduemilaquattrocentocinquantacinque
		)
	Test_ToWordsItalian()
		{
		for n in .tests.Members()
			Assert(n.ToWordsItalian() is: .tests[n])
		}
	}