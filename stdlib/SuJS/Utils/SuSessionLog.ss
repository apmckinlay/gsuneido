// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Singleton
	{
	New()
		{
		.Ensure()
		}

	num: false
	Login(user, ip, userAgent)
		{
		Assert(.num is: false)
		.num = Timestamp()
		QueryOutput('su_session_logs', [
			susl_num: .num,
			:user,
			susl_ip: ip,
			susl_user_agent: userAgent])
		}

	Logout()
		{
		if .num isnt false
			QueryDo('update su_session_logs
				where susl_num is ' $ Display(.num) $ '
				set susl_logout_time = ' $ Display(Date()))
		}

	Error(err, reconnected? = false)
		{
		Assert(.num isnt: false)
		QueryOutput('su_session_errors', [
			suse_num: Timestamp(),
			susl_num: .num,
			suse_error: err,
			suse_reconnected: reconnected?])
		}

	desc: 'Web logins (last 30 days)'
	Stats()
		{
		stat = Object()
		totalLogin = 0
		totalMinutes = 0
		totalReconnects = 0
		totalDisconnects = 0
		days = Object()

		if TableExists?('su_session_logs')
			{
			start = Date().Plus(days: -30).NoTime()
			QueryApply('su_session_logs where susl_num >= ' $ Display(start))
				{
				totalLogin++
				days[it.susl_num.NoTime()] = true
				if Date?(it.susl_logout_time)
					totalMinutes += it.susl_logout_time.MinusMinutes(it.susl_num)
				}

			QueryApply('su_session_errors
				where susl_num >= ' $ Display(start) $ '
				summarize suse_reconnected, count')
				{
				if it.suse_reconnected is true
					totalReconnects = it.count
				else
					totalDisconnects = it.count
				}
			}

		stat['Web logins (last 30 days)'] = totalLogin
		stat['Average logins per day'] = (totalLogin / days.Size()).Round(2)
		stat['Average minutes per login'] = (totalMinutes / totalLogin).Round(0)
		stat['Average minutes per error'] =
			(totalMinutes / (totalReconnects + totalDisconnects)).Round(0)
		stat['Disconnect percentage'] = (totalDisconnects / totalLogin).DecimalToPercent()
		return stat
		}

	Ensure()
		{
		Database('ensure su_session_logs
			(susl_num, user, susl_ip, susl_user_agent, susl_logout_time)
			key (susl_num)')
		Database('ensure su_session_errors
			(suse_num, susl_num, suse_error, suse_reconnected)
			key (suse_num)
			index (susl_num) in su_session_logs cascade')
		}

	// override super.Reset to avoid getting cleared
	Reset() {}
	}
