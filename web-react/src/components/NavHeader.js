import React from 'react';

import {
  Alignment,
  Button,
  Navbar,
  Classes
} from "@blueprintjs/core";

const NavHeader = () => {
  return (
    <Navbar className={Classes.DARK}>
      <Navbar.Group align={Alignment.LEFT}>
        <Navbar.Heading>
          Earthworks
        </Navbar.Heading>
      </Navbar.Group>
      <Navbar.Group align={Alignment.RIGHT}>
        <Button className={Classes.MINIMAL} icon="cog"/>
        <Button className={Classes.MINIMAL} icon="user"/>
      </Navbar.Group>
    </Navbar>
  );
};

export default NavHeader
