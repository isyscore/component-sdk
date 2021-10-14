package isyscore

import "encoding/json"

type LicenseCustomer struct {
	EnterpriseName string `json:"enterprise_name"`
	ContactEmail   string `json:"contact_email"`
	ContactName    string `json:"contact_name"`
	ContactPhone   string `json:"contact_phone"`
}

type LicenseData struct {
	LicenseCode string           `json:"license_code"`
	Customer    *LicenseCustomer `json:"customer"`
}

type ComponentRegister struct {
	Name             string `json:"name"`
	ShowName         string `json:"showName"`
	Description      string `json:"description"`
	VersionCode      int    `json:"versionCode"`
	VersionName      string `json:"versionName"`
	IsOpenSource     int    `json:"isOpenSource"`
	IsEnabled        int    `json:"isEnabled"`
	IsUnderCarriage  int    `json:"isUnderCarriage"`
	CompactOsVersion string `json:"compactOsVersion"`
	ProducerCompany  string `json:"producerCompany"`
	ProducerContact  string `json:"producerContact"`
	ProducerEmail    string `json:"producerEmail"`
	ProducerPhone    string `json:"producerPhone"`
	ProducerUrl      string `json:"producerUrl"`
}

type ResultComponentRegister struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    *ComponentRegister `json:"data"`
}

type ComponentLicensed struct {
	IsRevoked      int   `json:"isRevoked"`
	IsTrial        int   `json:"isTrial"`
	TrialStartDate int64 `json:"trialStartDate"`
	TrialEndDate   int64 `json:"trialEndDate"`
}

type ResultComponentLicensed struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    *ComponentLicensed `json:"data"`
}

type ComponentLicense struct {
	LicenseName string `json:"licenseName"`
	LicenseText string `json:"licenseText"`
}

func (l *ComponentLicense) String() string {
	b, _ := json.Marshal(l)
	return string(b)
}

type ResultComponentLicense struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *ComponentLicense `json:"data"`
}

type ComponentProducer struct {
	Company string `json:"company"`
	Contact string `json:"contact"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Url     string `json:"url"`
}

func (p *ComponentProducer) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

type ResultComponentProducer struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    *ComponentProducer `json:"data"`
}
