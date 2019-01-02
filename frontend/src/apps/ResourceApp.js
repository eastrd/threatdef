import React, { Component } from "react";
import "antd/dist/antd.css";
import NavBar from "../components/NavBar";
import CredTable from "../components/data/CredTable";
import { Row, Col } from "antd";

class ResourceApp extends Component {
  render() {
    return (
      <div>
        <NavBar />
        <Row>
          <Col span={12}>
            <h2>Brute Force Combinations</h2>
            <CredTable pagesize={10} secondsToWait={30} />
          </Col>
          <Col span={12}> Under Construction... </Col>
        </Row>
      </div>
    );
  }
}

export default ResourceApp;
