// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(hwnd = 0)
		{
		data = .load()
		if data.Empty?()
			{
			Alert("No AI conversation files found in .ai/")
			return false
			}
		return Dialog(hwnd, [this, data],
			keep_size: #AiAgentLoadControl, closeButton?:, title: "Load Conversation")
		}

	New(data)
		{
		super(.layout(data))
		}

	Startup()
		{
		.Vert.List.ScrollToBottom()
		}

	layout(data)
		{
		return [#Vert,
			[#ListStretch, #(date_time, prompt), data, columnsSaveName: #AiAgentLoad],
			#Skip,
			#(Horz, (Button, Browse), Fill, (Button, Load), Skip, (Button, Cancel))]
		}

	On_Load()
		{
		selected = .Vert.List.GetSelection()
		if selected.Empty?()
			{
			.AlertInfo("Load Conversation", "Please select a conversation to load")
			return
			}
		.Window.Result(.Vert.List.GetRow(selected[0]).path)
		}

	List_DoubleClick(row, col/*unused*/)
		{
		.Window.Result(.Vert.List.GetRow(row).path)
		}

	On_Cancel()
		{
		.Window.Result(false)
		}

	On_Browse()
		{
		filename = OpenFileName(filter: "Log Files (*.md)|*.md|All Files (*.*)|*.*")
		if filename isnt ""
			.Window.Result(filename)
		}

	load()
		{
		result = Object()
		base = Paths.Combine(GetCurrentDirectory(), ".ai")
		Dir(`.ai/ai*.md`, details:)
			{ |file|
			prompt = ""
			path = Paths.ToLocal(Paths.Combine(base, file.name))
			File(path)
				{ |f|
				while false isnt line = f.Readline()
					if line is `## {{ User }}`
						{
						f.Readline()
						while false isnt line = f.Readline()
							{
							if line.Prefix?(`## `)
								break
							prompt $= line $ '\n'
							}
						prompt = prompt.Trim().Tr('\n', '\t')
						break
						}
				}
			result.Add([:prompt, date_time: file.date, :path])
			}
		return result
		}
	}
