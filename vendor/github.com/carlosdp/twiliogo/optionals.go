package twiliogo

type Optional interface {
	GetParam() (string, string)
}

type Callback string

func (callback Callback) GetParam() (string, string) {
	return "Url", string(callback)
}

type ApplicationSid string

func (applicationSid ApplicationSid) GetParam() (string, string) {
	return "ApplicationSid", string(applicationSid)
}

type Method string

func (method Method) GetParam() (string, string) {
	return "Method", string(method)
}

type FallbackUrl string

func (fallbackUrl FallbackUrl) GetParam() (string, string) {
	return "FallbackUrl", string(fallbackUrl)
}

type FallbackMethod string

func (fallbackMethod FallbackMethod) GetParam() (string, string) {
	return "FallbackMethod", string(fallbackMethod)
}

type StatusCallback string

func (statusCallback StatusCallback) GetParam() (string, string) {
	return "StatusCallback", string(statusCallback)
}

type StatusCallbackMethod string

func (statusCallbackMethod StatusCallbackMethod) GetParam() (string, string) {
	return "StatusCallbackMethod", string(statusCallbackMethod)
}

type SendDigits string

func (sendDigits SendDigits) GetParam() (string, string) {
	return "SendDigits", string(sendDigits)
}

type IfMachine string

func (ifMachine IfMachine) GetParam() (string, string) {
	return "IfMachine", string(ifMachine)
}

type Timeout string

func (timeout Timeout) GetParam() (string, string) {
	return "Timeout", string(timeout)
}

type Record string

func (record Record) GetParam() (string, string) {
	return "Record", string(record)
}

type To string

func (to To) GetParam() (string, string) {
	return "To", string(to)
}

type From string

func (from From) GetParam() (string, string) {
	return "From", string(from)
}

type Status string

func (status Status) GetParam() (string, string) {
	return "Status", string(status)
}

type StartTime string

func (startTime StartTime) GetParam() (string, string) {
	return "StartTime", string(startTime)
}

type ParentCallSid string

func (parentCallSid ParentCallSid) GetParam() (string, string) {
	return "ParentCallSid", string(parentCallSid)
}

type DateSent string

func (dateSent DateSent) GetParam() (string, string) {
	return "DateSent", string(dateSent)
}

type Body string

func (body Body) GetParam() (string, string) {
	return "Body", string(body)
}

type MediaUrl string

func (mediaUrl MediaUrl) GetParam() (string, string) {
	return "MediaUrl", string(mediaUrl)
}

type FriendlyName string

func (friendlyName FriendlyName) GetParam() (string, string) {
	return "FriendlyName", string(friendlyName)
}

type PhoneNumber string

func (phoneNumber PhoneNumber) GetParam() (string, string) {
	return "PhoneNumber", string(phoneNumber)
}

type AreaCode string

func (areaCode AreaCode) GetParam() (string, string) {
	return "AreaCode", string(areaCode)
}

type MessagingServiceSid string

func (sid MessagingServiceSid) GetParam() (string, string) {
	return "MessagingServiceSid", string(sid)
}
