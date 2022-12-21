import {echoService} from './echo'
import {checkService} from './check'
import {echoService as echoServiceV2} from './v2/echo'

export const services = {echoService: echoService(), checkService: checkService(), echoServiceV2: echoServiceV2()}
