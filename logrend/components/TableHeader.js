import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  text: {
    fontSize: 12
  },
  row: {
    flex: 1,
    flexDirection: 'row',
    height: 100,
    borderTop: 1,
    borderBottom: 1
  },
  col: {
    flex: 1,
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center',
    borderRight: 1
  },
  colSm: {
    flexBasis: '20px',
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center',
    borderRight: 1,
    padding: 0,
    margin: 0
  },
  colEnd: {
    flex: 5,
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center'
  },
  col60: {
    flex: 1,
    flexBasis: '60%'
  },
  col40: {
    flex: 1,
    flexBasis: '40%'
  },
  col2: {
    flex: 2
  },
  col5: {
    flex: 5,
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center',
    borderRight: 1
  },
  colRight: {
    flex: 2,
    alignItems: 'right'
  },
  rotated: {
    transform: 'rotate(-90deg)',
    fontSize: 10,
    padding: 0,
    margin: 0
  }
});

export default (props) => (
  <View style={styles.text}>
    <View style={styles.row}>
      <View style={styles.colSm}><Text style={styles.rotated}>DEPTH</Text></View>
      <View style={styles.colSm}><Text style={styles.rotated}>USCS</Text></View>
      <View style={styles.colSm}><Text style={styles.rotated}>SYMBOL</Text></View>

      <View style={styles.col5}><Text>SOIL DESCRIPTION</Text></View>
      <View style={styles.colEnd}><Text>SAMPLES</Text></View>
    </View>
  </View>
)
