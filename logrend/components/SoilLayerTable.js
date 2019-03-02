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
    height: 451,
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

const paginateLayers = (layers) => {

  if (!layers) {
    return []
  }

  // for now, assume soils/tests/etc are properly ordered.  We can order this array on the calling service.
  // TODO: verify order.

  // pages holds an array of singlePages
  const pages = []

  // each singlePage is an array of objects
  const singlePage = []

  // the current method is to keep pages 0 - 10 m, so we'll keep track of which 10m interval we are on (i*10 to (i+1) * 10, starting at i=0)
  let i = 0

  for (let j = 0; j < layers.length; j++) {
    const item = layers[j]
    if (item.from >= ((i+1) * 10)) { i++ }
    if (pages.length < i+1) { pages.push([]) }

    const template = {
      from: item.from,
      to: item.to,
      desc: item.desc
    }

    // if we need to overflow this block, this object will be inserted into the list of layers.
    // declaring it here lets us store data and decide later if we need to use it,
    // rather than repeating logic elsewhere.
    const next = {}

    // strip text if it will not fit in a layer.
    const line_equivalent = 0.6
    if (item.to - item.from < line_equivalent ||  (i+1)*10 - item.from < line_equivalent) {

      // check if there is space on the next page for this description
      if (item.to - (i+1)*10 > 0.6) {
        next['desc'] = template.desc
      }
      template.desc = ''
    }


    if (item.to > ( (i+1) * 10)) {
      template.to = (i+1) * 10
      pages[i].push(template)
      next['from'] = (i+1) * 10,
      next['to'] = item.to,
      next['desc'] = next['desc'] || '(cont.)'
      
      layers.splice(j+1, 0, next)
      continue
    }
    pages[i].push(template)
  }
  console.log('# of pages:', pages.length)
  return pages || []
}

export default (props) => (
  <View style={styles.wrapper} wrap>
    <View style={styles.text} wrap>

      { paginateLayers(props.soils).map((page, i) => {
        return (
          <View style={styles.soilRowWrapper} break={ i != 0} key={`page${i}`}>
            <View style={styles.soilRow} nobreak>
              <View style={styles.soilColSmLeft}></View>
              <View style={styles.soilColSm}>
                  { Array.from({ length: 9 }).map((j, k) => {
                    return (
                      <View style={styles.rulerSegment} key={`ruler-page${i}-${k}`}><View style={styles.rulerText}><Text>{ (i * 10 + k + 1).toFixed(1) }</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
                    )
                  })}
                  <View style={styles.rulerSegmentLast}><View style={styles.rulerText}><Text>{ ((i+1) * 10).toFixed(1) }</Text></View><View style={styles.rulerHash}><View style={styles.hashMark}></View><View style={styles.hashMarkLast}></View></View></View>
              </View>
              <View style={styles.soilCol}>
              { page.map((soil, j) => {
                  const height = (((Number(soil.to) - Number(soil.from)) / 10 * 452))
    
                  return (
                    <View style={{height: height, display: 'flex', flexDirection: 'row', padding: 5, borderBottom:0.5}} key={`page${i}-layer${j}`}>
                      <Text style={styles.soilText}>{soil.desc}</Text>
                    </View>
                  )  
              }) }
            </View>
            <View style={styles.soilColSample} nobreak>
            </View>
          </View>
        </View>
        )

      }) }
    </View>
  </View>
)
