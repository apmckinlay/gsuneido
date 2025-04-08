// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestOne(t *testing.T) {
	assert := assert.T(t).This
	test := func(year int, month int, day int,
		hour int, minute int, second int, millisecond int) {
		d := NewDate(year, month, day, hour, minute, second, millisecond)
		assert(d.Year()).Is(year)
		assert(d.Month()).Is(month)
		assert(d.Day()).Is(day)
		assert(d.Hour()).Is(hour)
		assert(d.Minute()).Is(minute)
		assert(d.Second()).Is(second)
		assert(d.Millisecond()).Is(millisecond)
	}
	test(2014, 01, 15, 12, 34, 56, 789)
	test(1900, 01, 01, 0, 0, 0, 0)
	test(2499, 12, 31, 23, 59, 59, 999)
}

func TestDateLiteral(t *testing.T) {
	assert := assert.T(t).This
	good := func(s string) {
		t.Helper()
		d := DateFromLiteral(s)
		s = "#" + s
		assert(d.String()).Is(s)
		d = DateFromLiteral(s)
		assert(d.String()).Is(s)
	}
	bad := func(s string) {
		t.Helper()
		assert(DateFromLiteral(s)).Is(NilDate)
		assert(DateFromLiteral("#" + s)).Is(NilDate)
	}
	good("20140115")
	good("19000101")
	good("24991231")
	good("20140115.1234")
	good("20140115.123456")
	good("20140115.123456789")
	bad("2014123")
	bad("20141231.1")
	bad("20140115.123")
	bad("20140115.12345")
	bad("20140115.12345678")
	bad("20140230")
	bad("20130229")
	good("20120229") // leap year
	bad("20160523.0-15")

	s := "#20230316.123456777222"
	assert(DateFromLiteral(s).String()).Is(s)
}

func TestDatePack(t *testing.T) {
	pack := func(s string) {
		t.Helper()
		d := DateFromLiteral(s)
		buf := Pack(d)
		assert.T(t).This(buf[0]).Is(byte(PackDate))
		d2 := Unpack(buf)
		assert.T(t).This(d2).Is(d)
	}
	pack("20140115")
	pack("19000101")
	pack("24991231")
	pack("20140115.1234")
	pack("20140115.123456")
	pack("20140115.123456789")
}

func TestDateCompare(t *testing.T) {
	assert := assert.T(t)
	dates := []string{"19990101", "20001231"}
	times := []string{"080102333", "180102333"}
	extra := []string{"", "123"}
	vals := [8]Value{}
	i := 0
	for _, d := range dates {
		for _, t := range times {
			for _, e := range extra {
				vals[i] = DateFromLiteral("#" + d + "." + t + e)
				var y = DateFromLiteral("#" + d + "." + t + e)
				assert.That(vals[i].Equal(y))
				i++
			}
		}
	}
	for i := range len(vals) {
		for j := range len(vals) {
			assert.This(vals[i].Compare(vals[j])).Is(cmp.Compare(i, j))
			assert.This(vals[j].Compare(vals[i])).Is(cmp.Compare(j, i))
			assert.This(vals[i].Equal(vals[j])).Is(i == j)
			assert.This(vals[j].Equal(vals[i])).Is(i == j)
		}
	}
}

func TestDatePlus(t *testing.T) {
	plus := func(s string, year int, month int, day int,
		hour int, minute int, second int, ms int, expected string) {
		d := DateFromLiteral(s).(SuDate)
		e := DateFromLiteral(expected).(SuDate)
		assert.T(t).
			This(d.Plus(year, month, day, hour, minute, second, ms)).Is(e)
	}
	//						   y  m  d  h  m  s  ms

	// no overflow
	plus("20140115.123456789", 0, 0, 0, 0, 0, 0, 0, "20140115.123456789")
	plus("20140115.123456789", 0, 0, 0, 0, 0, 0, 1, "20140115.123456790")
	plus("20140115.123456789", 0, 0, 0, 0, 0, 0, -1, "20140115.123456788")
	plus("20140115.123456789", 0, 0, 0, 0, 0, 1, 0, "20140115.123457789")
	plus("20140115.123456789", 0, 0, 0, 0, 0, -1, 0, "20140115.123455789")
	plus("20140115.123456789", 0, 0, 0, 0, 1, 0, 0, "20140115.123556789")
	plus("20140115.123456789", 0, 0, 0, 0, -1, 0, 0, "20140115.123356789")
	plus("20140115.123456789", 0, 0, 0, 1, 0, 0, 0, "20140115.133456789")
	plus("20140115.123456789", 0, 0, 0, -1, 0, 0, 0, "20140115.113456789")
	plus("20140115.123456789", 0, 0, 1, 0, 0, 0, 0, "20140116.123456789")
	plus("20140115.123456789", 0, 0, -1, 0, 0, 0, 0, "20140114.123456789")
	plus("20140115.123456789", 0, 1, 0, 0, 0, 0, 0, "20140215.123456789")
	plus("20140215.123456789", 0, -1, 0, 0, 0, 0, 0, "20140115.123456789")
	plus("20140115.123456789", 1, 0, 0, 0, 0, 0, 0, "20150115.123456789")
	plus("20140115.123456789", -1, 0, 0, 0, 0, 0, 0, "20130115.123456789")

	// overflow
	plus("20140115.123456789", 0, 0, 0, 0, 0, 0, 300, "20140115.123457089")
	plus("20140115.123456789", 0, 0, 0, 0, 0, 0, 2300, "20140115.123459089")
	plus("20140115.123456789", 0, 0, 0, 0, 0, 0, -1800, "20140115.123454989")
	plus("20140115.235959999", 0, 0, 0, 0, 0, 0, 1, "20140116")
	plus("20120228", 0, 0, 1, 0, 0, 0, 0, "20120229") // leap year
	plus("20130228", 0, 0, 1, 0, 0, 0, 0, "20130301")
	plus("20140215", 0, 20, 0, 0, 0, 0, 0, "20151015")
	plus("20140115", 0, -2, 0, 0, 0, 0, 0, "20131115")
}

func TestDateWeekDay(t *testing.T) {
	weekday := func(s string, wd int) {
		d := DateFromLiteral(s).(SuDate)
		assert.T(t).This(d.WeekDay()).Is(wd)
	}
	weekday("20140112", 0)
	weekday("20140115", 3)
	weekday("20140118", 6)
}

func TestDateMinusDays(t *testing.T) {
	assert := assert.T(t).This
	minusdays := func(s1 string, s2 string, expected int) {
		d1 := DateFromLiteral(s1).(SuDate)
		d2 := DateFromLiteral(s2).(SuDate)
		assert(d1.MinusDays(d1)).Is(0)
		assert(d2.MinusDays(d2)).Is(0)
		assert(d1.MinusDays(d2)).Is(expected)
		assert(d2.MinusDays(d1)).Is(-expected)
	}
	minusdays("20140215", "20140214", 1)
	minusdays("20140215", "20140115", 31)
	minusdays("20140215", "20130215", 365)
	minusdays("20130215", "20120215", 366)
}

func TestDateMinusMs(t *testing.T) {
	assert := assert.T(t).This
	minusms := func(s1 string, s2 string, expected int64) {
		if len(s1) == 9 {
			s1 = "20140115." + s1
		}
		d1 := DateFromLiteral(s1).(SuDate)
		if len(s2) == 9 {
			s2 = "20140115." + s2
		}
		d2 := DateFromLiteral(s2).(SuDate)
		assert(d1.MinusMs(d1)).Is(int64(0))
		assert(d2.MinusMs(d2)).Is(int64(0))
		assert(d1.MinusMs(d2)).Is(expected)
		assert(d2.MinusMs(d1)).Is(-expected)
	}
	minusms("123456008", "123456005", 3)
	minusms("123456008", "123455005", 1003)
	minusms("123456008", "123356008", 60*1000)
	minusms("123456008", "113456008", 60*60*1000)

	minusms("20140115", "20140114.235959999", 1)
	minusms("20140115", "20140114.235958999", 1+1000)
	minusms("20140115", "20140114.225959999", 1+60*60*1000)
}

func TestDateFormat(t *testing.T) {
	format := func(date string, format string, expected string) {
		d := DateFromLiteral(date).(SuDate)
		assert.T(t).This(d.Format(format)).Is(expected)
	}
	format("20140108", "yy-M-d", "14-1-8")
	format("20140116", "yy-MM-dd", "14-01-16")
	format("20140116", "yyyy-MM-dd", "2014-01-16")
	format("20140116", "ddd MMM dd, yyyy", "Thu Jan 16, 2014")
	format("20140116", "xx dddd MMMM dd, yyyy zz",
		"xx Thursday January 16, 2014 zz")

	format("20140108.103855", "HH:mm:ss", "10:38:55")
	format("20140108.103855", "hh:mm:ss a", "10:38:55 a")
	format("20140108.103855", "hh:mm:ss aa", "10:38:55 am")
	format("20140108.103855", "hh:mm:ss A", "10:38:55 A")
	format("20140108.103855", "hh:mm:ss AA", "10:38:55 AM")
	format("20140108.233855", "HH:mm:ss", "23:38:55")
	format("20140108.233855", "hh:mm:ss a", "11:38:55 p")
	format("20140108.233855", "hh:mm:ss aa", "11:38:55 pm")
	format("20140108.233855", "hh:mm:ss A", "11:38:55 P")
	format("20140108.233855", "hh:mm:ss AA", "11:38:55 PM")
	format("20140108.093855", "hh:mm:ss", "09:38:55")
	format("20140108.093855", "h:mm:ss", "9:38:55")
	format("20140108.103855", "h:mm:ss", "10:38:55")
	format("20140108.093855", "h 'h:mm:ss' s", "9 h:mm:ss 55")
}

func TestParseDate(t *testing.T) {
	parse := func(ds string, fmt string, expected string) {
		d := ParseDate(ds, fmt)
		assert.T(t).This(ParseDate(ds, fmt)).Isnt(NilDate)
		assert.T(t).This(d.Format("yyyy MMM d")).Is(expected)
	}
	noparse := func(ds string, fmt string) {
		assert.T(t).This(ParseDate(ds, fmt)).Is(NilDate)
	}
	parse("090625", "yMd", "2009 Jun 25")
	parse("20090625", "yMd", "2009 Jun 25")
	parse("June 25, 2009", "yMd", "2009 Jun 25")
	parse("020304", "yMd", "2002 Mar 4")
	parse("020304", "Mdy", "2004 Feb 3")
	parse("032299", "yMd", "1999 Mar 22")
	parse("2009-06-25", "yMd", "2009 Jun 25")
	parse("Wed. 25 June '09", "yMd", "2009 Jun 25")
	parse("30000101", "yMd", "3000 Jan 1")

	noparse("19992525", "yMd")
	noparse("19991233", "yMd")
	noparse("30010303", "yMd")
}

func TestIncrement(t *testing.T) {
	assert := assert.T(t).This
	d := DateFromLiteral("#20230315").(SuDate)
	assert(d.time).Is(0)
	d = d.AddMs(1)
	assert(d.time).Is(1)
	d = d.AddMs(1)
	assert(d.time).Is(2)
	d.time = 999
	d = d.AddMs(1)
	assert(d.Millisecond()).Is(0)
	assert(d.Second()).Is(1)
	d.time = 1024 + 999
	d = d.AddMs(1)
	assert(d.Millisecond()).Is(0)
	assert(d.Second()).Is(2)

	d = NewDate(2021, 5, 21, 23, 59, 59, 999)
	assert(d.Millisecond()).Is(999)
	assert(d.Second()).Is(59)
	d2 := d.AddMs(1)
	assert(d2.time).Is(0)
	assert(d2.date).Is(d.date + 1)

	// spring daylight savings change
	d = DateFromLiteral("#20230312.015959999").(SuDate)
	d = d.AddMs(1)
	assert(d.Hour()).Is(2)
	assert(d.Minute()).Is(0)
}

func FuzzParseDate(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		ParseDate(s, "yMd")
	})
}

// to run: go test -fuzz=FuzzParseDate -run=FuzzParseDate
