// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
RepeatControl
	{
	New()
		{
		super(.layout())
		.addHeader()
		}

	SetChooseFieldList()
		{
		sourceFields = .Send('ChooseConditions_SourceFields')
		for row in .GetRows()
			{
			row.FindControl('condition_source').SetList(sourceFields.Members())
			conditionRec = row.Get()
			source = conditionRec.condition_source
			fields = sourceFields.GetDefault(source, #())
			fieldCtrl = row.FindControl('condition_field')
			fieldCtrl.SetMandatory(not fields.Empty?())
			fieldCtrl.SetFieldMap(fields)
			fieldCtrl.Set(conditionRec.condition_field)
			}
		}

	Set(data)
		{
		super.Set(data)
		.SetChooseFieldList()
		}

	layout()
		{
		ops_desc = Select2.TranslateOps()
		return Object(#Horz
			#(ChooseList, width: 8, list: #() name: 'condition_source', mandatory:)
			#(Skip 4)
			#(FieldPrompt, width: 15, name: 'condition_field')
			#(Skip 4)
			Object('ChooseList', ops_desc, width: 12, name: 'condition_op', mandatory:)
			#(Skip 4)
			Object('Field', name: 'condition_value') name: 'condition_horz')
		}

	addHeader()
		{
		horz = .FindControl('condition_horz')
		header = Object('Horz'
			Object('Static' 'Source', name: 'sourceStatic')
			Object('Static' 'Field', name: 'fieldStatic')
			Object('Static' 'Operator', name: 'opStatic')
			Object('Static' 'Value', name:'valueStatic'))

		content = .FindControl('content')
		headerCtrl = content.Insert(0, header)

		headers = headerCtrl.GetChildren()
		skip = horz.GetChildren()[1]
		headers[0/*=condition_source*/].CalcXminByControls(
			Object(horz.condition_source, skip))
		headers[1/*=condition_field*/].CalcXminByControls(
			Object(horz.condition_field, skip))
		headers[2/*=condition_op*/].CalcXminByControls(Object(horz.condition_op, skip))
		headers[3/*=condition_value*/].CalcXminByControls(horz.condition_value)
		}

	On_Plus(source)
		{
		super.On_Plus(source)
		.SetChooseFieldList()
		}

	Record_NewValue(field, value /*unused*/, source)
		{
		if field is 'condition_source'
			{
			rec = source.Get()
			rec.condition_field = ''
			}
		}

	RepeatRecord_Changed(member)
		{
		super.RepeatRecord_Changed(member)
		.SetChooseFieldList()
		}
	}
