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
  col: {
    flex: 1,
  },
  col2: {
    flex: 3
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

export default () => (
  <View style={styles.col}>
    <View style={styles.sectionWrapper}>
      <View style={styles.soilBreakdownWrapper}>
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
      <View style={styles.col}>
        <Text style={styles.sectionHeading}>Test results</Text>
        <View style={styles.headerRow}>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Opening (mm)</Text></View>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Percent passing</Text></View>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Spec</Text></View>
        </View>
        <View style={styles.tableRow}>
          <View style={styles.col}><Text style={styles.tableCell}>20mm</Text></View>
          <View style={styles.col}><Text style={styles.tableCell}>0%</Text></View>
          <View style={styles.col}><Text style={styles.tableCell}>0% - 5%</Text></View>
        </View>
        <View style={styles.tableRow}>
          <View style={styles.tableCol}><Text style={styles.tableCell}>16mm</Text></View>
          <View style={styles.tableCol}><Text style={styles.tableCell}>10%</Text></View>
          <View style={styles.col}><Text style={styles.tableCell}>0% - 12%</Text></View>
        </View>
      </View>
      <View style={styles.col}>
        <Text style={styles.sectionHeading}>Summary</Text>
        <View style={styles.tableRow}>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Sample:</Text></View>
          <View style={styles.tableCol}><Text>SA-1</Text></View>
        </View>
        <View style={styles.tableRow}>
          <View style={styles.tableCol}><Text style={styles.tableCell}>Source:</Text></View>
          <View style={styles.tableCol}><Text>BH19-1</Text></View>
        </View>
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
