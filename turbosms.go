package sachet

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

// sachet section
type TurbosmsConfig struct {
	Alogin    string `yaml:"login"`
	Apassword string `yaml:"password"`
}
type Turbosms struct {
	Login    string
	Password string
}

func NewTurbosms(config TurbosmsConfig) *Turbosms {
	Turbosms := &Turbosms{Login: config.Alogin, Password: config.Apassword}
	return Turbosms
}

// http  url for  turbosms
var urlSoap string = "http://turbosms.in.ua/api/soap.html"

type SoapBody struct {
	Contents []byte `xml:",innerxml"`
}

type SoapEnvelopeResponse struct {
	XMLName struct{} `xml:"Envelope"`
	Id1     string   `xml:"xmlns:SOAP-ENV,attr"`
	Id2     string   `xml:"xmlns:ns1,attr"`
	Body    SoapBody
}

type getAuthResponse struct {
	XMLName    struct{} `xml:"AuthResponse"`
	Result     string   `xml:"AuthResult"`
	StatusCode int
}

type SoapEnvelopeReqest struct {
	XMLName struct{} `xml:"SOAP-ENV:Envelope"`
	Id1     string   `xml:"xmlns:SOAP-ENV,attr"`
	Id2     string   `xml:"xmlns:ns1,attr"`
	Body    SoapBody `xml:"SOAP-ENV:Body"`
}

type getSendsmsRequest struct {
	XMLName     struct{} `xml:"ns1:SendSMS"`
	Sender      string   `xml:"ns1:sender"`
	Destination string `xml:"ns1:destination"`
	Text        string   `xml:"ns1:text"`
	Wappush     string   `xml:"ns1:wappush"`
}

type getAuthRequest struct {
	XMLName  struct{} `xml:"ns1:Auth"`
	User     string   `xml:"ns1:login"`
	Password string   `xml:"ns1:password"`
}

func SoapEncode(contents interface{}) ([]byte, error) {
	data, err := xml.MarshalIndent(contents, "    ", "  ")
	if err != nil {
		return nil, err
	}
	data = append([]byte("\n"), data...)
	env := SoapEnvelopeReqest{Id1: "http://schemas.xmlsoap.org/soap/envelope/", Id2: "http://turbosms.in.ua/api/Turbo", Body: SoapBody{Contents: data}}
	return xml.MarshalIndent(&env, "", "  ")
}

func SoapDecode(data []byte, contents interface{}) error {
	env := SoapEnvelopeResponse{Body: SoapBody{}}
	err := xml.Unmarshal(data, &env)
	if err != nil {
		return err
	}
	return xml.Unmarshal(env.Body.Contents, contents)
}

func Request(c *http.Client, url string, payload []byte) ([]byte, error, int) {

	resp, err := c.Post(url, "text/xml", bytes.NewBuffer(payload))
	statuscode := resp.StatusCode
	if err != nil {
		return nil, err, statuscode
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err, statuscode
	}
	return body, nil, statuscode
}

func (c *Turbosms) Send(message Message) (err error) {
	////// Encode Auth ////////////
	req := &getAuthRequest{User: c.Login, Password: c.Password}
	data, err := SoapEncode(&req)
	if err != nil {
		return err
	}
	cookieJar, _ := cookiejar.New(nil)
	clientConfig := &http.Client{
		Timeout: 15 * time.Second,
		Jar:     cookieJar,
	}
	reply, err, statusreply := Request(clientConfig, urlSoap, []byte(data))
	if err != nil {
		return err
	}
	var allrecipent string
	//	for _, recipent := range message.To {
	//	       allrecipent += recipent
	//	}

	/////  Encode SendSms /////////
	sms := &getSendsmsRequest{Sender: message.From, Destination: strings.Join(message.To, ","), Text: message.Text, Wappush: ""}
	datasms, err := SoapEncode(&sms)
	if err != nil {
		return err
	}
	////// Request ///////////
	replysms, err, statusreplysms := Request(clientConfig, urlSoap, []byte(datasms))
	if err != nil {
		return err

	}
	if statusreply == 200 && statusreplysms == 200 && err == nil {
		return nil
	}

	var resp getAuthResponse
	err = SoapDecode([]byte(reply), &resp)
	if err != nil {
		return err
	}

	return fmt.Errorf("Failed sending sms. Reason: %s, statusCode: %d", string(replysms), statusreplysms)
}
