import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  text: {
    fontSize: 10,
    flex: 1,
    flexDirection: 'column'
  },
  row: {
    flex: 1,
    flexDirection: 'row',
  },
  col: {
    flex: 1
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
    textTransform: 'uppercase',
    borderBottom: 1,
  },
});

export default (props) => (
  <View style={styles.row}>
    <View style={styles.text}>
      <View style={styles.titleRow}>
      <View style={styles.row}>
        <View style={styles.col}>
          <Text></Text>
        </View>
        <View style={styles.col3}>
          <Text>BOREHOLE RECORD</Text>
        </View>
        <View style={styles.col}>
          <Text>BH18-1</Text>
        </View>
      </View>

      </View>
      <View style={styles.row}>
        <View style={styles.col60}><Text>CLIENT: {props.client}</Text></View>
        <View style={styles.col40}><Text>PROJECT NO.: {props.projectNumber}</Text></View>
      </View>
      <View style={styles.row}>
        <View style={styles.col60}><Text>PROJECT: {props.project}</Text></View>
        <View style={styles.col40}><Text>LOCATION: {props.location}</Text></View>
      </View>
      <View style={styles.row}>
        <View style={styles.col60}><Text>EASTING: {props.easting}</Text></View>
        <View style={styles.col40}><Text>ELEV: {props.elevation}</Text></View>

      </View>
      <View style={styles.row}>
        <View style={styles.col60}><Text>NORTHING: {props.northing}</Text></View>
        <View style={styles.col40}><Text>DRILLED BY: {props.company}</Text></View>
      </View>
      <View style={styles.row}>
        <View style={styles.col60}><Text>ZONE: {props.zone}</Text></View>
        <View style={styles.col40}><Text>METHOD: {props.method}</Text></View>
      </View>
    </View>

  </View>
)
