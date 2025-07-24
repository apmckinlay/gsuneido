// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "About Suneido"
	CallClass()
		{
		return ToolDialog(0, this, border: 0)
		}
	New()
		{
		super(.makecontrols())
		.yscroll = 0
		.scrolldir = -3
		.timer = Delay(2.SecondsInMs(), .autoscroll)
		}
	makecontrols()
		{
		try
			suneido = Object('Image',
				Query1("suneidoc", name: 'suneido.gif').text)
		catch
			suneido = #(Horz
				(Static Suneido size: '+39' color: 255 weight: 'bold')
				(Static 'TM' size: '-1'))
		return Object('Vert'
				Object('Border',
				.layout(suneido)
				border: 20) xstretch: 0)
		}
	layout(suneido)
		{
		return Object('Vert'
			Object('Horz' 'Fill' suneido 'Fill'),
			#(Skip 5)
			Object('Horz' 'Fill'
				Object('Static' 'Integrated Application Platform' size: '+5') 'Fill')
			#(Skip 5)
			Object('Horz' 'Fill'
				Object('Static' 'Running from ' $ ExePath()) 'Fill')
			Object('Horz' 'Fill' Object('Static' 'Built: ' $ Built()) 'Fill')
			#(Skip 5)
			#(Horz Fill
				(Static 'Licensed under the GNU General Public License Version 2') Fill)
			Object('Horz' 'Fill'
				Object('Static'
					'Copyright \xa9 2000-' $ Date().Year() $ ' Suneido Software Corp.')
				'Fill')
			#(Horz Fill (Static 'All rights reserved worldwide.') Fill)
			#(Horz Fill (WebLink, "suneido.com", "suneido.com") Fill)
			#(Skip 5)
			#(Horz Fill (Static 'Scintilla source code editor') Fill)
			#(Horz Fill (WebLink, "scintilla.org", "scintilla.org") Fill)
			#(Skip 5)
			#(Horz Fill (Static 'Hunspell spelling checker') Fill)
			#(Horz Fill
				(WebLink, "hunspell.sourceforge.net", "hunspell.sourceforge.net/") Fill)
			#(Skip 5)).
			Add(.contributors())
		}
	contribs: #(
		Oliver_Ackermann
		Petr_Antos
		Oliver_Ackermann
		Petr_Antos
		Luis_Alfredo
		Tracy_Arams
		'Roberto_Artigas Jr'
		Paul_Blankenfeld
		Erik_Braaten
		'Gerardo Antonio Garza_Casso'
		'Jean-Luc_Chervais'
		Maxwell_Correya
		Colin_Coller
		Randy_Coulman
		Jeremy_Cowgar
		Tyler_Davidson
		_DeusTech
		Kim_Dong
		'Domingo Alvarez_Duarte'
		Jason_Elias
		'Helmut_Enck-Radana'
		Jeff_Ferguson
		Mark_Gabor
		Mauro_Giubileo
		Tony_Hallett
		Steve_Heyns
		Jennie_Hill
		'Jean-Charles_Hoarau'
		'Cor_de Jong'
		Kevin_Kinnell
		Martin_Ledoux
		'Bj\xf6rn_Lietz-Spendig'
		Mal_Malakov
		Claudio_Mascioni
		Andrew_McKinlay
		Valerio_Muzi
		Santiago_Ottonello
		Tomas_Polak
		Andrew_Price
		Francek_Prijatelj
		Li_Qian
		'Arne Christian_Riis'
		'Jos_van Roosmalen'
		Jssi_Salmela
		Johan_Samyn
		Jos_Schaars
		Victor_Schappert
		Stefan_Schmiedl
		Mateus_Vendramini
		Hao_Xie)
	contributors()
		{
		firstnames = ''
		lastnames = ''
		for name in .contribs
			{
			firstnames $= name.BeforeFirst('_') $ '\n'
			lastnames $= name.AfterFirst('_') $ '\n'
			}
		return Object('GroupBox' 'Contributors'
				Object('Vert'
					Object('Scroll',
						Object('Horz'
							'Fill'
							Object('Static' firstnames justify: 'RIGHT')
							'Skip'
							Object('Static' lastnames)
							'Fill')
						noEdge:, dyscroll: 1, ymin: 100)
					#(Skip 5))
			)
		}
	autoscroll()
		{
		if .Destroyed?()
			return
		c = .Vert.Border.Vert.GroupBox.Vert.Scroll
		if (.timer isnt false and .yscroll is c.GetYscroll())
			{
			c.Scroll(0, .scrolldir)
			if (.yscroll is c.GetYscroll()) // scrolling up
				{
				.scrolldir *= -1
				.timer = Delay(2.SecondsInMs(), .autoscroll)
				}
			else // scrolling down
				{
				.yscroll -= .scrolldir
				.timer = Delay(100, .autoscroll)  /*= scrolling down delay*/
				}
			}
		else
			.timer = false
		}
	timer: false
	Destroy()
		{
		if .timer isnt false
			.timer.Kill()
		.timer = false
		super.Destroy()
		}
	}
