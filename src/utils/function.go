package utils

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/sfa119f/backend_xyz/src/dictionary"
)

func JsonResp(w http.ResponseWriter, code int, data interface{}, err error) {
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: data, 
			Error: dictionary.NoError,
		})
	} else if code == 400 {
		fmt.Println("error message:", err)
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: nil, Error: err.Error(),
		})
	} else if err != nil {
		fmt.Println("error message:", err)
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: nil, Error: dictionary.UndisclosedError,
		})
	} else {
		fmt.Println("error message:", dictionary.UndisclosedError)
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: nil, Error: dictionary.UndisclosedError,
		})
	}
}
