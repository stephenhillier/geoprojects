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
    to: 2
  },
  {
    desc: 'SAND, some gravel, some silt',
    from: 2,
    to: 6
  },
  {
    desc: 'Silty CLAY',
    from: 6,
    to: 9
  },
  {
    desc: 'SAND',
    from: 9,
    to: 13
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
          client={summary.client}
          project={summary.project}
          projectNumber={summary.projectNo}
          location={summary.location}
          easting={summary.easting}
          northing={summary.northing}
          zone={summary.zone}
          elevation={summary.elevation}
          company={summary.company}
          method={summary.method}
          date={summary.date}
          boreholeNum={props.boreholeNum}
          fixed
        />
        <SoilTableHeader fixed/>
      </View>
      <View wrap>
        <SoilLayerTable soils={soils} wrap/>
      </View>

      <View style={styles.footerSection} fixed>
          <TitleBlockBottom
            date={summary.date}
            client={summary.client}
            clientAddress={summary.clientAddress}
            totalDepth={summary.totalDepth}
            fixed
          />
      </View>
    </Page>
  </Document>
);
ReactPDF.render(<BoreholeLog/>, `${__dirname}/example.pdf`);

export default BoreholeLog
