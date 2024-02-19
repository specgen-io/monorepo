import * as service from '../spec/echo'
import * as models from '../spec/models'

export let echoService = () => {
    let echoBodyString = async (params: service.EchoBodyStringParams): Promise<string> => {
        return params.body
    }

    let echoBodyModel = async (params: service.EchoBodyModelParams): Promise<models.Message> => {
        return {int_field: params.body.int_field, string_field: params.body.string_field}
    }

    let echoBodyArray = async (params: service.EchoBodyArrayParams): Promise<string[]> => {
        return params.body
    }

    let echoBodyMap = async (params: service.EchoBodyMapParams): Promise<Record<string, string>> => {
        return params.body
    }

    let echoFormData = async (params: service.EchoFormDataParams): Promise<models.FormParameters> => {
        return {
            int_field: params.intParam,
            long_field: params.longParam,
            float_field: params.floatParam,
            double_field: params.doubleParam,
            decimal_field: params.decimalParam,
            bool_field: params.boolParam,
            string_field: params.stringParam,
            string_opt_field: params.stringOptParam,
            string_defaulted_field: params.stringDefaultedParam,
            string_array_field: params.stringArrayParam,
            uuid_field: params.uuidParam,
            date_field: params.dateParam,
            date_array_field: params.dateArrayParam,
            datetime_field: params.datetimeParam,
            enum_field: params.enumParam,
        }
    }

    let echoFormUrlencoded = async (params: service.EchoFormUrlencodedParams): Promise<models.FormParameters> => {
        return {
            int_field: params.intParam,
            long_field: params.longParam,
            float_field: params.floatParam,
            double_field: params.doubleParam,
            decimal_field: params.decimalParam,
            bool_field: params.boolParam,
            string_field: params.stringParam,
            string_opt_field: params.stringOptParam,
            string_defaulted_field: params.stringDefaultedParam,
            string_array_field: params.stringArrayParam,
            uuid_field: params.uuidParam,
            date_field: params.dateParam,
            date_array_field: params.dateArrayParam,
            datetime_field: params.datetimeParam,
            enum_field: params.enumParam,
        }
    }

    let echoQuery = async (params: service.EchoQueryParams): Promise<models.Parameters> => {
        return {
            int_field: params.intQuery,
            long_field: params.longQuery,
            float_field: params.floatQuery,
            double_field: params.doubleQuery,
            decimal_field: params.decimalQuery,
            bool_field: params.boolQuery,
            string_field: params.stringQuery,
            string_opt_field: params.stringOptQuery,
            string_defaulted_field: params.stringDefaultedQuery,
            string_array_field: params.stringArrayQuery,
            uuid_field: params.uuidQuery,
            date_field: params.dateQuery,
            date_array_field: params.dateArrayQuery,
            datetime_field: params.datetimeQuery,
            enum_field: params.enumQuery,
        }
    }

    let echoHeader = async (params: service.EchoHeaderParams): Promise<models.Parameters> => {
        return {
            int_field: params.intHeader,
            long_field: params.longHeader,
            float_field: params.floatHeader,
            double_field: params.doubleHeader,
            decimal_field: params.decimalHeader,
            bool_field: params.boolHeader,
            string_field: params.stringHeader,
            string_opt_field: params.stringOptHeader,
            string_defaulted_field: params.stringDefaultedHeader,
            string_array_field: params.stringArrayHeader,
            uuid_field: params.uuidHeader,
            date_field: params.dateHeader,
            date_array_field: params.dateArrayHeader,
            datetime_field: params.datetimeHeader,
            enum_field: params.enumHeader,
        }
    }

    let echoUrlParams = async (params: service.EchoUrlParamsParams): Promise<models.UrlParameters> => {
        return {
            int_field: params.intUrl,
            long_field: params.longUrl,
            float_field: params.floatUrl,
            double_field: params.doubleUrl,
            decimal_field: params.decimalUrl,
            bool_field: params.boolUrl,
            string_field: params.stringUrl,
            uuid_field: params.uuidUrl,
            date_field: params.dateUrl,
            datetime_field: params.datetimeUrl,
            enum_field: params.enumUrl,
        }
    }

    let echoEverything = async (params: service.EchoEverythingParams): Promise<service.EchoEverythingResponse> => {
        return {
            status: "ok",
            data: {
                body_field: params.body,
                float_query: params.floatQuery,
                bool_query: params.boolQuery,
                uuid_header: params.uuidHeader,
                datetime_header: params.datetimeHeader,
                date_url: params.dateUrl,
                decimal_url: params.decimalUrl,
            }
        }
    }

    let sameOperationName = async (): Promise<service.SameOperationNameResponse> => {
        return {status: 'ok'}
    }

    return {echoBodyString, echoBodyModel, echoBodyArray, echoBodyMap, echoFormData, echoFormUrlencoded, echoQuery, echoHeader, echoUrlParams, echoEverything, sameOperationName}
}