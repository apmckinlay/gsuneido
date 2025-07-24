// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Qccm_RecordLength()
		{
		recordData = Record()
		warningThreshold = Qc_RecordLength.Qc_RecordLength_warningThreshold

		text = "class{" $ "\n x = 555".Repeat(899) $ "}"
		recordData.code = text
		recordLength = Qc_RecordLength(recordData, true)
		Assert(recordLength.desc
			is: "Record is 900 lines long. Limit to " $ 700 $ " lines")
		Assert(recordLength.rating is: 3)

		text = "class{" $ "\n x = 555".Repeat(900) $ "}"
		recordData.code = text
		recordLength = Qc_RecordLength(recordData, true)
		Assert(recordLength.desc
			is: "Record is 901 lines long. Limit to " $ 700 $ " lines")
		Assert(recordLength.rating is: 0)

		text = "class{" $ "\n x = 55".Repeat(700) $ "}"
		recordData.code = text
		recordLength = Qc_RecordLength(recordData, true)
		Assert(recordLength.desc
			is: "Record is 701 lines long. Limit to " $ 700 $ " lines")
		Assert(recordLength.rating is: 3)

		text = "class{" $ "\n x = 5".Repeat(698) $ "}"
		recordData.code = text
		recordLength = Qc_RecordLength(recordData, true)
		Assert(recordLength.desc is:
			"Record is 699 lines long. Remain under " $ warningThreshold $ " lines")
		Assert(recordLength.rating is: 5)

		text = "class{" $ "\n x = 5".Repeat(699) $ "}"
		recordData.code = text
		recordLength = Qc_RecordLength(recordData, true)
		Assert(recordLength.desc is:
			"Record is 700 lines long. Limit to " $ warningThreshold $ " lines")
		Assert(recordLength.rating is: 5)

		text = ""
		recordData.code = text
		recordLength = Qc_RecordLength(recordData, true)
		Assert(recordLength.desc is: "")
		recordLength = Qc_RecordLength(recordData, false)
		Assert(recordLength.desc
			is: "Record is 0 lines long. Limit to " $ warningThreshold $ " lines")
		}
	}