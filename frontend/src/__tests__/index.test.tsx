/* eslint-disable */
import React from "react";
import ReactDOM from "react-dom/client";
import ZoomApp from "../ZoomApp";

const rootMock = {
  render: jest.fn(),
};

jest.mock("react-dom/client", () => ({
  createRoot: jest.fn(() => rootMock),
}));

describe("Root DOM", () => {
  test("Renders App.tsx", () => {
    const root = document.createElement("div");
    root.id = "root";
    document.body.append(root);

    require("../index.tsx");

    expect(ReactDOM.createRoot).toHaveBeenCalledWith(root);
    expect(rootMock.render).toHaveBeenCalledWith(<ZoomApp />);
  });
});
