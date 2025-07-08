// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/pack"
)

/*
SuDate is a Suneido date/time Value

Represents a readable "local" date and time.
Does not take into account time zones or daylight savings.

It is designed to be efficient to pack and unpack
and to convert to human readable formats.
(Calculations are less common.)
*/
type SuDate struct {
	ValueBase[SuDate]
	// 21 bits for year, 4 bits for month (1-12), 5 bits for day (1-31)
	date uint32
	// 10 bits for hour, 6 bits for minute, 6 bits for second, 10 bits for ms
	time uint32
}

var NilDate SuDate
var DateBegin = NewDate(1700, 1, 1, 0, 0, 0, 0)
var DateEnd = NewDate(3000, 1, 1, 0, 0, 0, 0)

func DateTime(date uint32, time uint32) SuDate {
	d := SuDate{date: date, time: time}
	if !valid(d.Year(), d.Month(), d.Day(),
		d.Hour(), d.Minute(), d.Second(), d.Millisecond()) {
		return NilDate
	}
	return d
}

// NewDate returns a SuDate value, month is 1-12, day is 1-31
func NewDate(yr int, mon int, day int, hr int, min int, sec int, ms int) SuDate {
	if !valid(yr, mon, day, hr, min, sec, ms) {
		return NilDate
	}
	date := uint32(yr<<9) | uint32(mon<<5) | uint32(day)
	time := uint32(hr<<22) | uint32(min<<16) | uint32(sec<<10) | uint32(ms)
	return DateTime(date, time)
}

/* Now returns a SuDate for the current local date & time */
func Now() SuDate {
	return FromGoTime(time.Now())
}

// DateFromLiteral returns an SuDate or an SuTimestamp
// from the Suneido literal format yyyymmdd[.hhmm[ss[mmm[ccc]]]]
// It returns NilDate if the string is invalid.
func DateFromLiteral(s string) PackableValue {
	if s[0] == '#' {
		s = s[1:]
	}
	datelen := strings.IndexRune(s, '.')
	timelen := 0
	if datelen == -1 {
		datelen = len(s)
	} else {
		timelen = len(s) - datelen - 1
	}
	if datelen != 8 || (timelen != 0 && timelen != 4 && timelen != 6 &&
		timelen != 9 && timelen != 12) {
		return NilDate
	}

	year := nsub(s, 0, 4)
	month := nsub(s, 4, 6)
	day := nsub(s, 6, 8)

	hour := nsub(s, 9, 11)
	minute := nsub(s, 11, 13)
	second := nsub(s, 13, 15)
	millisecond := nsub(s, 15, 18)
	d := NewDate(year, month, day, hour, minute, second, millisecond)
	if timelen == 12 {
		extra := nsub(s, 18, 21)
		if extra <= 0 || extra >= 256 {
			return NilDate
		}
		return SuTimestamp{SuDate: d, extra: uint8(extra)}
	}
	return d
}

func nsub(s string, from int, to int) int {
	if to > len(s) {
		return 0
	}
	i, err := strconv.Atoi(s[from:to])
	if err != nil {
		return -1
	}
	return i
}

func valid(yr int, mon int, day int, hr int, min int, sec int, ms int) bool {
	if yr == mmYear.max &&
		(mon != 1 || day != 1 || hr != 0 || min != 0 || sec != 0 || ms != 0) {
		return false
	}
	if !mmYear.valid(yr) || !mmMonth.valid(mon) || !mmDay.valid(day) ||
		!mmHour.valid(hr) || !mmMinute.valid(min) ||
		!mmSecond.valid(sec) || !mmMillisecond.valid(ms) {
		return false
	}
	t := goTime(yr, mon, day, 0, 0, 0, 0)
	return t.Year() == yr && int(t.Month()) == mon && t.Day() == day
}

// getters

func (d SuDate) Year() int {
	return int(d.date >> 9)
}

func (d SuDate) Month() int {
	return int((d.date >> 5) & 0xf)
}

func (d SuDate) Day() int {
	return int(d.date & 0x1f)
}

func (d SuDate) Hour() int {
	return int(d.time >> 22)
}

func (d SuDate) Minute() int {
	return int((d.time >> 16) & 0x3f)
}

func (d SuDate) Second() int {
	return int((d.time >> 10) & 0x3f)
}

func (d SuDate) Millisecond() int {
	return int(d.time & 0x3ff)
}

func (d SuDate) Plus(yr int, mon int, day int, hr int, min int, sec int, ms int) SuDate {
	yr += d.Year()
	mon += d.Month()
	day += d.Day()
	hr += d.Hour()
	min += d.Minute()
	sec += d.Second()
	ms += d.Millisecond()
	nd := NormalizeDate(yr, mon, day, hr, min, sec, ms)
	if nd == NilDate {
		panic("bad date")
	}
	return nd
}

func NormalizeDate(yr int, mon int, day int, hr int, min int, sec int, ms int) SuDate {
	// use UTC to avoid timezone daylight savings issues
	t := time.Date(yr, time.Month(mon), day, hr, min, sec, ms*1000000, time.UTC)
	return FromGoTime(t)
}

// AddMs is used by Timestamp. It will usually be faster than Plus
func (d SuDate) AddMs(ms int) SuDate {
	assert.That(0 < ms && ms < 100)
	orig := d
	if int(d.Millisecond())+ms < 1000 {
		d.time += uint32(ms) // fast path
		return d
	}
	return orig.Plus(0, 0, 0, 0, 0, 0, 1) // slower fallback
}

func (d SuDate) WithoutMs() SuDate {
	d.time &^= 0x3ff
	return d
}

// WeekDay returns the day of the week - Sun is 0, Sat is 6
func (d SuDate) WeekDay() int {
	return int(d.ToGoTime().Weekday())
}

// MinusDays returns the difference between two Dates in days
func (d SuDate) MinusDays(other SuDate) int {
	return (int)(d.jday() - other.jday())
}

func (d SuDate) jday() int64 {
	return julianDayNumber(d.Year(), d.Month(), d.Day())
}

// julianDayNumber returns the time's Julian Day Number
// relative to the epoch 12:00 January 1, 4713 BC, Monday.
// NOTE: based on Go time package code
func julianDayNumber(year, month, day int) int64 {
	a := int64(14-month) / 12
	y := int64(year) + 4800 - a
	m := int64(month) + 12*a - 3
	return int64(day) + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

// MinusMs returns the difference between two Dates in milliseconds
//
// WARNING: doing this around daylight savings changes may be problematic
func (d SuDate) MinusMs(other SuDate) int64 {
	if d.date == other.date {
		return d.timeAsMs() - other.timeAsMs()
	}
	return d.UnixMilli() - other.UnixMilli()
}

func (d SuDate) timeAsMs() int64 {
	return int64(d.Millisecond()) +
		int64(1000)*int64(d.Second()+60*(d.Minute()+60*d.Hour()))
}

// UnixMilli returns the time in milliseconds since 1 Jan 1970
func (d SuDate) UnixMilli() int64 {
	return d.ToGoTime().UnixMilli()
}

func SuDateFromUnixMilli(t int64) SuDate {
	return FromGoTime(time.UnixMilli(t))
}

func (d SuDate) ToGoTime() time.Time {
	return goTime(d.Year(), d.Month(), d.Day(),
		d.Hour(), d.Minute(), d.Second(), d.Millisecond())
}

func goTime(yr int, mon int, day int, hr int, min int, sec int, ms int) time.Time {
	return time.Date(yr, time.Month(mon), day, hr, min, sec, ms*1000000, time.Local)
}

func FromGoTime(t time.Time) SuDate {
	return NewDate(t.Year(), int(t.Month()), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000)
}

// Format converts the date to a string in the specified format
func (d SuDate) Format(fmt string) string {
	fmtlen := len(fmt)
	var dst strings.Builder
	add := func(i int) {
		dst.WriteByte('0' + byte(i))
	}
	dst.Grow(fmtlen)
	for i := 0; i < fmtlen; i++ {
		n := 1
		if ascii.IsLetter(fmt[i]) {
			for c := fmt[i]; i+1 < fmtlen && fmt[i+1] == c; i++ {
				n++
			}
		}
		switch fmt[i] {
		case 'y':
			yr := d.Year()
			if n >= 4 {
				add(yr / 1000)
			}
			if n >= 3 {
				add((yr % 1000) / 100)
			}
			if n >= 2 || (yr%100) > 9 {
				add((yr % 100) / 10)
			}
			add(yr % 10)
		case 'M':
			mon := d.Month()
			if n > 3 {
				dst.WriteString(months[mon-1])
			} else if n == 3 {
				dst.WriteString(months[mon-1][0:3])
			} else {
				if n >= 2 || mon > 9 {
					add(mon / 10)
				}
				add(mon % 10)
			}
		case 'd':
			if n > 3 {
				dst.WriteString(days[d.WeekDay()])
			} else if n == 3 {
				dst.WriteString(days[d.WeekDay()][0:3])
			} else {
				if n >= 2 || d.Day() > 9 {
					add(d.Day() / 10)
				}
				add(d.Day() % 10)
			}
		case 'h': // 12 hour
			hr := d.Hour() % 12
			if hr == 0 {
				hr = 12
			}
			if n >= 2 || hr > 9 {
				add(hr / 10)
			}
			add(hr % 10)
		case 'H': // 24 hour
			if n >= 2 || d.Hour() > 9 {
				add(d.Hour() / 10)
			}
			add(d.Hour() % 10)
		case 'm':
			if n >= 2 || d.Minute() > 9 {
				add(d.Minute() / 10)
			}
			add(d.Minute() % 10)
		case 's':
			if n >= 2 || d.Second() > 9 {
				add(d.Second() / 10)
			}
			add(d.Second() % 10)
		case 'a':
			if d.Hour() < 12 {
				dst.WriteRune('a')
			} else {
				dst.WriteRune('p')
			}
			if n > 1 {
				dst.WriteRune('m')
			}
		case 'A', 't':
			if d.Hour() < 12 {
				dst.WriteRune('A')
			} else {
				dst.WriteRune('P')
			}
			if n > 1 {
				dst.WriteRune('M')
			}
		case '\'':
			for i++; i < fmtlen && (fmt[i] != '\''); i++ {
				dst.WriteByte(fmt[i])
			}
		case '\\':
			i++
			dst.WriteByte(fmt[i])
		default:
			for ; n > 0; n-- {
				dst.WriteByte(fmt[i])
			}
		}
	}
	return dst.String()
}

var months = []string{
	"January", "February", "March", "April", "May", "June", "July",
	"August", "September", "October", "November", "December"}
var days = []string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
	"Saturday"}

// ParseDate converts a human readable date to a SuDate.
//
// Returns NilDate if it fails.
func ParseDate(s string, order string) SuDate {
	NOTSET := 9999
	year := NOTSET
	month := 0
	day := 0
	hour := NOTSET
	minute := NOTSET
	second := NOTSET
	millisecond := 0

	datePatterns := []string{
		"", // set to supplied order
		"md",
		"dm",
		"dmy",
		"mdy",
		"ymd",
	}

	syspat := getSyspat(order, datePatterns)

	// scan
	const MAXTOKENS = 20
	var tokens [MAXTOKENS]int
	var typ [MAXTOKENS + 1]minmax
	ntokens := 0
	tok := func(t int) {
		if ntokens < MAXTOKENS {
			tokens[ntokens] = t
			ntokens++
		}
	}
	gotTime := false
	var prev byte
	for si := 0; si < len(s); {
		c := s[si]
		next := nextWord(s, si)
		if next != "" {
			si += len(next)
			i := 0
			for ; i < 12; i++ {
				if strings.HasPrefix(months[i], next) {
					break
				}
			}
			if i < 12 {
				typ[ntokens] = mmMonth
				tok(i + 1)
			} else if next == "Am" || next == "Pm" {
				if next[0] == 'P' {
					if hour < 12 {
						hour += 12
					}
				} else { // (word[0] == 'A')
					if hour == 12 {
						hour = 0
					}
					if hour > 12 {
						return NilDate
					}
				}
			} else {
				// ignore days of week
				for i = 0; i < 7; i++ {
					if strings.HasPrefix(days[i], next) {
						break
					}
				}
				if i >= 7 {
					return NilDate
				}
			}
		} else if next = nextNumber(s, si); next != "" {
			n, _ := strconv.Atoi(next)
			size := len(next)
			si += size
			c = get(s, si)
			if size == 6 || size == 8 {
				dig := digits{next, 0}
				switch size {
				case 6:
					// date with no separators with yy
					tok(dig.get(2))
					tok(dig.get(2))
					tok(dig.get(2))
				case 8:
					// date with no separators with yyyy
					for i := range 3 {
						if syspat[i] == 'y' {
							tok(dig.get(4))
						} else {
							tok(dig.get(2))
						}
					}
				}
				if c == '.' { // time
					si++
					time := nextNumber(s, si)
					tlen := len(time)
					si += tlen
					if tlen == 4 || tlen == 6 || tlen == 9 {
						dig = digits{time, 0}
						hour = dig.get(2)
						minute = dig.get(2)
						if tlen >= 6 {
							second = dig.get(2)
							if tlen == 9 {
								millisecond = dig.get(3)
							}
						}
					}
				}
			} else if prev == ':' || c == ':' || ampmAhead(s, si) {
				// time
				gotTime = true
				if hour == NOTSET {
					hour = n
				} else if minute == NOTSET {
					minute = n
				} else if second == NOTSET {
					second = n
				} else {
					return NilDate
				}
			} else {
				// date
				tok(n)
				if prev == '\'' {
					typ[ntokens-1] = mmYear
				}
			}
		} else {
			prev = c // ignore
			si++
		}
		if ntokens >= MAXTOKENS {
			return NilDate
		}
	}
	if hour == NOTSET {
		hour = 0
	}
	if minute == NOTSET {
		minute = 0
	}
	if second == NOTSET {
		second = 0
	}

	// search for date match
	pat := 0
	p := ""
	for ; pat < len(datePatterns); pat++ {
		p = datePatterns[pat]
		// try one pattern
		var t int
		for t = 0; t < len(p) && t < ntokens; t++ {
			var part minmax
			switch p[t] {
			case 'y':
				part = mmYear
			case 'm':
				part = mmMonth
			case 'd':
				part = mmDay
			default:
				assert.ShouldNotReachHere()
			}
			if (typ[t] != mmUnknown && typ[t] != part) ||
				tokens[t] < part.min || tokens[t] > part.max {
				break
			}
		}
		// stop at first match
		assert.That(p != "")
		if t == len(p) && t == ntokens {
			break
		}
	}
	assert.That(p != "")

	now := Now()

	if pat < len(datePatterns) {
		// use match
		for t := range len(p) {
			switch p[t] {
			case 'y':
				year = tokens[t]
			case 'm':
				month = tokens[t]
			case 'd':
				day = tokens[t]
			default:
				assert.ShouldNotReachHere()
			}
		}
	} else if gotTime && ntokens == 0 {
		year = now.Year()
		month = now.Month()
		day = now.Day()
	} else {
		return NilDate // no match
	}

	if year == NOTSET {
		if month >= max(now.Month()-5, 1) &&
			month <= min(now.Month()+6, 12) {
			year = now.Year()
		} else if now.Month() < 6 {
			year = now.Year() - 1
		} else {
			year = now.Year() + 1
		}
	} else if year < 100 {
		thisyr := now.Year()
		year += thisyr - thisyr%100
		if year-thisyr > 20 {
			year -= 100
		}
	}
	if !valid(year, month, day, hour, minute, second, millisecond) {
		return NilDate
	}
	return NewDate(year, month, day, hour, minute, second, millisecond)
}

func nextWord(s string, si int) string {
	dst := []byte{}
	for ; si < len(s) && ascii.IsLetter(s[si]); si++ {
		dst = append(dst, byte(ascii.ToLower(s[si])))
	}
	if len(dst) == 0 {
		return ""
	}
	dst[0] = byte(ascii.ToUpper(dst[0]))
	return string(dst)
}

func nextNumber(s string, si int) string {
	i := si
	for i < len(s) && ascii.IsDigit(s[i]) {
		i++
	}
	return s[si:i]
}

func getSyspat(order string, datePatterns []string) []byte {
	syspat := make([]byte, 3)
	i := 0
	oc := byte(0)
	prev := byte(0)
	for oi := 0; oi < len(order) && i < 3; oi++ {
		oc = order[oi]
		if oc != prev && (oc == 'y' || oc == 'M' || oc == 'd') {
			syspat[i] = byte(ascii.ToLower(oc))
			i++
		}
		prev = oc
	}
	if i != 3 {
		panic("invalid date format: '" + order + "'")
	}
	datePatterns[0] = string(syspat)

	// swap month-day patterns if system setting is day first
	for i = 0; i < 3; i++ {
		if syspat[i] == 'm' {
			break
		} else if syspat[i] == 'd' {
			datePatterns[1], datePatterns[2] = datePatterns[2], datePatterns[1]
		}
	}
	return syspat
}

func ampmAhead(s string, i int) bool {
	s0 := get(s, i)
	if s0 == ' ' {
		i++
		s0 = get(s, i)
	}
	s0 = byte(ascii.ToLower(s0))
	return (s0 == 'a' || s0 == 'p') &&
		ascii.ToLower(get(s, i+1)) == 'm'
}

func get(s string, i int) byte {
	if i >= len(s) {
		return 0
	}
	return s[i]
}

type digits struct {
	s string
	i int
}

func (d *digits) get(n int) int {
	d.i += n
	i, _ := strconv.Atoi(d.s[d.i-n : d.i])
	return i
}

type minmax struct {
	min int
	max int
}

func (m minmax) valid(n int) bool {
	return m.min <= n && n <= m.max
}

var (
	mmYear        = minmax{0, 3000}
	mmMonth       = minmax{1, 12}
	mmDay         = minmax{1, 31}
	mmHour        = minmax{0, 23}
	mmMinute      = minmax{0, 59}
	mmSecond      = minmax{0, 59}
	mmMillisecond = minmax{0, 999}
	mmUnknown     = minmax{0, 0}
)

// Value interface --------------------------------------------------

var _ Value = (*SuDate)(nil)

func (d SuDate) String() string {
	if d.time == 0 {
		return fmt.Sprintf("#%04d%02d%02d", d.Year(), d.Month(), d.Day())
	}
	s := fmt.Sprintf("#%04d%02d%02d.%02d%02d%02d%03d",
		d.Year(), d.Month(), d.Day(),
		d.Hour(), d.Minute(), d.Second(), d.Millisecond())
	if strings.HasSuffix(s, "00000") {
		return s[0:14]
	}
	if strings.HasSuffix(s, "000") {
		return s[0:16]
	}
	return s
}

func (d SuDate) Equal(other any) bool {
	return d == other
}

func (d SuDate) Hash() uint64 {
	h := uint64(17)
	h = 31*h + uint64(d.date)
	h = 31*h + uint64(d.time)
	return h
}

func (d SuDate) Hash2() uint64 {
	return d.Hash()
}

func (SuDate) Type() types.Type {
	return types.Date
}

func (d SuDate) Compare(other Value) int {
	if cmp := cmp.Compare(ordDate, Order(other)); cmp != 0 {
		return cmp * 2
	}
	if st, ok := other.(SuTimestamp); ok {
		return CompareSuTimestamp(SuTimestamp{SuDate: d}, st)
	}
	d2 := other.(SuDate)
	if d.date < d2.date {
		return -1
	} else if d.date > d2.date {
		return +1
	} else if d.time < d2.time {
		return -1
	} else if d.time > d2.time {
		return +1
	}
	return 0
}

func (SuDate) SetConcurrent() {
	// immutable so ok
}

// DateMethods is initialized by the builtin package
var DateMethods Methods

var gnDates = Global.Num("Dates")

func (SuDate) Lookup(th *Thread, method string) Value {
	return Lookup(th, DateMethods, gnDates, method)
}

// Packable interface -----------------------------------------------

var _ Packable = SuDate{}

// PackSize returns the packed size (Packable interface)
func (SuDate) PackSize(*uint64) int {
	return 9
}

func (SuDate) PackSize2(*uint64, packStack) int {
	return 9
}

// Pack packs into the supplied byte slice (Packable interface)
func (d SuDate) Pack(_ *uint64, buf *pack.Encoder) {
	buf.Put1(PackDate).Uint32(d.date).Uint32(d.time)
}

// UnpackDate unpacks a date from the supplied byte slice
func UnpackDate(s string) Value {
	d := pack.NewDecoder(s[1:])
	date := d.Uint32()
	time := d.Uint32()
	sd := SuDate{date: date, time: time}
	if d.Remaining() > 0 {
		return UnpackTimestamp(sd, d)
	}
	return sd
}
