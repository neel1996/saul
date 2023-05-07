import React from "react";
import "@testing-library/jest-dom";
import { act, fireEvent, render } from "@testing-library/react";
import { loginWithGoogle, loginWithGithub } from "@services/login";
import Login from "./Login";

jest.mock("@services/login");

describe("Login", () => {
    afterEach(() => {
        jest.clearAllMocks();
        jest.resetAllMocks();
    });

    it("should render login component correctly", async () => {
        let tree = await render(<Login />);
        const { findByTestId, container } = tree;

        expect(await findByTestId("login")).toBeInTheDocument();

        expect(container).toMatchSnapshot();
    });

    it("should login with google on click", async () => {
        let tree = await render(<Login />);
        const { findByTestId } = tree;

        const loginWithGoogleButton = await findByTestId("login-with-google");

        await act(async () => {
            await fireEvent.click(loginWithGoogleButton);
        });

        expect(loginWithGoogle).toHaveBeenCalledTimes(1);
    });

    it("should login with github on click", async () => {
        let tree = await render(<Login />);
        const { findByTestId } = tree;

        const loginWithGithubButton = await findByTestId("login-with-github");

        await act(async () => {
            await fireEvent.click(loginWithGithubButton);
        });

        expect(loginWithGithub).toHaveBeenCalledTimes(1);
    });
});
