// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(recordData, minimizeOutput? = false)
		{
		lineWarnings = Object()
		magicNums = .findMagicNumbers(recordData)
		warnings = .buildWarnings(magicNums, recordData, lineWarnings, minimizeOutput?)
		desc = .createMagicDescription(magicNums, minimizeOutput?)
		rating = Max(Qc_CalculateCodeRating.MaxRating - warnings.Size(), 0)

		return Object(warnings: warnings.UniqueValues(), :desc, :rating, :lineWarnings)
		}

	duplicateMagicNums(magicNums)
		{
		duplicateMagNums = Object().Set_default(Object())
		addToDesc = ""
		nums = magicNums.Map({ it.magicNum })
		for magNum in nums.UniqueValues()
			if 1 < magNumCount = nums.Count(magNum)
				duplicateMagNums[magNumCount].Add(magNum)
		for occurrenceNum in duplicateMagNums.Members().Sort!()
			{
			dups = ""
			duplicateMagNums[occurrenceNum].Each({ dups $= it $ ", " })
			dups = dups[..-2]
			addToDesc $= '\n' $ occurrenceNum $ " occurrences of magic number(s): " $ dups
			}
		return addToDesc
		}

	findMagicNumbers(recordData)
		{
		magicNums = Object()
		tokenTracker = Object(tokenFound?: false, token: "", prevToken: "",
			prev2Token: "", prev3Token: "", len: 0)
		.withScanner(recordData.code)
			{ |token, scan|
			tokenTracker.token = token
			if false isnt magicNum = .scanMagicNumber(recordData, scan, tokenTracker)
				magicNums.Add(magicNum)
			}
		return magicNums
		}

	withScanner(code, block)
		{
		scan = ScanCode(code)
		while scan isnt token = scan.Next()
			(block)(token, scan)
		}

	validTokens: ('=', ':', '#', '>>', '<<', '>>=', '<<=') // not magic if next to these
	numsToIgnore: ('0', '1', '2')// (The pos/neg of value here is ignored)
	methodCallsToIgnore: ('Round')

	scanMagicNumber(recordData, scan, tokenTracker)
		{
		type = scan.Type()
		if type in ("WHITESPACE", "COMMENT", "NEWLINE")
			return false
		magicNum = Object()
		prefix = .posMagicNumber(scan, tokenTracker, recordData.code)
			? ''
			: .negMagicNumber(scan, tokenTracker, recordData.code)
				? '-'
				: false

		if prefix isnt false
			{
			magicNum.pos = scan.Position() - tokenTracker.token.Size()
			magicNum.len = tokenTracker.token.Size()
			magicNum.lineNum = recordData.code[.. scan.Position()].Lines().Size()
			magicNum.magicNum = prefix $ tokenTracker.token
			}

		tokenTracker.prev3Token = tokenTracker.prev2Token
		tokenTracker.prev2Token = tokenTracker.prevToken
		tokenTracker.prevToken = tokenTracker.token
		return magicNum.Empty?() ? false : magicNum
		}

	negMagicNumber(scan, tokenTracker, code)
		{
		if tokenTracker.prevToken is '-' and
			.MagicNumber?(scan, tokenTracker.token, tokenTracker.prev2Token,
				tokenTracker.prev3Token)
			return not .magicNumComment?(scan, code)
		return false
		}

	magicNumComment?(scan, code)
		{
		line = code[scan.Position() ..].BeforeFirst('\n')
		return not line.Has?("=*/") and not line.Has?("= */") and
			(line.Has?("/* =") or line.Has?("/*="))
		}

	posMagicNumber(scan, tokenTracker, code)
		{
		if tokenTracker.prevToken isnt '-' and
			.MagicNumber?(scan, tokenTracker.token, tokenTracker.prevToken,
				tokenTracker.prev2Token)
			return not .magicNumComment?(scan, code)
		return false
		}
	MagicNumber?(scan, token, firstPrevToken, secondPrevToken)
		{
		return token.Number?() and
			scan.Context() is "code" and
			not .numsToIgnore.Has?(token) and
			not .validTokens.Has?(firstPrevToken) and
			not .methodCallsToIgnore.Has?(secondPrevToken)
		}

	buildWarnings(magicNums, recordData, lineWarnings, minimizeOutput?)
		{
		warnings = Object()
		for numOb in magicNums
			{
			if minimizeOutput?
				{
				lineWarnings.Add(Object(numOb.pos, numOb.len, warning:))
				warnings.Add(Record(name: recordData.lib $ ':' $ recordData.recordName $
					':' $ numOb.lineNum $ " - Magic Number"))
				}
			else
				warnings.Add(Record(name: "Magic number on line: " $ numOb.lineNum))
			}
		return warnings
		}

	createMagicDescription(magicNums, minimizeOutput?)
		{
		if not magicNums.Empty?()
			return "Magic numbers found" $ .duplicateMagicNums(magicNums)
		else if not minimizeOutput?
			return "No magic numbers were found"
		return ''
		}
	}