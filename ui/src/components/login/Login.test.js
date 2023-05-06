import React from "react";
import { render } from "@testing-library/react";
import Login from "./Login";

describe("Login", () => {
    it("should render login component correctly", async () => {
        let tree = await render(<Login />);
        const container = tree.container;

        expect(container).toMatchSnapshot();
    });
});
