package old

var httpStatusCode = map[string]string{
	"continue":           "100",
	"switching_protocol": "101",
	"processing":         "102",

	"ok":                            "200",
	"created":                       "201",
	"accepted":                      "202",
	"non_authoritative_information": "203",
	"no_content":                    "204",
	"reset_content":                 "205",
	"partial_content":               "206",
	"multi_status":                  "207",
	"already_reported":              "208",
	"im_used":                       "226",

	"multiple_choices":   "300",
	"moved_permanently":  "301",
	"found":              "302",
	"see_other":          "303",
	"not_modified":       "304",
	"use_proxy":          "305",
	"switch_proxy":       "306",
	"temporary_redirect": "307",
	"permanent_redirect": "308",

	"bad_request":                     "400",
	"unauthorized":                    "401",
	"payment_required":                "402",
	"forbidden":                       "403",
	"not_found":                       "404",
	"method_not_allowed":              "405",
	"not_acceptable":                  "406",
	"proxy_authentication_required":   "407",
	"request_timeout":                 "408",
	"conflict":                        "409",
	"gone":                            "410",
	"length_required":                 "411",
	"precondition_failed":             "412",
	"request_entity_too_large":        "413",
	"request_uri_too_long":            "414",
	"unsupported_media_type":          "415",
	"requested_range_not_satisfiable": "416",
	"expectation_failed":              "417",
	"im_a_teapot":                     "418",
	"unprocessable_entity":            "422",
	"too_early":                       "425",
	"upgrade_required":                "426",
	"precondition_required":           "428",
	"too_many_requests":               "429",
	"request_header_fields_too_large": "431",
	"unavailable_for_legal_reasons":   "451",

	"internal_server_error":           "500",
	"not_implemented":                 "501",
	"bad_gateway":                     "502",
	"service_unavailable":             "503",
	"gateway_timeout":                 "504",
	"http_version_not_supported":      "505",
	"variant_also_negotiates":         "506",
	"insufficient_storage":            "507",
	"loop_detected":                   "508",
	"bandwidth_limit_exceeded":        "509",
	"not_extend":                      "510",
	"network_authentication_required": "511",
	"connection_timed_out":            "522",
	"proxy_declined_request":          "523",
	"timeout_occurred":                "524",
}

func HttpStatusCode(statusName Name) string {
	return httpStatusCode[statusName.SnakeCase()]
}

const (
	HttpStatusInternalServerError string = "internal_server_error"
	HttpStatusNotFound            string = "not_found"
	HttpStatusBadRequest          string = "bad_request"
)
