import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  row: {
    height: 100
  },
  rowTitleBlock: {
    borderTop: 1,
    marginTop: 'auto',
    flex: 1,
    flexDirection: 'row'
  },
  summaryBlock: {
    flex: 1,
    flexDirection: 'column',
    borderLeft: 0.5,
    fontSize: 10
  },
  companyLogoBlock: {
    flex: 1,
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center',
  },
  companyBlock: {
    flex: 1,
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center',
    borderLeft: 0.5
  },
  companyDetailsRow: {
    flex: 1,
    flexDirection: 'row',
    borderBottom: 0.5
  },
  clientDetailsRow: {
    flex: 1,
    flexDirection: 'row'
  },
  checkedByRow: {
    padding: 5,
    borderBottom: 0.5,
    flex: 1
  },
  totalDepthPagesRow: {
    padding: 5,
    flex: 1
  }
})

export default (props) => (
  <View style={styles.row}>
    <View style={styles.rowTitleBlock}>
      <View style={styles.companyLogoBlock}>
        <Text>ISLAND CIVIL</Text>
      </View>
      <View style={styles.companyBlock}>
        <View style={styles.companyDetailsRow}><View><Text>Island Civil Ltd.</Text><Text>123 Main St, Victoria, BC</Text></View></View>
        <View style={styles.clientDetailsRow}><View><Text>Prepared for:</Text><Text>Bigtime Engineering Inc</Text><Text>Esquimalt, BC</Text></View></View>

      </View>
      <View style={styles.summaryBlock}>
        <View style={styles.checkedByRow}>
          <Text>LOGGED BY: STH</Text>
          <Text>CHECKED BY: STH</Text>
          <Text>APPROVED BY: STH</Text>
        </View>
        <View style={styles.totalDepthPagesRow}>
          <Text>DRILLING DATE: {props.date}</Text>
          <Text>TOTAL DEPTH: {props.totalDepth}</Text>
          <Text>PAGE: 1 of 1</Text>
        </View>
      </View>
    </View>
  </View>
)
