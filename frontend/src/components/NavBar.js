import React, { Component } from "react";
import { Layout, Menu, Icon } from "antd";
import "antd/dist/antd.css";

const { Header } = Layout;

class NavBar extends Component {
  render() {
    return (
      <Header>
        <Menu
          mode="horizontal"
          theme="dark"
          style={{ lineHeight: "60px" }}
          onClick={({ key }) => {
            switch (key) {
              case "map":
                window.location.href = "/";
                break;
              case "data":
                window.location.href = "/data";
                break;
              case "resource":
                window.location.href = "/resource";
                break;
              default:
                console.log("Error on navbar routing.");
            }
          }}
        >
          <Menu.Item key="map">
            <Icon type="radar-chart" />
            Cyber Map
          </Menu.Item>
          <Menu.Item key="data">
            <Icon type="table" />
            Data
          </Menu.Item>
          <Menu.Item key="resource">
            <Icon type="profile" />
            Resource
          </Menu.Item>
        </Menu>
      </Header>
    );
  }
}

export default NavBar;
