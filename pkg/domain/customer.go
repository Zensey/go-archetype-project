package domain

type Customer struct {
	ID                   int    `json:"customerID" pg:"id,pk" sql:",pk"`
	PayerID              int    `json:"payerID,"`
	TypeID               string `json:"type_id,"`
	FullName             string `json:"fullName,"`
	CompanyName          string `json:"companyName,"`
	FirstName            string `json:"firstName,"`
	LastName             string `json:"lastName,"`
	GroupID              int    `json:"groupID,"`
	EDI                  string `json:"EDI,"`
	IsPOSDefaultCustomer int    `json:"isPOSDefaultCustomer,"`
	CountryID            string `json:"countryID,"`
	Phone                string `json:"phone,"`
	EInvoiceEmail        string `json:"eInvoiceEmail,"`
	Email                string `json:"email,"`
	Fax                  string `json:"fax,"`
	Code                 string `json:"code,"`
	ReferenceNumber      string `json:"referenceNumber,"`
	VatNumber            string `json:"vatNumber,"`
	BankName             string `json:"bankName,"`
	BankAccountNumber    string `json:"bankAccountNumber,"`
	BankIBAN             string `json:"bankIBAN,"`
	BankSWIFT            string `json:"bankSWIFT,"`
	PaymentDays          int    `json:"paymentDays,"`
	Notes                string `json:"notes,"`
	LastModified         int    `json:"lastModified,"`
	CustomerType         string `json:"customerType,"`

	Address string `json:"address,"`
	//CustomerAddresses    sharedCommon.Addresses `json:"addresses,"`
	Street     string `json:"street,"`
	Address2   string `json:"address2,"`
	City       string `json:"city,"`
	PostalCode string `json:"postalCode,"`
	Country    string `json:"country,"`
	State      string `json:"state,"`

	//ContactPersons       ContactPersons         `json:"contactPersons"`
	// Web-shop related fields
	//Username  string `json:"webshopUsername"`
	//LastLogin string `json:"webshopLastLogin"`

	Dirty bool `json:"-"`
}
