import React from 'react'
import ReactPDF from '@react-pdf/renderer'
import BoreholeLog from './boreholeLog.js'
import express from 'express'
import http from 'http'
import axios from 'axios'
const jwt = require('express-jwt');
// const jwtAuthz = require('express-jwt-authz');
const jwksRsa = require('jwks-rsa');
const cookieParser = require('cookie-parser')
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');

const app = express();
app.use(cookieParser())

// Adapted from original code by Inventid https://github.com/inventid/pedeffy/blob/master/src/index.js
// and https://github.com/diegomura/react-pdf/issues/269

const renderReact = async (component, data) => {
	const rootElemComponent = React.createElement(component, data);
	return await ReactPDF.renderToStream(rootElemComponent);
};

const createPdf = async (reactTemplate, data, response) => {
  const started = new Date();
  try {
    response.set('Content-Type', "application/pdf");
    const readStream = await renderReact(reactTemplate, data);
    readStream.pipe(response);
    // When the stream end the response is closed as well
    readStream.on('end', () => console.log(`[logrend] Rendered in ${new Date() - started}ms`));
  } catch (e) {
    console.log(`Error occurred while rendering: "${e}"`);
    response.status(500).end();
  }
};


const getTokenFromRequest = (req) => {
  if (req.headers.authorization && req.headers.authorization.split(' ')[0] === 'Bearer') {
      return req.headers.authorization.split(' ')[1];
  } else if (req.query && req.query.token) {
    return req.query.token;
  } else if (req.cookies && req.cookies.access_token) {
    return req.cookies.access_token
  }
  return null;
}

const checkJwt = jwt({
  // Dynamically provide a signing key
  // based on the kid in the header and 
  // the signing keys provided by the JWKS endpoint.
  secret: jwksRsa.expressJwtSecret({
    cache: true,
    rateLimit: true,
    jwksRequestsPerMinute: 1,
    jwksUri: `https://earthworks.auth0.com/.well-known/jwks.json`
  }),


  getToken: getTokenFromRequest,

  // Validate the audience and the issuer.
  audience: 'https://earthworks.islandcivil.com',
  issuer: `https://earthworks.auth0.com/`,
  algorithms: ['RS256']
});

app.get('/logs/:projectID/boreholes/:boreholeID/:boreholeSlug.pdf', checkJwt, async function(req, res) {
  const { boreholeID, projectID } = req.params;
  const token = getTokenFromRequest(req)
  let project
  let borehole
  let soils

  if (!token) {
    res.status(500).send('Unable to parse authentication token');
  }

  const client = axios.create({
    baseURL: `http://${process.env.PROJECTS_SERVICE}/api/v1`,
    headers: { 'Authorization': `Bearer ${token}` }
  })

  const bhReq = client.get(`/boreholes/${boreholeID}`)
  const prReq = client.get(`/projects/${projectID}`)
  const soilReq = client.get(`/boreholes/${boreholeID}/strata`)

  Promise.all([bhReq, prReq, soilReq]).then((values) => {
    borehole = values[0].data
    project = values[1].data
    soils = values[2].data

    if (borehole.project != project.id) {
      res.status(401).send('Invalid borehole for this project')
    }

    const data = {
      project: project.name,
      client: project.client,
      clientAddress: '',
      projectNo: project.number,
      location: project.location,
      easting: borehole.location[1],
      northing: borehole.location[0],
      zone: '7',
      elevation: '',
      company: 'Acme Drilling Co.',
      method: 'Auger',
      date: borehole.end_date,
      totalDepth: '',
      boreholeName: borehole.name,
      soils: soils.map((item) => {
        return {
          from: item.start,
          to: item.end,
          desc: item.description
        }
      })
    }
  
    createPdf(BoreholeLog, data, res).then((output) => {
      return output
    }).catch((e) => {
      console.error(e)
      res.status(500).send('error rendering borehole log')
    })
  }).catch((e) => {
    res.status(401).send('error collecting borehole data')
  })
});

app.get('/logs/:projectID/sieves/:sieveID/:sieveSlug.pdf', checkJwt, async function(req, res) {

  const PROTO_PATH = __dirname + '/../plotsvc/proto/plotsvc/plotsvc.proto';

  const packageDefinition = protoLoader.loadSync(
      PROTO_PATH,
      {keepCase: true,
       longs: String,
       enums: String,
       defaults: true,
       oneofs: true
      });
  const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
  const sieveplot = protoDescriptor.plotsvc;


  const { sieveID, projectID } = req.params;
  const token = getTokenFromRequest(req)
  let project
  let testData
  let figure

  if (!token) {
    res.status(500).send('Unable to parse authentication token');
  }

  const client = axios.create({
    baseURL: `http://${process.env.PROJECTS_SERVICE}/api/v1`,
    headers: { 'Authorization': `Bearer ${token}` }
  })

  const prReq = client.get(`/projects/${projectID}`)
  const labTestReq = client.get(`/projects/${projectID}/lab/tests/${sieveID}`)

  const plotClient = new sieveplot.SievePlot(`${process.env.PLOTS_SERVICE || 'localhost'}:50051`, grpc.credentials.createInsecure())

  Promise.all([prReq, labTestReq]).then((values) => {
    project = values[0].data
    testData = values[1].data
  

    const reportData = {
      date: Date.now(),
    }

    console.log(`${process.env.PLOTS_SERVICE || 'localhost'}:50051`)

    // this isn't pretty, but the node grpc implementation does not have good async support
    plotClient.plotSieve({}, (err, fig) => {
      if (err || !fig || !fig.ok || !fig.figure ) {
        console.error(err, fig)
        console.log((err || !fig || !fig.ok || !fig.figure ))
        res.status(500).send('error rendering plot')    
      } else {

        reportData.figure = fig.figure
        createPdf(BoreholeLog, reportData, res).then((output) => {
          return output
        }).catch((e) => {
          console.error(e)
          res.status(500).send('error rendering report')
        })
    
      }
    })


  }).catch((e) => {
    console.error(e)
    res.status(401).send('error collecting report data')
  })
});

http.createServer(app).listen(8081);
