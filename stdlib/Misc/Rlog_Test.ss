// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getRotationNum()
		{
		getRotationNum = Rlog.Rlog_getRotationNum
		daysPer = 1
		Assert(getRotationNum(#20061231, daysPer) is: -1)
		Assert(getRotationNum(#20070102, daysPer) is: 1)
		Assert(getRotationNum(#20080101, daysPer) is: 365)
		Assert(getRotationNum(#20070101, daysPer) is: 0)
		daysPer = 7
		Assert(getRotationNum(#20061231, daysPer) is: -1)
		Assert(getRotationNum(#20070101, daysPer) is: 0)
		Assert(getRotationNum(#20070102, daysPer) is: 0)
		Assert(getRotationNum(#20080101, daysPer) is: 52)
		}

	Test_currentLog()
		{
		Assert(Rlog.CurrentLog("log", -15) is: "log-5.log")
		Assert(Rlog.CurrentLog("log", 0) is: "log0.log")
		Assert(Rlog.CurrentLog("log", 0, -1) is: "log-1.log")
		Assert(Rlog.CurrentLog("log", 9) is: "log9.log")
		Assert(Rlog.CurrentLog("log", 9, 1) is: "log0.log")
		Assert(Rlog.CurrentLog("log", 10) is: "log0.log")
		Assert(Rlog.CurrentLog("log", 10, -1) is: "log9.log")
		Assert(Rlog.CurrentLog("log", 89) is: "log9.log")
		Assert(Rlog.CurrentLog("log", 89, -1) is: "log8.log")
		Assert(Rlog.CurrentLog("log", 89, 1) is: "log0.log")
		}

	Test_deleteOldFile()
		{
		Suneido.Delete(#Rlog)

		deleteOldFile = Rlog.Rlog_deleteOldFile
		mock = Mock(Rlog)
		mock.Eval(deleteOldFile, "log", "log1.log", 1)
		Assert(Suneido.Rlog["log"] is: 1)
		mock.Verify.deleteFile([anyArgs:])

		mock = Mock(Rlog)
		mock.Eval(deleteOldFile, "log", "log1.log", 1)
		mock.Verify.Never().deleteFile([anyArgs:])

		mock = Mock(Rlog)
		mock.Eval(deleteOldFile, "log", "log1.log", 2)
		Assert(Suneido.Rlog["log"] is: 2)
		mock.Verify.deleteFile([anyArgs:])

		Suneido.Delete(#Rlog)
		}

	Test_createNewFile?()
		{
		// 1 day per log
		Assert(Rlog.Rlog_createNewFile?("20070101.123456780", 0) is: false)
		Assert(Rlog.Rlog_createNewFile?("20070101.123456780", 1) is: true)

		cl = Rlog { DaysPerLog: 7 }
		Assert(cl.Rlog_createNewFile?("20070101.123456780", 0) is: false)
		Assert(cl.Rlog_createNewFile?("20070107.123456780", 0) is: false)
		Assert(cl.Rlog_createNewFile?("20070101.123456780", 1) is: true)
		Assert(cl.Rlog_createNewFile?("20070107.123456780", 1) is: true)
		Assert(cl.Rlog_createNewFile?("20070108.123456780", 1) is: false)
		Assert(cl.Rlog_createNewFile?("garbage", 0) is: true)
		}

	Test_firstPossibleRotationDate()
		{
		// 1 day per log
		Assert(Rlog.Rlog_firstPossibleRotationDate(0) is: #20070101)
		Assert(Rlog.Rlog_firstPossibleRotationDate(1) is: #20070102)
		Assert(Rlog.Rlog_firstPossibleRotationDate(31) is: #20070201)

		cl = Rlog { DaysPerLog: 7 }
		Assert(cl.Rlog_firstPossibleRotationDate(0) is: #20070101)
		Assert(cl.Rlog_firstPossibleRotationDate(1) is: #20070108)
		Assert(cl.Rlog_firstPossibleRotationDate(5) is: #20070205)
		}

	Test_formatMessage()
		{
		formatMessage = Rlog.Rlog_formatMessage
		Assert(
			formatMessage(#20001122.1234, 'message')
			is: '20001122.1234\tmessage\r\n')
		Assert(
			formatMessage(#20001122.1234, 'message\r\nmessage2\rmessage3\nmessage4')
			is: '20001122.1234\tmessage message2 message3 message4\r\n')

		rlog = Rlog
			{
			MemoryArena(unused) { return '118mb'}
			Rlog_sessionId(unused) { return 'rlogtest@127.0.0.1' }
			}
		formatMessage = rlog.Rlog_formatMessage
		Assert(
			formatMessage(#20001122.1234, '%m\t%s\tmessage\tparams')
			is: '20001122.1234\t118mb\trlogtest@127.0.0.1\tmessage\tparams\r\n')
		}

	Test_getLastLog()
		{
		testLog = 'testing_last_rlog'
		mock = Mock(Rlog)
		mock.When.dirLogFiles(testLog).Return(#())
		mock.When.getLastLog([anyArgs:]).CallThrough()
		Assert(mock.getLastLog(testLog) is: false)

		mock.When.Rlog_dirLogFiles(testLog).Return(
			Object(
				Object(date: #20170101, name: testLog $ "6.log"),
				Object(date: #20170102, name: testLog $ "7.log")))
		Assert(mock.getLastLog(testLog) is: testLog $ '7.log')

		mock.When.Rlog_dirLogFiles(testLog).Return(
			Object(
				Object(date: #20170101, name: testLog $ "9.log"),
				Object(date: #20170102, name: testLog $ "0.log")))
		Assert(mock.getLastLog(testLog) is: testLog $ '0.log')
		}

	Test_hoursPastMidnight()
		{
		daysPerLog = 1
		hoursPastMidnight = 2
		getRotationNum = Rlog.Rlog_getRotationNum
		Assert(getRotationNum(#20201217.2234, daysPerLog) is: 5099)
		Assert(getRotationNum(#20201218.0010, daysPerLog) is: 5100)
		Assert(getRotationNum(#20201218.0010, daysPerLog, hoursPastMidnight) is: 5099)
		Assert(getRotationNum(#20201218.015959, daysPerLog, hoursPastMidnight) is: 5099)
		Assert(getRotationNum(#20201218.0200, daysPerLog, hoursPastMidnight) is: 5100)
		Assert(getRotationNum(#20201218.0600, daysPerLog, hoursPastMidnight) is: 5100)
		}
	}
