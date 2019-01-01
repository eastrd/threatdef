import React from "react";
import ReactDOM from "react-dom";
import DataApp from "./apps/DataApp";
import { BrowserRouter, Route, Switch } from "react-router-dom";

ReactDOM.render(
  <BrowserRouter>
    <Switch>
      <Route exact path="/data" component={DataApp} />
    </Switch>
  </BrowserRouter>,
  document.getElementById("root")
);
