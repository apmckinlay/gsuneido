// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
/*
Transforms a number into words in Dutch.
Both the integer and the fractional part are transformed,
seperated by the word 'komma'. If there is no fractional part,
or if it is zero, the fractional part is not mentioned.
If there is no integer part, the integer part is transformed as zero.
The method also accepts negative numbers,
returning 'min ' as the start of the string.
It can handle numbers from 0.000000000001 upto 9,999,999,999,999,999.9999
(max 16 integer digits (= 9 trillion ...), and max 12 fractional digits).
You cannot have a number with all 16 integer digits and all 12 fractional digits.
A fractional number seems to be limited to 16 positions (which is confirmed in
sunumber.h), that is the integer digits, the decimal sign,
and the fractional digits alltogether.
These are in fact limits (I found) for Suneido.
*/
class
	{
	words_dutch_long: #(0: 'nul', 1: 'een', 2: 'twee', 3: 'drie', 4: 'vier',
		5: 'vijf', 6: 'zes', 7: 'zeven', 8: 'acht', 9: 'negen',
		10: 'tien', 11: 'elf', 12: 'twaalf', 13: 'dertien',
		14: 'veertien', 15: 'vijftien', 16: 'zestien', 17: 'zeventien',
		18: 'achttien', 19: 'negentien', 20: 'twintig', 30: 'dertig',
		40: 'veertig', 50: 'vijftig', 60: 'zestig',
		70: 'zeventig', 80: 'tachtig', 90: 'negentig')
	sectionSize: 3
	NumberToWords(number)
		{
		if (number is 0)	// always a bit special ...
			return .words_dutch_long[number]

		_res = Object(sInt: "", sFrac: "")
		theNumber = number.Abs()
		for (nLoopcount = 1; nLoopcount <= 2; ++nLoopcount)
			{
			sNum = .getNumStr(nLoopcount, theNumber)
			nRank = 1
			while sNum isnt ""
				{
				if sNum.Size() > .sectionSize
					{
					sNumRank = sNum[-.sectionSize :: .sectionSize]
					sNum = sNum[.. sNum.Size() - .sectionSize]
					}
				else
					{
					sNumRank = sNum
					sNum = ""
					}
				.buildThreeDigits(sNumRank, nRank, nLoopcount)
				++nRank
				}
			}
		return (number < 0 ? "min " : "") $ _res.sInt $
			(_res.sFrac.Size() > 0 ? " komma " $ _res.sFrac : "")
		}
	getNumStr(nLoopcount, theNumber)
		{
		sNum = ""
		if (nLoopcount is 1)
			// It seems that we can handle integers larger than
			// 2147483647 because we immediately convert to a string.
			// Discovered this by accident. Nice side effect for this function !
			sNum = "" $ theNumber.Int()
		else if (theNumber.FracDigits() > 0)
			{
			if (theNumber.FracDigits() <= 4) /*= Pitty the Frac() method
				only returns the first 4 digits */
				sNum = ("" $ theNumber.Frac())[1 ..]
			else
				{
				// Seems a maximum of 12 digits in the fractional part are preserved.
				// Could be a restriction of the conversion to string.
				sNum = "" $ theNumber
				sNum = sNum[sNum.Find(".") + 1 ..]
				}
			}
		return sNum
		}
	ten: 10
	hundred: 100
	buildThreeDigits(sNumRank, nRank, nLoopcount)
		{
		nNum = Number(sNumRank)
		nHundreds = (nNum / .hundred).Int()
		nTens = ((nNum - (nHundreds * .hundred)) / .ten).Int()
		nUnits = (nNum - (nHundreds * .hundred) - (nTens * .ten))
		nLessThanHundred = nNum - (nHundreds * .hundred)
		sRes = (nHundreds > 0 ?
			(nHundreds > 1 ? .words_dutch_long[nHundreds] : "") $
				((nRank is 1 and nLessThanHundred > 0) ? "honderd " : "honderd")
			: "")
		sRes $= .buildTwoDigits(nLessThanHundred, nUnits, nTens)
		if (sRes isnt "nul" and sRes isnt "")
			{
			sRes = .handle2To6(nRank, sRes)
			if (nLoopcount is 1)
				_res.sInt = sRes $ .space(_res.sInt) $ _res.sInt
			else
				_res.sFrac = sRes $ .space(_res.sFrac) $ _res.sFrac
			}
		}
	buildTwoDigits(nLessThanHundred, nUnits, nTens)
		{
		s = ''
		if (nLessThanHundred > 0)
			{
			if nUnits is 0
				s = .words_dutch_long[nTens * .ten]
			else if (nLessThanHundred <= 20) /*= less than 20, read from map*/
				s = .words_dutch_long[nLessThanHundred]
			else
				s = .words_dutch_long[nUnits] $
					((nUnits is 2 or nUnits is 3) ? "ën" : "en") $ /*= separator */
					.words_dutch_long[nTens * .ten]
			}
		return s
		}
	handle2To6(nRank, sRes)
		{
		switch nRank
			{
		case 2 :
			sRes = (sRes is 'een' ? 'duizend' : sRes $ 'duizend')
		case 3 : /*= 3 */
			sRes $= 'miljoen'
		case 4 : /*= 4 */
			sRes $= 'miljard'
		case 5 : /*= 5 */
			sRes $= 'biljoen'
		case 6 : /*= 6 */
			sRes $= 'triljoen'
		default:
			}
		return sRes
		}
	space(val)
		{
		return val.Size() > 0 ? " " : ""
		}
	}
