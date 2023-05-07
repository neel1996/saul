import React from "react";
import "@testing-library/jest-dom";
import { render } from "@testing-library/react";
import Login from "./Login";

describe("Login", () => {
    it("should render login component correctly", async () => {
        let tree = await render(<Login />);
        const { findByTestId, container } = tree;

        expect(await findByTestId("login")).toBeInTheDocument();

        expect(container).toMatchSnapshot();
    });
});
