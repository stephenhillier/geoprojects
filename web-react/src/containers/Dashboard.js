import React, { Component } from 'react';

import { FocusStyleManager, Card } from "@blueprintjs/core";
import { Route } from 'react-router-dom';

import Sidenav from '../components/Sidenav';
import ProjectSearch from '../components/ProjectSearch';
import Project from './Project/Project'


 
FocusStyleManager.onlyShowFocusOnTabs();

class App extends Component {
  constructor(props) {
    super(props)

    this.changeProjectHandler = this.changeProjectHandler.bind(this)
  }

  state = {
    projectID: null,
    projectName: null,
  }

  changeProjectHandler(id, name) {
    this.setState({
      projectID: id,
      projectName: name
    })
  }

  render() {
    return (
      <div className="dashboard">
        <Sidenav
            className="sidenav"
            projectID={this.state.projectID}
            projectName={this.state.projectName}
            changeProjectHandler={this.changeProjectHandler}
        />
        <Card className={`gutter main`}>
            <div>
              <Route path="/" exact render={() => {
                  return <ProjectSearch changeProjectHandler={this.changeProjectHandler}/>
                }
              }>
              </Route>
              <Route path="/projects/new" exact/>
              <Route path="/projects/:id" render={(props) => {
                  return <Project {...props}  changeProjectHandler={this.changeProjectHandler}/>
                }
              }/>
            </div>
        </Card>
      </div>
    );
  }; 
};

export default App;
 