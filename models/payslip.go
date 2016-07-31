package models

import "time"

type PayslipSummary struct {
	Gross      float64
	Taxable    float64
	Taxes      float64
	Deductions float64
	Net        float64
}

type PayslipHead struct {
	Name    string
	Current float64
	YTD     float64
}

type Payslip struct {
	DocumentID string
	Date       time.Time
	StartDate  time.Time
	EndDate    time.Time
	NetPay     float64
	PayslipSummary
	IncomeHeads    []*PayslipHead
	DeductionHeads []*PayslipHead
	TaxHeads       []*PayslipHead
}
