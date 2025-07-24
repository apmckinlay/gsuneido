// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(printMode, reporterMode = 'simple')
		{
		ob = Object('Vert', #(Skip 4))
		ob.Add(Object('HorzEqual',
			LastContribution('Reporter').TopLeft_Extra,
			Object(#Horz, #(Skip 5) Object('Static', '', size: '+2', weight: 'bold',
				color: CLR.amber name: 'schedWarn'))
			'Fill', 'Skip', #(Button 'New')
			'Skip', #(Button 'Open...')
			'Skip', #(Button 'Save')
			'Skip', #(Button 'Save As...')
			'Skip', #(Button 'Help'), pad: 5))
		ob.Add(#(Skip 4))
		ob.Add(Object('Tabs'
				Object('Scroll'
					Object('Border', .input_tab()), noEdge:,
					Tab: 'Input')
				Object('Border', .formulas_tab(),
					Tab: 'Formulas'),
				.enumerateDesignTab(reporterMode),
				Object('Scroll'
					Object('Border', .sort_tab(reporterMode)), noEdge:,
					Tab: 'Sort')
				Object('Scroll'
					Object('Border', .select_tab(reporterMode)), noEdge:,
					Tab: 'Select')
				constructAll:
				))
		ob.Add(#(Skip 4))
		.add_extra_buttons(ob, reporterMode)

		if printMode is true
			ob.Add(#(Skip 4))

		return Object('Record', ob)
		}

	enumerateDesignTab(reporterMode)
		{
		tab = reporterMode is 'form'
			? .form_tab()
			: reporterMode is 'enhanced'
				? .design_tab()
				: .old_design_tab()

		return Object('Scroll', Object('Border', tab), noEdge:, Tab: 'Design')
		}

	add_extra_buttons(ob, reporterMode)
		{
		horz = Object('HorzEqual'
			#Skip
			#(Button 'Print')
			#Skip
			#(Button 'Preview')
			#Skip
			#PDFButton
			#Skip
			#(Button 'Page Setup')
			#Skip)
		if reporterMode isnt #form
			horz.Add(#(Button 'Export'),
				#Skip,
				#(CheckBox 'Print Lines' name: print_lines))

		horz.Add(
			#Fill #Skip
			#(LinkButton 'View Data')
			#Fill #Skip
			#(LinkButton 'Export Report Design')
			#Skip #Skip
			#(LinkButton 'Import Report Design')
			#Skip
			)
		horz.pad = 5
		ob.Add(horz)
		}

	input_tab()
		{
		layout = Object('Vert',
			Object('Horz'
				Object('Vert'
					#(Heading 'Data Source'),
					#(Skip 4),
					ReporterDataSource),
				#(Fill)
				Object('Static', '', name: 'reportHistory', size: '-1', justify: 'RIGHT'))
			'Skip',
			#(Heading 'Description') #(Skip 4) #(Editor, xstretch: 0,
			ystretch: 0, name: 'report_description'),
			'Skip',
			.summarizeAccordion())
		return layout
		}

	summarizeAccordion()
		{
		layout = Object('Vert'
			#(Static size: '+2'
				'Only use this section if you want to summarize the data.')
			#(Skip 8),
			#(Pair (Static 'By') (ChooseTwoList listField: reporter_summarizeby_cols,
				name: summarize_by, width: 30)))
		for i in .. Reporter.MaxSummarizeFields
			layout.Add(#(Skip 4), Object('Horz'
				Object('Pair' #(Static 'Function'),
					Object('ChooseList' #(total, maximum, minimum, average, count)
						name: 'summarize_func' $ i))
				'Skip'
				Object('Pair' #(Static 'Field'),
					Object('ChooseList' listField: 'reporter_summarizeby_cols',
						name: 'summarize_field' $ i))))
		return Object('Accordion' Object('Summarize', layout))
		}

	formulas_tab()
		{
		formulaLayout = Object('Horz',
				Object('Vert'
					Object('Field', name: 'calc'),
					Object('Horz', 'Fill',
						Object('ChooseFieldType', name: 'type', reporter:),
						xstretch: 0))
				Object(.formula_editor, name: 'formula', font: '@mono')
				#(Skip 30)
				Object('CheckBox', name: 'form_val')
				#(Skip 30)
				)

		return Object('Vert'
			Object('Horz'
				#(Static 'Field')
				#(Fill .7)
				#(Static 'Formula')
				#(Skip 30),
				FormulaEditor.AddFormulaButtons(),
				#Fill
				#(Static 'Menu Option')
				#(Skip 20))
			#(Skip 4)
			Object('Scroll', Object('Repeat', formulaLayout, minusAtFront:,
				plusAtBottom:, maxRecords: Reporter.Ncalcs, name: 'formulas',
				saveKey:))
			#(Skip 3)
			Object('Static',
				'Choose a Field Name that is different from the existing ones.   ' $
				'Number Formulas can use: + - * /     Text Formulas can use: $')
			#(Skip 3)
			Object('Static',
				'For example:   Bonus Commission			Amount * .075\n' $
				'		     Number, 2 decimals'
				xstretch: 1)
			)
		}

	formula_editor: FormulaEditorControl
		{
		Xstretch: 1
		Ystretch: 0
		Height: 3

		EN_KILLFOCUS()
			{
			idx = .Send('FindRowByRecordControl', .Controller)
			.Send('FormulaKillFocus', idx)
			return super.EN_KILLFOCUS()
			}
		}

	old_design_tab()
		{
		return Object('Vert'
			Object('WndPane' Object('Border'
				Object('Vert'
					Object('Horz'
						Object('Static' Date().ShortDateTime(), xmin: 100,
							whitebgnd:)
						'Fill'
						#(Field, name: heading1, width: 35, justify: 'CENTER',
							size: '+2', weight: 'bold')
						'Fill'
						#(Static 'Page #', xmin: 100, justify: 'RIGHT',
							whitebgnd:))
					#(Horz
						Fill
						(Field, name: heading2, width: 35, justify: 'CENTER',
							size: '+2', weight: 'bold')
						Fill)
					#(Horz
						Fill
						(Editor, xmin: 550, height: 3, xstretch: 0, ystretch: 0,
							name: header)
						Fill)
					'Skip'
					'ReporterColumns'
					'Fill'
					#(Horz
						Fill
						(Editor, xmin: 550, height: 3, xstretch: 0, ystretch: 0,
							name: footer)
						Fill)
					)))
			#(Skip 6)
			#(Horz
				(Button 'Add/Remove Columns...' pad: 30)
				(Skip 20)
				(Static 'Drag columns to resize or rearrange.')
				(Skip 20)
				(Static 'Click on a column to change the heading or to total.')
				)
			)
		}

	design_tab()
		{
		return Object('Vert'
				Object('MultiCanvas', name: 'reporterCanvas')
			)
		}

	form_tab()
		{
		return #('Vert', #('ReporterCanvasColumns' name: 'ReporterColumns')
				#('ReporterCanvas', name: 'reporterCanvas')
			)
		}

	sort_form_grid()
		{
		rows = Object()
		rows.Add(#( // Horz's improve spacing
			(Horz (Skip 3) (Static 'Field'))))
		for i in .. Reporter.SortRows
			rows.Add(Object(Object('ChooseList', listField: 'reporter_sortcolumns',
				name: 'sort' $ i)))
		return rows
		}

	sort_grid()
		{
		rows = Object()
		rows.Add(#( // Horz's improve spacing
			(Horz)
			(Horz Skip)
			(Horz (Static 'Enable'))
			(Horz Skip)
			(Horz (Static 'Show'))
			(Horz Skip)
			(Horz (Static 'New Page'))
			))
		rows.Add(#( // Horz's improve spacing
			(Horz (Skip 3) (Static 'Field'))
			(Horz Skip)
			(Horz (Static 'Summaries'))
			(Horz Skip)
			(Horz (Static 'Before'))
			(Horz Skip)
			(Horz (Static 'For Each'))
			))
		for i in .. Reporter.SortRows
			rows.Add(Object(
				Object('ChooseList', listField: 'reporter_sortcolumns', name: 'sort' $ i),
				'Skip',
				Object('CheckBox', name: 'total' $ i)
				'Skip',
				Object('CheckBox', name: 'show' $ i)
				'Skip',
				Object('CheckBox', name: 'page' $ i)
				))
		return rows
		}

	sort_tab(reporterMode)
		{
		rows = reporterMode is #form ? .sort_form_grid() : .sort_grid()
		return Object('Vert'
			Object('Grid', rows)
			'Skip'
			#(CheckBox 'Reverse Order', name: reverse)
			'Skip'
			#(Static 'Select the fields you want to order your report by.')
			'Skip'
			#(Static 'Note: Totals must be set up on the Design tab.')
			)
		}

	select_tab(reporterMode)
		{
		// need extra vert to limit stretch
		return Object('Vert', Object('Vert',
			#(Skip 4)
			Object('Select2', SelectFields(), Record(),
				printParams: reporterMode isnt #form, menuOptions:)
			#(Skip 20)
			#(Button 'Clear Select' xmin: 80)
//			'Skip'
//			#(Heading 'Report Interface Options')
//			#(ChooseTwoList listField: reporter_summarizeby_cols,
//				name: report_filters, width: 40)
			xstretch: false))
		}
	}
