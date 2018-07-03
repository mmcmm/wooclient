package woocommerce 

import (
	"net/url"
)

func addParams(paramsIn url.Values, paramsOut url.Values) url.Values {
	if paramsIn.Get("page") != "" {
		paramsOut.Add("page", paramsIn.Get("page"))
	}
	if paramsIn.Get("per_page") != "" {
		paramsOut.Add("per_page", paramsIn.Get("per_page"))
	}
	if paramsIn.Get("after") != "" {
		paramsOut.Add("after", paramsIn.Get("after"))
	}
	if paramsIn.Get("status") != "" {
		paramsOut.Add("status", paramsIn.Get("status"))
	}
	if paramsIn.Get("order") != "" {
		paramsOut.Add("order", paramsIn.Get("order"))
	}

	return paramsOut
}

func addURLParams(params *url.Values, specURL string) string {
	if params.Get("page") != "" {
		specURL += "&page=" + params.Get("page") 
	}
	if params.Get("per_page") != "" {
		specURL += "&per_page=" + params.Get("per_page") 
	}
	if params.Get("after") != "" {
		specURL += "&after=" + params.Get("after") 
	}
	if params.Get("status") != "" {
		specURL += "&status=" + params.Get("status") 
	}
	if params.Get("order") != "" {
		specURL += "&order=" + params.Get("order") 
	}

	return specURL
}