package resources

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ginserver/dao"
	"ginserver/util"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// PriceResource ...
type PriceResource struct {
}

type resData struct {
	Code    int         `json:"code"`
	Body    interface{} `json:"body"`
	Message string      `json:"message"`
}

// NewPriceResource ...
func NewPriceResource(e *gin.Engine) {
	u := PriceResource{}
	// Setup Routes
	e.POST("/post/priceData", u.getPriceData)
}

// TestData ...
type TestData struct {
	Symbol    string `form:"symbol"`
	Frequence string `form:"frequence"`
	Subtime   string `form:"subtime"`
}

func (r *PriceResource) getPriceData(c *gin.Context) {
	var testData TestData
	if c.ShouldBind(&testData) == nil {
		// 参数验证
		symbol := testData.Symbol
		symbolValidationArr := []string{"BTC", "ETH"}
		symbolValidation := util.StringInSlice(symbol, symbolValidationArr)

		frequence := testData.Frequence
		frequenceValidationArr := []string{"hour", "day", "week", "month"}
		frequenceValidation := util.StringInSlice(frequence, frequenceValidationArr)

		subtime := testData.Subtime
		subtimeValidation := false
		subtimeInt, err := strconv.Atoi(subtime)
		if err == nil && subtimeInt >= 0 {
			switch frequence {
			case "hour":
				if subtimeInt <= 59 {
					subtimeValidation = true
				}
			case "day":
				if subtimeInt <= 23 {
					subtimeValidation = true
				}
			case "week":
				if subtimeInt <= 6 {
					subtimeValidation = true
				}
			case "month":
				if subtimeInt <= 28 && subtimeInt >= 1 {
					subtimeValidation = true
				}
			default:
			}
		}

		var resStr string
		var res resData
		if !symbolValidation {
			resStr += "{\"code\":40001,"
			resStr += "\"body\":\"\","
			resStr += "\"message\": \"Symbol validation failed!\"}"
			err := json.Unmarshal([]byte(resStr), &res)
			if err != nil {
				fmt.Println(err)
			}
			c.JSON(400, res)
			return
		}

		if !frequenceValidation {
			resStr += "{\"code\":40002,"
			resStr += "\"body\":\"\","
			resStr += "\"message\": \"Frequence validation failed!\"}"
			err := json.Unmarshal([]byte(resStr), &res)
			if err != nil {
				fmt.Println(err)
			}
			c.JSON(400, res)
			return
		}

		if !subtimeValidation {
			resStr += "{\"code\":40003,"
			resStr += "\"body\":\"\","
			resStr += "\"message\": \"Subtime validation failed!\"}"
			err := json.Unmarshal([]byte(resStr), &res)
			if err != nil {
				fmt.Println(err)
			}
			c.JSON(400, res)
			return
		}

		// 匹配数据
		if symbolValidation && frequenceValidation && subtimeValidation {
			resStr += "{\"code\":200,"

			if len(subtime) == 1 {
				subtime = "0" + subtime
			}
			query := bson.M{"symbol": symbol}

			if frequence == "hour" {
				regStr := ":" + subtime + ":00"
				query["time"] = bson.M{
					"$regex": regStr,
				}
			}
			if frequence == "day" {
				regStr := " " + subtime + ":00:00"
				query["time"] = bson.M{
					"$regex": regStr,
				}
			}
			if frequence == "week" {
				//  计算周日期
				var weekArr [26]string
				var WeekDayMap = map[string]int{
					"Monday":    1,
					"Tuesday":   2,
					"Wednesday": 3,
					"Thursday":  4,
					"Friday":    5,
					"Saturday":  6,
					"Sunday":    0,
				}
				currentDateUnix := time.Now().Unix()
				currentWeekday := WeekDayMap[time.Now().Weekday().String()]
				tracebackCnt := 0
				if subtimeInt < currentWeekday {
					tracebackCnt = currentWeekday - subtimeInt
				}
				if subtimeInt > currentWeekday {
					tracebackCnt = currentWeekday + 7 - subtimeInt
				}
				targetWeekdayUnix := currentDateUnix - int64(tracebackCnt)*86400

				targetWeekday := time.Unix(targetWeekdayUnix, 0).Format("2006-01-02") + " 00:00:00"
				weekArr[25] = targetWeekday

				for i := 0; i < 25; i++ {
					targetWeekdayUnix -= 604800
					weekArr[25-i-1] = time.Unix(targetWeekdayUnix, 0).Format("2006-01-02") + " 00:00:00"
				}

				query["time"] = bson.M{
					"$in": weekArr,
				}
			}
			if frequence == "month" {
				var monthArr [6]string

				currentDateArr := strings.Split(time.Now().Format("2006-01-02"), "-")
				currentYear, _ := strconv.Atoi(currentDateArr[0])
				currentMonth, _ := strconv.Atoi(currentDateArr[1])
				currentDay, _ := strconv.Atoi(currentDateArr[2])
				monthCnt := 6
				if subtimeInt < currentDay {
					monthCnt = 5
					monthArr[5] = time.Now().Format("2006-01-") + subtime + " 00:00:00"
				}

				for i := 0; i < monthCnt; i++ {
					currentMonth--
					if currentMonth <= 0 {
						currentYear--
						currentMonth += 12
					}
					currentYearStr := strconv.Itoa(currentYear)
					currentMonthStr := strconv.Itoa(currentMonth)
					if len(currentMonthStr) == 1 {
						currentMonthStr = "0" + currentMonthStr
					}
					monthArr[monthCnt-i-1] = currentYearStr + "-" + currentMonthStr + "-" + subtime + " 00:00:00"
				}
				fmt.Println(monthArr)
				query["time"] = bson.M{
					"$in": monthArr,
				}
			}
			dbSession := c.MustGet("db").(*mgo.Session).Copy()
			defer dbSession.Close()

			priceData := dao.GetPriceData(dbSession, query)
			dataByte, _ := json.Marshal(priceData)
			resStr += "\"body\":" + string(dataByte) + ","

			resStr += "\"message\": \"Get price data Successed!\"}"
			json.Unmarshal([]byte(resStr), &res)

			c.JSON(200, res)
		}

	}

}
