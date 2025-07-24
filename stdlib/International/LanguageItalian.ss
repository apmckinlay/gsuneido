// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// FIXME: 28 should be ventotto
// FIXME: 10245554 should be
//	dieci milioni duecentoquarantacinquemilacinquecentocinquantaquattro
class
	{
	three: 3
	NumberToWords(numero)
		{
		numero = numero.Int().Abs()
		if numero > 999999999999999 /*= max number*/
			return 'numero troppo grande'

		numlett= #(#('unbilione','bilioni'), #('unmiliardo','miliardi'),
			#('unmilione','milioni'), #('mille','mila'))
		numstr= numero.Pad(15) /*= max 15 cifre */
		s = ''
		n = -15

		for (i = 0; i < 4; ++i) /*= bilioni,miliardi,milioni,migliaia */
			{
			threeDigits = numstr[n :: .three]
			if (threeDigits isnt '000')
				s $= (threeDigits is '001')
					? numlett[i][0]
					: .centinaia(threeDigits) $ numlett[i][1]
			n += .three
			}
		//centinaia
		lastThreeDigits = numstr[-.three ..]
		s $= (lastThreeDigits is '000')
			? (s is '' ) ? 'zero' : ''
			: .centinaia(lastThreeDigits)
		return s
		}

	centinaia(numstr)
		{
		numlett = #('zero', 'uno', 'due', 'tre', 'quattro','cinque', 'sei', 'sette',
				'otto', 'nove','dieci', 'undici', 'dodici', 'tredici','quattordici',
				'quindici', 'sedici', 'diciassette','diciotto', 'diciannove', 'venti',
				30: 'trenta', 40: 'quaranta', 50: 'cinquanta', 60: 'sessanta',
				70: 'settanta', 80: 'ottanta', 90: 'novanta')
		//centinaia
		s = ''
		if (numstr[0] isnt '0')
			{
			s = (numstr[0] isnt '1') ? numlett[Number(numstr[0])] : ''
			s $='cento'
			}
		numstr= numstr[-2 ..]
		if (numstr is '00')
			return s
		//decine
		if (Number(numstr) < 21) /*= in the map */
			s $= numlett[Number(numstr)]
		else{
			s $= numlett[Number(numstr[-2] $ '0')]
			//unita'
			if (numstr[-1] isnt '0')
				s $= numlett[Number(numstr[-1])]
			}
		s = s.Replace('ao' 'o').Replace('au' 'u').Replace('iu' 'u').Replace('oo' 'o')
		return s
		}
	}