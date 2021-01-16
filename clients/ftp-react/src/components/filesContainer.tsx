import React, { Component } from "react";
import FilesTree from "./filesTree";
export interface FilesContainerProps {}

export interface FilesContainerState {}

class FilesContainer extends React.Component<
  FilesContainerProps,
  FilesContainerState
> {
  state = {};
  render() {
    return (
      <div>
        <h1>Container</h1>
        <div>
          <FilesTree />
        </div>
      </div>
    );
  }
}

export default FilesContainer;
