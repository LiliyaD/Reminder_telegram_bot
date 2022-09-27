package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type funcType func(string) (interface{}, error)

var (
	errBadParameters     = errors.New("Parameters are incorrect")
	errUnavailableServer = errors.New("Server is unavailable. Please, try later")
	errDate              = errors.New("Date can't be parsed")
)

func ParseDate(d string) (time.Time, error) {
	date, err := time.Parse("02.01.2006", d)
	if err != nil {
		return time.Time{}, errDate
	}

	return date, nil
}

func parseDateBot(d string) (interface{}, error) {
	if d == "_" {
		return time.Time{}, nil
	}
	return ParseDate(d)
}

func parseTimesBot(d string) (interface{}, error) {
	if d == "_" {
		return uint8(0), nil
	}
	t, err := strconv.ParseUint(d, 10, 8)
	if err != nil {
		journal.LogError(d, err)
		return uint8(0), err
	}

	return uint8(t), nil
}

func parseQuantityBot(d string) (interface{}, error) {
	if d == "_" {
		return float32(0.0), nil
	}
	q, err := strconv.ParseFloat(d, 32)
	if err != nil {
		journal.LogError(d, err)
		return float32(0), err
	}

	return float32(q), nil
}

func parseData(s string, allowSkips bool) ([]interface{}, error) {
	p := strings.Split(s, " ")
	if len(p) != 5 {
		journal.LogError(s, errBadParameters)
		return nil, errBadParameters
	}

	if !allowSkips {
		for i := 0; i < 5; i++ {
			if p[i] == "_" {
				journal.LogError(s, errBadParameters)
				return nil, errBadParameters
			}
		}
	}

	res := make([]interface{}, 5)
	res[0] = p[0]

	f := [4]funcType{parseDateBot, parseDateBot, parseTimesBot, parseQuantityBot}

	var err error
	for i := 0; i < 4; i++ {
		res[i+1], err = f[i](p[i+1])
		if err != nil {
			journal.LogError(p[i+1], err)
			return nil, err
		}
	}

	return res, nil
}

func formText(act, beginDate, endDate string, timesPerDay uint32, quantityPerTime float32) string {
	t := strconv.FormatUint(uint64(timesPerDay), 10)
	q0 := quantityPerTime
	n := int32(q0)
	k := q0 - float32(n)
	var prec int = 0
	if k > 0 {
		prec = 2
	}
	q := strconv.FormatFloat(float64(q0), 'f', prec, 32)

	return t + " " + act + " " + q + " times per day from " + beginDate + " to " + endDate + "\n"
}

func prepareErrorTextForUser(err error) string {
	switch status.Code(err) {
	case codes.Unavailable:
		return errUnavailableServer.Error()
	case codes.InvalidArgument:
		return multiPkg.ErrInvalidParameter.Error()
	case codes.NotFound:
		return multiPkg.ErrNotExist.Error()
	case codes.Internal:
		return multiPkg.ErrTimeOut.Error()
	case codes.AlreadyExists:
		return multiPkg.ErrExist.Error()
	default:
		return err.Error()
	}
}
