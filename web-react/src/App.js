import React, { Component } from 'react';

import { FocusStyleManager } from "@blueprintjs/core";
import { BrowserRouter } from 'react-router-dom';

import '@blueprintjs/core/lib/css/blueprint.css';
import './App.css';

import Dashboard from './containers/Dashboard'


 
FocusStyleManager.onlyShowFocusOnTabs();

class App extends Component {
  render() {
    return (
      <div>
        <BrowserRouter>
          <Dashboard></Dashboard>
        </BrowserRouter>
      </div>
    );
  }; 
};

export default App;
 