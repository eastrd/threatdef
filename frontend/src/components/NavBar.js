import React, { Component } from "react";
import { Layout, Menu, Icon } from "antd";
import "antd/dist/antd.css";
import "../styles/App.css";

const { Header } = Layout;

class NavBar extends Component {
  render() {
    return (
      <Header>
        <Menu mode="horizontal" theme="dark" style={{ lineHeight: "60px" }}>
          <Menu.Item>
            <Icon type="radar-chart" />
            Cyber Map
          </Menu.Item>
          <Menu.Item>
            <Icon type="table" />
            Data
          </Menu.Item>
          <Menu.Item>
            <Icon type="download" />
            Downloads
          </Menu.Item>
        </Menu>
      </Header>
    );
  }
}

export default NavBar;
