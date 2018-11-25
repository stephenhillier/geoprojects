import React, { Component } from 'react';
import axios from '../../axios.js';
import { Route } from 'react-router-dom';

import {
  H1
} from "@blueprintjs/core";

class Project extends Component {

  state = {
    project: {}
  }

  componentDidMount() {
    axios.get('projects/' + this.props.match.params.id).then((response) => {
      this.setState({project: response.data})
      this.updateProject(response.data.id, response.data.name)
    }).catch((e) => {
      console.error(e)
    });
  };

  updateProject(id, name) {
    this.props.changeProjectHandler(id, name)
  }

  render() {
    return (
      <div>
        <H1>{ this.state.project.name }</H1>
        <Route path="/boreholes" render={() => (<div>Boreholes</div>)}/>
      </div>
    );
  }
};

export default Project