// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Subject Template"
	Xmin: 600
	CallClass(hwnd, emailSubject)
		{
		rptName = emailSubject.type
		data = emailSubject.data
		cols = emailSubject.cols

		return OkCancel(Object(this, data, rptName, cols), .Title, hwnd)
		}

	New(.data, .rptName, .cols)
		{
		super(.layout())
		.subjectCtrl = .FindControl('sci_ctrl')
		.previewCtrl = .FindControl('preview')
		.colsList = .FindControl('colsList')
		if .readonly
			.FindControl('Add_Field_to_Subject').SetEnabled(false)
		if false isnt defaultSubjectTemplate = .GetSavedSubject(.rptName)
			.subjectCtrl.Set(defaultSubjectTemplate)
		}

	layout()
		{
		.readonly = not .HasPermission?()
		field_dicts = .cols
		list = field_dicts.Map(SelectFields.GetFieldPrompt)
		return Object('Vert',
			Object('Title', .Title $ ' - ' $ .rptName),
			'Skip',
			Object('HorzSplit',
				Object('ListBox', list, sort:, name: 'colsList'),
				Object('Vert',
					#(Static 'Subject:'),
					.editor(name: 'sci_ctrl', readonly: .readonly,
						Addon_highlight_customFields: [:field_dicts, useMarkdown:]),
					'Skip',
					#(Static 'Preview:'),
					.editor(readonly:, name: 'preview'))),
			'Skip',
			#(Horz (Button, 'Add Field to Subject', command: 'AddField'),
				Fill (Static 'Subject Template will apply to all users'))
			)
		}

	HasPermission?()
		{
		return OptContribution('Email_DefaultSubject_HasPermission?', { false })()
		}

	editor(@args)
		{
		return Object(.editorCtrl, Addon_detab:, height: 5, xstretch: 10,
			ystretch: 1, margin: 0).Merge(args)
		}

	editorCtrl: ScintillaAddonsEditorControl
		{
		Addons: #()
		}

	ListBoxDoubleClick(unused)
		{
		.On_AddField()
		}

	EditorChange(source)
		{
		if source.Name is 'sci_ctrl'
			.previewCtrl.Set(.BuildSubject(source.Get(), .cols, .data))
		}

	On_AddField()
		{
		if .readonly
			return
		.subjectCtrl.Paste('<' $ .colsList.Get() $ '> ')
		}

	OK()
		{
		if .readonly
			{
			.On_Cancel()
			return false
			}

		report = .rptName $ ' ~ default email subject'
		subject = .subjectCtrl.Get().Trim()
		if subject.Has?('\r') or subject.Has?('\n')
			{
			.AlertInfo(.Title, 'Subject should be a single line')
			return false
			}

		QueryDo('delete params where report is ' $ Display(report))
		if subject isnt ''
			QueryOutput('params', [:report , params: subject])
		else
			{
			.AlertInfo(.Title, 'Please enter a Subject')
			return false
			}

		translatedSubject = .BuildSubject(.subjectCtrl.Get(), .cols, .data)
		return translatedSubject
		}

	BuildSubject(subject, cols, data)
		{
		formatterFunc = OptContribution('Email_DefaultSubjectFormatter', FormatValue)
		for col in cols
			{
			fieldPrompt = SelectFields.GetFieldPrompt(col)
			subject = subject.Replace("(?q)<" $ fieldPrompt $ ">", {|unused|
				value = data[col]
				Object?(value) ? value.Join(', ') : formatterFunc(value, col)
				})
			}

		subject = subject.Trim()
		return subject.Tr('\r\n', ' ')
		}

	GetSavedSubject(rptName)
		{
		report = rptName $ ' ~ default email subject'
		return false isnt (p = Query1('params', :report)) ? p.params : false
		}

	AddEmailSubjectInfo(rpt, type, data, cols)
		{
		if rpt.Params.GetDefault('NoEmailSubject', false) is true
			return

		if rpt.Params.Member?('EmailSubject')
			{
			if rpt.Params.EmailSubject.data is data // keep if same data
				return
			rpt.Params.Delete('EmailSubject')
			rpt.Params.NoEmailSubject = true
			return
			}

		emailSubject = rpt.Params.EmailSubject = Object()
		emailSubject.type = type
		emailSubject.data = data
		emailSubject.cols = .filterCols(cols)
		}
	filterCols(cols)
		{
		return cols.Copy().RemoveIf(Internal?).RemoveIf({
			Datadict(it).Control[0] =~ 'ScintillaRichWordAddons' })
		}
	// disable alt-f4 instead of using OkCancel wrapper because we have custom layout
	// for position of OkCancel button
	ConfirmDestroy()
		{
		return false
		}
	}