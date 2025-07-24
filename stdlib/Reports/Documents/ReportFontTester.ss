// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	Params.On_Preview(#(Vert
		#(Text 'This is a Helvetica')
		#(Text 'This is a Helvetica 20',
			font: (size: 20))
		#(Text 'This is a Helvetica 20 Italic',
			font: (size: 20, italic:))
		#(Text 'This is a Helvetica 20 Bold',
			font: (size: 20, weight: 'bold'))
		#(Text 'This is a Helvetica 20 Bold Italic',
			font: (size: 20, italic:, weight: 'bold'))
		#(Text 'This is a Courier',
			font: (name: 'Courier New'))
		#(Text 'This is a Courier 20',
			font: (name: 'Courier New', size: 20))
		#(Text 'This is a Courier 20 Italic',
			font: (name: 'Courier New', size: 20, italic:))
		#(Text 'This is a Courier 20 Bold',
			font: (name: 'Courier New', size: 20, weight: 'bold'))
		#(Text 'This is a Courier 20 Bold Italic',
			font: (name: 'Courier New', size: 20, italic:, weight: 'bold'))
		#(Text 'This is a Times New Roman',
			font: (name: 'Times New Roman'))
		#(Text 'This is a Times New Roman 20',
			font: (name: 'Times New Roman', size: 20))
		#(Text 'This is a Times New Roman 20 Italic',
			font: (name: 'Times New Roman', size: 20, italic:))
		#(Text 'This is a Times New Roman 20 Bold',
			font: (name: 'Times New Roman', size: 20, weight: 'bold'))
		#(Text 'This is a Times New Roman 20 Bold Italic',
			font: (name: 'Times New Roman', size: 20, italic:, weight: 'bold'))
		), previewWindow: GetFocus())
	}