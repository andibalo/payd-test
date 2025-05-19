package response

const (
	Success           Code = "RMS0000"
	ServerError       Code = "RMS0001"
	BadRequest        Code = "RMS0002"
	InvalidRequest    Code = "RMS0004"
	Failed            Code = "RMS0073"
	Pending           Code = "RMS0050"
	InvalidInputParam Code = "RMS0032"
	DuplicateUser     Code = "RMS0033"
	NotFound          Code = "RMS0034"

	Unauthorized   Code = "RMS0502"
	Forbidden      Code = "RMS0503"
	GatewayTimeout Code = "RMS0048"
)

type Code string

var codeMap = map[Code]string{
	Success:           "success",
	Failed:            "failed",
	Pending:           "pending",
	BadRequest:        "bad or invalid request",
	Unauthorized:      "Unauthorized Token",
	GatewayTimeout:    "Gateway Timeout",
	ServerError:       "Internal Server Error",
	InvalidInputParam: "Other invalid argument",
	DuplicateUser:     "duplicate user",
	NotFound:          "Not found",
}

func (c Code) AsString() string {
	return string(c)
}

func (c Code) GetStatus() string {
	switch c {
	case Success:
		return "SUCCESS"

	default:
		return "FAILED"
	}
}

func (c Code) GetMessage() string {
	return codeMap[c]
}

func (c Code) GetVersion() string {
	return "1"
}
