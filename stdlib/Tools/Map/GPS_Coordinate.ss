// Copyright (C) 2014 Axon Development Corporation All rights reserved worldwide.
class
	{
	lat: false
	long: false

	// expects format like "49.454629N,112.784675W" or "49.454629,-112.784675"
	New(latLongStr)
		{
		if not String?(latLongStr)
			return

		ob = latLongStr.Split(",").Map!(#Trim)
		if ob.Size() isnt 2
			return

		.lat = .checkMaxValue(.convertToNumber(ob[0]), 90 /*= maximum latitude*/)
		.long = .checkMaxValue(.convertToNumber(ob[1]), 180 /*= maximum longitude */)
		}

	checkMaxValue(val, max)
		{
		return val isnt false and Abs(val) <= max ? val : false
		}

	convertToNumber(latOrLong)
		{
		directionChar = latOrLong[-1]
		if not directionChar.Alpha?()
			return latOrLong.Number?() ? Number(latOrLong) : false

		coord = latOrLong[.. -1]
		if not coord.Number?()
			return false

		sign = directionChar is 'W' or directionChar is 'S' ? -1 : 1
		return sign * Number(coord)
		}

	Latitude()
		{
		return .lat
		}

	Longitude()
		{
		return .long
		}

	Valid?()
		{
		return .lat isnt false and .long isnt false
		}

	// using haversine formula (http://www.movable-type.co.uk/scripts/latlong.html)
	// is not completely accurate since earth is not a perfect sphere
	// Distance is returned in km
	DistanceFrom(gpsCoord)
		{
		if not .Valid?() or not gpsCoord.Valid?()
			return false

		earthsRadiusInKM = 6371
		dLat = .degreeToRadians(gpsCoord.Latitude() - .Latitude())
		dLon = .degreeToRadians(gpsCoord.Longitude() - .Longitude())
		a =	(dLat/2).Sin() * (dLat/2).Sin() +
			(.degreeToRadians(.Latitude())).Cos() *
			(.degreeToRadians(gpsCoord.Latitude())).Cos() *
			(dLon/2).Sin() * (dLon/2).Sin()
		c = 2 * .atan2(a.Sqrt(), (1-a).Sqrt())
		return earthsRadiusInKM * c  // km
		}

	degreeToRadians(deg)
		{
		return deg * (PI / 180)
		}

	atan2(y, x)
		{
		return 2 * ( y / ( (x * x + y * y).Sqrt() + x ) ).ATan()
		}

	ToString()
		{
		return .Latitude().Format('-###.######') $ ',' $
			.Longitude().Format('-###.######')
		}

	DMS_Latitude(round = 0)
		{
		return .convertToDMS(.lat, round)
		}

	DMS_Longitude(round = 0)
		{
		return .convertToDMS(.long, round)
		}

	convertToDMS(latOrLong, round)
		{
		if latOrLong is false
			return false
		degrees = String(latOrLong.Int())
		minutesBase = Abs(latOrLong.Frac())*60
		minutes = String(minutesBase.Int()).LeftFill(2, '0')
		seconds = String((minutesBase.Frac()*60).Round(round)).LeftFill(2, '0')
		return Object(:degrees, :minutes, :seconds)
		}
	}