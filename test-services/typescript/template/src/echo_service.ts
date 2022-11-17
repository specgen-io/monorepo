import * as service from './spec/echo'
import * as models from './spec/models'

export let echoService = (): service.EchoService => {
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

    let echoQuery = async (params: service.EchoQueryParams): Promise<models.Parameters> => {
        return {
            int_field: params.int_query,
            long_field: params.long_query,
            float_field: params.float_query,
            double_field: params.double_query,
            decimal_field: params.decimal_query,
            bool_field: params.bool_query,
            string_field: params.string_query,
            string_opt_field: params.string_opt_query,
            string_defaulted_field: params.string_defaulted_query,
            string_array_field: params.string_array_query,
            uuid_field: params.uuid_query,
            date_field: params.date_query,
            date_array_field: params.date_array_query,
            datetime_field: params.datetime_query,
            enum_field: params.enum_query,
        }
    }

    let echoHeader = async (params: service.EchoHeaderParams): Promise<models.Parameters> => {
        return {
            int_field: params['Int-Header'],
            long_field: params['Long-Header'],
            float_field: params['Float-Header'],
            double_field: params['Double-Header'],
            decimal_field: params['Decimal-Header'],
            bool_field: params['Bool-Header'],
            string_field: params['String-Header'],
            string_opt_field: params['String-Opt-Header'],
            string_defaulted_field: params['String-Defaulted-Header'],
            string_array_field: params['String-Array-Header'],
            uuid_field: params['Uuid-Header'],
            date_field: params['Date-Header'],
            date_array_field: params['Date-Array-Header'],
            datetime_field: params['Datetime-Header'],
            enum_field: params['Enum-Header'],
        }
    }

    let echoUrlParams = async (params: service.EchoUrlParamsParams): Promise<models.UrlParameters> => {
        return {
            int_field: params.int_url,
            long_field: params.long_url,
            float_field: params.float_url,
            double_field: params.double_url,
            decimal_field: params.decimal_url,
            bool_field: params.bool_url,
            string_field: params.string_url,
            uuid_field: params.uuid_url,
            date_field: params.date_url,
            datetime_field: params.datetime_url,
            enum_field: params.enum_url,
        }
    }

    let echoEverything = async (params: service.EchoEverythingParams): Promise<service.EchoEverythingResponse> => {
        return {
            status: "ok",
            data: {
                body_field: params.body,
                float_query: params.float_query,
                bool_query: params.bool_query,
                uuid_header: params['Uuid-Header'],
                datetime_header: params['Datetime-Header'],
                date_url: params.date_url,
                decimal_url: params.decimal_url,
            }
        }
    }

    let sameOperationName = async (): Promise<service.SameOperationNameResponse> => {
        return {status: 'ok'}
    }

    return {echoBodyString, echoBodyModel, echoBodyArray, echoBodyMap, echoQuery, echoHeader, echoUrlParams, echoEverything, sameOperationName}
}