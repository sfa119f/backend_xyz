package dictionary

const (
	UndisclosedError	= "Something went wrong"
	NoError				= "None"
	NotFoundError		= "Not found"
	InvalidParamError	= "Invalid parameter"
)

type APIResponse struct {
	Data	 interface{}	`json:"data"`
	Detail interface{}	`json:"detail,omitempty"`
	Error	 string				`json:"error"`
}
