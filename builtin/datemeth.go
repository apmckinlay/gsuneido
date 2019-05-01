package builtin

import (
	"strings"
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	DateMethods = Methods{
		"MinusDays": method1("(date)", func(this Value, val Value) Value {
			t1 := this.(SuDate)
			if t2, ok := val.(SuDate); ok {
				return IntVal(t1.MinusDays(t2))
			}
			panic("date.MinusDays requires date")
		}),
		"MinusSeconds": method1("(date)", func(this Value, val Value) Value {
			t1 := this.(SuDate)
			if t2, ok := val.(SuDate); ok {
				ms := t1.MinusMs(t2)
				return fromFloat(float64(ms) / 1000)
			}
			panic("date.MinusSeconds requires date")
		}),
		"FormatEn": method1("(format)", func(this, arg Value) Value {
			return SuStr(this.(SuDate).Format(IfStr(arg)))
		}),
		"GetLocalGMTBias": method0(func(this Value) Value { // should be static
			_, offset := time.Now().Zone()
			return IntVal(-offset / 60)
		}),
		"Plus": methodRaw("(n, default)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&paramSpecPlus, as)
				return this.(SuDate).Plus(ToInt(args[0]), ToInt(args[1]),
					ToInt(args[2]), ToInt(args[3]), ToInt(args[4]),
					ToInt(args[5]), ToInt(args[6]))
			}),
		"WeekDay": method1("(firstDay='Sun')", func(this, arg Value) Value {
			i := dayOfWeek(arg)
			return IntVal((this.(SuDate).WeekDay() - i) % 7)
		}),

		"Year": method0(func(this Value) Value {
			return IntVal(this.(SuDate).Year())
		}),
		"Month": method0(func(this Value) Value {
			return IntVal(this.(SuDate).Month())
		}),
		"Day": method0(func(this Value) Value {
			return IntVal(this.(SuDate).Day())
		}),
		"Hour": method0(func(this Value) Value {
			return IntVal(this.(SuDate).Hour())
		}),
		"Minute": method0(func(this Value) Value {
			return IntVal(this.(SuDate).Minute())
		}),
		"Second": method0(func(this Value) Value {
			return IntVal(this.(SuDate).Second())
		}),
		"Millisecond": method0(func(this Value) Value {
			return IntVal(this.(SuDate).Millisecond())
		}),
	}
}

var paramSpecPlus = params("(years=0, months=0, days=0, " +
	"hours=0, minutes=0, seconds=0, milliseconds=0)")

func dayOfWeek(x Value) int {
	if i, ok := x.IfInt(); ok {
		return i
	}
	s := strings.ToLower(ToStr(x))
	days := []string{"sunday", "monday", "tuesday",
		"wednesday", "thursday", "friday", "saturday"}
	for i, d := range days {
		if strings.HasPrefix(d, s) {
			return i
		}
	}
	panic("usage: date.WeekDay(day name or number)")
}
