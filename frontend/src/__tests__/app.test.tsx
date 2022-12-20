import React from "react";
import { render } from "@testing-library/react";

import ZoomApp from "../ZoomApp";

test("renders app component", () => {
  const { container, getByRole } = render(<ZoomApp />);

  expect(getByRole("main")).toBeInTheDocument();
  expect(container).toMatchSnapshot();
});
