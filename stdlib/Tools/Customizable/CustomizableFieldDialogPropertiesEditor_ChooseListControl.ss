// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// Must be named like this in order to be found by CustomizableFieldDialog
	Name: 'peditor'

	RangeFrom: 10
	RangeTo: 40
	MaxItems: 300
	DefaultWidth: 10
	Format: 'Text'

	New()
		{
		.field_width = .FindControl('field_width')
		.list = .FindControl('items')
		.tooltip = .FindControl('tooltip')
		.status = .FindControl('Status')
		}

	Controls()
		{
		return Object('Vert',
			#Skip,
			Object('Static' 'Select Text Width (from ' $ Display(.RangeFrom) $
				' to ' $ Display(.RangeTo) $ ')'),
			#Skip,
			Object('Number' rangefrom: .RangeFrom, rangeto: .RangeTo,
				set: .DefaultWidth, name: 'field_width'),
			#Skip,
			Object('Static' 'Please enter items to choose from'),
			Object('List' columns: #('Items'), data: #(), name: 'items',
				defWidth: 380, headerSelectPrompt: 'no_prompts'),
			#Status,
			#Skip,
			Object('Static' 'Tooltip'),
			#Skip,
			Object('Field' name: 'tooltip'))
		}

	Valid?()
		{
		if not res = .validateListItems()
			return res
		value = .field_width.Get()
		return Number?(value) and .RangeFrom <= value and value <= .RangeTo
		}

	validateListItems()
		{
		list = .getList()
		msg = .GetValidMsg(list)
		.setStatusBar(msg)
		return msg is ""
		}

	GetValidMsg(list)
		{
		if list.Size() < 2
			return 'Should have at least two options to choose from'
		else if list.Size() > .MaxItems
			return 'Cannot have more than ' $ .MaxItems $ ' items\n\n' $
				'Use the "Text, from custom table" type if you need more items\n\n' $
				'Please contact Axon for assistance'
		return ""
		}

	setStatusBar(msg)
		{
		normal = msg is ''
		invalid = msg isnt ''
		.status.Set(msg, :invalid, :normal)
		}

	getList()
		{
		return .list.Get().Map({ it.Items.Trim() }).Remove("")
		}

	// should return an object with options i.e. (list: ('a' 'b') width: 50 )
	Get()
		{
		width = .field_width.Get()
		status = .tooltip.Get()
		return Object(control: Object(list: .getList(), :width, :status)
			format: Object(.Format, :width))
		}

	Set(object)
		{
		list = object.GetDefault('Control_list', #())
		data = list.Map({ Object(Items: it) })
		.list.Set(data)
		.field_width.Set(object.GetDefault('Control_width', .DefaultWidth))
		.tooltip.Set(object.GetDefault('Control_status', ''))
		}

	ConvertFieldType_CustomKey(fieldOb, data)
		{
		try
			.outputTableAndValues(fieldOb.colnme, fieldOb.custpe, data.control.list)
		catch(err)
			{
			SuneidoLog('ERROR: (CAUGHT) - Converting Choose List to Custom Key: ' $ err,
				params: Object(fieldOb, data),
				caughtMsg: 'no msg given; field did not convert; may need attention')
			return false
			}
		return Object(control: Object(customField: fieldOb.colnme), format: Object())
		}

	outputTableAndValues(field, prompt, list)
		{
		if not TableExists?(table = field $ '_table')
			CustomFieldControl_CustomTable(field, prompt, 'configlib')
		for desc in list
			if false is Query1(table, name: desc)
				QueryOutput(table, Record(name: desc))
		}

	List_Deletions()
		{
		return .validateListItems()
		}

	List_AfterEdit()
		{
		return .validateListItems()
		}
	}
