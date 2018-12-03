import React from 'react';
import { Document, Page, View } from '@react-pdf/renderer'
import ReactPDF from '@react-pdf/renderer'


const Doc = () => (
  <Document>
    <Page wrap>
      <View fixed>
        Project: Steve's Mansion
      </View>
    </Page>
  </Document>
);
ReactPDF.render(<Doc/>, `${__dirname}/example.pdf`);
