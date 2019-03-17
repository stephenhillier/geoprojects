import React from 'react';
import ReactPDF, { Document, Page, View, Text, Image, Font, StyleSheet } from '@react-pdf/renderer'

import Header from './Header'
import Summary from './Summary'
import TitleBlock from './TitleBlock';

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
  bodySection: {
    borderLeft: 1,
    borderRight: 1,
    borderTop: 1,
    flex: 0
  },
  summarySection: {
    borderLeft: 1,
    borderRight: 1,
    borderTop: 1,
    flex: 1
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
    right: 0,
    flex: 0
  }
});

Font.register(`${__dirname}/NotoSans-Regular.ttf`, { family: 'Noto Sans' });

const SieveReport = (props) => (
  <Document>
    <Page style={{padding: 25, paddingBottom: 125, fontSize: 10, flexDirection: 'column'}} size="Letter" wrap>
      <View fixed style={{flex: 0}}>
        <Header></Header>
      </View>
      <View style={styles.bodySection}>
        <Image src={`data:image/png;base64, ${props.figure}`}/>
      </View>
      <View style={styles.summarySection}>
        <Summary>
        </Summary>
      </View>
      <View style={styles.footerSection} fixed>
        <TitleBlock
          date={props.date}
          client={props.client}
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