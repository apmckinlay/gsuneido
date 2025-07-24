// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_importFile()
		{
		test_cl = ImportExportObject
			{
			New()
				{
				.messages = Object()
				}
			openFileNameResults: ()
			ImportExportObject_doWithFile(file, block)
				{
				FakeFile(file, block)
				}
			ImportExportObject_displayMessage(title /*unused*/, message)
				{
				.messages.Add(message)
				}

			DisplayMessageCalls()
				{
				return .messages
				}

			SetOpenFileNameResults(results)
				{
				.openFileNameResults = results
				}

			ImportExportObject_openFileName(@unused)
				{
				return .openFileNameResults.PopFirst()
				}
			}

		str = ImportExportObject.ImportExportObject_exportKey
		fn = ImportExportObject.ImportExportObject_toText
		futureDate = Date().Plus(days: 8)
		futureExpiryDateOb =
			Object(import_export_expiry_date: futureDate, import_export_type: "test")
		openFileNameResults = Object(
			"",
			"invalid file",
			1,
			str $ fn(Object(import_export_type: "presets",
				import_export_expiry_date: futureDate)),
			str $ fn(#(import_export_expiry_date: #19000101)),
			str $ fn(futureExpiryDateOb)
			)
		cl = new test_cl
		cl.SetOpenFileNameResults(openFileNameResults)
		displayMessages = cl.DisplayMessageCalls()
		Assert(cl.Import('test import', 'test', false, false) is: false)

		Assert(cl.Import('test import', 'test', false, false) is: false)
		Assert(displayMessages[0] is: 'Unable to import. Invalid file.')
		Assert(cl.Import('test import', 'test', false, false) is: false)
		Assert(displayMessages[1] is: "Unable to import from:\n\n\t1")

		Assert(cl.Import('test import', 'test', false, false) is: false)
		Assert(displayMessages[2] has: 'Sorry, this is not a valid test file.\n\n')

		Assert(cl.Import('test import', 'test', false, false) is: false)
		Assert(displayMessages[3] has: 'Import not allowed')
		Assert(displayMessages[3] has: 'has expired (')

		Assert(cl.Import('test import', 'test', false, false) is: futureExpiryDateOb)
		}

	Test_exportImportFile()
		{
		text = 'sample exported file'
		packedContent = ImportExportObject.ImportExportObject_toText(text)

		unpackedContent = ImportExportObject.ImportExportObject_fromText(packedContent)
		Assert(unpackedContent is: text)
		}
	}
