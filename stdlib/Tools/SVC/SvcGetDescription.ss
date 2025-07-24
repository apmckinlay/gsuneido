// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Send Local Changes"
	CallClass(hwnd, list, message)
		{
		return OkCancel(Object(this, list, message), .Title, hwnd)
		}
	New(list, message)
		{
		super(.makecontrols(message))
		.Vert.list.Set(.buildChangeList(list))
		.setWarningMsg(message)
		.comment = .Vert.desc
		recentBtn = .FindControl('recent')
		if Suneido.GetDefault('svc_comments', #()).Empty?()
			recentBtn.SetEnabled(false)
		else
			{
			.comment.Set(.getComment(Suneido.svc_comments.Last()))
			recentBtn.SetMenu(Suneido.svc_comments.
				Map({ it.Ellipsis(100, atEnd:) }).Reverse!()) /*= character limit */
			}
		.comment.SelectAll()
		.comment.SetFocus()
		}
	makecontrols(message)
		{
		controls = ['Vert'
			#(Editor readonly:, height: 5, name: list,
				font: "@mono" size: 10)
			'Skip'
			#(Static 'Description of Changes')
			#(ScintillaAddons wrap:, margin: 0, height: 7, name: desc,
				Addon_speller: #(ignore: (refactor, refactored, refactoring)))
			#(Skip 3),
			#(Horz Fill (Static
				'Please say "WHY" not just "what", the diff will show us "what"')
				Fill)
			'Skip'
			#(Horz
				(MenuButton 'Recent', name: 'recent', command: 'Recent')
				Skip
				(MenuButton 'Standard' (cosmetic, renamed, 'minor refactor'))
				(Skip 50)
				Fill
				'OkCancel')
			xmin: 300
			]
		if message isnt ""
			controls.Add(
				['ScintillaAddonsEditor', fontSize: '+2', height: 5, readonly:, wrap:,
					name: 'msg' ]
				#(Skip 6), at: 1)
		return controls
		}
	buildChangeList(list)
		{
		return list.Map({ it.type.Tr(" ") $ it.lib $ ":" $ it.name }).Join('\r\n')
		}
	setWarningMsg(message)
		{
		if false isnt msg = .FindControl('msg')
			{
			msg.DefineStyle(0, CLR.RED, bold:)
			msg.Set(message)
			}
		}
	On_Recent(item)
		{
		.comment.Set(.getComment(item))
		}
	On_Standard(item)
		{
		.comment.Set(item)
		}
	getComment(item)
		{
		return Suneido.svc_comments.
			FindOne({ it.Prefix?(item.BeforeFirst(' -')) }).AfterFirst('- ')
		}
	On_OK()
		{
		.Send('On_OK')
		}
	OK()
		{
		.save_comment()
		if '' is desc = .comment.Get().Trim()
			{
			.AlertError('Svc Get Description', "Please enter a description")
			return false
			}
		return desc
		}
	Cancel()
		{
		.save_comment()
		return false
		}
	save_comment()
		{
		if '' is comment = .comment.Get().Trim()
			return

		svcComments = Suneido.GetInit('svc_comments', { Object() })

		// remove and add to bring to front
		svcComments.RemoveIf({ it.AfterFirst('- ') is comment })
		svcComments.Add(Date().Format('HH:mm:ss') $ ' - ' $ comment)

		if svcComments.Size() > 10 /*= max number of comments saved */
			svcComments.Delete(0)
		}
	}
