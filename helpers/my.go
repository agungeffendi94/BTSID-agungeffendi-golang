package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cast"
)

func TimeInLocal(local string) (time.Time, time.Time) {
	t := time.Now()

	location, err := time.LoadLocation(local)
	if err != nil {
		fmt.Println(err)
	}
	resp := t.In(location)

	return t, resp
}

func Pass2Hash(plaintext string) string {
	pjg_char := len(plaintext)
	holderstr := []string{}
	for i := 0; i < pjg_char; i++ {
		start := i
		end := i + 1
		strnya := plaintext[start:end]
		h := sha1.New()
		h.Write([]byte(strnya))
		sha1_hash := hex.EncodeToString(h.Sum(nil))
		holderstr = append(holderstr, sha1_hash)
	}
	newstring := strings.Join(holderstr, "")
	hash := sha1.New()
	hash.Write([]byte(newstring))
	hashedstring := hex.EncodeToString(hash.Sum(nil))

	return hashedstring
}

func FormattedDate(stringdate string) string {
	myDate, err := time.Parse("2006-01-02 15:04:05", stringdate)
	if err != nil {
		fmt.Println("Second Parse : " + err.Error())
	}
	newdate := myDate.Format("2 Jan 2006")
	return newdate
}

func FormattedDateTZ(stringdate string) string {
	myDate, err := time.Parse("2006-01-02T15:04:05Z", stringdate)
	if err != nil {
		fmt.Println("Second Parse : " + err.Error())
	}
	newdate := myDate.Format("2 Jan 2006")
	return newdate
}

func FormattedDateWtHour(stringdate string) string {
	myDate, err := time.Parse("2006-01-02 15:04:05", stringdate)
	if err != nil {
		fmt.Println("Second Parse : " + err.Error())
	}
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println(err)
	}
	local := myDate.In(location)

	newdate := cast.ToString(local.Hour()) + ":" + cast.ToString(local.Minute()) + ", " + local.Format("2 Jan 2006")
	return newdate
}

func FormattedDateWtHourTZ(stringdate string) string {
	myDate, err := time.Parse("2006-01-02T15:04:05Z", stringdate)
	if err != nil {
		fmt.Println("Second Parse : " + err.Error())
	}

	newdate := cast.ToString(myDate.Hour()) + ":" + cast.ToString(myDate.Minute()) + ", " + myDate.Format("2 Jan 2006")
	return newdate
}

func GetJWTClaims(header []string, key string) (string, error) {
	// tokenString :=
	if header[0] == "" {
		return "", errors.New("header empty, request ignored")
	}
	jwtsplit := strings.Split(header[0], " ")
	token, _, err := new(jwt.Parser).ParseUnverified(jwtsplit[1], jwt.MapClaims{})
	if err != nil {
		fmt.Printf("Error %s", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// obtains claims
		ret := fmt.Sprint(claims[key])
		return ret, nil
	}

	return "", nil
}

func DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
