import Koa from 'koa'
import bodyParser from 'koa-bodyparser'
import Router from '@koa/router'
{{#swagger.value}}
import {koaSwagger} from 'koa2-swagger-ui'
import yamljs from 'yamljs'
{{/swagger.value}}
import {specRouter} from './spec/spec_router'
import {services} from './services'

const app = new Koa()
const port = 8081

app.use(bodyParser({enableTypes: ['json', 'form', 'text']}))

{{#swagger.value}}
let docsRouter = new Router()
docsRouter.get('/docs', koaSwagger({swaggerOptions: {spec: yamljs.load("./docs/swagger.yaml")}}))
app.use(docsRouter.routes())
{{/swagger.value}}

let spec = specRouter(services)
app.use(spec.routes()).use(spec.allowedMethods())

app.listen(port, () => {
    console.log( `server started at http://localhost:${ port }` )
})
