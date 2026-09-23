// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	setup()
		{
		.SpyOn(SvcCore.SvcCore_svcHooks).Return("")
		.SpyOn(SvcTable.Publish).Return(0)
		.svc = .Svc()
		.svcTable = .SvcTable(.lib = .MakeLibrary())
		}

	code(n)
		{
		return "function ()\r\n{\r\nreturn " $ n $ "\r\n}"
		}

	fmt(text)
		{
		return SvcAutoFormat.Format(#Test, text)
		}

	modify(name, text)
		{
		Transaction(update:)
			{|t|
			.svcTable.Update(.svcTable.Get(name, t), :t, newText: text)
			}
		}

	send(name, desc, asof = false, type = ' ')
		{
		return .svc.SendLocalChanges([[lib: .lib, :type, :name]], desc, #testuser, asof)
		}

	recText(name)
		{
		if false is rec = Query1(.lib, :name)
			throw "record not found: " $ name
		return rec.text
		}

	sent(since)
		{
		return QueryAll(.svcTable.Table() $ "_master where lib_committed > " $
			Display(since) $ " sort lib_committed")
		}

	Test_autoFormat?()
		{
		m = SvcAutoFormat.SvcAutoFormat_autoFormat?

		svcTable = Object(Type: #lib)
		Assert(m(svcTable, #RecordName))
		Assert(m(svcTable, #RecordName__other))
		Assert(m(svcTable, #RecordName__webgui))
		Assert(m(svcTable, #RecordName__trial) is: false)
		Assert(m(svcTable, #RecordName__alpha) is: false)

		svcTable = Object(Type: #book)
		Assert(m(svcTable, #RecordName) is: false)
		Assert(m(svcTable, #RecordName__other) is: false)
		Assert(m(svcTable, #RecordName__webgui) is: false)
		Assert(m(svcTable, #RecordName__trial) is: false)
		Assert(m(svcTable, #RecordName__alpha) is: false)
		}

	Test_unformatted_original_is_committed_separately()
		{
		.setup()
		.CommitAdd(.svc, .svcTable, name = #Rec, original = .code(1), #add)
		.modify(name, modified = .code(2))

		since = Timestamp()
		Assert(.send(name, "real change"))

		rows = .sent(since)
		Assert(rows isSize: 2)
		Assert(rows[0].comment is: SvcAutoFormat.Desc)
		Assert(rows[0].id is: #testuser)
		Assert(rows[0].text is: .fmt(original))
		Assert(rows[1].comment is: "real change")
		Assert(rows[1].text is: .fmt(modified))

		Assert(.recText(name) is: .fmt(modified))
		.AssertSvcEmpty(.svc, .lib)
		}

	Test_hashes_chain_through_both_commits()
		{
		.setup()
		.CommitAdd(.svc, .svcTable, name = #Rec, original = .code(1), #add)
		.modify(name, modified = .code(2))

		since = Timestamp()
		Assert(.send(name, "real change"))

		rows = .sent(since)
		Assert(rows isSize: 2)
		Assert(rows[0].lib_before_hash is: Svc.Hash(original))
		Assert(rows[0].lib_after_hash is: Svc.Hash(.fmt(original)))
		Assert(rows[1].lib_before_hash is: Svc.Hash(.fmt(original)))
		Assert(rows[1].lib_after_hash is: Svc.Hash(.fmt(modified)))
		}

	Test_formatted_original_gets_no_extra_commit()
		{
		.setup()
		.CommitAdd(.svc, .svcTable, name = #Rec, .fmt(.code(1)), #add)
		.modify(name, modified = .fmt(.code(2)))

		since = Timestamp()
		Assert(.send(name, "real change"))

		rows = .sent(since)
		Assert(rows isSize: 1)
		Assert(rows[0].comment is: "real change")
		Assert(rows[0].text is: modified)
		}

	Test_new_record_is_formatted_in_one_commit()
		{
		.setup()
		name = #NewRec
		unformatted = .code(1)
		.svcTable.Output([parent: 0, :name, text: unformatted])

		since = Timestamp()
		Assert(.send(name, "add it", type: '+'))

		rows = .sent(since)
		Assert(rows isSize: 1)
		Assert(rows[0].comment is: "add it")
		Assert(rows[0].text is: .fmt(unformatted))
		}

	Test_new_record_moved_in_is_formatted()
		{
		.setup()
		name = #MovedIn
		unformatted = .code(1)
		.svcTable.Output([parent: 0, :name, text: unformatted])
		QueryDo("update " $ .lib $ " where name is " $ Display(name) $
			" set lib_before_path = " $ Display(.lib $ "/old"))

		since = Timestamp()
		Assert(.send(name, "moved and changed", type: '+'))

		rows = .sent(since)
		Assert(rows isSize: 1)
		Assert(rows[0].comment is: "moved and changed")
		Assert(rows[0].text is: .fmt(unformatted))
		}

	Test_uncompilable_record_is_sent_as_is()
		{
		.setup()
		.CommitAdd(.svc, .svcTable, name = #Rec, "this is not code", #add)
		.modify(name, modified = "this is still not code")

		since = Timestamp()
		Assert(.send(name, "real change"))

		rows = .sent(since)
		Assert(rows isSize: 1)
		Assert(rows[0].text is: modified)
		}

	Test_move_is_not_formatted()
		{
		.setup()
		.CommitAdd(.svc, .svcTable, name = #Rec, original = .code(1), #add)
		QueryDo("update " $ .lib $ " where name is " $ Display(name) $
			" set lib_before_path = " $ Display(.lib $ "/old") $ ", lib_modified = " $
			Display(Timestamp()))

		since = Timestamp()
		Assert(.send(name, #moved))

		rows = .sent(since)
		Assert(rows isSize: 1)
		Assert(rows[0].text is: original)
		}

	Test_batch_send_chains_asof()
		{
		.setup()
		names = #(One, Two, Three)
		for n in names
			.CommitAdd(.svc, .svcTable, n, .code(1), "add " $ n)
		changes = Object()
		for n in names
			{
			.modify(n, .code(2))
			changes.Add(Object(lib: .lib, type: ' ', name: n))
			}

		since = Timestamp()
		Assert(.svc.SendLocalChanges(changes, "batch change", #testuser))

		rows = .sent(since)
		Assert(rows isSize: 6)
		Assert(rows.Map({ it.comment })
			is: [SvcAutoFormat.Desc, "batch change",
				SvcAutoFormat.Desc, "batch change",
				SvcAutoFormat.Desc, "batch change"])
		for n in names
			Assert(.recText(n) is: .fmt(.code(2)))
		.AssertSvcEmpty(.svc, .lib)
		}

	Test_failed_send_leaves_the_record_alone()
		{
		.setup()
		.CommitAdd(.svc, .svcTable, name = #Rec, .code(1), #add)
		.modify(name, modified = .code(2))

		since = Timestamp()

		Assert(.send(name, "real change", Date.Begin()) is: false)
		Assert(.sent(since) isSize: 0)
		Assert(.recText(name) is: modified)
		}

	Test_Format()
		{
		Assert(SvcAutoFormat.Format('X', "this is not code") is: false)

		formatted = SvcAutoFormat.Format('X', .code(1))
		Assert(formatted isnt: false)
		Assert(formatted isnt: .code(1))
		Assert(SvcAutoFormat.Format('X', formatted) is: formatted) // idempotent

		logs = .SpyOn(SuneidoLog).Return("").CallLogs()
		.SpyOn(AstFmtEquals?).Return(false)
		Assert(SvcAutoFormat.Format('X', .code(1)) is: false)
		Assert(logs isSize: 1)
		}
	}
