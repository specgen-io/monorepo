import * as echo from './test-service/echo';
import * as check from './test-service/check';

declare global {
  interface Window {
    echoClient: (config: {baseURL: string}) => echo.EchoClient    
    checkClient: (config: {baseURL: string}) => check.CheckClient    
  }
}

window.echoClient = echo.client
window.checkClient = check.client