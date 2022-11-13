import {echoService} from './echo_service'
import {checkService} from './check_service'
import {echoService as echoServiceV2} from './v2/echo_service'

export const services = {echoService: echoService(), checkService: checkService(), echoServiceV2: echoServiceV2()}
