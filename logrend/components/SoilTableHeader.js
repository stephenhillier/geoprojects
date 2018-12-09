import React from 'react';
import { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  wrapper: {
    flex: 1,
    flexDirection: 'row',
  },
  text: {
    fontSize: 12,
    flex: 1,
    flexDirection: 'column',
  },
  row: {
    flex: 1,
    flexDirection: 'row',
    borderTop: 1,
    borderBottom: 1
  },
  tableHeaderRow: {
    height: 65,
    width: '100%',
    borderLeft: 1,
    borderRight: 1,
  },
  soilDescCol: {
    justifyContent: 'center',
    width: 300,
    textAlign: 'center',
    borderLeft: 1
  },
  col5: {
    flex: 5,
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center',
    borderLeft: 0.5
  },
  colSm: {
    width: 55,
    justifyContent: 'center',
    textAlign: 'center',
    padding: 0,
    margin: 0,
  },
  rotated: {
    transform: 'rotate(-90deg)',
    fontSize: 10,
    flexDirection: 'column',
    padding: 0,
    margin: 0,
  },
})

export default (props) => (
  <View fixed style={styles.tableHeaderRow}>
    <View fixed style={styles.row}>
      <View fixed style={styles.colSm}>
        <View fixed style={styles.rotated}>
          <Text fixed>DEPTH</Text>
        </View>
      </View>
      <View fixed style={styles.soilDescCol}><Text fixed>SOIL DESCRIPTION</Text></View>
      <View fixed style={styles.col5}><Text fixed>SAMPLES</Text></View>
    </View>
  </View>
)