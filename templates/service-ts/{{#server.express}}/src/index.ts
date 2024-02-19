import express from 'express'
import cors from 'cors'
var multipart = require('connect-multiparty')
{{#swagger.value}}
import swaggerUi from 'swagger-ui-express'
import yamljs from 'yamljs'
{{/swagger.value}}
import {specRouter} from './spec/spec_router'
import {services} from './services'

const app = express()
const port = 8081

app.use(cors())
app.use(express.json())
app.use(express.text())
app.use(express.urlencoded())
app.use(multipart())

{{#swagger.value}}
app.use("/docs/swagger.yaml", express.static("docs/swagger.yaml"))
app.use("/docs", swaggerUi.serve, swaggerUi.setup(yamljs.load("./docs/swagger.yaml")))
{{/swagger.value}}

app.use("/", specRouter(services))

app.listen(port, () => {
    console.log( `server started at http://localhost:${ port }` )
})
