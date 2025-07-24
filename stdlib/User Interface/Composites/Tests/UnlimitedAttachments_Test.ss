// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_GetLastRowPosition()
		{
		Assert(UnlimitedAttachments.GetLastRowPos("") is: 0)
		Assert(UnlimitedAttachments.GetLastRowPos(Object()) is: 0)
		ob = Object(#(attachment0: 'fred', attachment1: 'barney'))
		Assert(UnlimitedAttachments.GetLastRowPos(ob) is: 0)

		ob = Object(#(attachment0: 'fred', attachment1: 'barney'
			attachment2: 'wilma', attachment3: 'betty', attachment4: 'dino'))
		Assert(UnlimitedAttachments.GetLastRowPos(ob) is: 0)

		ob = Object(#(attachment0: 'fred', attachment1: 'barney'
			attachment3: 'betty', attachment4: 'dino'))
		Assert(UnlimitedAttachments.GetLastRowPos(ob) is: 0)


		ob = Object(
			#(attachment0: 'fred', attachment1: 'barney'
				attachment2: 'wilma', attachment3: 'betty', attachment4: 'dino')
			#(attachment0: 'bambam', attachment4: 'pebbles'))
		Assert(UnlimitedAttachments.GetLastRowPos(ob) is: 1)
		}

	Test_GetLastPosUsed()
		{
		Assert(UnlimitedAttachments.GetNextAvailRowPos("") is: 0)
		Assert(UnlimitedAttachments.GetNextAvailRowPos(Object()) is: 0)

		ob = Object(attachment0: 'fred')
		Assert(UnlimitedAttachments.GetNextAvailRowPos(ob) is: 1)

		ob = Object(attachment0: 'fred', attachment3: 'barney')
		Assert(UnlimitedAttachments.GetNextAvailRowPos(ob) is: 4)

		ob = Object(attachment0: 'fred', attachment4: 'barney')
		Assert(UnlimitedAttachments.GetNextAvailRowPos(ob) is: false)

		ob = Object(attachment4: 'barney')
		Assert(UnlimitedAttachments.GetNextAvailRowPos(ob) is: false)
		}

	Test_GetNextAvailPosition()
		{
		Assert(UnlimitedAttachments.GetNextAvailPos("") is: Object(row: 0, pos: 0))
		Assert(UnlimitedAttachments.GetNextAvailPos(Object()) is: Object(row: 0, pos: 0))

		ob = Object(#(attachment0: 'fred', attachment1: 'barney'))
		Assert(UnlimitedAttachments.GetNextAvailPos(ob) is: Object(row: 0, pos: 2))

		ob = Object(#(attachment0: 'fred', attachment3: 'barney'))
		Assert(UnlimitedAttachments.GetNextAvailPos(ob) is: Object(row: 0, pos: 4))

		ob = Object(#(attachment4: 'barney'))
		Assert(UnlimitedAttachments.GetNextAvailPos(ob) is: Object(row: 1, pos: 0))

		ob = Object(#(attachment0: 'a', attachment1: 'b', attachment2: 'c',
				attachment3: 'd', attachment4: 'e')
			#(attachment0: 'a', attachment1: 'b', attachment2: 'c',
				attachment3: 'd', attachment4: 'e')
			#(attachment0: 'a', attachment1: 'b'))
		Assert(UnlimitedAttachments.GetNextAvailPos(ob) is: Object(row: 2, pos: 2))
		}

	Test_GetValue()
		{
		Assert(UnlimitedAttachments.GetValue("", 0, 0) is: '')
		Assert(UnlimitedAttachments.GetValue(Object(), 0, 0) is: '')

		ob = Object(#(attachment0: 'a', attachment1: 'b'))
		Assert(UnlimitedAttachments.GetValue(ob, 0, 1) is: 'b')
		Assert(UnlimitedAttachments.GetValue(ob, 0, 2) is: '')
		Assert(UnlimitedAttachments.GetValue(ob, 1, 0) is: '')
		Assert(UnlimitedAttachments.GetValue(ob, 0, 'attachment0') is: 'a')
		}

	Test_listAttachment()
		{
		fn = UnlimitedAttachments.UnlimitedAttachments_listAttachment
		Assert(fn("", false) is: Object())
		Assert(fn(#(), false) is: Object())

		ob = #(#(attachment0: 'attach0', attachment1: 'attach1'))
		Assert(fn(ob, false) is: #(
			#(file: 'attach0', subfolder: '', fullPath: 'attach0',
				labels: '', row: 0, pos: 0),
			#(file: 'attach1', subfolder: '', fullPath: 'attach1',
				labels: '', row: 0, pos: 1)))

		ob = #(
			#(attachment0: '/test_folder/fred.csv  Axon Label: label1',
				attachment1: 'barney Axon Label: label1, label2'
				attachment2: ' Axon Label: label1',
				attachment3: 'betty Axon Label: label3',
				attachment4: 'dino')
			#(attachment0: 'bambam',
				attachment4: 'pebbles Axon Label: label1'))
		Assert(fn(ob, 'label1') is: #(
			#(file: 'fred.csv', subfolder: '/test_folder',
				fullPath: '/test_folder/fred.csv', labels: 'label1', row: 0, pos: 0)
			#(file: 'barney', subfolder: '', fullPath: 'barney',
				labels: 'label1, label2', row: 0, pos: 1)
			#(file: 'pebbles', subfolder: '', fullPath: 'pebbles',
				labels: 'label1', row: 1, pos: 4)))
		}

	Test_CompareAttachments()
		{
		fn = UnlimitedAttachments.CompareAttachments
		oldOb = ''
		newOb = #(#(attachment0: 'attach0', attachment1: 'attach1'))
		Assert(fn(oldOb, newOb, 'label1') is: #(changed?: false, newAttachments: #()))

		// old attachments is empty
		// new attachments is not empty
		newOb = #(
			#(attachment0: 'fred  Axon Label: label1',
				attachment1: 'barney Axon Label: label1, label2'
				attachment2: ' Axon Label: label1',
				attachment3: 'betty Axon Label: label3',
				attachment4: 'dino')
			#(attachment0: 'bambam',
				attachment4: 'pebbles Axon Label: label1'))
		Assert(fn(oldOb, newOb, 'label1') is: #(changed?:, newAttachments: #(
			#(file: 'fred', subfolder: '', fullPath: 'fred',
				labels: 'label1', row: 0, pos: 0)
			#(file: 'barney', subfolder: '', fullPath: 'barney',
				labels: 'label1, label2', row: 0, pos: 1)
			#(file: 'pebbles', subfolder: '', fullPath: 'pebbles',
				labels: 'label1', row: 1, pos: 4))))

		// old attachments is not empty
		// new attachments removed 'label3' from attachment1 at row 0
		oldOb = #(
			#(attachment0: 'fred  Axon Label: label1',
				attachment1: 'barney Axon Label: label1, label2, label3'
				attachment2: ' Axon Label: label1',
				attachment3: 'betty Axon Label: label3',
				attachment4: 'dino')
			#(attachment0: 'bambam',
				attachment4: 'pebbles Axon Label: label1'))
		Assert(fn(oldOb, newOb, 'label1').changed? is: false)

		// old attachments is not empty
		// new attachments removed 'label1' from attachment1 at row 0
		newOb = #(
			#(attachment0: 'fred  Axon Label: label1',
				attachment1: 'barney Axon Label: label2, label3'
				attachment2: ' Axon Label: label1',
				attachment3: 'betty Axon Label: label3',
				attachment4: 'dino')
			#(attachment0: 'bambam',
				attachment4: 'pebbles Axon Label: label1'))
		Assert(fn(oldOb, newOb, 'label1').changed?)

		// old attachments is not empty
		// new attachments changed the file name of attachment1 at row 0
		newOb = #(
			#(attachment0: 'fred  Axon Label: label1',
				attachment1: 'barney1 Axon Label: label1, label2, label3'
				attachment2: ' Axon Label: label1',
				attachment3: 'betty Axon Label: label3',
				attachment4: 'dino')
			#(attachment0: 'bambam',
				attachment4: 'pebbles Axon Label: label1'))
		Assert(fn(oldOb, newOb, 'label1').changed?)

		// old attachments is not empty
		// new attachments changed the file name of attachment3 at row 0
		newOb = #(
			#(attachment0: 'fred  Axon Label: label1',
				attachment1: 'barney Axon Label: label1, label2, label3'
				attachment2: ' Axon Label: label1',
				attachment3: 'betty1 Axon Label: label3',
				attachment4: 'dino')
			#(attachment0: 'bambam',
				attachment4: 'pebbles Axon Label: label1'))
		Assert(fn(oldOb, newOb, 'label1').changed? is: false)

		// old attachments is not empty
		// new attachments is empty
		newOb = ''
		Assert(fn(oldOb, newOb, 'label1') is: #(changed?: true, newAttachments: #()))
		}
	}
