package models

import (
	"bitcoin/config"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type SignalEvent struct {
	Time        time.Time `json:"time"`
	ProductCode string    `json:"product_code"`
	Side        string    `json:"side"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
}

func (s *SignalEvent) Save() bool {
	cmd := fmt.Sprintf("INSERT INTO %s (time, product_code, side, price, size) VALUES (?, ?, ?, ?, ?)", tableNameSignalEvents)
	_, err := DbConnection.Exec(cmd, s.Time.Format(time.RFC3339), s.ProductCode, s.Side, s.Price, s.Size)
	if err != nil {
		// UNIQUE constraint failedは同じ時間のレコードが入った時にエラーが起きる（すでにレコードあるためtrueを返す）
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			log.Println(err)
			return true
		}
		return false
	}
	return true
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
}

func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

// 渡した引数の個数分のevents取得する関数
func GetSignalEventsByCount(loadEvents int) *SignalEvents {
	cmd := fmt.Sprintf(`SELECT * FROM (
        SELECT time, product_code, side, price, size FROM %s WHERE product_code = ? ORDER BY time DESC LIMIT ? )
        ORDER BY time ASC;`, tableNameSignalEvents)
	rows, err := DbConnection.Query(cmd, config.Config.ProductCode, loadEvents)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var signalEvents SignalEvents
	for rows.Next() {
		var signalEvent SignalEvent
		rows.Scan(&signalEvent.Time, &signalEvent.ProductCode, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
		signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	return &signalEvents
}

// 引数で指定した時間以降のデータを取得する関数
func GetSignalEventsAfterTime(timeTime time.Time) *SignalEvents {
	cmd := fmt.Sprintf(`SELECT * FROM (
                SELECT time, product_code, side, price, size FROM %s
                WHERE DATETIME(time) >= DATETIME(?)
                ORDER BY time DESC
            ) ORDER BY time ASC;`, tableNameSignalEvents)
	rows, err := DbConnection.Query(cmd, timeTime.Format(time.RFC3339))
	if err != nil {
		return nil
	}
	defer rows.Close()

	var signalEvents SignalEvents
	for rows.Next() {
		var signalEvent SignalEvent
		rows.Scan(&signalEvent.Time, &signalEvent.ProductCode, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
		signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	}
	return &signalEvents
}

func (s *SignalEvents) CanBuy(time time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return true
	}
	lastSignal := s.Signals[lenSignals-1]
	// 最後の取引がSellかつ今の時間よりも前の取引であることが条件
	if lastSignal.Side == "SEll" && lastSignal.Time.Before(time) {
		return true
	}
	return false
}

func (s *SignalEvents) CanSell(time time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		// canbuyと違いfalseなのは買いが1つもなければ売れないから
		return false
	}
	lastSignal := s.Signals[lenSignals-1]
	// 最後の取引がBuyかつ今の時間よりも前の取引であることが条件
	if lastSignal.Side == "BUY" && lastSignal.Time.Before(time) {
		return true
	}
	return false
}

func (s *SignalEvents) Buy(ProductCode string, time time.Time, price, size float64, save bool) bool {
	if s.CanBuy(time) {
		return false
	}
	SignalEvent := SignalEvent{
		ProductCode: ProductCode,
		Time: time,
		Side: "BUY",
		Price: price,
		Size: size,
	}
	if save {
		SignalEvent.Save()
	}
	s.Signals = append(s.Signals, SignalEvent)
	return true
}

func (s *SignalEvents) Sell(ProductCode string, time time.Time, price, size float64, save bool) bool {
	if s.CanSell(time) {
		return false
	}
	SignalEvent := SignalEvent{
		ProductCode: ProductCode,
		Time: time,
		Side: "Sell",
		Price: price,
		Size: size,
	}
	if save {
		SignalEvent.Save()
	}
	s.Signals = append(s.Signals, SignalEvent)
	return true
}

// 利益計算する関数
func (s *SignalEvents) Profit() float64 {
	total := 0.0
	beforeSell := 0.0
	// isHoldingは買ったものをまだ持っていたらtrue
	isHolding := false
	for i, signalEvent := range s.Signals {
		if i == 0 && signalEvent.Side == "SELL" {
			continue
		}
		if signalEvent.Side == "BUY" {
			total -= signalEvent.Price * signalEvent.Size
			isHolding = true
		}
		if signalEvent.Side == "SELL" {
			total += signalEvent.Price * signalEvent.Size
			isHolding = false
			beforeSell = total
		}
	}
	if isHolding {
		return beforeSell
	}
	return total
}

// SignalEventsのフィールドにprofitを追加
func (s SignalEvents) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&struct {
		Signals []SignalEvent `json:"signals,omitempty"`
		Profit  float64       `json:"profit,omitempty"`
	}{
		Signals: s.Signals,
		Profit:  s.Profit(),
	})
	if err != nil {
		return nil, err
	}
	return value, err
}

func (s *SignalEvents) CollectAfter(time time.Time) *SignalEvents {
	for i, signal := range s.Signals {
		if time.After(signal.Time) {
			continue
		}
		return &SignalEvents{Signals: s.Signals[i:]}
	}
	return nil
}