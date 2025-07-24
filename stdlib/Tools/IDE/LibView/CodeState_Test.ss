// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.lib = .MakeLibrary(
			.validFunc = [
				name: #ValidFunc,
				text: 'function () { Print("Always valid") }'],
			.invalidFunc = [
				name: #InvalidFunc,
				text: 'function () { Print("Before it was invalid") }'
				lib_invalid_text: 'function () { } Print("After it was invalid")'],
			.invalidClass = [
				name: #InvalidClass,
				text: 'class { }',
				lib_invalid_text: 'class { InvalidText }'],
			.modifiedFunc = [
				name: #ModifiedFunc,
				text: 'function () { Print("Will be made invalid") }'])
		}

	Test_Valid?()
		{
		// Test getting the applicable code
		// 	- Valid Records: Will return their code
		//	- Invalid Records: Will return the code from their "hidden" records (*)
		//	- If a record is not found: Will return false
		rec = .validFunc.Copy()
		Assert(rec.lib_current_text is: 'function () { Print("Always valid") }')
		Assert(CodeState.Valid?(.lib, rec))

		rec.valid? = false // will preceed compilation checking
		Assert(rec.lib_current_text is: 'function () { Print("Always valid") }')
		Assert(CodeState.Valid?(.lib, rec) is: false)

		rec = .invalidFunc.Copy()
		Assert(rec.lib_current_text is: 'function () { } Print("After it was invalid")')
		rec.text = rec.lib_invalid_text // simulate text being made in valid
		Assert(CodeState.Valid?(.lib, rec) is: false)
		}

	Test_main()
		{
		.SpyOn(LibTreeModel.Libs).Return(Object(.lib))
		codeState = CodeState{ CodeState_log(@unused) { } }
		.inst = new codeState(.libTreeModel = new LibTreeModel)

		// Record is updated, code compiles
		rec = .libTreeRec(name = #ModifiedFunc, 'function () { Print("Still valid") }')
		.inst.Save(rec)
		Assert(rec.lib_current_text is: 'function () { Print("Still valid") }')

		// Record is updated, code no longer compiles, invalid field is filled
		// The loaded code will differ from the most recent invalid changes
		rec = .libTreeRec(name, code = 'function () { } Print("No longer valid")')
		.inst.Save(rec)
		Assert(rec.text is: '')
		Assert(rec.lib_invalid_text is: code)
		Assert(rec.lib_current_text is: code)

		// Record is updated, code still does not compile, invalid field is updated
		// The unloaded code will be the same as before, while the invalid code is updated
		rec = .libTreeRec(name, code = 'function () { } Print("Still invalid")')
		.inst.Save(rec)
		// text will not be present as the record was invalid. Text is removed to prevent
		// TreeModel from unnecessarily updating it
		Assert(rec hasntMember: #text)
		Assert(rec.lib_invalid_text is: code)
		Assert(rec.lib_current_text is: code)

		// The code is made valid, resulting in the clearing of the invalid column
		// and the updating text
		rec = .libTreeRec(name, code = 'function () { Print("Made valid once more") }')
		.inst.Save(rec)
		Assert(.inst.InvalidRec(.lib, name) is: false)
		Assert(rec.text is: code)
		Assert(rec.lib_invalid_text is: '')
		Assert(rec.lib_current_text is: code)
		}

	// Simulate a record being passed from ExplorerMultiControl to LibTreeModel
	libNumFactor: 100000
	libTreeRec(name, text = false)
		{
		rec = .libTreeModel.Get(Query1(.lib, :name).num += .libNumFactor)
		rec.parent = .libNumFactor
		if text isnt false
			rec.text = text
		return rec
		}

	Test_InvalidRecs()
		{
		m = CodeState.InvalidRecs

		invalidRecs = m(.lib)
		Assert(invalidRecs isSize: 2)
		Assert(invalidRecs has: .lib $ ':InvalidClass')
		Assert(invalidRecs has: .lib $ ':InvalidFunc')

		// Create record to be deleted, should be seen as invalid
		.MakeLibraryRecord(
			[name: #ToBeDeleted,
				text: 'class { }',
				lib_invalid_text: 'class { InvalidText }']
			table: .lib)
		invalidRecs = m(.lib)
		Assert(invalidRecs isSize: 3)
		Assert(invalidRecs has: .lib $ ':InvalidClass')
		Assert(invalidRecs has: .lib $ ':InvalidFunc')
		Assert(invalidRecs has: .lib $ ':ToBeDeleted')

		// Delete record, should no longer appear as an invalid record
		QueryDo('update ' $ .lib $ ' where name is #ToBeDeleted and group is -1
			set group = -2')
		invalidRecs = m(.lib)
		Assert(invalidRecs isSize: 2)
		Assert(invalidRecs has: .lib $ ':InvalidClass')
		Assert(invalidRecs has: .lib $ ':InvalidFunc')
		Assert(invalidRecs hasnt: .lib $ ':ToBeDeleted')
		}
	}
