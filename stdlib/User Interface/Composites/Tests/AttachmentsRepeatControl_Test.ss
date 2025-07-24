// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_shiftAttachments()
		{
		// Doesnt't handle if dest is "" or if dest = src
		// Those are handled before a call to this method
		method = AttachmentsRepeatControl.AttachmentsRepeatControl_shiftAttachments
		// When src is lower than dest
		data = Object([attachment0: "0", attachment1: "", attachment2: "2",
			 attachment3: "3", attachment4: "4"])
		method(1, 4, data)
		expectedOutput = Object([attachment0: "0", attachment1: "2", attachment2: "3",
			 attachment3: "4", attachment4: "4"])
		Assert(data is: expectedOutput)
		// When src is lower than dest and there is an empty space in between
		data = Object([attachment0: "0", attachment1: "", attachment2: "2",
			 attachment3: "", attachment4: "4"])
		method(1, 4, data)
		expectedOutput = Object([attachment0: "0", attachment1: "", attachment2: "2",
			 attachment3: "4", attachment4: "4"])
		Assert(data is: expectedOutput)
		// When src is lower than dest multirow
		data = Object([attachment0: "0", attachment1: "", attachment2: "2",
			 attachment3: "3", attachment4: "4"],
			 [attachment0: "5", attachment1: "6", attachment2: "7",attachment3: "8",
			 attachment4: "9"])
		method(1, 7, data)
		expectedOutput = Object([attachment0: "0", attachment1: "2", attachment2: "3",
			 attachment3: "4", attachment4: "5"],
			 [attachment0: "6", attachment1: "7", attachment2: "7",attachment3: "8",
			 attachment4: "9"])
		Assert(data is: expectedOutput)
		 // When src is lower than dest and there is an empty space in between multirow
		data = Object([attachment0: "0", attachment1: "", attachment2: "2",
			 attachment3: "3", attachment4: "4"],
			 [attachment0: "", attachment1: "6", attachment2: "7",attachment3: "8",
			 attachment4: "9"])
		method(1, 7, data)
		expectedOutput = Object([attachment0: "0", attachment1: "", attachment2: "2",
			 attachment3: "3", attachment4: "4"],
			 [attachment0: "6", attachment1: "7", attachment2: "7",attachment3: "8",
			 attachment4: "9"])
		Assert(data is: expectedOutput)
		}
	Test_shiftAttachments2()
		{
		method = AttachmentsRepeatControl.AttachmentsRepeatControl_shiftAttachments
		// When src is higher than dest
		data = Object([attachment0: "0", attachment1: "1", attachment2: "2",
			 attachment3: "3", attachment4: ""])
		method(4, 1, data)
		expectedOutput = Object([attachment0: "0", attachment1: "1", attachment2: "1",
			 attachment3: "2", attachment4: "3"])
		Assert(data is: expectedOutput)
		// When src is higher than dest and there is a space in between
		data = Object([attachment0: "0", attachment1: "1", attachment2: "2",
			 attachment3: "", attachment4: ""])
		method(4, 1, data)
		expectedOutput = Object([attachment0: "0", attachment1: "1", attachment2: "1",
			 attachment3: "2", attachment4: ""])
		Assert(data is: expectedOutput)
		// When src is higher than dest multirow
		data = Object([attachment0: "0", attachment1: "1", attachment2: "2",
			 attachment3: "3", attachment4: "4"],
			 [attachment0: "5", attachment1: "6", attachment2: "7",attachment3: "8",
			 attachment4: ""])
		method(9, 1, data)
		expectedOutput = Object([attachment0: "0", attachment1: "1", attachment2: "1",
			 attachment3: "2", attachment4: "3"],
			 [attachment0: "4", attachment1: "5", attachment2: "6",attachment3: "7",
			 attachment4: "8"])
		Assert(data is: expectedOutput)
		 // When src is higher than dest and there is an empty space in between multirow
		data = Object([attachment0: "0", attachment1: "1", attachment2: "2",
			 attachment3: "3", attachment4: "4"],
			 [attachment0: "5", attachment1: "", attachment2: "7",attachment3: "8",
			 attachment4: ""])
		method(9, 1, data)
		expectedOutput = Object([attachment0: "0", attachment1: "1", attachment2: "1",
			 attachment3: "2", attachment4: "3"],
			 [attachment0: "4", attachment1: "5", attachment2: "7",attachment3: "8",
			 attachment4: ""])
		Assert(data is: expectedOutput)
		}

	Test_firstEmptySlot()
		{
		method = AttachmentsRepeatControl.AttachmentsRepeatControl_firstEmptySlot
		data = Object([attachment0: "0", attachment1: "1", attachment2: "2",
			 attachment3: "3", attachment4: ""])
		start = 0
		desc = false
		Assert(method(data, start, desc) is: 4)
		// Multiple Rows with multiple empty spaces
		data = Object([attachment0: "0", attachment1: "1", attachment2: "2",
			attachment3: "3", attachment4: "4"],
			[attachment0: "", attachment1: "6", attachment2: "7",attachment3: "8",
			attachment4: ""])
		start = 3
		Assert(method(data, start, desc) is: 5)
		// Starting in between two empty spaces
		start = 6
		Assert(method(data, start, desc) is: 9)
		// Going backwards
		desc = true
		Assert(method(data, start, desc) is: 5)
		// Going backwards
		data = Object([attachment0: "0", attachment1: "", attachment2: "2",
			attachment3: "3", attachment4: "4"],
			[attachment0: "", attachment1: "6", attachment2: "7",attachment3: "8",
			attachment4: ""])
		Assert(method(data, start, desc) is: 5)
		start = 0
		desc = false
		// Row is full
		data = Object([attachment0: "0", attachment1: "1", attachment2: "2",
			attachment3: "3", attachment4: "4"])
		 Assert(method(data, start, desc) is: 5)
		}
	}