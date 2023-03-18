package dictionary

const (
	UndisclosedError		= "something went wrong"
	NoError							= "none"
	NotFoundError				= "not found"
	InvalidParamError		= "invalid parameter"
	InvalidRequestError	= "invalid request"
)

type APIResponse struct {
	Data	 interface{}	`json:"data"`
	Detail interface{}	`json:"detail,omitempty"`
	Error	 string				`json:"error"`
}
