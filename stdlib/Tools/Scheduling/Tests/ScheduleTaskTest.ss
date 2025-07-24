// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	testTable1:	"schedulecontroltesttable1"
	getUIDTests:	(645431, -8, false)	// DO NOT CHANGE WITHOUT EXAMINING CODE
	Setup()
		{
		try
			Database("destroy " $ .testTable1)
		catch (x /*unused*/)
		ScheduleNextEvent.EnsureTaskTable(.testTable1)
		}
	Teardown()
		{
		// destroy test table; this should be destroyed after test; if an
		// exception is raised here, a test has failed, since this table
		// *should* exist subsequent to the running of the tests
		Database("destroy " $ .testTable1)
		}
	Test_General_Methods()
		{
		// test getUID
		for (counter = 0; counter < 100; counter++)
			{
			Assert(ScheduleTask.ScheduleTask_getUID(.testTable1)
				is: counter msg: "getUID")
			QueryOutput(.testTable1, Object(uid: counter))
			}
		for (counter = 0; counter < .getUIDTests.Size(); counter++)
			{
			QueryOutput(.testTable1, Object(uid: .getUIDTests[counter]))
			Assert(ScheduleTask.ScheduleTask_getUID(.testTable1)
				is: (.getUIDTests[0] + 1) msg: "getUID")
			}
		}
	Test_UpdateTaskRecord()
		{
		testOb = Object(taskTable: .testTable1, taskname: "this is", task: "a test")
		testOb.Set_readonly()

		testTask = ScheduleTask(testOb)
		Assert(testTask['uid']
			is: (.getUIDTests[0] + 1)
			msg: "New/getUID/UpdateTaskRecord")

		testOb2 = Object(taskTable:, taskname:, task:)
		testTask.UpdateTaskRecord(testOb2)	// Update testOb2 from testTask
		Assert(testOb is: testOb2)

		// test MinUID
		Assert((testTask.MinUID() is 0) and (ScheduleTask.MinUID() is 0)
			is: true msg: "MinUID")
		}
	}
