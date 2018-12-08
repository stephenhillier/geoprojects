import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';


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
    to: 10
  }
]

const styles = StyleSheet.create({
  wrapper: {
    flex: 1,
    flexDirection: 'row'
  },
  text: {
    fontSize: 12,
    flex: 1,
    flexDirection: 'column'
  },
  row: {
    flex: 1,
    flexDirection: 'row',
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
    flexBasis: '55px',
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center',
    padding: 0,
    margin: 0,
    width: 20
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
    borderLeft: 1
  },
  colRight: {
    flex: 2,
    alignItems: 'right'
  },
  rotated: {
    transform: 'rotate(-90deg)',
    fontSize: 10,
    flexDirection: 'column',
    padding: 0,
    margin: 0,
  },
  rulerSegment: {
    flexGrow: 1,
    flexDirection: 'row',
    borderBottom: 0.5
  },
  rulerSegmentLast: {
    flexGrow: 1,
    flexDirection: 'row'
  },
  rulerText: {
    fontSize: 10,
    marginTop: 'auto',
    paddingRight: 5,
    flex: 3,
    flexDirection: 'column'
  },
  tableHeaderRow: {
    height: 65,
  },
  soilRow: {
    display: 'flex',
    flexGrow: 1,
    flexDirection: 'row'
  },
  soilColSmLeft: {
    flexBasis: '15px',
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'right',
    padding: 0,
    margin: 0,
  },
  soilColSm: {
    flexBasis: '40px',
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'right',
    padding: 0,
    margin: 0,
  },
  soilCol: {
    flex: 5,
    borderLeft: 1,
    flexDirection: 'column'
  },
  soilText: {
    paddingLeft: 10
  },
  rulerHash: {
    flex: 1,
    flexDirection: 'column'
  },
  hashMark: {
    flex: 1,
    flexDirection: 'row',
    borderBottom: 0.5
  },
  hashMarkLast: {
    flex: 1,
    flexDirection: 'row'
  },
});

export default (props) => (
  <View style={styles.text}>
    <View style={styles.tableHeaderRow}>
      <View style={styles.row}>
        <View style={styles.colSm}>
          <View style={styles.rotated}>
            <Text>DEPTH</Text>
          </View>
        </View>
        <View style={styles.col5}><Text>SOIL DESCRIPTION</Text></View>
        <View style={styles.col5}><Text>SAMPLES</Text></View>
      </View>
    </View>
    <View style={styles.soilRow}>
      <View style={styles.soilColSmLeft}></View>
      <View style={styles.soilColSm}>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>1.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>2.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>3.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>4.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>5.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>6.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>7.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>8.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegment}><View style={styles.rulerText}><Text>9.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
          <View style={styles.rulerSegmentLast}><View style={styles.rulerText}><Text>10.0</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
      </View>
      <View style={styles.soilCol}>
      { soils.map((soil) => {
          const height = (Number(soil.to) - Number(soil.from)) * 10 - 4

          return (
            <View style={{flex: height, flexDirection: 'row', borderBottom:0.5}} key={`${soil.from} ${soil.to}`}>
              <Text style={styles.soilText}>{soil.desc}</Text>
            </View>
          )  
      }) }
      </View>
      <View style={styles.soilCol}>
      </View>
    </View>
  </View>
)
