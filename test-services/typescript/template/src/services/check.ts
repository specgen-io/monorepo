import * as service from '../spec/check'
import * as models from '../spec/models'
import * as errors from '../spec/errors'

export let checkService = (): service.CheckService => {
    let checkEmpty = async (): Promise<void> => {}

    let checkEmptyResponse = async (params: service.CheckEmptyResponseParams): Promise<void> => {}
    
    let checkForbidden = async (): Promise<service.CheckForbiddenResponse> => {
        return {status: 'forbidden'}
    }

    let sameOperationName = async (): Promise<service.SameOperationNameResponse> => {
        return {status: 'ok'}
    }

    let checkBadRequest = async (): Promise<service.CheckBadRequestResponse> => {
        return {status: 'bad_request', data: {message: 'Error returned from service implementation', location: errors.ErrorLocation.UNKNOWN}}
    }

    return {checkEmpty, checkEmptyResponse, checkForbidden, sameOperationName, checkBadRequest}
}