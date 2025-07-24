// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(paramData = #())
		{
		printParams = paramData.Empty?() and _report.Params.Member?('printParams')
			? _report.Params.printParams
			: paramData.Members()
		if paramData.Empty?()
			paramData = _report.Params

		return .buildParamsFmt(printParams, paramData)
		}

	buildParamsFmt(printParams, paramData)
		{
		paramsfmt = Object('WrapItems')
		for param in printParams
			{
			if param is 'ReportDestination'
				continue

			if String?(param) and param.Suffix?('Filters') and Object?(paramData[param])
				{
				if paramData.GetDefault('printFilterNames', false) is true
					paramsfmt.Add(Object('Horz',
						Object('Text', Prompt(param) $ ': '),
						font: #(name: 'Arial', size: 8, weight: 'bold')))
				for ob in paramData[param]
					.handleParamData(ob.condition_field, ob, paramsfmt)
				}
			else
				.handleParamData(param, paramData, paramsfmt)
			}
		return paramsfmt
		}

	handleParamData(param, paramData, paramsfmt)
		{
		if Object?(paramData[param]) and paramData.Member?(param)
			GetParamsWhere.ConvertDateCodes(paramData[param])
		.ctrlFormat(param, paramData, paramsfmt)
		}

	ctrlFormat(param, paramData, paramsfmt)
		{
		// param is an object for formula fields from menu reporter reports
		prompt = .handlePrompt(param)
		paramformat = .handleParamFormat(param)
		param = .handleParamField(param)

		.updateFromRule(param, paramData)
		if .skipPrint?(paramData, param)
			return

		if false is format = .getFormat(paramData, param, paramformat)
			return

		if .skipPrintCheckmark?(format, paramData, param)
			return

		.convertWrapFormat(format)
		format.Delete('width')
		format.data = paramData[param]

		paramsfmt.Add(Object('Horz',
			Object('Text', Opt(prompt, ": ")), format,
			font: #(name: 'Arial', size: 8)))
		}

	handlePrompt(param)
		{
		return Object?(param) and param.Member?('paramPrompt')
			? param.paramPrompt
			: Prompt.WithInfo(param, fallbackMethod: 'Heading')
		}

	handleParamFormat(param)
		{
		return Object?(param) and param.Member?('paramFormat')
			? param.paramFormat
			: false
		}

	handleParamField(param)
		{
		if Object?(param) and param.Member?('paramField')
			param = param.paramField
		return param
		}

	updateFromRule(param, paramData)
		{
		if String?(param)
			{
			dd = Datadict(param, getMembers: #(ParamsNoSave))
			if dd.GetDefault('ParamsNoSave', false) is true
				paramData[param]
			}
		}

	skipPrint?(paramData, param)
		{
		return not paramData.Member?(param) or paramData[param] is ""
		}

	getFormat(paramData, param, paramformat)
		{
		// recognize ParamsSelect data
		if (Object?(paramData[param]) and
			paramData[param].Member?("operation") and
			paramData[param].Member?("value"))
			{
			if paramData[param].operation is ""
				return false
			format = Object("ParamsSelect" name: param, format: paramformat)
			}
		else
			format = Datadict(param).Format.Copy()
		return format
		}

	skipPrintCheckmark?(format, paramData, param)
		{
		return format[0] is "CheckMark" and paramData[param] is false
		}

	convertWrapFormat(format)
		{
		if String?(format[0]) and format[0].Has?('Wrap')
			format[0] = 'Text'
		}
	}
