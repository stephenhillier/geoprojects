import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  titleRow: {
    marginBottom: 10,
    fontSize: 18,
    textAlign: 'center',
    textTransform: 'uppercase',
    borderBottom: 1,
    flex: 1,
    flexDirection: 'row'
  },
  col: {
    flex: 1
  },
  col2: {
    flex: 3
  }
});

export default () => (
  <View>
    <View style={styles.titleRow}>
      <View style={styles.col}>
        <Text></Text>
      </View>
      <View style={styles.col2}>
        <Text>BOREHOLE RECORD</Text>
      </View>
      <View style={styles.col}>
        <Text>BH18-1</Text>
      </View>
    </View>
  </View>
)
