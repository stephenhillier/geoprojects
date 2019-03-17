import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  titleRow: {
    marginBottom: 10,
    fontSize: 18,
    height: 25,
    textAlign: 'center',
    textTransform: 'uppercase'
  },
  sectionWrapper: {
    borderBottom: 1,
    flex: 0,
    flexDirection: 'row'
  },
  sectionWrapper: {
    borderBottom: 1,
    flex: 0,
  },
  soilBreakdownWrapper: {
    flex: 1,
    flexDirection: 'row'
  },
  soilCol: {
    flex: 1,
    textAlign: 'center',
    marginVertical: '5px'
  },
  soilColSampleSummary: {
    flex: 3,
    textAlign: 'left',
    marginLeft: '5px',
    marginVertical: '5px'
  },
  resultsCol: {
    paddingTop: '8px',
    flex: 1,
  },
  col: {
    flex: 1,
  },
  col2: {
    flex: 2
  },
  sectionHeading: {
    textAlign: 'center',
    marginTop: '4px',
    marginBottom: '6px'
  },
  tableRow: {
    flexDirection: 'row'
  },
  tableCol: {
    flex: 1,
  },
  tableCell: {
    textAlign: 'center'
  },
  headerRow: {
    flexDirection: 'row',
    textDecoration: 'underline',
    marginBottom: '5px'
  },
  summaryWrapper: {
    flex: 1,
    flexDirection: 'row'
  },
  row: {
    flexDirection: 'row'
  }
});

const sieves = [
  {
    size: 20,
    units: 'mm',
    passing: 100,
    specUpper: 100,
    specLower: 95
  },
  {
    size: 16,
    units: 'mm',
    passing: 90,
    specUpper: 100,
    specLower: 85
  },
  {
    size: 12,
    units: 'mm',
    passing: 80,
    specUpper: 90,
    specLower: 65
  },
  {
    size: 10,
    units: 'mm',
    passing: 62,
    specUpper: 85,
    specLower: 50
  },
  {
    size: 5,
    units: 'mm',
    passing: 45,
    specUpper: 65,
    specLower: 30
  },
  {
    size: 2,
    units: 'mm',
    passing: 33,
    specUpper: 45,
    specLower: 20
  },
  {
    size: 1,
    units: 'mm',
    passing: 26,
    specUpper: 35,
    specLower: 10
  },
  {
    size: 0.630,
    units: 'mm',
    passing: 16,
    specUpper: 25,
    specLower: 5
  },
  {
    size: 0.315,
    units: 'mm',
    passing: 12,
    specUpper: 20,
    specLower: 0
  },
  {
    size: 0.160,
    units: 'mm',
    passing: 8,
    specUpper: 15,
    specLower: 0
  },
  {
    size: 0.08,
    units: 'mm',
    passing: 4,
    specUpper: 10,
    specLower: 0
  }
]

export default () => (
  <View style={styles.col}>
    <View style={styles.sectionWrapper}>
      <View style={styles.soilBreakdownWrapper}>
        <View style={styles.soilColSampleSummary}>
          <Text>Sample: SA-1</Text>
          <Text>Source: BH19 @ 2m - 5m</Text>
        </View>
        <View style={styles.soilCol}>
          <Text>Fines (silt/clay)</Text>
          <Text>10.0%</Text>
        </View>
        <View style={styles.soilCol}>
          <Text>Sand</Text>
          <Text>10.0%</Text>
        </View>
        <View style={styles.soilCol}>
          <Text>Gravel</Text>
          <Text>10.0%</Text>
        </View>
        <View style={styles.soilCol}>
          <Text>Cobbles</Text>
          <Text>0.0%</Text>
        </View>
      </View>
    </View>
    <View style={styles.summaryWrapper}>
      <View style={styles.resultsCol}>
        <View style={styles.headerRow}>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Opening (mm)</Text></View>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Percent passing</Text></View>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Spec (percent range)</Text></View>
        </View>


        { sieves.map((sieve, i) => {
          return (
            <View style={styles.tableRow} key={`sieveResultRow${i}`}>
              <View style={styles.col}><Text style={styles.tableCell}>{sieve.size}</Text></View>
              <View style={styles.col}><Text style={styles.tableCell}>{sieve.passing}</Text></View>
              <View style={styles.col}><Text style={styles.tableCell}>{sieve.specLower} - {sieve.specUpper}</Text></View>
            </View>
          )
        })}
      </View>
      <View style={styles.resultsCol}>
        <View style={styles.tableRow}>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Cc</Text></View>
          <View style={styles.tableCol}><Text>0.104</Text></View>
        </View>
        <View style={styles.tableRow}>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Cu</Text></View>
          <View style={styles.tableCol}><Text>1.1</Text></View>
        </View>
      </View>
    </View>
  </View>
)
