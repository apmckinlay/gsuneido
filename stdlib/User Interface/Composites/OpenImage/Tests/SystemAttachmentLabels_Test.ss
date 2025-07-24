// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		testLib1 = .makeOpenImageClass(
			[label: 'Label1', condition: true, name: 'A'],
			[label: 'Label2', condition: true, name: 'B'],
			[label: 'Label2', condition: true, name: 'C'],
			[label: 'Label3', condition: true, name: 'D'],
			[label: 'Label4', condition: false, name: 'E']
			)
		// Testing 1 library
		Assert(labels = .func([testLib1]) isSize: 3)
		Assert(labels has: 'Label1')
		Assert(labels has: 'Label2')
		Assert(labels has: 'Label3')
		Assert(labels hasnt: 'Label4') // Condition not met

		testLib2 = .makeOpenImageClass(
			[label: 'Label3', condition: true, name: 'F'],
			[label: 'Label4', condition: true, name: 'G'],
			[label: 'Label5', condition: false, name: 'H'],
			[label: 'Label6', condition: true, name: 'I'],
			[label: 'Label6', condition: true, name: 'J'],
			[label: 'Label7', condition: false, name: 'K']
			)
		// Testing 2 libraries
		Assert(labels = .func([testLib1, testLib2]) isSize: 5)
		Assert(labels has: 'Label1')
		Assert(labels has: 'Label2')
		Assert(labels has: 'Label3')
		Assert(labels has: 'Label4') // Added from library testLib2
		Assert(labels has: 'Label6')
		Assert(labels hasnt: 'Label5') // Condition not met
		Assert(labels hasnt: 'Label7') // Condition not met

		// Testing 1 library with a secondary library present with OpenImage label classes
		Assert(labels = .func([testLib2]) isSize: 3)
		Assert(labels has: 'Label3')
		Assert(labels has: 'Label4')
		Assert(labels has: 'Label6')
		Assert(labels hasnt: 'Label1') // Other library (testLib1)
		Assert(labels hasnt: 'Label2') // Other library (testLib1)
		Assert(labels hasnt: 'Label5') // Condition not met
		Assert(labels hasnt: 'Label7') // Condition not met

		// Empty libraries passed in
		Assert(.func([]) isSize: 0)

		// Testing Libraries() (this will include our two test libraries).
		// Should have a minimum size of 5 due to the records this test outputs
		Assert(.func(false).Size() greaterThanOrEqualTo: 5)
		}

	makeOpenImageClass(@labels)
		{
		records = [table: .TempTableName()]
		labels.Each()
			{
			name = 'OpenImage_' $ it.name $ 'Label'
			records.Add([text: .text(it), :name, group: -1])
			}
		return .MakeLibraryRecord(@records)
		}

	text(label)
		{
		return `OpenImageLabelBase
	{
	Label: "` $ label.label $ `"
	Condition()
		{
		return ` $ label.condition $ `
		}
	}
	`
		}

	func(libraries) // Simulate SystemAttachmentLabels with Test_lib only
		{
		labels = Object()
		SystemAttachmentLabels.
			SystemAttachmentLabels_loop({ labels.AddUnique(it.Label) }, libraries)
		return labels.Remove('')
		}
	}
