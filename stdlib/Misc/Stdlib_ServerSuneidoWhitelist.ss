// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
#(
	exact: (
		CSDev?:,
		CSDevServerWindowProc:,
		csDevPrint:,
		DevExeHandled:,
		DevHwnd:,
		DisabledOption:,
		init_times:,
		NetworkServiceIPAddress:,
		PopulateEventDisabled:,
		RackPerfStats:,
		RunAsStandalone:,
		RunningHttpServer:,
		SafeShutdown:,
		ServerContinuousTestInfo:,
		SchedulerHeartBeat:,
		SchedulerLastProcessStarted:,
		SchedulerLastProcessCompleted:,
		SchedulerStopEmail:,
		SessionsToKill:,
		SuneidoLog_MessageCount:, // SuneidoLog.messageCount
		SuneidoLog_MessageDate:, // SuneidoLog.messageDate
		start_time:,
		sujswebtest_done:,
		TestRunner:,
		TestRunningExpectedErrors:,
		TestRunningLogs:,
		ValidateQueryAny1?:,
		WebSocketExpiredConnection:,
		WebSocketReconnect:,
		wiki_regenerate:,
		workSpaceFindStop:
		),

	pattern: (
		"^CircuitBreaker_", // CircuitBreaker.memberPrefix
		"^ConnectionHandler-", // ConnectionErrorHandler.idPrefix
		"^Eta_EnterDepartProximity.Process",
		"^SchedulerLastLog_",
		"^TestWatchTable_",
		"^Tests_\d\d\d\d\d\d\d\d_\d\d\d\d",
		"^TimeoutTester_"
		)
	)
