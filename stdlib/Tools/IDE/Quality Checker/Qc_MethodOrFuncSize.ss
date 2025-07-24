// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	starRatings: #( 5: #(.05, .15, .30),
					4: #(.07, .22, .44),
					3: #(.08, .26, .50),
					2: #(.10, .30, .55),
					1: #(.13, .35, .60))
	thresholds: #(60, 50, 40)
	//	Empty lines do not contribute to line count
	//	Counts function header, '{', and '}' lines
	CallClass(recordData, minimizeOutput? = false)
		{
		lineWarnings = Object()
		methods = .calculateMethodSizes(recordData, minimizeOutput?)
		rating = Qc_CalculateCodeRating(methods.methodSizes, .thresholds, .starRatings)
		warnings = .generateWarnings(recordData, methods.largeMethodList, lineWarnings,
			minimizeOutput?)
		desc = .createDescription(warnings, minimizeOutput?, rating)
		return Object(:rating, :warnings, :desc, :lineWarnings)
		}

	calculateMethodSizes(recordData, minimizeOutput?)
		{
		largeMethodList = Object()
		methodSizes = Object()
		for method in recordData.qc_method_sizes
			{
			methodLength = method.lines
			methodSizes.Add(methodLength)
			if methodLength > .thresholds.Min() or not minimizeOutput?
				{
				methodName = method.name is '' ? 'function' : method.name
				lineNum = recordData.code[.. method.from].LineCount()
				largeMethodList.Add(Object(:methodName, :methodLength, :lineNum))
				}
			}
		largeMethodList.Sort!({|x,y| x.methodLength > y.methodLength })
		return Object(:methodSizes, :largeMethodList)
		}

	generateWarnings(recordData, largeMethodList, lineWarnings, minimizeOutput?)
		{
		warnings = Object()
		largeMethodList.Each()
			{
			name = it.methodLength $ ' lines in ' $ it.methodName
			lineNum = it.methodName is 'function' ? it.lineNum + 1 : it.lineNum
			if minimizeOutput?
				{
				name = recordData.lib $ ':' $ recordData.recordName $ ':' $ lineNum $
					" - " $ name
				lineWarnings.Add(Object(lineNum))
				}
			warnings.Add([:name])
			}
		return warnings
		}

	createDescription(warnings, minimizeOutput?, rating)
		{
		if warnings.Empty?()
			return minimizeOutput? ? '' : 'No methods found to check method sizes'

		affected = rating < Qc_CalculateCodeRating.MaxRating ? '' : ' not'
		return 'Method sizes - Rating' $ affected $ ' affected -> Limit to 40 lines'
		}
	}