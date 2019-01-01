import React, { Component } from "react";
import { Layout } from "antd";
import "antd/dist/antd.css";
import TunnelTable from "../components/tables/TunnelTable";
import "../styles/App.css";
import NavBar from "../components/NavBar";
import CmdTable from "../components/tables/CmdTable";

const { Content } = Layout;

class DataApp extends Component {
  render() {
    return (
      <div>
        <Layout>
          <NavBar />
          <Content>
            <h1>Tunnel Data</h1>
            <TunnelTable pagesize={5} secondsToWait={5} />
          </Content>
          <Content>
            <h1>Command Input Data</h1>
            <CmdTable pagesize={5} secondsToWait={10} />
          </Content>
        </Layout>
      </div>
    );
  }
}

export default DataApp;
