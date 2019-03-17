import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  titleBlock: {
    width: '100%',
    height: '100%'
  },
  rowTitleBlock: {
    borderTop: 1,
    flex: 1,
    flexDirection: 'row',
    height: '100%'
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
    borderBottom: 0.5,
    padding: 5

  },
  clientDetailsRow: {
    flex: 1,
    flexDirection: 'row',
    padding: 5
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
  <View fixed style={styles.titleBlock}>
    <View fixed style={styles.rowTitleBlock}>
      <View fixed style={styles.companyLogoBlock}>
        <Text fixed>ISLAND CIVIL</Text>
      </View>
      <View fixed style={styles.companyBlock}>
        <View fixed style={styles.companyDetailsRow}><View fixed><Text fixed>Island Civil Ltd.</Text><Text fixed>Victoria, BC</Text></View></View>
        <View fixed style={styles.clientDetailsRow}><View fixed><Text fixed>Prepared for:</Text><Text fixed>{props.client}</Text><Text fixed>{props.clientAddress}</Text></View></View>

      </View>
      <View fixed style={styles.summaryBlock}>
        <View fixed style={styles.checkedByRow}>
          <Text fixed>TESTED BY: STH</Text>
          <Text fixed>CHECKED BY: STH</Text>
          <Text fixed>APPROVED BY: STH</Text>
        </View>
        <View fixed style={styles.totalDepthPagesRow}>
          <Text fixed>START DATE: {props.date}</Text>
          <Text fixed>COMPLETED DATE: {props.date}</Text>
          <Text render={({ pageNumber, totalPages }) => (
            `PAGE: ${pageNumber} OF ${totalPages}`
          )} fixed />
        </View>
      </View>
    </View>
  </View>
)
