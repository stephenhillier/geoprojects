import React from 'react';
import NavHeader from './NavHeader';
import { Link } from 'react-router-dom';

import {
  Card,
  Button,
  Classes
} from "@blueprintjs/core";

const Sidenav = (props) => {
  return (
    <div className="gutter">
      <NavHeader/>
      <Card>
        <div>
          <Link to="/"><Button minimal={true} large={true} icon="projects" text="Projects Home"/></Link>
          <h3 className={Classes.TEXT_MUTED}>{props.projectName}</h3>
        </div>

        { props.projectID &&
        <div>
          <div>
            <Link className="dashboard-menu-item" to={`/projects/${props.projectID}`}><Button minimal={true} large={true} text="Summary"/></Link>
          </div>
          <div>
            <Link to={`/projects/${props.projectID}/boreholes`}><Button minimal={true} large={true} text="Boreholes"/></Link>
          </div>
          <div>
            <Link to={`/projects/${props.projectID}`}><Button minimal={true} large={true} text="Instrumentation"/></Link>
          </div>
          <div>
            <Link to={`/projects/${props.projectID}`}><Button minimal={true} large={true} text="Samples"/></Link>
          </div>
          <div>
            <Link to={`/projects/${props.projectID}`}><Button minimal={true} large={true} text="Lab testing"/></Link>
          </div>
        </div>
        }
      </Card>
    </div>
  );
};

export default Sidenav
