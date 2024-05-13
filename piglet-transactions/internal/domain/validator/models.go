package validation

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ValTrans struct {
	Date      *timestamppb.Timestamp `validate:"required"`
	TransType int32                  `validate:"required,min=1,max=4"`
	Sum       float32                `validate:"required,min=0.001"`
	Comment   string
}

type ValIncome struct {
	IdCategory string `validator:"uuid4"`
	IdBillTo   string `validator:"required,uuid4"`
	Sender     string
	Repeat     bool
}

type ValExpense struct {
	IdCategory string `validator:"uuid4"`
	IdBillFrom string `validator:"required,uuid4"`
	Recipient  string
	Repeat     bool
}

type ValDebt struct {
	DebtType       bool
	IdBillFrom     string `validator:"uuid4"`
	IdBillTo       string `validator:"uuid4"`
	CreditorDebtor string
}

type ValTransfer struct {
	IdBillFrom string `validator:"required,uuid4"`
	IdBillTo   string `validator:"required,uuid4"`
}

type ValCategory struct {
	CategoryType bool
	Name         string `validator:"required"`
	Mandatory    bool
}
