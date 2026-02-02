// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (index)
	{
	switch (index)
		{
	case SM.CXMAXIMIZED:
		return 1930 /*=default screen width*/
	case SM.CYMAXIMIZED:
		return 1080 /*=default screen height*/
	case SM.CYEDGE:
		return 2
	case SM.CXVSCROLL:
		return 16
	case SM.CXHSCROLL:
		return 16
		}
	}
