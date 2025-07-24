// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
// Runs on the SVC server
function ()
	{
	// check SVC scheduler is running
	heartbeat = ServerSuneido.Get('SchedulerHeartBeat', Date.Begin())
	return heartbeat < Date().Plus(days: -1) ? 'SVC Scheduler Not Running' : ''
	}
