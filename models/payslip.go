package models

import "time"

type PayslipSummary struct {
	Gross      float64 `json:"gross" bson:"gross"`
	Taxable    float64 `json:"taxable" bson:"taxable"`
	Taxes      float64 `json:"taxes" bson:"taxes"`
	Deductions float64 `json:"deductions" bson:"deductions"`
	Net        float64 `json:"net" bson:"net"`
}

type PayslipHead struct {
	Name    string  `json:"name" bson:"name"`
	Current float64 `json:"current" bson:"current"`
	YTD     float64 `json:"ytd" bson:"ytd"`
}

type Payslip struct {
	DocumentID string    `json:"document_id" bson:"document_id"`
	Currency   string    `json:"currency" bson:"currency"`
	Date       time.Time `json:"date" bson:"date"`
	StartDate  time.Time `json:"start_date" bson:"start_date"`
	EndDate    time.Time `json:"end_date" bson:"end_date"`
	NetPay     float64   `json:"net_pay" bson:"net_pay"`
	PayslipSummary
	IncomeHeads    []*PayslipHead `json:"income_heads" bson:"income_heads"`
	DeductionHeads []*PayslipHead `json:"deduction_heads" bson:"deduction_heads"`
	TaxHeads       []*PayslipHead `json:"tax_heads" bson:"tax_heads"`
}
