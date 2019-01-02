import React, { Component } from "react";
import "antd/dist/antd.css";
import TunnelTable from "../components/data/TunnelTable";
import "../styles/tables.css";
import NavBar from "../components/NavBar";
import CmdTable from "../components/data/CmdTable";
import { Row, Col } from "antd";

class DataApp extends Component {
  render() {
    return (
      <div>
        <NavBar />
        <Row>
          <Col span={12}>
            <h2>Tunnel Data</h2>
            <TunnelTable pagesize={5} secondsToWait={5} />
          </Col>
          <Col span={12}>
            <h2>Command Input Data</h2>
            <CmdTable pagesize={5} secondsToWait={10} />
          </Col>
        </Row>
      </div>
    );
  }
}

export default DataApp;
