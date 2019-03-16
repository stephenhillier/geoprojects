import React from 'react';
import ReactPDF, { Document, Page, View, Text, Image, Font, StyleSheet } from '@react-pdf/renderer'

import Summary from '../components/Summary'
import TitleBlockBottom from '../components/TitleBlockBottom';

const PROTO_PATH = __dirname + '/../../plotsvc/proto/plotsvc/plotsvc.proto';
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');
// Suggested options for similarity to existing grpc.load behavior
const packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
// The protoDescriptor object has the full package hierarchy
const sieveplot = protoDescriptor.plotsvc;

const summary = {
  project: "Esquimalt Towers",
  client: "EXAMPLE INDUSTRIES, INC.",
  clientAddress: "Esquimalt, BC",
  projectNo: "A4321",
  location: "Esquimalt, BC",
  sample: "BH1 - SA1",
  easting: '777777',
  northing: '5555555',
  zone: '7',
  elevation: '2 m',
  date: 'Jan 1, 2018',
  totalDepth: '10 m'
}

const sieves = [
  {
    size: 20,
    passing: 10,
  },
]

const styles = StyleSheet.create({
  body: {
    fontSize: 10,
  },
  text: {
  },
  headerSection: {
    height: 190,
  },
  footerSection: {
    height: 100,
    position: 'absolute',
    bottom: 25,
    width: '100%',
    borderLeft: 1,
    borderRight: 1,
    borderBottom: 1,
    left: 25,
    right: 0
  }
});

Font.register(`${__dirname}/NotoSans-Regular.ttf`, { family: 'Noto Sans' });

const SieveReport = (props) => (
  <Document>
    <Page style={{padding: 25, paddingBottom: 125, fontSize: 10}} size="Letter" wrap>
      <View style={styles.headerSection} fixed>
        <Summary
          client={props.client}
          project={props.project}
          projectNumber={props.projectNo}
          location={props.location}
          easting={props.easting}
          northing={props.northing}
          zone={props.zone}
          elevation={props.elevation}
          date={props.date}
          boreholeNum={props.boreholeName}
          fixed
        />
      </View>
      <View wrap>
      </View>
      <Image src={`data:image/png;base64, ${props.figure}`}/>

      <View style={styles.footerSection} fixed>
        <TitleBlockBottom
          date={props.date}
          client={props.client}
          clientAddress={props.clientAddress}
          totalDepth={props.totalDepth}
          fixed
        />
      </View>
    </Page>
  </Document>
);

const client = new sieveplot.SievePlot('localhost:50051', grpc.credentials.createInsecure())

client.plotSieve({}, (err, fig) => {
  if (err || !fig || !fig.ok || !fig.figure ) {
    console.log('error!')
  } else {
    ReactPDF.render(<SieveReport figure={fig.figure} />, `${__dirname}/example.pdf`);
  }
})


export default SieveReport