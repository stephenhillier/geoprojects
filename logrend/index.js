import React from 'react';
import ReactPDF, { Document, Page, View, Text, Font, StyleSheet } from '@react-pdf/renderer'

import Heading from './components/Heading'
import Summary from './components/Summary'
import TableHeader from './components/TableHeader'

const summary = {
  project: "Esquimalt Towers",
  client: "Bigtime Engineering",
  projectNo: "A4321",
  location: "Esquimalt, BC",
  easting: '777777',
  northing: '5555555',
  zone: '7',
  elevation: '2 m',
  company: 'Acme Drilling Co.',
  method: 'Auger',
  date: 'Jan 1, 2018'
}

const styles = StyleSheet.create({
  body: {
    fontFamily: 'Noto Sans',
    fontSize: 12,
    height: '100%',
    border: 1
  },
  page: {
    padding: 20,
  },
  text: {
    marginLeft: 20,
    marginRight: 20,
    marginTop: 20,
    marginBottom: 10
  },
  headerSection: {
    margin: 10
  }
});

Font.register(`${__dirname}/NotoSans-Regular.ttf`, { family: 'Noto Sans' });

const Doc = () => (
  <Document>
    <Page size="Letter" style={styles.page}>
      <View style={styles.body}>
        <View style={styles.headerSection}>
          <Heading/>
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
          />
        </View>
        <TableHeader/>
      </View>
    </Page>
  </Document>
);
ReactPDF.render(<Doc/>, `${__dirname}/example.pdf`);
