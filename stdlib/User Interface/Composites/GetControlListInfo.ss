// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(field, ctrl = false)
		{
		if not String?(field) or field.Blank?()
			return #()

		if ctrl is false
			ctrl = .GetControlFromField(field)
		return .getControlListInfo(ctrl)
		}

	GetControlFromField(field)
		{
		if not String?(field)
			return field

		dd = Datadict(field)
		control = dd.Member?('SelectControl') ? dd.SelectControl : dd.Control
		control = control.Copy()
		control.name = field
		control.readonly = false
		return control
		}

	getControlListInfo(ctrl)
		{
		if ctrl.Member?('listField')
			{
			list = [][ctrl.listField]
			if ctrl.Member?('displayField')
				list = list.Map({ it.ProjectValues(Object(ctrl.displayField)).Join(',') })
			if String?(list)
				list = list.Split(ctrl.GetDefault('splitValue', ','))
			return list.Map({ .formatItem(it, ctrl) }).Instantiate()
			}
		else if ctrl.Member?('list')
			return ctrl.list.Map({ .formatItem(it, ctrl) }).Instantiate()
		else if ctrl[0] is 'ChooseList' and ctrl.Member?(1)
			return ctrl[1].Map({ .formatItem(it, ctrl) })
		else if ctrl[0] is 'ChooseDate'
			return ctrl
		else
			{
			ctrlClass = Global(ctrl[0] $ 'Control')
			if ctrlClass.Method?('GetList')
				{
				list = ctrlClass.GetList().Copy()
				if not list.Member?('query')
					return list.Map({ .formatItem(it, ctrl) })
				list.keys = ctrl.GetDefault('keys', false)
				return list
				}
			else if ctrl[0] is 'CustomKey'
				return .customKeyGetListInfo(ctrl)
			else if ctrlClass.Method?('Key_BuildQuery')
				return .keyGetListInfo(ctrl, ctrlClass)
			}
		return #()
		}

	keyGetListInfo(ctrl, ctrlClass)
		{
		ctrlq = ctrl.Member?('query')
			? ctrl.query
			: ctrl.GetDefault(1, '')
		if Function?(ctrlq)
			ctrlq = ctrlq()
		query = ctrlClass.Key_BuildQuery(ctrlq,
			.getRestrictions('restrictions', ctrl),
			.getRestrictions('invalidRestrictions', ctrl),
			noSend:)
		optionalRestrictions = ctrl.GetDefault(#optionalRestrictions, #())

		numField = ctrl.Member?('field')
			? ctrl.field
			: ctrl.Member?('numField')
				? ctrl.numField
				: ctrl.GetDefault(2, false)

		displayField = ctrl.Member?('displayField')
			? ctrl.displayField
			: ctrl.GetDefault('nameField',
				numField.Replace("(_num|_name|_abbrev)$", "_name"))
		columns = .getColumns(ctrl, query)
		customizeQueryCols = ctrl.GetDefault(#customizeQueryCols, false)
		excludeSelect = ctrl.GetDefault(#excludeSelect, #())
		if String?(excludeSelect)
			excludeSelect = Global(excludeSelect)()
		keys = ctrl.GetDefault('keys', false)
		return Object(:query, :columns, field: numField, :displayField,
			:customizeQueryCols, :keys, :optionalRestrictions, :excludeSelect)
		}

	customKeyGetListInfo(ctrl)
		{
		query = ctrl.customField $ '_table'
		field = QueryStrategy(query).Extract('\(.*\)').Tr('()')
		columns = QueryColumns(query)
		customizeQueryCols = ctrl.GetDefault(#customizeQueryCols, true)
		keys = ctrl.GetDefault('keys', false)
		return Object(:query, :columns, :field,
			displayField: field, :customizeQueryCols, :keys)
		}

	getRestrictions(mem, control)
		{
		if not control.Member?(mem) or control[mem] is false
			return false
		return control[mem].Has?(' ') // query string, not rule
			? control[mem]
			: [][control[mem]]
		}

	getColumns(ctrl, query)
		{
		return ctrl.Member?('cols')
			? ctrl.cols
			: ctrl.Member?('columns')
				? ctrl.columns
				: ctrl.GetDefault(3, QueryColumns(query))
		}

	formatItem(item, ctrl)
		{
		sep = ctrl.GetDefault('listSeparator', ' - ')
		if sep isnt ''
			item = item.BeforeFirst(sep)
		return item.Trim()
		}
	}
