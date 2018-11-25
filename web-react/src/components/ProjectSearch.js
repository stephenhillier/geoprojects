import React, { Component } from 'react';
import axios from '../axios.js';
import { Link } from 'react-router-dom'

import {
  H1,
  FormGroup,
  InputGroup,
  Button,
  ButtonGroup,
  Classes
} from "@blueprintjs/core";

class ProjectSearch extends Component {

  state = {
    projects: []
  }

  componentDidMount() {
    this.props.changeProjectHandler(null, null)    
    axios.get('projects').then((response) => {
      this.setState({projects: response.data.results})
    }).catch((e) => {
      console.error(e)
    });
  };

  selectProjectHandler(id, name) {
    this.props.changeProjectHandler(id, name)
  }

  render() {
    const projectTableRows = this.state.projects.map((project) => {
      return (
        <tr key={project.id}>
          <td><Link to={'/projects/' + project.id}>{project.id} - {project.name}</Link></td>
          <td>{project.location}</td>
        </tr>
      )
    });
    return (
      <div>
        <H1>Projects</H1>
        <div className="form-row">
          <FormGroup
              helperText="Search for Projects"
              label="Search"
              labelFor="project-search-input"
          >
              <InputGroup id="project-search-input"/>
          </FormGroup>
        </div>
        <table className="bp3-html-table">
          <thead>
            <tr>
              <th>Project</th>
              <th>Location</th>
            </tr>
          </thead>
          <tbody>
            {projectTableRows}
          </tbody>
        </table>
        {/* <ButtonGroup minimal={true}>
          <Button intent="primary" minimal={true}>Previous</Button>
          <Button intent="primary" minimal={true}>Next</Button>
        </ButtonGroup> */}
        <div>
          <Button intent="primary" text="New Project"/>
        </div>
      </div>
    );
  }
};

export default ProjectSearch