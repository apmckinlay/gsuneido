// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	(mapSource)
	)
Contributions:
	(
	// after we can make all maps work with lat long,
	// plugin doesn't need allowLatLong option,
	// and MapButtonControl doesn't need onlyLatLong? argument
	('Map', 'mapSource', name: 'Google', func: MapGoogle, multiLoc:, allowLatLong:)
	('Map', 'mapSource', name: 'MapQuest', func: MapQuest)
	('Map', 'mapSource', name: 'Bing', func: MapBing, allowLatLong:)
	)
)
