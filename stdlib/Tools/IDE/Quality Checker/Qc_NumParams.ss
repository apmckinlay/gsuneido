// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	thresholds: #(9, 7, 5) //Number of params
	starRatings: #( 5: #(.005, .03, .10),
					4: #(.010, .05, .14),
					3: #(.030, .07, .20),
					2: #(.050, .10, .25),
					1: #(.080, .15, .30))

	CallClass(recordData, minimizeOutput? = false)
		{
		lineWarnings = Object()
		paramInfo = .calculateParamList(recordData)
		rating = Qc_CalculateCodeRating(checkCounts: paramInfo.Map({ it.numVars }),
			thresholds: .thresholds, starRatings: .starRatings)
		warnings = .generateWarnings(recordData, minimizeOutput?, lineWarnings, paramInfo)
		desc = .createDescription(warnings, minimizeOutput?, rating)
		return Object(:rating, :desc, :warnings, :lineWarnings)
		}

	generateWarnings(recordData, minimizeOutput?, lineWarnings, paramInfo)
		{
		warnings = Object()
		paramInfo.Each()
			{
			if it.numVars <= .thresholds.Min() and minimizeOutput?
				continue

			name = it.numVars $ ' parameters in ' $ it.methodName
			if minimizeOutput?
				{
				lineNum = it.methodName is 'function' ? it.lineNum + 1 : it.lineNum
				name = recordData.lib $ ':' $ recordData.recordName $ ':' $ lineNum $
					' - ' $ name
				lineWarnings.Add(Object(lineNum))
				}
			warnings.Add([:name])
			}
		return warnings
		}

	createDescription(warnings, minimizeOutput?, rating)
		{
		if minimizeOutput? and warnings.Empty?()
			return ''

		affected = rating < Qc_CalculateCodeRating.MaxRating ? '' : ' not'
		return 'Number of parameters - Rating' $ affected $ ' affected -> Limit to ' $
			.thresholds.Min()
		}


	calculateParamList(recordData)
		{
		paramInfo = Object()
		for method in recordData.qc_method_sizes
			{
			isClass? = method.to isnt ""
			numVars = .calculateNumParams(isClass?, method, recordData.code)
			methodName = isClass? ? method.name : 'function'
			lineNum = recordData.code[.. method.from].LineCount()
			paramInfo.Add(Object(:methodName, :numVars, :lineNum))
			}
		paramInfo.Sort!({ |x,y| x.numVars > y.numVars })
		return paramInfo
		}

	calculateNumParams(isClass?, method, code)
		{
		methodVars = Object()
		codeToProcess = isClass? ? code[method.from .. method.to] : code[method.from ..]
		lineAdjust = isClass? ? 0 : 1

		ClassHelp.RetrieveParamsList(codeToProcess, methodVars)
		return methodVars.Size() - lineAdjust - .calcNumVarsWithDefaults(codeToProcess)
		}

	calcNumVarsWithDefaults(code)
		{
		numDefaultVars = 0
		code = code.AfterFirst('(')
		scan = Scanner(code)
		nest = 1
		while scan isnt token = scan.Next()
			{
			if scan.Type() in (#WHITESPACE, #NEWLINE, #COMMENT)
				continue
			if token is '('
				nest++
			else if token is ')'
				nest--
			else if token is '='
				numDefaultVars++
			if nest is 0
				break
			}
		return numDefaultVars
		}
	}