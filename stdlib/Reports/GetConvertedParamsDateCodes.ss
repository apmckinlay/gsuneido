// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(params, reporterModel = false)
		{
		if not Object?(params)
			return false
		paramsConvert = Record()
		for mem in params.Members()
			{
			paramVal = params[mem]
			paramsConvert[mem] = Object?(paramVal) ? paramVal.Copy() : paramVal

			if .convertableParamVal?(paramVal)
				.convertParamIfDate(mem, reporterModel, paramVal, paramsConvert)
			else
				paramsConvert[mem] =
					GetConvertedParamsDateCodes(paramVal, reporterModel)
			}
		return paramsConvert
		}
	convertableParamVal?(paramVal)
		{
		return not Object?(paramVal) or
			(paramVal.Member?('operation') and paramVal.Member?('value'))
		}

	convertParamIfDate(mem, reporterModel, paramVal, paramsConvert)
		{
		if false is dateType = .dateType(mem, reporterModel)
			return

		if not Object?(paramVal)
			.convertDate(paramVal, dateType, { paramsConvert[mem] = it })
		else
			for val in #(value, value2)
				if paramVal.Member?(val)
					.convertDate(paramVal[val], dateType, {paramsConvert[mem][val] = it})
		}
	convertDate(paramVal, dateType, setFn)
		{
		if paramVal is "" or Date?(paramVal) or Object?(paramVal)
			return

		showTime = dateType is 'DateTime'
		date = DateControl.ConvertToDate(paramVal, convertDateCodes?:, :showTime)
		date = Date?(date) and not showTime ? date.NoTime() : date
		setFn(date)
		}

	dateType(field, reporterModel = false)
		{
		ctrl = Datadict(String(field)).Control
		if false isnt type = .getDateType(ParamsSelectControl.ControlName(ctrl))
			return type

		if reporterModel isnt false and
			false isnt type = .getDateType(reporterModel.GetFormulaType(field))
				return type

		return false
		}

	getDateType(type)
		{
		return type.Has?('Date')
			? type.Has?('Time')
				? 'DateTime'
				: 'Date'
			: false
		}
	}