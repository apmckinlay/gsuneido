// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Show Members"

	New(num)
		{
		.listMembers = .FindControl('channelMembers')
		.listMembers.SetReadOnly(true)
		.buildList(num)
		}

	Controls()
		{
		return Object('List',
			name: 'channelMembers',
			columns: #(user, bizuser_name),
			defaultColumns: #(user, bizuser_name),
			columnsSaveName: .Title,
			title: .Title,
			resetColumns:)
		}

	buildList(num)
		{
		QueryApply('user_settings where key is "IM_SubscribedChannels"')
			{ |rec|
			if rec.value.Empty?()
				continue

			if rec.value.Has?(num)
				QueryApply1('users', user: rec.user)
					{ |user|
					.listMembers.AddRow(
						Object(user: user.user, bizuser_name: user.bizuser_name))
					}
			}
		}

	GetMembersFromChannelNum(num)
		{
		if num is #() or num is '' or num is '""'
			return #()

		obj = Object()
		return_obj = Object()
		QueryApply('user_settings where key is "IM_SubscribedChannels"')
			{ |rec|
			if rec.value.Empty?()
				continue

			if rec.value.Has?(num)
				QueryApply1('users', user: rec.user)
					{ |user|
					obj.Add(Object(:user, bizuser_name: user.bizuser_name))
					}
			}
		for (i = 0; i < obj.Size(); i++)
			{
			return_obj.Add(obj[i].user.user)
			}
		return return_obj
		}
	}
