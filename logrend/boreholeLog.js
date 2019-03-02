import React from 'react';
import ReactPDF, { Document, Page, View, Text, Font, StyleSheet } from '@react-pdf/renderer'

import Summary from './components/Summary'
import SoilLayerTable from './components/SoilLayerTable'
import TitleBlockBottom from './components/TitleBlockBottom';
import SoilTableHeader from './components/SoilTableHeader';

const summary = {
  project: "Esquimalt Towers",
  client: "EXAMPLE INDUSTRIES, INC.",
  clientAddress: "Esquimalt, BC",
  projectNo: "A4321",
  location: "Esquimalt, BC",
  easting: '777777',
  northing: '5555555',
  zone: '7',
  elevation: '2 m',
  company: 'Acme Drilling Co.',
  method: 'Auger',
  date: 'Jan 1, 2018',
  totalDepth: '10 m'
}

const soils = [
  {
    desc: 'SAND and GRAVEL, some silt',
    from: 0,
    to: 2.5
  },
  {
    desc: 'SAND, some gravel, some silt',
    from: 2.5,
    to: 9.5
  },
  {
    desc: 'Silty CLAY',
    from: 9.5,
    to: 11
  },
  {
    desc: 'SAND',
    from: 11,
    to: 15
  },
  {
    desc: 'CLAY',
    from: 15,
    to: 45
  }
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

const BoreholeLog = (props) => (
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
          company={props.company}
          method={props.method}
          date={props.date}
          boreholeNum={props.boreholeName}
          fixed
        />
        <SoilTableHeader fixed/>
      </View>
      <View wrap>
        <SoilLayerTable soils={props.soils} wrap/>
      </View>

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
ReactPDF.render(<BoreholeLog/>, `${__dirname}/example.pdf`);

export default BoreholeLog
