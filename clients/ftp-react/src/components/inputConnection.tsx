import React, { Component, ChangeEvent } from "react";
import { TextField, Grid, Button } from "@material-ui/core";

export interface InputCnProps {}

export interface InputCnState {}

class InputCn extends Component<InputCnProps, InputCnState> {
  state = {
    host: "54.39.115.191",
    port: 21,
    username: "renechiquete@gmail.com.95750",
    password: "holamundo123.#",
  };
  render() {
    return (
      <Grid direction="row" alignItems="center" justify="center" container>
        <TextField
          id="host"
          variant="filled"
          label="Host"
          type="search"
          value={this.state.host}
          onChange={this.handleChange}
        />

        <TextField
          id="port"
          variant="filled"
          label="Port"
          type="number"
          value={this.state.port}
          onChange={this.handleChange}
        />

        <TextField
          id="username"
          variant="filled"
          label="Username"
          value={this.state.username}
          onChange={this.handleChange}
        />

        <TextField
          id="password"
          variant="filled"
          label="Password"
          value={this.state.password}
          type="password"
          onChange={this.handleChange}
        />

        <Button
          variant="contained"
          color="primary"
          onClick={() => this.connect()}
        >
          Connect
        </Button>
      </Grid>
    );
  }

  handleChange = (event: ChangeEvent<{ value: string; id: string }>) => {
    this.setState({ [event.target.id]: event.target.value });
  };

  connect = () => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    const raw = JSON.stringify({
      host: "54.39.115.191",
      username: "renechiquete@gmail.com.95750",
      password: "holamundo123.#",
    });

    const requestOptions = {
      method: "POST",
      headers: [["Content-Type", "application/json"]],
      body: JSON.stringify(this.state),
    };

    fetch("http://localhost:8088/api/connect", requestOptions)
      .then((response) => response.text())
      .then((result) => console.log(result))
      .catch((error) => console.log("error", error));
  };
}

export default InputCn;
