// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
Test
	{
	field1: class
		{
		New(.readOnly = false)
			{ }
		GetReadOnly()
			{ return .readOnly }
		Dirty?(unused = '') // Based on stdlib:Control.Dirty?
			{ return false }
		}
	field2: class
		{
		New(.readOnly = false)
			{ }
		GetReadOnly()
			{ return .readOnly }
		dirty?: false
		Dirty?(dirty? = '')
			{
			if Boolean?(dirty?)
				.dirty? = dirty?
			return .dirty?
			}
		Get()
			{ return .value }
		NewValue(value)
			{ .value = value.Lower() }
		}
	field3: class
		{
		New(.readOnly = false)
			{ }
		GetReadOnly()
			{ return .readOnly }
		dirty?: false
		Dirty?(dirty? = '')
			{
			if Boolean?(dirty?)
				.dirty? = dirty?
			return .dirty?
			}
		Field: class
			{
			Process_newvalue()
				{
				Assert(true) // Should see code coverage here
				}
			}
		}
	Test_main()
		{
		mock = Mock(ControlFillin)
		mock.When.FillinFields([anyArgs:]).CallThrough()
		mock.ControlFillin_collectedFields = #()
		mock.ControlFillin_tabsControl = false
		mock.When.ensureConstructedControl([anyArgs:]).CallThrough()

		fieldMap = Object(
			boolean: 		new .field1(),
			boolean_yesno: 	new .field1(readOnly:),
			address1: 		new .field2(),
			address2: 		new .field2(readOnly:),
			id:				new .field3(),
			desc: 			false)
		for field, control in fieldMap
			mock.When.findControl(field).Return(control)

		// Empty fillinData
		mock.FillinFields(fieldMap.Members(), [], controlData = [])
		Assert(controlData isSize: 0)

		fillinData = [
			boolean: true,
			boolean_yesno: false,
			address1: 'Street Address 1',
			address2: 'Street Address 2',
			desc: 'Generic Description',
			id: '<id value>']
		mock.FillinFields(fieldMap.Members(), fillinData, controlData = [])
		Assert(controlData isSize: 3)
		Assert(controlData.boolean)
		Assert(controlData.address1 is: 'street address 1')
		Assert(controlData.id is: '<id value>')

		// Simulating base Control.Dirty?
		Assert(fieldMap.boolean.Dirty?() is: false)
		Assert(fieldMap.boolean_yesno.Dirty?() is: false)
		// Controls are updated and marked as dirty
		Assert(fieldMap.address1.Dirty?())
		Assert(fieldMap.id.Dirty?())
		// Control is readonly so it is not updated or marked as dirty
		Assert(fieldMap.address2.Dirty?() is: false)
		}

	Test_ensureConstructedControl()
		{
		mock = Mock(ControlFillin)
		mock.When.ensureConstructedControl([anyArgs:]).CallThrough()
		mock.When.EnsureTab([anyArgs:]).CallThrough()
		mock.When.findControl([anyArgs:]).Return(control = new .field1, false)

		// First call to .findControl returns a control
		Assert(mock.ensureConstructedControl('field1') is: control)

		// First call to .findControl returns false, tabsControl is false
		mock.ControlFillin_tabsControl = false
		Assert(mock.ensureConstructedControl('field1') is: false)
		// Ensure the public method tabsControl check is covered
		Assert(mock.EnsureTab('field1') is: false)

		// First call to .findControl returns false, CollectFields.FindTab fails to find
		tabsMock = Mock()
		tabsMock.When.FindTab([anyArgs:]).Return(false, 'tabIdx')
		tabsMock.When.ConstructAndSetTab([anyArgs:]).Do({ })
		mock.ControlFillin_tabsControl = tabsMock

		// field1 does not exist in collectedFields
		mock.ControlFillin_collectedFields = #(((section: 'Tab1'), (name: 'field2')))
		Assert(mock.ensureConstructedControl('field1') is: false)

		// First call to .findControl returns false, fails to find tab, returns false
		mock.When.findControl('field2').Return(false)
		Assert(mock.ensureConstructedControl('field2') is: false)
		tabsMock.Verify.Never().Constructed?([anyArgs:])

		// First call to .findControl returns false, fails to find tab
		// second call to .findControl returns false
		mock.When.findControl('field2').Return(false)
		tabsMock.When.Constructed?([anyArgs:]).Return(false, true)
		Assert(mock.ensureConstructedControl('field2') is: false)

		// First call to .findControl returns false, tab is constructed
		// second call to .findControl returns false
		mock.When.findControl('field2').Return(false)
		tabsMock.When.Constructed?([anyArgs:]).Return(false, true)
		Assert(mock.ensureConstructedControl('field2') is: false)

		// First call to .findControl returns false, tab is already constructed
		// second call to .findControl returns false
		mock.When.findControl('field2').Return(false, control = new .field2)
		tabsMock.When.Constructed?([anyArgs:]).Return(true)
		Assert(mock.ensureConstructedControl('field2') is: control)
		}
	}