// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	thresholds: (15, 12, 8)
	starRatings: #( 5: #(.025, .075, .20),
					4: #(.040, .100, .25),
					3: #(.050, .130, .30),
					2: #(.060, .150, .33),
					1: #(.080, .180, .35))

	//The McCabe complexity - The number of branching points + 1
	CallClass(recordData, minimizeOutput? = false)
		{
		lineWarnings = Object()
		complexityInfo = .calculateComplexities(recordData)

		counts = complexityInfo.Map({ it.complexity })
		rating = Qc_CalculateCodeRating(counts, .thresholds, .starRatings)

		warnings = .generateWarnings(complexityInfo, minimizeOutput?, recordData,
			lineWarnings)
		desc = .createDescription(warnings, minimizeOutput?, rating)
		return Object(:desc, :rating, :warnings, :lineWarnings)
		}

	calculateComplexities(recordData)
		{
		complexityInfo = Object()
		for method in recordData.qc_method_sizes
			{
			complexity = .calculateMcCabeComplexity(recordData.code, method)
			lineNum = recordData.code[.. method.from].LineCount()
			complexityInfo.Add(Object(methodName: method.name, :complexity, :lineNum))
			}
		complexityInfo.Sort!({|x,y| x.complexity > y.complexity })

		return complexityInfo
		}

	generateWarnings(complexityInfo, minimizeOutput?, recordData, lineWarnings)
		{
		warnings = Object()
		for method in complexityInfo
			{
			if .libRecIsFunction(method.methodName)
				{
				gotoLineNum = method.lineNum + 1
				warningSuffix = "Function McCabe complexity: " $ method.complexity
				}
			else
				{
				gotoLineNum = method.lineNum
				warningSuffix = "McCabe complexity is "  $ method.complexity $ " for " $
					method.methodName
				}

			warning = ''
			if not minimizeOutput?
				warning = warningSuffix
			else if method.complexity > .thresholds.Min()
				warning = recordData.lib $ ':' $ recordData.recordName $ ':' $
					gotoLineNum $ " - " $ warningSuffix

			if warning isnt ''
				{
				if minimizeOutput?
					lineWarnings.Add(Object(gotoLineNum))
				warnings.Add([name: warning])
				}
			}
		return warnings
		}

	createDescription(warnings, minimizeOutput?, rating)
		{
		if warnings.Empty?()
			return minimizeOutput? ? '' : 'No methods found to run McCabe complexity on'

		ratingAffect = rating < Qc_CalculateCodeRating.MaxRating ? '' : ' not'
		return "McCabe function complexity - Rating" $ ratingAffect $
			" affected -> Limit to " $ .thresholds.Min()
		}

	calculateMcCabeComplexity(code, method)
		{
		mcComplexity = 1
		if not .libRecIsFunction(method.to)
			code = code[method.from..method.to]
		scan = Scanner(code)
		while scan isnt token = scan.Next()
			if token in ('if', 'case', 'default', 'while', 'catch',
				'for', 'and', 'or', '&&', '||', '?')
				mcComplexity++

		return mcComplexity
		}

	libRecIsFunction(methodName)
		{
		return methodName is ''
		}
	}
