// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// Contributed by Jeremy Cowgar (jcowgar@bhsys.com)
Hwnd
	{
	Name: "ProgressBar"
	Xmin: 100
	Ymin: 20
	wordmax: 65535

	New(style = 0)
		{
		.CreateWindow( 'msctls_progress32', '', WS.VISIBLE | style)
		}

	SetRange(min, max)
		{
		if min > max or min < 0 or max > .wordmax
			return false
		.SendMessage(PBM.SETRANGE, 0, MAKELONG(min, max))
		}

	// Set position (if range is 0 -> 100, 50 would be right in the middle!
	SetPos(pos)
		{
		if pos < 0 or pos > .wordmax
			return false
		.SendMessage(PBM.SETPOS, pos)
		}

	// Add's onto the current position, for instance if the range is 0 -> 100
	// and the current position is 50,
	// if you were to call SetDeltaPos(25) the range would now be 75.
	SetDeltaPos(pos)
		{
		if pos < 0 or pos > .wordmax
			return false
		.SendMessage(PBM.DELTAPOS, pos)
		}

	// This is used to determine the step taken when calling the StepIt function.
	// For instance, if you SetStep(25), and your current pos is 50, calling StepIt()
	// would move the current pos to 75, the next, 100.
	SetStep(step)
		{
		.SendMessage(PBM.SETSTEP, step)
		}

	// Increment the pos by the value declared in SetStep. The intial setting is 10.
	StepIt()
		{
		.SendMessage(PBM.STEPIT)
		}

	GetReadOnly()			// read-only not applicable to progressbar
		{
		return true
		}
	}
