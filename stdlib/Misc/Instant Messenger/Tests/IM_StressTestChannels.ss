// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(param)
		{
		setupParam = .setup(param)
		validity = .stress(setupParam)
		if validity is false
			ServerPrint('test failed at stress()')
		if setupParam.cleanOnFinish
			.cleanup(param)
		}

	setup(param)
		{
		// record to imchannels
		ServerPrint("setting up for stress testing")
		param.imchannel_num = Timestamp()
		QueryDo('insert {imchannel_num: ' $ Display(param.imchannel_num) $
			', imchannel_name: "stress_test2", imchannel_abbrev: "stresstest2",
				imchannel_status: "active"}
			into im_channels')
		return param
		}

	stress(param)
		{
		ServerPrint("beginning stress test")
		ServerPrint("test parameter : ",param)
		for (i = 0; i < param.messageCount; i++)
			{
			if false is IM_MessengerControl.SendDummyMessages(
				'GLOBAL','[stress_test] stress test external | Lorem Ipsum',
				param.imchannel_num)
				{
				ServerPrint('received fall response on iteration : ',i)
				return false
				}
			if not param.noSleep
				Thread.Sleep(param.sleepInterval)
			}
		ServerPrint('successfully added ',param.messageCount,'messages')
		return true
		}

	cleanup(param)
		{
		ServerPrint("cleaning up...")
		QueryDo('delete im_history where imchannel_num is ' $
			Display(param.imchannel_num))
		QueryDo('delete im_channels where imchannel_num is ' $
			Display(param.imchannel_num))
		}

	}
