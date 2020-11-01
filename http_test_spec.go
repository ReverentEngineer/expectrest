package main

/*
#include <rdg.h>
#include <stdlib.h>
#cgo LDFLAGS: -lrdg
*/
import "C"

import (
	"errors"
	"unsafe"
)

type HTTPTestSpec struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    *string           `json:"body,omitempty"`

	Expect struct {
		Code    *int               `json:"code,omitempty"`
		Headers *map[string]string `json:"headers,omitempty"`
		Body    *string            `json:"body,omitempty"`
	} `json:"expect"`
}

func generatePermutations(expression string) ([]string, error) {
	permutations := make([]string, 0)

	cString := C.CString(expression)
	rdg := C.rdg_new(cString)

	if rdg == nil {
		C.free(unsafe.Pointer(cString))
		return nil, errors.New("Invalid expression")
	}

	cResultSize := C.size_t(0)
	cResult := (*C.uchar)(unsafe.Pointer(nil))

	for C.rdg_generate(&cResult, &cResultSize, rdg) != 0 {
		result := C.GoString((*C.char)(unsafe.Pointer(cResult)))
		permutations = append(permutations, result)
	}

	C.rdg_free(rdg)
	C.free(unsafe.Pointer(cString))
	return permutations, nil
}

func ExpandHTTPTestSpecs(testSpecs []HTTPTestSpec) ([]HTTPTestSpec, error) {
	results := make([]HTTPTestSpec, 0)

	for _, testSpec := range testSpecs {

		methods, err := generatePermutations(testSpec.Method)

		if err != nil {
			return nil, err
		}

		var bodies []string
		if testSpec.Body != nil {
			bodies, err = generatePermutations(*testSpec.Body)
			if err != nil {
				return nil, err
			}
		}

		for _, method := range methods {
			if bodies != nil {
				for _, body := range bodies {
					result := HTTPTestSpec{
						Url:    testSpec.Url,
						Method: method,
						Body:   &body,
					}
					results = append(results, result)
				}
			} else {
				result := HTTPTestSpec{
					Url:    testSpec.Url,
					Method: method,
				}
				results = append(results, result)
			}
		}
	}

	return results, nil
}
