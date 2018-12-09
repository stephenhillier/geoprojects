import React from 'react'
import ReactPDF from '@react-pdf/renderer'
import BoreholeLog from './boreholeLog.js'
import express from 'express'
import http from 'http'
const jwt = require('express-jwt');
// const jwtAuthz = require('express-jwt-authz');
const jwksRsa = require('jwks-rsa');

const app = express();

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
    readStream.on('end', () => console.log(`[logrend][borehole] Rendered ${data.boreholeNum} in ${new Date() - started}ms`));
  } catch (e) {
    console.log(`Error occurred while rendering: "${e}"`);
    response.status(500).end();
  }
};


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

  // Validate the audience and the issuer.
  audience: 'https://earthworks.islandcivil.com',
  issuer: `https://earthworks.auth0.com/`,
  algorithms: ['RS256']
});


app.get('/logs/boreholes/:boreholeNum.pdf', checkJwt, async function(req, res) {
  const { boreholeNum } = req.params;
  const data = { boreholeNum };
  return await createPdf(BoreholeLog, data, res);
});


http.createServer(app).listen(8081);
