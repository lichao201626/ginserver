package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ginserver/db"
	"ginserver/models"
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2/bson"
)

// CollectPriceData ...
func CollectPriceData() {
	// TODO: Deal with Error
	// defer
	go job()
}

func job() {
	c := cron.New()
	c.AddFunc("59 * * * * *", func() {
		// 查询数据库中最新价格日期
		session := db.Session()
		database := session.DB("quanta_lab_aip")
		collection := database.C("prices")
		latestPrice := models.NewPrice()
		collection.Find(nil).Sort("-time").One(&latestPrice)
		session.Close()

		currentTimeStr := time.Now().Format("2006-01-02")
		lastestDataTimeStr := strings.Split(latestPrice.Time, " ")[0]

		if currentTimeStr != lastestDataTimeStr {
			filter := generateDataFilter()
			data := requestPriceData(filter)
			savePriceData(data)
		}
	})
	c.Start()
	select {}
}

func generateDataFilter() map[string]string {
	// 查询数据库中最新价格日期
	session := db.Session()
	database := session.DB("quanta_lab_aip")
	collection := database.C(models.CollectionPrice)

	latestPrice := models.NewPrice()
	collection.Find(nil).Sort("-time").One(&latestPrice)
	session.Close()

	// 计算价格查询条件
	currentTime := time.Now()
	lastestDataTimeStr := latestPrice.Time
	if lastestDataTimeStr == "" {
		lastestDataTimeStr = currentTime.Add(time.Duration(-24*180) * time.Hour).Format("2006-01-02 15:04:05")
	}
	lastestDataTime, err := time.Parse("2006-01-02 15:04:05", lastestDataTimeStr)
	if err != nil {
		fmt.Println(err)
	}
	nextDataTime := lastestDataTime.Add(time.Duration(24) * time.Hour)
	nextDataTimeStr := nextDataTime.Format("2006-01-02 15:04:05")

	filter := make(map[string]string)
	// 28800: 时间戳转换后快了8小时
	filter["start"] = strconv.FormatInt(lastestDataTime.Unix()-28800, 10)
	filter["end"] = strconv.FormatInt(nextDataTime.Unix()-28800, 10)

	fmt.Println("start time: ", lastestDataTimeStr)
	fmt.Println("end time: ", nextDataTimeStr)

	return filter
}

type priceData struct {
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
}

func requestPriceData(filter map[string]string) []priceData {
	// TODO: 查询日期精确化
	preURL := "http://hist-quote.1tokentrade.cn/candles?contract=okex/btc.usdt&since="
	url := preURL + filter["start"] + "&until=" + filter["end"] + "&duration=1m&format=json"
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("ot-key", "7jAtk7Yh-EshqCHr0-NG8qMHnE-tDRdHbZe")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data []priceData
	json.Unmarshal([]byte(string(body)), &data)

	return data
}

func savePriceData(rawData []priceData) {
	for i := 0; i < len(rawData); i++ {
		price := models.NewPrice()
		price.Symbol = "BTC"
		price.Open = rawData[i].Open
		price.High = rawData[i].High
		price.Low = rawData[i].Low
		price.Close = rawData[i].Close
		price.Volume = rawData[i].Volume
		price.Time = time.Unix(rawData[i].Timestamp, 0).Format("2006-01-02 15:04:05")
		if i == len(rawData)-1 {
			fmt.Println("Prices update to: ", price.Time)
		}
		price.Timestamp = rawData[i].Timestamp

		session := db.Session()
		database := session.DB("quanta_lab_aip")
		collection := database.C("prices")
		// 数据库去重检验
		var testPrice []priceData
		collection.Find(bson.M{"time": price.Time}).All(&testPrice)

		if len(testPrice) == 0 {
			collection.Insert(price)
		} else {
			fmt.Println("Duplicated data:", price.Time)
		}

		session.Close()
	}
}
