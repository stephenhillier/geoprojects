import React from 'react';
import ReactPDF, { Document, Page, View, Text } from '@react-pdf/renderer'


const Doc = () => (
  <Document>
    <Page>
      <View>
        <Text>
          Project: Steve's Mansion
        </Text>
      </View>
    </Page>
  </Document>
);
ReactPDF.render(<Doc/>, `${__dirname}/example.pdf`);
