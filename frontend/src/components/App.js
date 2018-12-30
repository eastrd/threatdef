import React, { Component } from "react";
import { Layout, Menu, Icon, Row, Col } from "antd";
import "antd/dist/antd.css";
import TunnelTable from "./api/TunnelTable";
import CmdTable from "./api/CmdTable";
import "../styles/App.css";

const { Header, Footer, Content } = Layout;

class App extends Component {
  render() {
    return (
      <div>
        <Layout>
          <Header>
            <Menu mode="horizontal" theme="dark" style={{ lineHeight: "60px" }}>
              <Menu.Item>
                <Icon type="mail" />
                Intel
              </Menu.Item>
              <Menu.Item>
                <Icon type="appstore" />
                Resource
              </Menu.Item>
            </Menu>
          </Header>
          <Layout style={{ maxHeight: 200 }}>
            <img id="map" src={require("../res/map.jpg")} />
          </Layout>
          <Layout style={{ minHeight: 200 }}>
            <Row gutter={8}>
              <Col span={12}>
                <TunnelTable />
              </Col>
              <Col span={12}>
                <CmdTable />
              </Col>
            </Row>
          </Layout>
        </Layout>
      </div>
    );
  }
}

export default App;
