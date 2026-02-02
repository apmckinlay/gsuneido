// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_beforeRecord()
		{
		mock = Mock(IM_ChannelsControl)
		mock.When.beforeRecord([anyArgs:]).CallThrough()
		mock.When.getUserSettingOb('IM_MentionedChannels').Return(#(1, 2))
		mock.When.getUserSettingOb('IM_UnreadChannels').Return(#(1, 3))
		mock.When.getLastActivity(1).Return(false)
		mock.When.getLastActivity(2).Return([im_num: Date()])
		mock.When.getLastActivity(3).Return([im_num: Date().Minus(days: 2)])

		mock.beforeRecord(rec = [imchannel_num: 1], false)
		Assert(rec.imchannel_num is: 1)
		Assert(rec.imchannel_notification is: CLR.orange)
		Assert(rec.imchannel_lastActivity is: 'No Activity')

		mock.beforeRecord(rec = [imchannel_num: 2], #())
		Assert(rec.imchannel_num is: 2)
		Assert(rec.imchannel_notification is: CLR.orange)
		Assert(rec.imchannel_lastActivity is: '0 day(s) ago')

		mock.beforeRecord(rec = [imchannel_num: 3], #(4, 3, 2))
		Assert(rec.imchannel_num is: 3)
		Assert(rec.imchannel_notification is: CLR.GREEN)
		Assert(rec.imchannel_lastActivity is: '2 day(s) ago')
		Assert(rec.imchannel_joinedStatus)
		}
	}