package flexi

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
)

const (
	// PaymentURL - URL платежной страницы
	PaymentURL = "https://pay.flexidengi.ru/"

	// PaymentMethodWebmoneyWmr - WebMoney WMR
	PaymentMethodWebmoneyWmr = 38
	// PaymentMethodQiwi - QIWI
	PaymentMethodQiwi = 44
	// PaymentMethodTest - Тестовый(возвращает произвольный статус оплаты, не передается в реестрах)
	PaymentMethodTest = 45
	// PaymentMethodMobile - Мобильная коммерция
	PaymentMethodMobile = 60
	// PaymentMethodCreditCards -Банковские карты VISA / VISA Electron / MasterCard / Maestro
	PaymentMethodCreditCards = 61

	// PaymentStatusProcessed - транзакция выполнена
	PaymentStatusProcessed = "PROCESSED"
	// PaymentStatusFailed - транзакция провалена
	PaymentStatusFailed = "FAILED"
)

// Flexi - базовый объект работы с платежами
type Flexi struct {
	ServiceID int
	SecretKey string

	OrderID    int
	ProductID  int
	CustomerID string
	Summ       int
	Currency   string
	Count      int

	PaymentMethodID int
}

// NewFlexi - создание объекта работы с оплатой
func NewFlexi(serviceID int, secretKey string) *Flexi {

	return &Flexi{
		ServiceID: serviceID,
		SecretKey: secretKey,
	}
}

// Sign - формиование хеша параметров
func (f *Flexi) Sign() string {

	//md5(service_id+customer_id+order_id+product_id+ summ+currency+count+payment_method_id+ secret_key)

	baseString := []string{
		// обязательные полядля хеша
		strconv.Itoa(f.ServiceID),
		f.CustomerID,
		strconv.Itoa(f.OrderID),
		strconv.Itoa(f.ProductID),
	}

	{ // дополнительные элементы для хеша
		if f.Summ != 0 {
			baseString = append(baseString, strconv.Itoa(f.Summ))
		}
		if f.Currency != "" {
			baseString = append(baseString, f.Currency)
		}
		if f.Count != 0 {
			baseString = append(baseString, strconv.Itoa(f.Count))
		}
		if f.PaymentMethodID != 0 {
			baseString = append(baseString, strconv.Itoa(f.PaymentMethodID))
		}
	}

	baseString = append(baseString, f.SecretKey)

	data := []byte(strings.Join(baseString, ""))
	return fmt.Sprintf("%x", md5.Sum(data))
}

// SetOrderID -
func (f *Flexi) SetOrderID(orderID int) *Flexi {

	f.OrderID = orderID
	return f
}

// SetProductID -
func (f *Flexi) SetProductID(productID int) *Flexi {

	f.ProductID = productID
	return f
}

// SetCustomerID -
func (f *Flexi) SetCustomerID(sustomerID string) *Flexi {

	f.CustomerID = sustomerID
	return f
}

// SetSumm -
func (f *Flexi) SetSumm(summ int) *Flexi {

	f.Summ = summ
	return f
}

// SetCurrency -
func (f *Flexi) SetCurrency(currency string) *Flexi {

	f.Currency = currency
	return f
}

// SetCount -
func (f *Flexi) SetCount(count int) *Flexi {

	f.Count = count
	return f
}

// SetPaymentMethodID -
func (f *Flexi) SetPaymentMethodID(paymentMethodID int) *Flexi {

	f.PaymentMethodID = paymentMethodID
	return f
}

// MakeForm - собирает форму, автоматически отправляющую данные на оплату
func (f *Flexi) MakeForm() string {

	element := `<input type="hidden" name="%s" value="%v" />`

	elementsList := []string{}

	for name, value := range f.GetElements() {
		elementsList = append(elementsList, fmt.Sprintf(element, name, value))
	}

	html := `<!DOCTYPE html>
<html>
    <head><title>Redirecting...</title></head>
    <body onload="document.forms[0].submit();">
        <p>Redirecting to payment page...</p>
        <form action="%s" method="post">%s</form>
    </body>
</html>`

	return fmt.Sprintf(html, PaymentURL, strings.Join(elementsList, "\n"))
}

// GetElements - выдаёт массив полей и их значений с для генерации форм
func (f *Flexi) GetElements() map[string]string {

	elementsList := map[string]string{
		// обязательные элементы
		"service_id":  strconv.Itoa(f.ServiceID),
		"customer_id": f.CustomerID,
		"order_id":    strconv.Itoa(f.OrderID),
		"product_id":  strconv.Itoa(f.ProductID),
		"hash":        f.Sign(),
	}

	// необязательные элементы
	{
		if f.Summ != 0 {
			elementsList["sum"] = strconv.Itoa(f.Summ)
		}
		if f.Currency != "" {
			elementsList["currency"] = f.Currency
		}
		if f.Count != 0 {
			elementsList["count"] = strconv.Itoa(f.Count)
		}
		if f.PaymentMethodID != 0 {
			elementsList["payment_method_id"] = strconv.Itoa(f.PaymentMethodID)
		}
	}

	return elementsList
}
