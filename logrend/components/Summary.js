import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  text: {
    fontSize: 10,
    flex: 1,
    flexDirection: 'column',
    paddingHorizontal: 10,
    padding: 10,
    borderLeft: 1,
    borderRight: 1,
    borderTop: 1
  },
  row: {
    flex: 1,
    flexDirection: 'row',
  },
  col: {
    flex: 1
  },
  col60: {
    flex: 6,
  },
  col40: {
    flex: 4,
  },
  col3: {
    flex: 3
  },
  colRight: {
    flex: 2,
    alignItems: 'right'
  },
  titleRow: {
    marginBottom: 10,
    fontSize: 18,
    height: 30,
    textAlign: 'center',
    textTransform: 'uppercase'
  },
});

export default (props) => (
  <View style={styles.row} fixed>
    <View style={styles.text} fixed>
      <View style={styles.titleRow} fixed>
      <View style={styles.row} fixed>
        <View style={styles.col} fixed>
          <Text fixed></Text>
        </View>
        <View style={styles.col3} fixed>
          <Text fixed>BOREHOLE RECORD</Text>
        </View>
        <View fixed style={styles.col}>
          <Text fixed>{props.boreholeNum}</Text>
        </View>
      </View>

      </View>
      <View fixed style={styles.row}>
        <View fixed style={styles.col60}><Text fixed>CLIENT: {props.client}</Text></View>
        <View fixed style={styles.col40}><Text fixed>PROJECT NO.: {props.projectNumber}</Text></View>
      </View>
      <View fixed style={styles.row}>
        <View fixed style={styles.col60}><Text fixed>PROJECT: {props.project}</Text></View>
        <View fixed style={styles.col40}><Text fixed>LOCATION: {props.location}</Text></View>
      </View>
      <View fixed style={styles.row}>
        <View fixed style={styles.col60}><Text fixed>EASTING: {props.easting}</Text></View>
        <View fixed style={styles.col40}><Text fixed>ELEV: {props.elevation}</Text></View>

      </View>
      <View fixed style={styles.row}>
        <View fixed style={styles.col60}><Text fixed>NORTHING: {props.northing}</Text></View>
        <View fixed style={styles.col40}><Text fixed>DRILLED BY: {props.company}</Text></View>
      </View>
      <View fixed style={styles.row}>
        <View fixed style={styles.col60}><Text fixed>ZONE: {props.zone}</Text></View>
        <View fixed style={styles.col40}><Text fixed>METHOD: {props.method}</Text></View>
      </View>
    </View>

  </View>
)
