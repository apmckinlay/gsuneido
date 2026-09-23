// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin:  350
	Ymin:  350
	Title: #Preferences
	CallClass()
		{
		IDESettings.Ensure()
		OkCancel([this], .Title)
		}

	New()
		{
		super([#Vert, .controls(), name: #settings])
		}

	controls()
		{
		settings = [#Tabs,
			selectedTabColor: IDESettings.Get(#ide_selected_tab_color, false),
			selectedTabBold: IDESettings.Get(#ide_selected_tab_bold, true)]
		.tabs().Each()
			{|tab|
			accordion = [#Accordion, Tab: tab.Extract(#tabName)]
			tab.Members().Sort!().Each()
				{|accordionName|
				accordion.Add([accordionName, .buildCtrls(tab[accordionName])])
				}
			settings.Add(accordion)
			}
		return settings
		}

	// Level 1 Values 	= 	tabName, and accordion controls
	// Level 2 Member 	= 	accordion name
	// Level 2 Values 	= 	control or controls
	// Level 3 Member 	= 	control static text
	// Level 3 Value 	= 	control
	//		Control		=	Added as is
	// 		Controls	=	Wrapped in a From and joined seperated by nl
	tabs()
		{
		return [.cosmeticTab(), .developerTab()]
		}

	cosmeticTab()
		{
		return Object(tabName: #Cosmetic,
			"Customize Font": Object(
				Controls: [#Font, noTitle?:, name: #ide_logfont,
					subHeading: "Font for control static text",
					fontFaceName: StdFonts.Ui(),
					initFont: IDESettings.Get(#ide_logfont, false)],
				Editor: [#Font, noTitle?:, name: #ide_scifont,
					subHeading: "Font for editor text controls",
					fontHeight: ScintillaControl.DefaultFontSize,
					initFont: IDESettings.Get(#ide_scifont, false)]),
			"Color Scheme": #ColorScheme)
		}

	developerTab()
		{
		return Object(tabName: #Developer,
			Tabs: Object(
				"Persistent Save Limit": #(IDESettingSpinner, ide_tab_persistent_limit),
				"Keep Active Tab Left": #(IDESettingToggle, ide_move_tab, defaultVal:),
				"Scroll Tabs": #(IDESettingToggle, ide_scroll_tabs, defaultVal:),
				"Selected Tab Bold": #(IDESettingToggle, ide_selected_tab_bold,
					defaultVal:),
				"Selected Tab Color": [#ChooseColor,
					name: #ide_selected_tab_color,
					color: IDESettings.Get(#ide_selected_tab_color, CLR.SpotBlue)]),
			Scintilla: #(
				"Show Line Numbers": (IDESettingToggle, ide_show_line_numbers),
				"Show Whitespace": (IDESettingToggle, ide_show_whitespace),
				"Show Folding Margin": (IDESettingToggle, ide_show_fold_margin,
					defaultVal:),
				"Show Annotations": (IDESettingToggle, ide_show_annotations)
				),
			BookEdit: #("Auto Refresh": (IDESettingToggle, ide_book_auto_refresh,
					defaultVal:)))
		}

	buildCtrls(layout)
		{
		ctrls = Object()
		if Object?(layout)
			{
			ctrls.Add(#Form, overlap: false)
			for mem in layout.Members().Sort!()
				ctrls.Add(
					[#StaticText, mem $ ':', group: 0],
					[#Vert, layout[mem], group: 1],
					#nl)
			}
		else
			ctrls = layout
		return ctrls
		}

	OK()
		{
		postSaveCtrls = Object()
		resetCache? = false
		for col in IDESettings.Columns
			{
			// IDE Settings which don't have controls added yet. Allows for pre-staging
			if false is ctrl = .FindControl(col)
				continue
			if IDESettings.Get(col) isnt value = ctrl.Get()
				{
				resetCache? = true
				IDESettings.Set(col, value)
				if ctrl.Method?(#PostSave)
					postSaveCtrls.Add(ctrl)
				}
			}
		if resetCache?
			IDESettings.ResetCache()
		for ctrl in postSaveCtrls
			ctrl.PostSave()
		PubSub.Publish(#Redir_SendToEditors, #SyncPreferences, init:)
		return true
		}
	}
