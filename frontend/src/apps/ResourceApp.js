import React, { Component } from "react";
import "antd/dist/antd.css";
import NavBar from "../components/NavBar";
import CredTable from "../components/data/CredTable";
import { Row, Col } from "antd";

// Craft URLs for static username & password text file
const dictionaryFile = require("../res/dictionary.zip");
const dictionaryFileUrl =
  window.location.protocol +
  "//" +
  window.location.hostname +
  ":" +
  window.location.port +
  dictionaryFile;

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
          <Col span={12}>
            <h3 style={{ textAlign: "center" }}>
              <a href={dictionaryFileUrl} download>
                Download SSH Brute Force Dictionary
              </a>
            </h3>
          </Col>
        </Row>
      </div>
    );
  }
}

export default ResourceApp;
