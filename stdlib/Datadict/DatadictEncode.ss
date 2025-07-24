// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// Arg 0: 	The field you are attempting to get the encode for
// Arg 1:	The value you are attempting to encode
// Arg > 1:	Additional parameters used by the encode
//				IE: Field_date.Encode('01122017', fmt: 'dMy') returns #20171201
//					DatadictEncode('date', '01122017', fmt: 'dMy') returns #20171201
function (@args) /*usage: (field, value, optionalParams)*/
	{
	return (Datadict(args[0], #(Encode)).Encode)(@+1args)
	}
