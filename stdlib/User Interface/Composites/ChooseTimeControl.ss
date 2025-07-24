// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Name: ChooseTime
	New(.mandatory = false, width = 5, hidden = false, tabover = false, readonly = false)
		{
		super(Object('Field', :mandatory, :width), :hidden, :tabover, :readonly)
		.Send("Data")
		}
	Getter_DialogControl()
		{ return Object(.timecontrol, (.Valid?() ? .Field.Get() : '')) }

	timecontrol: Controller
		{
		Name: 'TimeList'
		New(time = 0)
			{
			super(.layout())
			.ctrls = .FindControl('controls')
			.Set(time)
			}

		layout()
			{
			return Object('Record'
				Object('Vert',
					Object('Horz',
						.hourList(),
						.minuteList(),
						name: 'controls')
				))
			}

		hourList()
			{
			list = Settings.Get('TimeFormat').Suffix?('tt')
				? #('12 am', '1 am', '2 am', '3 am', '4 am', '5 am', '6 am',
					'7 am', '8 am', '9 am','10 am','11 am','12 pm','1 pm','2 pm',
					'3 pm','4 pm','5 pm','6 pm','7 pm', '8 pm', '9 pm', '10 pm',
					'11 pm')
				: #("0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12",
					"13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23")

			return Object('ListBox' list name: 'Hour' xmin: 75 ymin: 160)
			}

		minuteList()
			{
			return #(ListBox #('00', '05', '10', '15', '20', '25', '30', '35', '40',
				'45', '50', '55')
				name: 'Min' xmin: 50 ymin: 160)
			}

		OK()
			{
			return .Get()
			}
		ListBoxDoubleClick(sel/*unused*/)
			{
			.Send('On_OK')
			}

		Get()
			{
			hour = .ctrls.Hour.GetText(.ctrls.Hour.GetSelected()).Split(' ')
			if hour is #()
				return false

			suffix = hour.Size() > 1 ? hour[1] : ''
			return hour[0] $ ':' $
				.ctrls.Min.GetText(.ctrls.Min.GetSelected()) $ ' ' $ suffix
			}
		roundMinutes: 5
		minutesInHour: 60
		round_to_five(minutes)
			{
			minute = (Number(minutes) / .roundMinutes).Round(0) * .roundMinutes
			if minute is .minutesInHour
				minute = 0
			return ("0" $ minute)[-2..]
			}
		Set(time)
			{
			results = .splitTime(time)
			.ctrls.Hour.SetCurSel(.ctrls.Hour.FindString(.hourSelect(results.hour,
				results.suffix,	Settings.Get('TimeFormat'))))
			.ctrls.Min.SetCurSel(.ctrls.Min.FindString(.round_to_five(results.minute)))
			}
		timeSplitRx: '(\d?\d):?(\d\d)\s?(am|AM|pm|PM)?'
		splitTime(time)
			{
			if false is results = time.ExtractAll(.timeSplitRx)
				return Object(hour: 0, minute: '', suffix: '')
			return Object(hour: Number(results.GetDefault(1 /*=hour*/, '')),
				minute: results.GetDefault(2 /*=minute*/, ''),
				suffix: results.GetDefault(3 /*=suffix*/, ''))
			}
		hourSelect(hour, suffix, timeFormat)
			{
			if not timeFormat.Suffix?('tt')
				return String(hour)
			if hour is 0
					{
					hour = 12
					suffix = 'am'
					}
				else if hour > 12 /* = hours */
					{
					hour = hour - 12 /* = hours */
					suffix = 'pm'
					}
			return Display(hour) $ ' ' $ suffix
			}
		}
	FieldKillFocus()
		{
		if not super.Dirty?() or not .Valid?()
			return
		time = .Get()
		if time is ''
			return
		s = .format(String(time))
		if s is time
			return
		SetWindowText(.Field.Hwnd, s)
		.Dirty?(true)
		}
	FieldReturn()
		{
		dirty? = .Dirty?()
		.FieldKillFocus()
		.Field.KillFocus()
		if (dirty?)
			.NewValue(.Get())
		}
	Set(time)
		{
		super.Set(.format(time))
		}
	format(time)
		{
		if time is ''
			return ''

		timeFmtLength = 4
		formatTime = String(time).LeftFill(timeFmtLength, '0')
		if (not .Valid?() or (false is (ob = .SplitTime(formatTime))))
			return time

		if ob.suffix is 'pm' and ob.hour isnt 12 /*=noon*/
			ob.hour += 12 /*=conventTo24Hour*/
		if ob.suffix is 'am' and ob.hour is 12 /*=midnight*/
			ob.hour = 0
		// set second and millisecond to zero for safety
		// Date() wants hour in 24hour format, don't convert it
		d = Date(hour: ob.hour, minute: ob.minute, second: 0, millisecond: 0)
		return d.Time()
		}
	SplitTime(time)
		{
		if time.Has?(':')
			{
			if false is hourMinute = .splitTimeOnColon(time)
				return false
			hour = hourMinute.hour
			minute = hourMinute.minute
			}
		else
			{
			if not time.Numeric?()
				return false
			hour = Number(time[.. -2])
			minute = Number(time[-2..])
			}
		suffix = .suffixFromTime(time)
		return Object(:hour, :minute, :suffix)
		}

	splitTimeOnColon(time)
		{
		hour = time.BeforeFirst(':')
		if not hour.Number?()
			return false
		hour = Number(hour)
		minute = time.AfterFirst(':').Lower()
		if minute.Has?('am') or minute.Has?('pm')
			minute = minute.Replace('am|pm', '').Trim()
		if not minute.Number?()
			return false
		minute = Number(minute)
		return Object(:hour, :minute)
		}

	suffixFromTime(time)
		{
		return time.Lower().Has?('am')
			? 'am'
			: time.Lower().Has?('pm')
				? 'pm' : ''
		}

	Valid?()
		{
		time = String(.Get())
		if false is .validCheck?(time, .mandatory)
			return false
		return super.Valid?()
		}
	validCheck?(time, mandatory)
		{
		if time is ''
			return mandatory isnt true

		if not time.Has?(':') and not time.Number?()
			return false

		if false is ob = .SplitTime(time)
			return false

		if ob.suffix is 'am' and ob.hour > 12 /* = hours */
			return false

		if not Time?(String(.ConvertToMilitary(ob)))
			return false

		return true
		}

	ValidData?(@args)
		{
		value = args[0]
		mandatory = args.GetDefault('mandatory', false)
		if value is ''
			return not mandatory
		return Time?(value)
		}

	NewValue(value /*unused*/)
		{ .Send('NewValue', .Get()) }
	Get()
		{
		time = .Field.Get()
		if time is ''
			return ''
		if false is ob = .SplitTime(time)
			return time
		return .ConvertToMilitary(ob)
		}
	ConvertToMilitary(ob)
		{
		hour = ob.hour
		if ob.suffix is 'pm' and hour < 12 /* = hours */
			hour += 12 /* = hours */
		if hour is 12 /* = hours */ and ob.suffix is 'am'
			hour = 0
		return Number(hour $ ob.minute.Pad(2, '0'))
		}
	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
