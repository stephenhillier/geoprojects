import React from 'react';
import ReactPDF, { Text, View, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
  titleRow: {
    marginBottom: 10,
    fontSize: 18,
    height: 25,
    textAlign: 'center',
    textTransform: 'uppercase'
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
    <Text style={styles.titleRow}>Particle Size Distribution</Text>
  </View>
)
