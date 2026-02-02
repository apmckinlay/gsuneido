// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Process(errMsg, errorOb, msg, classId, threshold = 10)
		{
		logType = "ERROR: (CAUGHT) "
		caughtMsg = "generic connection error handling; may need attention"
		if .errorMatched?(errMsg, errorOb)
			{
			numberOfConnectionErrs = .CalcNumberOfConnectionErrs(classId, errorOb, errMsg)
			if numberOfConnectionErrs < threshold
				{
				logType = "INFO: "
				caughtMsg = ''
				}
			}
		err = errMsg.Has?('<Message>')
			? errMsg.AfterFirst('<Message>').BeforeFirst('</Message>')
			: errMsg
		.AddToLog(logType $ classId $ ': ' $ msg $ ' - ' $ err, classId, :caughtMsg)
		.updateLogCount(classId)
		}

	errorMatched?(errMsg, errorOb)
		{
		for error in errorOb
			if errMsg.Has?(error)
				return true

		return false
		}

	CalcNumberOfConnectionErrs(classId, errorOb, errMsg)
		{
		if .errorMatched?(errMsg, errorOb)
			return ServerSuneido.Get(.idPrefix $ classId, #(0))[0]
		else
			return 0
		}

	AddToLog(msg, caughtMsg = '')
		{
		SuneidoLog(msg, :caughtMsg)
		}

	idPrefix: 'ConnectionHandler-'
	updateLogCount(classId)
		{
		id = .idPrefix $ classId
		.resetIfNeeded(id)
		log = ServerSuneido.Get(id)
		ServerSuneido.Set(id, Object(log[0] + 1, createdAt: log.createdAt))
		}

	resetIfNeeded(id)
		{
		if not ServerSuneido.HasMember?(id) or .getCountDuration(id) > 24 /*=hours*/
			ServerSuneido.Set(id, Object(0, createdAt: Date()))
		}

	getCountDuration(id)
		{
		return Date().MinusHours(ServerSuneido.Get(id, #()).GetDefault('createdAt',
			Date()))
		}
	}