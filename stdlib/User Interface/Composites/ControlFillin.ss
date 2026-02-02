// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
class
	{
	// README: Only call the public methods after the parent control has been constructed.
	// Else you risk having .tabsControl initialized to: false.
	// This is done to reduce the number of times we call: .findControl('Tabs')
	New(.parent, .layout)
		{
		}

	FillinFields(fields, fillinData, controlData)
		{
		// Process boolean fields first as they can control other field's edit states
		for field in booleanFields = fields.Filter({ DatadictType(it) is 'boolean' })
			.FillinField(field, fillinData, controlData)

		// Process remaining fields
		for field in fields.Remove(@booleanFields)
			.FillinField(field, fillinData, controlData)
		}

	FillinField(field, fillinData, controlData)
		{
		if not fillinData.Member?(field)
			return false

		if false is control = .ensureConstructedControl(field)
			return false

		if control.GetReadOnly()
			return false

		controlData[field] = fillinData[field]
		control.Dirty?(true)

		// Trigger Key control lookups / fillins
		if control.Member?('Field') and control.Field.Method?('Process_newvalue')
			control.Field.Process_newvalue()

		if not control.Method?('NewValue')
			return true

		control.NewValue(controlData[field], source: control)
		controlData[field] = control.Get()
		return true
		}

	// NOTE: This does not currently handle nested tab controls
	ensureConstructedControl(field)
		{
		if false isnt control = .findControl(field)
			return control

		if false is .tabsControl
			return false // Control should be constructed already

		if false is tab = CollectFields.FindTab(field, .collectedFields)
			return false // Failed to find the field's tab

		return .EnsureTab(tab)
			? .findControl(field)
			: false
		}

	findControl(field)
		{
		return .parent.FindControl(field)
		}

	getter_tabsControl()
		{
		return .tabsControl = .findControl('Tabs')
		}

	getter_collectedFields()
		{
		return .collectedFields = CollectFields(.layout, path?:)
		}

	EnsureTab(tab)
		{
		if false is .tabsControl
			return false

		if false is tabIdx = .tabsControl.FindTab(tab)
			return false

		if .tabsControl.Constructed?(tabIdx)
			return true

		.tabsControl.ConstructAndSetTab(tabIdx)
		return .tabsControl.Constructed?(tabIdx)
		}
	}
