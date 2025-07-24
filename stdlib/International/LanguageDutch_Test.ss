// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ToWordsDutch()
		{
		data_ToWordsDutch = #(
			(0, "nul") (1, "een") (2, "twee") (3, "drie") (4, "vier") (5, "vijf")
			(6, "zes") (7, "zeven") (8, "acht") (9, "negen") (10, "tien") (11, "elf")
			(12, "twaalf") (13, "dertien") (14, "veertien") (15, "vijftien")
			(16, "zestien") (17, "zeventien") (18, "achttien") (19, "negentien")
			(20, "twintig") (30, "dertig") (40, "veertig") (50, "vijftig") (60, "zestig")
			(70, "zeventig") (80, "tachtig") (90, "negentig") (100, "honderd")
			(1000, "duizend") (10000, "tienduizend") (100000, "honderdduizend")
			(1000000, "eenmiljoen") (1000000000, "eenmiljard")
			(1000000000000, "eenbiljoen") (1000000000000000, "eentriljoen")
			(122, "honderd tweeëntwintig") (223, "tweehonderd drieëntwintig")
			(123456789012345, "honderddrieëntwintigbiljoen " $
				"vierhonderdzesenvijftigmiljard zevenhonderdnegenentachtigmiljoen " $
				"twaalfduizend driehonderd vijfenveertig")
			(999999999999999, "negenhonderdnegenennegentigbiljoen " $
				"negenhonderdnegenennegentigmiljard negenhonderdnegenennegentigmiljoen " $
				"negenhonderdnegenennegentigduizend negenhonderd negenennegentig")
			(1.999999999999, "een komma negenhonderdnegenennegentigmiljard " $
				"negenhonderdnegenennegentigmiljoen negenhonderdnegenennegentigduizend " $
				"negenhonderd negenennegentig")
			(1.123456789012, "een komma honderddrieëntwintigmiljard " $
				"vierhonderdzesenvijftigmiljoen zevenhonderdnegenentachtigduizend twaalf")
			(99999999999.9999, "negenennegentigmiljard " $
				"negenhonderdnegenennegentigmiljoen negenhonderdnegenennegentigduizend " $
				"negenhonderd negenennegentig komma negenduizend negenhonderd " $
				"negenennegentig")
			(999999999999.999, "negenhonderdnegenennegentigmiljard " $
				"negenhonderdnegenennegentigmiljoen negenhonderdnegenennegentigduizend " $
				"negenhonderd negenennegentig komma negenhonderd negenennegentig")
			(123456789012.123, "honderddrieëntwintigmiljard " $
				"vierhonderdzesenvijftigmiljoen zevenhonderdnegenentachtigduizend " $
				"twaalf komma honderd drieëntwintig")
			(9999.99999999999, "negenduizend negenhonderd negenennegentig komma " $
				"negenennegentigmiljard negenhonderdnegenennegentigmiljoen " $
				"negenhonderdnegenennegentigduizend negenhonderd negenennegentig")
			(1234.12345678901, "duizend tweehonderd vierendertig komma twaalfmiljard " $
				"driehonderdvijfenveertigmiljoen zeshonderdachtenzeventigduizend " $
				"negenhonderd een")
			)
		for x in data_ToWordsDutch
			{
			// positive number
			Assert((x[0]).ToWordsDutch() is: x[1])
			// negative number
			Assert((-1 * x[0]).ToWordsDutch() is: (x[0] is 0 ? "" : "min ") $ x[1])
			}
		}
	}