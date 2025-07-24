// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	Export: false
	Print(x, y, w, h, data = #{})
		{
		if .Data isnt false
			data = .Data
		super.Print(x, y, w, h, (.layout(data)))
		}
	layout(data)
		{
		return .ParamsSelectObjectToString(data)
		}
	ParamsSelectObjectToString(ob)
		{
		if not Object?(ob)
			return ob
		str = ''
		value1 = .format_data(ob.value)
		value2 = .format_data(ob.value2)

		if ob.operation is "range"
			str $= TranslateLanguage('From') $ ' ' $ value1 $ ' ' $
				TranslateLanguage('To') $ ' ' $	value2
		else if ob.operation.Has?("in list")
			{
			fmtValues = Object()
			for o in ob.value
				fmtValues.Add(.format_data(o))
			str $= TranslateLanguage(ob.operation.CapitalizeWords()) $ ' (' $
				fmtValues.Join(', ') $ ')'
			}
		else if ob.operation is "empty" or ob.operation is "not empty"
			str $= TranslateLanguage(ob.operation)
		else
			{
			value = value1
			if ob.operation is "equals" and value is ""
				str $= TranslateLanguage(ob.operation) $ ' ""'
			else
				str $= TranslateLanguage(ob.operation) $ " " $ value
			}
		return str
		}
	format_data(value)
		{
		if Date?(value)
			return value.ShortDate()
		if not String?(value)
			return Display(value)
		return value
		}
	}
