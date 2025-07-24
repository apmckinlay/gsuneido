// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(@args)
		{
		super(.controls(args))
		}
	controls(args)
		{
		data = args[0]
		if String?(data)
			data = args[0] = SparklineControl.GetData(args[0], args[1])
		red = CLR.RED
		blue = CLR.BLUE
		return ['Horz',
			['Static', data.First().Format('-###,###,###'), color: red],
			args.Add('Sparkline', at: 0),
			['Static', data.Last().Format('-###,###,###'), color: red],
			#(Skip 4),
			['Static', data.Min().Format('-###,###,###'), color: blue],
			#(Skip 4),
			['Static', data.Max().Format('-###,###,###'), color: blue]
			]
		}
	}