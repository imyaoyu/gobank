package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/imyaoyu/goat/app"
)

func main() {

	// Midware: check userId
	app.Func(func(c *app.ApiCtx) {
		if c.UserId == "" {
			c.Panic(401, "No Auth Info", nil)
		}
	})

	//Add Apis
	app.Api("b.1", QueryProd)
	app.Api("b.2", BuyProd)
	app.Api("b.3", SellProd)
	app.Api("b.4", QueryAcct)

	//Start Server
	app.Run()

}

type IQueryProd struct {
	ProdType string `validate:"min=1,max=1"`
}

type OQueryProd struct {
	Num   int
	Prods []*TProdInfo
}

type TProdInfo struct {
	Id                                   int64
	ProdId, ProdType, ProdName, ProdNote string
	ProdRate                             int64
	ProdRateStr                          string
	ProdTerm                             int
}

// 产品查询
func QueryProd(c *app.ApiCtx) {
	i, o := new(IQueryProd), new(OQueryProd)
	c.Init(i, o)
	var prods []*TProdInfo
	c.SelectS(&prods, 10, 0, "prod_id", "prod_type=?", i.ProdType)
	o.Num = len(prods)
	o.Prods = prods

}

type IBuyProd struct {
	ProdId string `validate:"min=1,max=10"`
	Amount string `validate:"regexp=^(0|[1-9]\d*)(\.\d{1,2})?$"`
}

type OBuyProd struct {
	ProdName, RateStr, OpenDate string
}

type TAcctInfo struct {
	Id                           int64
	UserId                       string
	ProdId, ProdType, ProdName   string
	Rate                         int64
	Amount, Balance              Money
	OpenDate, EndDate, CloseDate string
	Status                       string
	//Interest                     int64
}

// 购买存款
func BuyProd(c *app.ApiCtx) {
	i, o := new(IBuyProd), new(OBuyProd)
	c.Init(i, o)
	var prod TProdInfo
	if has := c.Select(&prod, "prod_id=?", i.ProdId); !has {
		c.Panic(500, "No Found Prod", nil)
	}
	acct := new(TAcctInfo)
	acct.UserId = c.UserId
	acct.ProdId = i.ProdId
	acct.ProdType = prod.ProdType
	acct.ProdName = prod.ProdName
	acct.Rate = prod.ProdRate
	acct.Amount = MoneyOf(i.Amount)
	acct.Balance = acct.Amount
	acct.OpenDate = Today().String()
	acct.EndDate = Today().AddMonth(prod.ProdTerm).String()
	acct.Status = "A"

	c.Insert(acct)

	o.ProdName = prod.ProdName
	o.RateStr = prod.ProdRateStr
	o.OpenDate = acct.OpenDate

}

type ISellProd struct {
	AcctId string
}

type OSellProd struct {
	CloseDate, Amount, Interest string
}

// 支取存款
func SellProd(c *app.ApiCtx) {

	i, o := new(ISellProd), new(OSellProd)
	c.Init(i, o)

	var acct TAcctInfo

	has := c.Select(&acct, "id=? and user_id=? and status='A'", i.AcctId, c.UserId)

	if !has {
		c.Panic(1001, "Invalid Acct", nil)
	}

	acct.CloseDate = Today().String()
	acct.Status = "C"
	acct.Balance = 0

	c.Update(acct, acct.Id, map[string]any{
		"close_date": acct.CloseDate,
		"status":     acct.Status,
		"balance":    acct.Balance,
	})

	o.CloseDate = acct.CloseDate
	o.Amount = acct.Amount.String()
	// o.Interest = TODO

}

type IQueryAcct struct {
	Status             string
	StartNum, LimitNum int
}

type OQueryAcct struct {
	TotalNum int64
	Accts    []*TAcctInfo
}

// 我的账户查询
func QueryAcct(c *app.ApiCtx) {
	i, o := new(IQueryAcct), new(OQueryAcct)
	c.Init(i, o)
	var accts []*TAcctInfo
	c.SelectS(&accts, i.LimitNum, i.StartNum, "id", "user_id=? and status=?", c.UserId, i.Status)
	total, _ := app.DB().Where("user_id=? and status=?", c.UserId, i.Status).Count(new(TAcctInfo))

	o.TotalNum = total
	o.Accts = accts

}

type Money int64

func MoneyOf(s string) Money {
	ss := strings.Split(s, ".")
	if len(ss) != 2 {
		panic(fmt.Errorf("%s is not a valid amount number", s))
	}
	if len(ss[1]) != 2 {
		panic(fmt.Errorf("%s is not a valid amount number", s))
	}

	numstr := ss[0] + ss[1]
	num, err := strconv.ParseInt(numstr, 10, 64)

	if err != nil {
		panic(fmt.Errorf("ParseInt %s error:%w", numstr, err))
	}

	return Money(num)
}

func (m Money) String() string {

	s := strconv.FormatInt(int64(m), 10)

	return s[:len(s)-2] + "." + s[len(s)-2:]
}

type Date struct {
	Year, Month, Day int
}

func (d *Date) String() string {
	return fmt.Sprintf("%04d%02d%02d", d.Year, d.Month, d.Day)
}

func (d *Date) AddMonth(num int) *Date {

	year, month := num/12, num%12
	d.Year += year
	d.Month += month
	if d.Month > 12 {
		d.Year += 1
		d.Month -= 12
	}
	switch d.Month {
	case 2:
		if (d.Year%4 == 0 && d.Year%100 != 0) || (d.Year%400 == 0) {
			if d.Day > 29 {
				d.Day = 29
			}
		} else {
			if d.Day > 28 {
				d.Day = 28
			}
		}
	case 4, 6, 9, 11:
		if d.Day > 30 {
			d.Day = 30
		}
	}

	return d

}

func Today() *Date {

	datetime := app.Now()

	year, _ := strconv.Atoi(datetime[0:4])
	month, _ := strconv.Atoi(datetime[5:7])
	day, _ := strconv.Atoi(datetime[8:10])

	return &Date{
		Year:  year,
		Month: month,
		Day:   day,
	}

}
