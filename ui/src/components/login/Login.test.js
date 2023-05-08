import React from "react";
import "@testing-library/jest-dom";
import { fireEvent, render, screen } from "@testing-library/react";
import { loginWithGoogle, loginWithGithub } from "@services/login";
import Login from "./Login";

jest.mock("@services/login");

describe("Login", () => {
    afterEach(() => {
        jest.clearAllMocks();
        jest.resetAllMocks();
    });

    it("should render login component correctly", async () => {
        const { container } = await render(<Login />);

        expect(await screen.findByTestId("login")).toBeInTheDocument();

        expect(container).toMatchSnapshot();
    });

    it("should login with google on click", async () => {
        await render(<Login />);

        const loginWithGoogleButton = await screen.findByTestId(
            "login-with-google"
        );

        await fireEvent.click(loginWithGoogleButton);

        expect(loginWithGoogle).toHaveBeenCalledTimes(1);
    });

    it("should login with github on click", async () => {
        await render(<Login />);

        const loginWithGithubButton = await screen.findByTestId(
            "login-with-github"
        );

        await fireEvent.click(loginWithGithubButton);

        expect(loginWithGithub).toHaveBeenCalledTimes(1);
    });
});
