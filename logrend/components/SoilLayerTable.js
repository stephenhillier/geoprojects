import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';



const styles = StyleSheet.create({
  wrapper: {
    width: '100%',
    fontSize: 12,
    borderLeft: 1,
    borderRight: 1

  },
  row: {
    flex: 1,
    flexDirection: 'row',
    width: '100%',
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
    flex: 1,
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center',
    padding: 0,
    margin: 0,
    width: 20
  },
  col60: {
    flex: 6,
  },
  col40: {
    flex: 4,
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
  rulerSegment: {
    display: 'flex',
    height: 45,
    flexDirection: 'row',
    borderBottom: 0.5
  },
  rulerSegmentLast: {
    display: 'flex',
    height: 45,
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
  soilRowWrapper: {
    height: 452,
    width: '100%',

  },
  soilRow: {
    display: 'flex',
    flexGrow: 1,
    flexDirection: 'row',
    width: '100%',
  },
  soilColSmLeft: {

    justifyContent: 'center',
    textAlign: 'right',
    padding: 0,
    margin: 0,
    width: 15,
  },
  soilColSm: {
    justifyContent: 'center',
    width: 40,
    textAlign: 'right',
    padding: 0,
    margin: 0,
  },
  soilCol: {
    width: 300,
    borderLeft: 1,
  },
  soilColSample: {
    flex: 5,
    borderLeft: 0.5,
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
  <View style={styles.wrapper} wrap>
    <View style={styles.text} wrap>

      <View style={styles.soilRowWrapper}>
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
          { props.soils.map((soil) => {
              const height = (((Number(soil.to) - Number(soil.from)) / 10 * 452))

              return (
                <View style={{height: height, display: 'flex', flexDirection: 'row', padding: 5, borderBottom:0.5}} key={`${soil.from} ${soil.to}`}>
                  <Text style={styles.soilText}>{soil.desc}</Text>
                </View>
              )  
          }) }
          </View>
          <View style={styles.soilColSample}>
          </View>
        </View>
      </View>
    </View>
  </View>
)
