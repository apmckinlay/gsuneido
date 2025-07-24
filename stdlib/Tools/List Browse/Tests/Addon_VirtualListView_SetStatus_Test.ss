// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cl: Addon_VirtualListView_SetStatus
		{
		Send(msg /*unused*/, rec) {
			if rec.extraValid isnt ""
				return rec.extraValid
			return 0 /* need to return based on rec */ }
		}
	Test_main()
		{
		view = FakeObject(GetGrid: FakeObject(RepaintRecord: #()),
			GetModel: Object(EditModel: .editModel()))
		statusBar = MockObject(#(
			#(Set, "", warn: false, normal:, invalid: false)
			))
		.cl.RefreshValid(view, statusBar, false)

		rec = FakeObject(New?: false)
		rec.extraValid = ""
		statusBar = MockObject(#(
			#(Set, "", warn: false, normal: false, invalid: false),
			#(SetValid, true)
			))
		.cl.RefreshValid(view, statusBar, rec)

		// Extra Validation
		rec = FakeObject(New?: false)
		rec.extraValid = Object("Test Extra Message", warn:)
		statusBar = MockObject(#(
			#(Set, "Test Extra Message", warn: true, normal: false, invalid: false),
			#(SetValid, true)
			))
		.cl.RefreshValid(view, statusBar, rec)

		// Warning on Modified Record
		view = FakeObject(GetGrid: FakeObject(RepaintRecord: #()),
			GetModel: Object(EditModel: .editModel("", "Test Warning Message")))
		rec = FakeObject(New?: false)
		rec.extraValid = ""
		statusBar = MockObject(#(
			#(Set, "Test Warning Message", warn:, normal: false, invalid: false),
			#(SetValid, true)
			))
		.cl.RefreshValid(view, statusBar, rec)

		// Warning on New Record
		view = FakeObject(GetGrid: FakeObject(RepaintRecord: #()),
			GetModel: Object(EditModel: .editModel("", "Test Warning Message")))
		rec = FakeObject(New?: true)
		rec.extraValid = ""
		statusBar = MockObject(#(
			#((Get), result: "",)
			#(Set, "", warn: false, normal: false, invalid: false),
			#(SetValid, true)
			))
		.cl.RefreshValid(view, statusBar, rec)

		// Warning but also Invalid
		view = FakeObject(GetGrid: FakeObject(RepaintRecord: #()),
			GetModel: Object(EditModel: .editModel("Test Valid Message",
				"Test Warning Message")))
		rec = FakeObject(New?: false)
		rec.extraValid = ""
		statusBar = MockObject(#(
			#(Set, "Test Valid Message", warn: false, normal: false, invalid: false),
			#(SetValid, false)
			))
		.cl.RefreshValid(view, statusBar, rec)

		// Invalid Modified Record
		view = FakeObject(GetGrid: FakeObject(RepaintRecord: #()),
			GetModel: Object(EditModel: .editModel("Test Valid Message")))
		rec = FakeObject(New?: false)
		rec.extraValid = ""
		statusBar = MockObject(#(
			#(Set, "Test Valid Message", warn: false, normal: false, invalid: false),
			#(SetValid, false)
			))
		.cl.RefreshValid(view, statusBar, rec)

		// New record Invalid Columns
		rec = FakeObject(New?: true)
		rec.extraValid = ""
		statusBar = MockObject(#(
			#((Get), result: "",)
			#(Set, "", warn: false, normal: false, invalid: false),
			#(SetValid, true)
			))
		.cl.RefreshValid(view, statusBar, rec)
		}

	editModel(invalidMsg = "", warningMsg = "")
		{
		edit = FakeObject(GetInvalidMsg: invalidMsg, GetWarningMsg: warningMsg)
		edit.ValidField = ""
		return edit
		}
	}
