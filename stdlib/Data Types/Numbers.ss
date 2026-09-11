// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	sig_Ceiling: "() :number"
	Ceiling()
		{
		i = .Int()
		return this > i ? i + 1 : i
		}

	sig_PercentToDecimal: "() :number"
	PercentToDecimal()
		{
		return (this * .01) /* = convert % to decimal */
		}

	sig_DecimalToPercent: "(round = 0) :number"
	DecimalToPercent(round = 0)
		{
		return (this * 100).Round(round) /* = convert decimal to % */
		}

	sig_Floor: "() :number"
	Floor()
		{
		i = .Int()
		return this < i ? i - 1 : i
		}

	sig_Abs: "() :number"
	Abs()
		{
		return this < 0 ? -this : this
		}

	sig_ToWords: "() :string"
	ToWords()
		{
		written = #(ZERO, ONE, TWO, THREE, FOUR,
			FIVE, SIX, SEVEN, EIGHT, NINE,
			TEN, ELEVEN, TWELVE, THIRTEEN,
			FOURTEEN, FIFTEEN, SIXTEEN, SEVENTEEN,
			EIGHTEEN, NINETEEN, TWENTY, 30: THIRTY,
			40: FORTY, 50: FIFTY, 60: SIXTY,
			70: SEVENTY, 80: EIGHTY, 90: NINETY)

		result = function(num, divisable, desc)
			{
			joinwith = divisable is 100 ? "AND " : ""
			result = (num / divisable).Int()
			remainder = num % divisable
			return result.ToWords() $ desc $
				(remainder > 0 ? ' ' $ joinwith $ remainder.ToWords() : "")
			}

		num = .Int()
		if num <= 20
			return written[num]
		else if num < 100
			{
			remainder = num % 10
			return written[num-remainder] $
				(remainder > 0 ? ' ' $ written[remainder] : "")
			}
		else if num < 1000
			return result(num, 100, " HUNDRED")
		else if num < 1_000_000
			return result(num, 1000, " THOUSAND")
		else
			return result(num, 1_000_000, " MILLION")
		}

	sig_EnFrancais: "() :string"
	EnFrancais()
		{
		return LanguageFrench.NumberToWords(this)
		}

	sig_ToWordsFrench: "() :string"
	ToWordsFrench()
		{
		return LanguageFrench.NumberToWords(this)
		}

	sig_ToWordsDutch: "() :string"
	ToWordsDutch()
		{
		return LanguageDutch.NumberToWords(this)
		}

	sig_ToWordSpanish: "() :string"
	ToWordSpanish()
		{
		return LanguageSpanish.NumberToWords(this)
		}

	sig_ToWordsItalian: "() :string"
	ToWordsItalian()
		{
		return LanguageItalian.NumberToWords(this)
		}

	sig_ToWordsSimple: "() :string"
	ToWordsSimple()
		{
		/* contributed by Johan Samyn
		Returns a string (in the current language), delimited by "* " and " *",
		summing up every digit in the number, from left to right,
		seperated by a space, and mentioning the decimal sign if applicable.
		The method also accepts negative numbers, and then returns a result
		starting with '* minus '.
		This simple notation is accepted in Belgium on official documents
		as a transformation in words of a number (e.g. invoice, check, ...).
		However, this method is ready to be used in any language available in
		translatelanguage.
		It can handle any number Suneido can handle, just like ToWordsDutch().
		*/
		words = #(zero, one, two, three, four,
			five, six, seven, eight, nine)
		decimal = TranslateLanguage(#decimalpoint)
		sNum = String(.Abs())
		if sNum[0] is '.'
			sNum = '0' $ sNum
		sResult = '*'
		if this < 0
			sResult $= ' ' $ TranslateLanguage("minus")
		for (i = 0; i < sNum.Size(); ++i)
			sResult $= ' ' $
				(sNum[i] is '.' ? decimal : TranslateLanguage(words[Number(sNum[i])]))
		return sResult $ " *"
		}

	sig_Pad: "(minSize, char = '0') :string"
	Pad(minSize, char = '0')
		{
		Assert(char.Size() is 1)
		str = String(.Abs()) // ignores sign
		n = minSize - str.Size()
		return (n > 0) ? char.Repeat(n) $ str : str
		}

	sig_Factorial: "() :number"
	Factorial()
		{
		if this > 0
			return this * (this - 1).Factorial()
		else if this is 0
			return 1
		else
			throw "Factorial: can't handle negative values"
		}

	sig_Int?: "() :boolean"
	Int?()
		{
		return this is .Int()
		}

	sig_IntDigits: "() :number"
	IntDigits()
		{
		num = .Abs().Int()
		return num is 0 ? 0 : String(num).Size()
		}

	sig_FracDigits: "() :number"
	FracDigits()
		{
		return String(.Abs().Frac()).Size() - 1
		}

	sig_EuroFormat: "(mask) :string"
	EuroFormat(mask)
		{
		return .Format(mask.Tr(",.", ".,")).Tr(",.", ".,")
		}

	sig_DollarFormat: "(mask) :string"
	DollarFormat(mask)
		{
		if String?(mask) and mask.Prefix?('-')
			mask = '(' $ mask.Replace("^-", "") $ ')'
		return '$' $ .Format(mask)
		}

	sig_Even?: "() :boolean"
	Even?()
		{
		return this%2 is 0
		}

	sig_Odd?: "() :boolean"
	Odd?()
		{
		return this%2 isnt 0
		}

	sig_Sign: "() :number"
	Sign()
		{
		return this < 0 ? -1 : 1
		}

	sig_RoundToPrecision: "(p) :number"
	RoundToPrecision(p)
		{
		if this is 0 or this is (1 / 0) or this is (-1 / 0)
			return this
		Assert(p greaterThan: 0)
		return .Round(-(.Abs().Log10().Floor() - (p - 1)))
		}

	sig_RoundToNearest: "(nearest = 1) :number"
	RoundToNearest(nearest = 1)
		{
		number = .Int()
		if nearest <= 0
			return number
		rem = number % nearest
		result = rem < nearest/2 ? number - rem : number + (nearest - rem)
		return result.Int()
		}

	sig_ToRGB: "() :object"
	ToRGB()
		{
		r = this & 0xff
		g = (this & 0xff00) >> 8
		b = (this & 0xff_0000) >> 16
		return [r, g, b]
		}

	sig_Times: "(block) :unknown"
	Times(block)
		{
		for (i = 0; i < this; ++i)
			block()
		}

	sig_MinutesInMs: "() :number"
	MinutesInMs()
		{
		return this * 60 * 1000
		}

	sig_SecondsInMs: "() :number"
	SecondsInMs()
		{
		return this * 1000
		}

	sig_SecondsInHours: "() :number"
	SecondsInHours()
		{
		return (this / 3600).Round(1)
		}

	sig_HoursInSeconds: "() :number"
	HoursInSeconds()
		{
		return (this * 3600).Round(0)
		}

	sig_HoursInMinutes: "() :number"
	HoursInMinutes()
		{
		return (this * 60).Round(0)
		}

	sig_SecondsInMinutes: "() :number"
	SecondsInMinutes()
		{
		return (this / 60).Round(0)
		}

	sig_InchesInTwips: "() :number"
	InchesInTwips()
		{
		return this * 1440 /* = twips in an inch */
		}

	sig_TwipsInInch: "() :number"
	TwipsInInch()
		{
		return this / 1440 /* = inch in twips */
		}

	sig_InchesInCanvasUnit: "() :number"
	InchesInCanvasUnit()
		{
		return this * 1440 / 17 /*= 1440 is twips in an inch,
			15 twips is 1 pixel in 96 DPI,
			not sure why we use 17, but it has been used in the calculatioin*/
		}

	sig_Kb: "() :number"
	Kb()
		{
		byteFactor = 1024.0
		return this * byteFactor
		}

	sig_Mb: "() :number"
	Mb()
		{
		byteFactor = 1024.0
		return this * byteFactor * byteFactor
		}

	sig_Gb: "() :number"
	Gb()
		{
		byteFactor = 1024.0
		return this * byteFactor * byteFactor * byteFactor
		}

	sig_SafeEval: "() :number"
	SafeEval()
		{
		return this
		}

	sig_Of: "(block) :object"
	Of(block)
		{
		if not Function?(block)
			{
			ob = Object()
			for ..this
				ob.Add(block)
			return ob
			}
		return Nof(this, block)
		}

	sig_OfStr: "(block) :string"
	OfStr(block)
		{
		s = ""
		for ..this
			s $= block()
		return s
		}
	}
