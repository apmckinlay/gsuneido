// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(.prompts)
		{
		super(.layout())
		.addHeader()
		.Send(#Data)
		}

	skip: 4
	layout()
		{
		row = Object('Horz'
			Object('ChooseList', .prompts, name: 'condition_field')
			Object('Skip', .skip)
			Object('ChooseList', Select2.TranslateOps(), name: 'condition_op')
			Object('Skip', .skip)
			Object(.multiTypeField, width: 12, name: 'condition_value')
			Object('Skip', .skip)
			Object(.multiTypeField, width: 12, name: 'condition_return')
			name: 'condition_horz')
		return Object('Vert'
			Object('Repeat', row, maxRecords: FormulaIf.MaxBranches, name: 'repeat')
			Object('Static' 'Else', name:'elseStatic')
			Object(.multiTypeField, width: 12, name: 'condition_else'))
		}


	multiTypeField: FieldControl
		{
		Get()
			{
			try
				return Display(super.Get().SafeEval())
			return Display(super.Get())
			}
		}

	addHeader()
		{
		horz = .FindControl('condition_horz')
		header = Object('Horz'
			Object('Static' 'Field', name: 'fieldStatic')
			Object('Skip', .skip)
			Object('Static' 'Operator', name: 'opStatic')
			Object('Skip', .skip)
			Object('Static' 'Value', name:'valueStatic')
			Object('Skip', .skip)
			Object('Static' 'Then', name:'returnStatic')
			)
		content = .FindControl('content')
		headerCtrl = content.Insert(0, header)
		headers = headerCtrl.GetChildren()
		headers[0/*=Field*/].CalcXminByControls(horz.condition_field)
		headers[2/*=Operator*/].CalcXminByControls(horz.condition_op)
		headers[4/*=Value*/].CalcXminByControls(horz.condition_value)
		headers[6/*=Then*/].CalcXminByControls(horz.condition_return)
		}

	NewValue(@unused)
		{
		.Send("NewValue", .Get())
		}

	Get()
		{
		repeatCtrl = .FindControl('repeat')
		elseCtrl = .FindControl('condition_else')
		lines = repeatCtrl.Get()
		elseValue = elseCtrl.Get()
		sep = lines.Size() > 1 ? '\r\n\t' : ''
		res = ''
		for line in lines
			{
			lineStr = .translate(line) $ ', ' $
				line.GetDefault('condition_return', '""') $ ', '
			res $= lineStr $ sep
			}
		res $=  elseValue is '' ? Display('') : elseValue
		return res
		}

	Valid?()
		{
		return .FindControl('repeat').Valid?()
		}

	SetValid(@unused)
		{
		}

	stringOps: (
		('contains',			fn: 'CONTAINS')
		('does not contain',	fn: 'not CONTAINS')
		('starts with',			fn: 'STARTSWITH')
		('ends with',			fn: 'ENDSWITH')
		('matches',				fn: 'MATCHES')
		('does not match',		fn: 'not MATCHES'))
	translate(condition)
		{
		pos = .stringOps.FindIf({ it[0] is condition.condition_op })
		conditionVal = condition.GetDefault('condition_value', '""')
		if pos isnt false
			return .stringOps[pos].fn $ '(' $
				condition.condition_field $ ', ' $ conditionVal $ ')'

		pos = Select2.Ops.FindIf({|x| x[0] is condition.condition_op })
		if pos is false
			return condition.condition_field
		op = Select2.Ops[pos]
		operator = op[1] is '=' ? 'is' : op[1] is '!=' ? 'isnt' : op[1]
		return condition.condition_field $ " " $ operator $ " " $
			(op[0].Suffix?('empty') ? Display("") : conditionVal)
		}
	}