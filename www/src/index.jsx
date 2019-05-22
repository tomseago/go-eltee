import React from "react";
import ReactDOM from "react-dom";

import App from "./app";

import "typeface-roboto";
// import RpcTest from "./rpc_test";

// class HelloMessage extends React.Component {
//   render() {
//     return <div>Hello {this.props.name}</div>;
//   }
// }

const mountNode = document.getElementById("app");
ReactDOM.render(<App />, mountNode);
